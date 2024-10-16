package executor

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/light"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"

	relayercommon "github.com/zkMeLabs/mechain-relayer/common"
	"github.com/zkMeLabs/mechain-relayer/config"
	"github.com/zkMeLabs/mechain-relayer/contract/crosschain"
	"github.com/zkMeLabs/mechain-relayer/contract/greenfieldlightclient"
	"github.com/zkMeLabs/mechain-relayer/contract/relayerhub"
	"github.com/zkMeLabs/mechain-relayer/logging"
	"github.com/zkMeLabs/mechain-relayer/metric"
	rtypes "github.com/zkMeLabs/mechain-relayer/types"
)

type BSCClient struct {
	rpcClient             *rpc.Client // for apis eth_getFinalizedBlock and eth_getFinalizedHeader usage, supported by BSC
	ethClient             *ethclient.Client
	crossChainClient      *crosschain.Crosschain
	greenfieldLightClient *greenfieldlightclient.Greenfieldlightclient
	relayerHub            *relayerhub.Relayerhub
	provider              string
	height                uint64
	updatedAt             time.Time
}

func newBSCClients(config *config.Config) []*BSCClient {
	bscClients := make([]*BSCClient, 0)
	for _, provider := range config.BSCConfig.RPCAddrs {
		rpcClient, err := rpc.DialContext(context.Background(), provider)
		if err != nil {
			panic("new rpc client error")
		}
		ethClient, err := ethclient.Dial(provider)
		if err != nil {
			panic("new eth client error")
		}
		greenfieldLightClient, err := greenfieldlightclient.NewGreenfieldlightclient(
			common.HexToAddress(config.RelayConfig.GreenfieldLightClientContractAddr),
			ethClient)
		if err != nil {
			panic("new greenfield light client error")
		}
		crossChainClient, err := crosschain.NewCrosschain(
			common.HexToAddress(config.RelayConfig.CrossChainContractAddr),
			ethClient)
		if err != nil {
			panic("new crossChain client error")
		}
		relayerHub, err := relayerhub.NewRelayerhub(
			common.HexToAddress(config.RelayConfig.RelayerHubContractAddr),
			ethClient)
		if err != nil {
			panic("new relayer hub error")
		}
		bscClients = append(bscClients, &BSCClient{
			rpcClient:             rpcClient,
			ethClient:             ethClient,
			crossChainClient:      crossChainClient,
			greenfieldLightClient: greenfieldLightClient,
			relayerHub:            relayerHub,
			provider:              provider,
			updatedAt:             time.Now(),
		})
	}
	return bscClients
}

type BSCExecutor struct {
	mutex              sync.RWMutex
	GreenfieldExecutor *GreenfieldExecutor
	clientIdx          int
	bscClients         []*BSCClient
	config             *config.Config
	privateKey         *ecdsa.PrivateKey
	txSender           common.Address
	relayers           []rtypes.Validator // cached relayers
	metricService      *metric.MetricService
}

func getBscPrivateKey(cfg *config.BSCConfig) string {
	var privateKey string
	if cfg.KeyType == config.KeyTypeAWSPrivateKey {
		result, err := config.GetSecret(cfg.AWSSecretName, cfg.AWSRegion)
		if err != nil {
			panic(err)
		}
		type AwsPrivateKey struct {
			PrivateKey string `json:"private_key"`
		}
		var awsPrivateKey AwsPrivateKey
		err = json.Unmarshal([]byte(result), &awsPrivateKey)
		if err != nil {
			panic(err)
		}
		privateKey = awsPrivateKey.PrivateKey
	} else {
		privateKey = cfg.PrivateKey
	}
	return privateKey
}

func NewBSCExecutor(cfg *config.Config, metricService *metric.MetricService) *BSCExecutor {
	privKey := viper.GetString(config.FlagConfigPrivateKey)
	if privKey == "" {
		privKey = getBscPrivateKey(&cfg.BSCConfig)
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		panic(err)
	}
	publicKey := ecdsaPrivKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("get public key error")
	}
	txSender := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &BSCExecutor{
		clientIdx:     0,
		bscClients:    newBSCClients(cfg),
		privateKey:    ecdsaPrivKey,
		txSender:      txSender,
		config:        cfg,
		metricService: metricService,
	}
}

func DecodeConsensusState(input []byte) (ConsensusState, error) {
	minimumLength := chainIDLength + heightLength + validatorSetHashLength
	inputLen := uint64(len(input))
	if (inputLen-minimumLength)%singleValidatorBytesLength != 0 {
		return ConsensusState{}, fmt.Errorf("expected input size %d+%d*N, actual input size: %d", minimumLength, singleValidatorBytesLength, inputLen)
	}

	pos := uint64(0)
	chainID := string(bytes.Trim(input[pos:pos+chainIDLength], "\x00"))
	pos += chainIDLength

	height := binary.BigEndian.Uint64(input[pos : pos+heightLength])
	pos += heightLength

	nextValidatorSetHash := input[pos : pos+validatorSetHashLength]
	pos += validatorSetHashLength

	validatorSetLength := (inputLen - minimumLength) / singleValidatorBytesLength
	validatorSetBytes := input[pos:]
	validatorSet := make([]*tmtypes.Validator, 0, validatorSetLength)
	for index := uint64(0); index < validatorSetLength; index++ {
		validatorBytes := validatorSetBytes[singleValidatorBytesLength*index : singleValidatorBytesLength*(index+1)]

		pos = 0
		pubkey := ed25519.PubKey(make([]byte, ed25519.PubKeySize))
		copy(pubkey[:], validatorBytes[:validatorPubkeyLength])
		pos += validatorPubkeyLength

		votingPower := int64(binary.BigEndian.Uint64(validatorBytes[pos : pos+validatorVotingPowerLength]))
		pos += validatorVotingPowerLength

		relayerAddress := make([]byte, relayerAddressLength)
		copy(relayerAddress[:], validatorBytes[pos:pos+relayerAddressLength])
		pos += relayerAddressLength

		relayerBlsKey := make([]byte, relayerBlsKeyLength)
		copy(relayerBlsKey[:], validatorBytes[pos:])

		validator := tmtypes.NewValidator(pubkey, votingPower)
		validator.SetRelayerAddress(relayerAddress)
		validator.SetBlsKey(relayerBlsKey)
		validatorSet = append(validatorSet, validator)
	}

	consensusState := ConsensusState{
		ChainID:              chainID,
		Height:               height,
		NextValidatorSetHash: nextValidatorSetHash,
		ValidatorSet: &tmtypes.ValidatorSet{
			Validators: validatorSet,
		},
	}

	return consensusState, nil
}

func DecodeLightBlockValidationInput(input []byte) (*ConsensusState, error) {
	singleValidatorBytesLength := validatorPubkeyLength + validatorVotingPowerLength + relayerAddressLength + relayerBlsKeyLength
	singleValidatorConsensusBytesLength := chainIDLength + heightLength + validatorSetHashLength + singleValidatorBytesLength
	if uint64(len(input)) < singleValidatorConsensusBytesLength {
		return nil, errors.New("invalid input")
	}

	cs, err := DecodeConsensusState(input)
	if err != nil {
		return nil, err
	}

	return &cs, nil
}

func ApplyLightBlock(cs *ConsensusState, block *tmtypes.LightBlock) (bool, error) {
	if uint64(block.Height) <= cs.Height {
		return false, fmt.Errorf("block height <= consensus height (%d < %d)", block.Height, cs.Height)
	}

	if err := block.ValidateBasic(cs.ChainID); err != nil {
		return false, err
	}

	if cs.Height == uint64(block.Height-1) {
		if !bytes.Equal(cs.NextValidatorSetHash, block.ValidatorsHash) {
			return false, fmt.Errorf("validators hash mismatch, expected: %s, real: %s", cs.NextValidatorSetHash, block.ValidatorsHash)
		}
		err := block.ValidatorSet.VerifyCommitLight(cs.ChainID, block.Commit.BlockID, block.Height, block.Commit)
		if err != nil {
			return false, err
		}
	} else {
		// Ensure that +`trustLevel` (default 1/3) or more of last trusted validators signed correctly.
		err := cs.ValidatorSet.VerifyCommitLightTrusting(cs.ChainID, block.Commit, light.DefaultTrustLevel)
		if err != nil {
			return false, err
		}

		// Ensure that +2/3 of new validators signed correctly.
		//
		// NOTE: this should always be the last check because untrustedVals can be
		// intentionally made very large to DOS the light client. not the case for
		// VerifyAdjacent, where validator set is known in advance.
		err = block.ValidatorSet.VerifyCommitLight(cs.ChainID, block.Commit.BlockID, block.Height, block.Commit)
		if err != nil {
			return false, err
		}
	}

	valSetChanged := !(bytes.Equal(cs.ValidatorSet.Hash(), block.ValidatorsHash))

	// update consensus state
	cs.Height = uint64(block.Height)
	cs.NextValidatorSetHash = block.NextValidatorsHash
	cs.ValidatorSet = block.ValidatorSet

	return valSetChanged, nil
}

// output:
// | validatorSetChanged | empty      | consensusStateBytesLength |  new consensusState |
// | 1 byte              | 23 bytes   | 8 bytes                   |                     |
func EncodeLightBlockValidationResult(validatorSetChanged bool, consensusStateBytes []byte) []byte {
	lengthBytes := make([]byte, validateResultMetaDataLength)
	if validatorSetChanged {
		copy(lengthBytes[:1], []byte{0x01})
	}

	consensusStateBytesLength := uint64(len(consensusStateBytes))
	binary.BigEndian.PutUint64(lengthBytes[validateResultMetaDataLength-uint64TypeLength:], consensusStateBytesLength)

	result := append(lengthBytes, consensusStateBytes...)
	return result
}

func (e *BSCExecutor) SetGreenfieldExecutor(ge *GreenfieldExecutor) {
	e.GreenfieldExecutor = ge
}

func (e *BSCExecutor) GetRpcClient() *rpc.Client {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.bscClients[e.clientIdx].rpcClient
}

func (e *BSCExecutor) GetEthClient() *ethclient.Client {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.bscClients[e.clientIdx].ethClient
}

func (e *BSCExecutor) getCrossChainClient() *crosschain.Crosschain {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.bscClients[e.clientIdx].crossChainClient
}

func (e *BSCExecutor) GetGreenfieldLightClient() *greenfieldlightclient.Greenfieldlightclient {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.bscClients[e.clientIdx].greenfieldLightClient
}

func (e *BSCExecutor) getRelayerHub() *relayerhub.Relayerhub {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.bscClients[e.clientIdx].relayerHub
}

func (e *BSCExecutor) SwitchClient() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.clientIdx++
	if e.clientIdx >= len(e.bscClients) {
		e.clientIdx = 0
	}
	logging.Logger.Infof("switch to provider: %s", e.config.BSCConfig.RPCAddrs[e.clientIdx])
}

func (e *BSCExecutor) GetLatestFinalizedBlockHeightWithRetry() (latestHeight uint64, err error) {
	return e.getLatestBlockHeightWithRetry(e.GetEthClient(), e.GetRpcClient(), true)
}

func (e *BSCExecutor) GetLatestBlockHeightWithRetry() (latestHeight uint64, err error) {
	return e.getLatestBlockHeightWithRetry(e.GetEthClient(), e.GetRpcClient(), false)
}

func (e *BSCExecutor) getLatestBlockHeightWithRetry(ethClient *ethclient.Client, rpcClient *rpc.Client, finalized bool) (latestHeight uint64, err error) {
	return latestHeight, retry.Do(func() error {
		latestHeight, err = e.getLatestBlockHeight(ethClient, rpcClient, finalized)
		return err
	}, relayercommon.RtyAttem,
		relayercommon.RtyDelay,
		relayercommon.RtyErr,
		retry.OnRetry(func(n uint, err error) {
			logging.Logger.Errorf("failed to query latest height, attempt: %d times, max_attempts: %d", n+1, relayercommon.RtyAttNum)
		}))
}

func (e *BSCExecutor) getLatestBlockHeight(client *ethclient.Client, rpcClient *rpc.Client, finalized bool) (uint64, error) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	if finalized {
		return e.getFinalizedBlockHeight(ctxWithTimeout, rpcClient)
	}
	header, err := client.HeaderByNumber(ctxWithTimeout, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

func (e *BSCExecutor) UpdateClientLoop() {
	ticker := time.NewTicker(SleepSecondForUpdateClient * time.Second)
	for range ticker.C {
		logging.Logger.Infof("start to monitor bsc data-seeds healthy")
		for _, bscClient := range e.bscClients {
			if time.Since(bscClient.updatedAt).Seconds() > DataSeedDenyServiceThreshold {
				msg := fmt.Sprintf("data seed %s is not accessable", bscClient.provider)
				logging.Logger.Error(msg)
				config.SendTelegramMessage(e.config.AlertConfig.Identity, e.config.AlertConfig.TelegramBotId,
					e.config.AlertConfig.TelegramChatId, msg)
			}
			var (
				height uint64
				err    error
			)
			if e.config.BSCConfig.IsOpCrossChain() {
				height, err = e.getLatestBlockHeightWithRetry(bscClient.ethClient, bscClient.rpcClient, false)
			} else {
				height, err = e.getLatestBlockHeightWithRetry(bscClient.ethClient, bscClient.rpcClient, true)
			}
			if err != nil {
				logging.Logger.Errorf("get latest block height error, err=%s", err.Error())
				continue
			}
			bscClient.height = height
			bscClient.updatedAt = time.Now()
		}
		highestHeight := uint64(0)
		highestIdx := 0
		for idx := 0; idx < len(e.bscClients); idx++ {
			if e.bscClients[idx].height > highestHeight {
				highestHeight = e.bscClients[idx].height
				highestIdx = idx
			}
		}
		// current client block sync is fall behind, switch to the client with the highest block height
		if e.bscClients[e.clientIdx].height+FallBehindThreshold < highestHeight {
			e.mutex.Lock()
			e.clientIdx = highestIdx
			e.mutex.Unlock()
		}
	}
}

func (e *BSCExecutor) GetBlockHeaderAtHeight(height uint64) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	header, err := e.GetEthClient().HeaderByNumber(ctx, big.NewInt(int64(height)))
	if err != nil {
		return nil, err
	}
	return header, nil
}

// GetNextReceiveSequenceForChannelWithRetry gets the next receive sequence for specified channel from BSC
func (e *BSCExecutor) GetNextReceiveSequenceForChannelWithRetry(channelID rtypes.ChannelId) (sequence uint64, err error) {
	return sequence, retry.Do(func() error {
		sequence, err = e.getNextReceiveSequenceForChannel(channelID)
		return err
	}, relayercommon.RtyAttem,
		relayercommon.RtyDelay,
		relayercommon.RtyErr,
		retry.OnRetry(func(n uint, err error) {
			logging.Logger.Errorf("failed to query receive sequence for channel %d, attempt: %d times, max_attempts: %d", channelID, n+1, relayercommon.RtyAttNum)
		}))
}

func (e *BSCExecutor) getNextReceiveSequenceForChannel(channelID rtypes.ChannelId) (sequence uint64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	callOpts := &bind.CallOpts{
		Pending: true,
		Context: ctx,
	}
	return e.getCrossChainClient().ChannelReceiveSequenceMap(callOpts, uint8(channelID))
}

// GetNextSendSequenceForChannelWithRetry gets the next send oracle sequence from  BSC
func (e *BSCExecutor) GetNextSendSequenceForChannelWithRetry() (sequence uint64, err error) {
	return sequence, retry.Do(func() error {
		sequence, err = e.getNextSendOracleSequence()
		return err
	}, relayercommon.RtyAttem,
		relayercommon.RtyDelay,
		relayercommon.RtyErr,
		retry.OnRetry(func(n uint, err error) {
			logging.Logger.Errorf("failed to query send oracle sequence, attempt: %d times, max_attempts: %d", n+1, relayercommon.RtyAttNum)
		}))
}

func (e *BSCExecutor) getNextSendOracleSequence() (sequence uint64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	callOpts := &bind.CallOpts{
		Pending: true,
		Context: ctx,
	}
	sentOracleSeq, err := e.getCrossChainClient().OracleSequence(callOpts)
	if err != nil {
		return 0, err
	}
	return uint64(sentOracleSeq + 1), nil
}

// GetNextDeliveryOracleSequenceWithRetry gets the next delivery Oracle sequence from Greenfield
func (e *BSCExecutor) GetNextDeliveryOracleSequenceWithRetry(chainId sdk.ChainID) (sequence uint64, err error) {
	return sequence, retry.Do(func() error {
		sequence, err = e.getNextDeliveryOracleSequence(chainId)
		return err
	}, relayercommon.RtyAttem,
		relayercommon.RtyDelay,
		relayercommon.RtyErr,
		retry.OnRetry(func(n uint, err error) {
			logging.Logger.Errorf("failed to query oracle sequence, attempt: %d times, max_attempts: %d", n+1, relayercommon.RtyAttNum)
		}))
}

func (e *BSCExecutor) getNextDeliveryOracleSequence(chainId sdk.ChainID) (uint64, error) {
	sequence, err := e.GreenfieldExecutor.GetNextReceiveOracleSequence(chainId)
	if err != nil {
		return 0, err
	}
	return sequence, nil
}

func (e *BSCExecutor) getTransactor(nonce uint64) (*bind.TransactOpts, error) {
	txOpts, err := bind.NewKeyedTransactorWithChainID(e.privateKey, big.NewInt(int64(e.config.BSCConfig.ChainId)))
	if err != nil {
		return nil, err
	}
	gasPrice, err := e.getGasPrice()
	if err != nil {
		return nil, err
	}
	txOpts.Nonce = big.NewInt(int64(nonce))
	txOpts.Value = big.NewInt(0)
	txOpts.GasLimit = e.config.BSCConfig.GasLimit
	txOpts.GasPrice = big.NewInt(gasPrice.Int64() + 1)
	return txOpts, nil
}

func (e *BSCExecutor) SyncTendermintLightBlock(height uint64) (common.Hash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	lightBlock, err := e.QueryTendermintLightBlockWithRetry(int64(height))
	if err != nil {
		return common.Hash{}, err
	}
	oldcsbts, err := e.GetGreenfieldLightClient().ConsensusStateBytes(nil)
	if err != nil {
		return common.Hash{}, err
	}
	// logging.Logger.Debugf("mechain-contracts ConsensusStateBytes: %s", hex.EncodeToString(oldcsbts))
	cs, err := DecodeLightBlockValidationInput(oldcsbts)
	if err != nil {
		return common.Hash{}, err
	}
	validatorSetChanged, err := ApplyLightBlock(cs, &lightBlock)
	if err != nil {
		return common.Hash{}, err
	}

	consensusStateBytes, err := cs.encodeConsensusState()
	if err != nil {
		return common.Hash{}, err
	}
	// logging.Logger.Debugf("validatorSetChanged: %t, new ConsensusStateBytes: %s", validatorSetChanged, hex.EncodeToString(consensusStateBytes))
	result := EncodeLightBlockValidationResult(validatorSetChanged, consensusStateBytes)
	nonce, err := e.GetEthClient().PendingNonceAt(ctx, e.txSender)
	if err != nil {
		return common.Hash{}, err
	}
	txOpts, err := e.getTransactor(nonce)
	if err != nil {
		return common.Hash{}, err
	}
	tx, err := e.GetGreenfieldLightClient().SyncLightBlock(txOpts, result, height)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (e *BSCExecutor) QueryTendermintLightBlockWithRetry(height int64) (lightBlock tmtypes.LightBlock, err error) {
	return lightBlock, retry.Do(func() error {
		lightBlock, err = e.GreenfieldExecutor.QueryTendermintLightBlock(height)
		return err
	}, relayercommon.RtyAttem,
		relayercommon.RtyDelay,
		relayercommon.RtyErr,
		retry.OnRetry(func(n uint, err error) {
			logging.Logger.Errorf("failed to query tendermint header, attempt: %d times, max_attempts: %d", n+1, relayercommon.RtyAttNum)
		}))
}

func (e *BSCExecutor) QueryLatestTendermintHeaderWithRetry() (lightBlockBts []byte, err error) {
	latestHeigh, err := e.GreenfieldExecutor.GetLatestBlockHeight()
	if err != nil {
		return nil, err
	}
	return lightBlockBts, retry.Do(func() error {
		lightBlock, err := e.GreenfieldExecutor.QueryTendermintLightBlock(int64(latestHeigh))
		protoBlock, err := lightBlock.ToProto()
		if err != nil {
			return err
		}
		lightBlockBts, err = protoBlock.Marshal()
		return err
	}, relayercommon.RtyAttem,
		relayercommon.RtyDelay,
		relayercommon.RtyErr,
		retry.OnRetry(func(n uint, err error) {
			logging.Logger.Errorf("failed to query tendermint header, attempt: %d times, max_attempts: %d", n+1, relayercommon.RtyAttNum)
		}))
}

func (e *BSCExecutor) GetNonce() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	return e.GetEthClient().PendingNonceAt(ctx, e.txSender)
}

func (e *BSCExecutor) CallBuildInSystemContract(blsSignature []byte, validatorSet *big.Int, msgBytes []byte, nonce uint64) (common.Hash, error) {
	txOpts, err := e.getTransactor(nonce)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := e.getCrossChainClient().HandlePackage(txOpts, msgBytes, blsSignature, validatorSet)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

// QueryLatestValidators used for gnfd -> bsc
func (e *BSCExecutor) QueryLatestValidators() ([]rtypes.Validator, error) {
	relayerAddresses, err := e.GetGreenfieldLightClient().GetRelayers(nil)
	if err != nil {
		return nil, err
	}
	blsKeys, err := e.GetGreenfieldLightClient().BlsPubKeys(nil)
	if err != nil {
		return nil, err
	}
	relayers := make([]rtypes.Validator, len(relayerAddresses))
	nextRelayerBtsStartIdx := 0

	for i, addr := range relayerAddresses {
		r := rtypes.Validator{
			RelayerAddress: addr,
			BlsPublicKey:   blsKeys[nextRelayerBtsStartIdx : nextRelayerBtsStartIdx+RelayerBytesLength][:],
		}
		nextRelayerBtsStartIdx = nextRelayerBtsStartIdx + RelayerBytesLength
		relayers[i] = r
	}
	return relayers, nil
}

// QueryCachedLatestValidators Used for gnfd -> bsc
func (e *BSCExecutor) QueryCachedLatestValidators() ([]rtypes.Validator, error) {
	if len(e.relayers) != 0 {
		return e.relayers, nil
	}
	return e.QueryLatestValidators()
}

func (e *BSCExecutor) UpdateCachedLatestValidatorsLoop() {
	ticker := time.NewTicker(UpdateCachedValidatorsInterval)
	for range ticker.C {
		relayers, err := e.QueryLatestValidators()
		if err != nil {
			logging.Logger.Errorf("update latest bsc relayers error, err=%s", err)
			continue
		}
		e.relayers = relayers
	}
}

func (e *BSCExecutor) GetLightClientLatestHeight() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	callOpts := &bind.CallOpts{
		Pending: true,
		Context: ctx,
	}
	latestHeight, err := e.GetGreenfieldLightClient().GnfdHeight(callOpts)
	if err != nil {
		return 0, err
	}
	return latestHeight, err
}

func (e *BSCExecutor) GetValidatorsBlsPublicKey() ([]string, error) {
	validators, err := e.QueryCachedLatestValidators()
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, v := range validators {
		keys = append(keys, hex.EncodeToString(v.BlsPublicKey[:]))
	}
	return keys, nil
}

func (e *BSCExecutor) GetInturnRelayer() (*rtypes.InturnRelayer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	callOpts := &bind.CallOpts{
		Pending: true,
		Context: ctx,
	}
	r, err := e.GetGreenfieldLightClient().GetInturnRelayer(callOpts)
	if err != nil {
		return nil, err
	}

	return &rtypes.InturnRelayer{
		BlsPublicKey: hex.EncodeToString(r.BlsKey),
		Start:        r.Start.Uint64(),
		End:          r.End.Uint64(),
	}, nil
}

func (e *BSCExecutor) getRelayerBalance() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	return e.GetEthClient().BalanceAt(ctx, e.txSender, nil)
}

func (e *BSCExecutor) claimReward() (common.Hash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	nonce, err := e.GetEthClient().PendingNonceAt(ctx, e.txSender)
	if err != nil {
		return common.Hash{}, err
	}
	txOpts, err := e.getTransactor(nonce)
	if err != nil {
		return common.Hash{}, err
	}
	txResp, err := e.getRelayerHub().ClaimReward(txOpts, e.txSender)
	if err != nil {
		return common.Hash{}, err
	}
	return txResp.Hash(), nil
}

func (e *BSCExecutor) getRewardBalance() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	callOpts := &bind.CallOpts{
		Pending: true,
		Context: ctx,
	}
	return e.getRelayerHub().RewardMap(callOpts, e.txSender)
}

// ClaimRewardLoop relayer would claim the reward if its balance is below 1BNB and the Reward is over 0.1BNB.
// if after refilled with the rewards, its balance is still lower than 1BNB, it will keep alerting
func (e *BSCExecutor) ClaimRewardLoop() {
	ticker := time.NewTicker(ClaimRewardInterval)
	for range ticker.C {
		logging.Logger.Info("starting claiming rewards loop")
		balance, err := e.getRelayerBalance()
		if err != nil {
			logging.Logger.Errorf("failed to get relayer balance err=%s", err.Error())
			continue
		}
		logging.Logger.Infof("current relayer balance is %v", balance)
		balance.Div(balance, BNBDecimal)

		e.metricService.SetBSCBalance(float64(balance.Int64()))

		// should not claim if balance > 1 BNB
		if balance.Cmp(BSCBalanceThreshold) > 0 {
			continue
		}
		rewardBalance, err := e.getRewardBalance()
		if err != nil {
			logging.Logger.Errorf("failed to get relayer reward balance err=%s", err.Error())
			continue
		}
		logging.Logger.Infof("current relayer reward balance is %v", balance)
		if rewardBalance.Cmp(BSCRewardThreshold) <= 0 {
			continue
		}
		// > 0.1 BNB
		txHash, err := e.claimReward()
		if err != nil {
			logging.Logger.Errorf("failed to claim reward, txHash=%s, err=%s", txHash, err.Error())
		}
		logging.Logger.Infof("claimed rewards, txHash is %s", txHash)
	}
}

// getFinalizedBlockHeight gets the finalizedBlockHeight, which is the larger one between (fastFinalizedBlockHeight, NumberOfBlocksForFinality from config).
func (e *BSCExecutor) getFinalizedBlockHeight(ctx context.Context, rpcClient *rpc.Client) (uint64, error) {
	var head *types.Header
	if err := rpcClient.CallContext(ctx, &head, "eth_getFinalizedHeader", e.config.BSCConfig.NumberOfBlocksForFinality); err != nil {
		return 0, err
	}
	if head == nil || head.Number == nil {
		return 0, ethereum.NotFound
	}
	return head.Number.Uint64(), nil
}

func (e *BSCExecutor) getGasPrice() (*big.Int, error) {
	var (
		gasPrice *big.Int
		err      error
	)
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	if e.config.BSCConfig.GasPrice == 0 {
		gasPrice, err = e.GetEthClient().SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		gasPrice = big.NewInt(int64(e.config.BSCConfig.GasPrice))
	}
	return gasPrice, nil
}
