package listener

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cometbft/cometbft/votepool"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/zkMeLabs/mechain-relayer/common"
	"github.com/zkMeLabs/mechain-relayer/config"
	"github.com/zkMeLabs/mechain-relayer/contract/crosschain"
	"github.com/zkMeLabs/mechain-relayer/db/dao"
	"github.com/zkMeLabs/mechain-relayer/db/model"
	"github.com/zkMeLabs/mechain-relayer/executor"
	"github.com/zkMeLabs/mechain-relayer/logging"
	"github.com/zkMeLabs/mechain-relayer/metric"
)

type BSCListener struct {
	config             *config.Config
	bscExecutor        *executor.BSCExecutor
	greenfieldExecutor *executor.GreenfieldExecutor
	DaoManager         *dao.DaoManager
	crossChainAbi      abi.ABI
	monitorService     *metric.MetricService
}

func NewBSCListener(cfg *config.Config, bscExecutor *executor.BSCExecutor, gnfdExecutor *executor.GreenfieldExecutor, dao *dao.DaoManager, ms *metric.MetricService) *BSCListener {
	crossChainAbi, err := abi.JSON(strings.NewReader(crosschain.CrosschainMetaData.ABI))
	if err != nil {
		panic("marshal abi error")
	}
	return &BSCListener{
		config:             cfg,
		bscExecutor:        bscExecutor,
		greenfieldExecutor: gnfdExecutor,
		DaoManager:         dao,
		crossChainAbi:      crossChainAbi,
		monitorService:     ms,
	}
}

func (l *BSCListener) StartLoop() {
	for {
		if err := l.poll(); err != nil {
			logging.Logger.Errorf("encounter err, err=%s", err.Error())
			time.Sleep(common.ErrorRetryInterval)
			continue
		}
	}
}

func (l *BSCListener) poll() error {
	latestPolledBlock, err := l.getLatestPolledBlock()
	if err != nil {
		return fmt.Errorf("failed to get latest polled block from db, error: %s", err.Error())
	}
	nextHeight := l.config.BSCConfig.StartHeight
	if (*latestPolledBlock != model.BscBlock{}) {
		latestPolledBlockHeight := latestPolledBlock.Height
		if nextHeight <= latestPolledBlockHeight {
			nextHeight = latestPolledBlockHeight + 1
		}
		var latestBlockHeight uint64
		if l.isOpCrossChain() {
			// currently Get finalized block is not support by OPBNB yet
			latestBlockHeight, err = l.bscExecutor.GetLatestBlockHeightWithRetry()
		} else {
			latestBlockHeight, err = l.bscExecutor.GetLatestFinalizedBlockHeightWithRetry()
		}
		if err != nil {
			logging.Logger.Errorf("failed to get latest finalized blockHeight, error: %s", err.Error())
			return err
		}
		if int64(latestPolledBlockHeight) >= int64(latestBlockHeight)-1 {
			time.Sleep(common.ListenerPauseTime)
			return nil
		}
	}
	if err = l.monitorCrossChainPkgAt(nextHeight, latestPolledBlock); err != nil {
		return err
	}
	return nil
}

func (l *BSCListener) getLatestPolledBlock() (*model.BscBlock, error) {
	return l.DaoManager.BSCDao.GetLatestBlock()
}

func (l *BSCListener) monitorCrossChainPkgAt(nextHeight uint64, latestPolledBlock *model.BscBlock) error {
	nextHeightBlockHeader, err := l.bscExecutor.GetBlockHeaderAtHeight(nextHeight)
	if err != nil {
		return fmt.Errorf("failed to get latest block header, error: %s", err.Error())
	}
	if nextHeightBlockHeader == nil {
		logging.Logger.Infof("BSC Block header at height %d not found", nextHeight)
		return nil
	}
	logging.Logger.Infof("retrieved BSC block header at height=%d", nextHeight)

	logs, err := l.queryCrossChainLogs(nextHeight)
	if err != nil {
		return fmt.Errorf("failed to get logs from block at height=%d, err=%s", nextHeight, err.Error())
	}
	relayPkgs := make([]*model.BscRelayPackage, 0)
	for _, log := range logs {
		logging.Logger.Infof("get log: %d, %s, %s", log.BlockNumber, log.Topics[0].String(), log.TxHash.String())
		relayPkg, err := ParseRelayPackage(&l.crossChainAbi,
			&log, nextHeightBlockHeader.Time,
			sdk.ChainID(l.config.GreenfieldConfig.ChainId),
			sdk.ChainID(l.config.BSCConfig.ChainId),
		)
		if err != nil {
			logging.Logger.Errorf("failed to parse event log, txHash=%s, err=%s", log.TxHash, err.Error())
			continue
		}
		if relayPkg == nil {
			continue
		}
		relayPkgs = append(relayPkgs, relayPkg)
	}

	if err = l.DaoManager.BSCDao.SaveBlockAndBatchPackages(
		&model.BscBlock{
			BlockHash:  nextHeightBlockHeader.Hash().String(),
			ParentHash: nextHeightBlockHeader.ParentHash.String(),
			Height:     nextHeight,
			BlockTime:  int64(nextHeightBlockHeader.Time),
		}, relayPkgs); err != nil {
		return fmt.Errorf("failed to persist block and tx to DB, err=%s", err.Error())
	}
	l.monitorService.SetBSCSavedBlockHeight(nextHeight)
	return nil
}

func (l *BSCListener) queryCrossChainLogs(height uint64) ([]types.Log, error) {
	client := l.bscExecutor.GetEthClient()
	topics := [][]ethcommon.Hash{{l.getCrossChainPackageEventHash()}}
	logs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(height)),
		ToBlock:   big.NewInt(int64(height)),
		Topics:    topics,
		Addresses: []ethcommon.Address{l.getCrossChainContractAddress()},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query cross chain logs, err=%s", err.Error())
	}
	return logs, nil
}

func (l *BSCListener) getCrossChainPackageEventHash() ethcommon.Hash {
	return ethcommon.HexToHash(CrossChainPackageEventHex)
}

func (l *BSCListener) getCrossChainContractAddress() ethcommon.Address {
	return ethcommon.HexToAddress(l.config.RelayConfig.CrossChainContractAddr)
}

func (l *BSCListener) PurgeLoop() {
	ticker := time.NewTicker(PurgeJobInterval)
	for range ticker.C {
		latestBscBlock, err := l.DaoManager.BSCDao.GetLatestBlock()
		if err != nil {
			logging.Logger.Errorf("failed to get latest DB BSC block, err=%s", err.Error())
			continue
		}
		blockHeightThreshHold := int64(latestBscBlock.Height) - NumOfHistoricalBlocks
		if blockHeightThreshHold <= 0 {
			continue
		}
		if err = l.DaoManager.BSCDao.DeleteBlocksBelowHeight(blockHeightThreshHold); err != nil {
			logging.Logger.Errorf("failed to delete Bsc blocks, err=%s", err.Error())
			continue
		}
		exists, err := l.DaoManager.BSCDao.ExistsUnprocessedPackage(blockHeightThreshHold)
		if err != nil || exists {
			continue
		}
		if err = l.DaoManager.BSCDao.DeletePackagesBelowHeightWithLimit(blockHeightThreshHold, DeletionLimit); err != nil {
			logging.Logger.Errorf("failed to delete bsc packages, err=%s", err.Error())
			continue
		}
		var eventType votepool.EventType
		if l.isOpCrossChain() {
			eventType = votepool.FromOpCrossChainEvent
		} else {
			eventType = votepool.FromBscCrossChainEvent
		}
		if err = l.DaoManager.VoteDao.DeleteVotesBelowHeightWithLimit(blockHeightThreshHold, uint32(eventType), DeletionLimit); err != nil {
			logging.Logger.Errorf("failed to delete votes, err=%s", err.Error())
		}
	}
}

func (l *BSCListener) isOpCrossChain() bool {
	return l.config.BSCConfig.IsOpCrossChain()
}
