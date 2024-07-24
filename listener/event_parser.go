package listener

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/bnb-chain/greenfield-relayer/common"
	"github.com/bnb-chain/greenfield-relayer/contract/zkmecrosschainupgradeable"
	"github.com/bnb-chain/greenfield-relayer/db"
	"github.com/bnb-chain/greenfield-relayer/db/model"
	rtypes "github.com/bnb-chain/greenfield-relayer/types"
)

func ParseRelayPackage(abi *abi.ABI, log *types.Log, timestamp uint64, greenfieldChainId, bscChainId sdk.ChainID) (*model.BscRelayPackage, error) {
	ev, err := parseCrossChainPackageEvent(abi, log)
	if err != nil {
		return nil, err
	}
	if sdk.ChainID(ev.SrcChainId) != bscChainId || sdk.ChainID(ev.DstChainId) != greenfieldChainId {
		return nil, fmt.Errorf("event log's chain id(s) not expected, SrcChainId=%d, DstChainId=%d", ev.SrcChainId, ev.DstChainId)
	}
	var p model.BscRelayPackage
	p.OracleSequence = ev.OracleSequence
	p.PackageSequence = ev.PackageSequence
	p.ChannelId = ev.ChannelId
	p.TxHash = log.TxHash.String()
	p.TxIndex = log.TxIndex
	p.TxTime = int64(timestamp)
	p.UpdatedTime = int64(timestamp)
	p.Height = log.BlockNumber
	p.Status = db.Saved
	p.PayLoad = hex.EncodeToString(ev.Payload)
	return &p, nil
}

func parseCrossChainPackageEvent(abi *abi.ABI, log *types.Log) (*rtypes.CrossChainPackageEvent, error) {
	var ev rtypes.CrossChainPackageEvent

	err := abi.UnpackIntoInterface(&ev, BSCCrossChainPackageEventName, log.Data)
	if err != nil {
		return nil, err
	}
	ev.OracleSequence = big.NewInt(0).SetBytes(log.Topics[1].Bytes()).Uint64()
	ev.PackageSequence = big.NewInt(0).SetBytes(log.Topics[2].Bytes()).Uint64()
	ev.ChannelId = uint8(big.NewInt(0).SetBytes(log.Topics[3].Bytes()).Uint64())
	return &ev, nil
}

func ParseZkmeSBTRelayPackage(abi *abi.ABI, log *types.Log, timestamp uint64, greenfieldChainId, bscChainId sdk.ChainID) (*model.GreenfieldRelayTransaction, error) {
	ev, err := parseZkmeSBTCrossChainPackageEvent(abi, log)
	if err != nil {
		return nil, err
	}
	if sdk.ChainID(ev.SrcChainId) != greenfieldChainId || sdk.ChainID(ev.DestChainId) != bscChainId {
		return nil, fmt.Errorf("zkmesbt event log's chain id(s) not expected, SrcChainId=%d, DstChainId=%d", ev.SrcChainId, ev.DestChainId)
	}
	var p model.GreenfieldRelayTransaction

	p.ChannelId = uint8(ev.ChannelId.Uint64())
	p.SrcChainId = ev.SrcChainId
	p.DestChainId = ev.DestChainId
	p.Sequence = ev.Sequence.Uint64()
	p.PackageType = uint32(sdk.SynCrossChainPackageType)
	p.TxHash = log.TxHash.String()
	p.TxTime = int64(timestamp)
	// Assign RelayerFee based on the chainid of blockchain.
	// if ev.DestChainId == "" {}
	p.RelayerFee = common.DefaultBscMirrorZkmeSBTRelayerFee
	p.AckRelayerFee = common.DefaultBscMirrorZkmeSBTAckRelayerFee
	p.UpdatedTime = time.Now().Unix()
	p.Height = log.BlockNumber
	p.Status = db.Saved
	p.PayLoad = hex.EncodeToString(ev.Payload)
	return &p, nil
}

func parseZkmeSBTCrossChainPackageEvent(abi *abi.ABI, log *types.Log) (*zkmecrosschainupgradeable.KYCDataLibEventData, error) {
	var ev zkmecrosschainupgradeable.KYCDataLibEventData

	err := abi.UnpackIntoInterface(&ev, ZkmeSBTCrossChainPackageEventName, log.Data)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}
