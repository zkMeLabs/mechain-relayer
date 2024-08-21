package common

import (
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/bnb-chain/greenfield-relayer/types"
)

var (
	RtyAttNum = uint(5)
	RtyAttem  = retry.Attempts(RtyAttNum)
	RtyDelay  = retry.Delay(time.Millisecond * 500)
	RtyErr    = retry.LastErrorOnly(true)
)

const (
	OracleChannelId              types.ChannelId = 0
	ZkmeSBTChannelId             types.ChannelId = 10
	SleepTimeAfterSyncLightBlock                 = 15 * time.Second

	OpBNBChainId    uint64 = 204
	PolygonChainId  uint64 = 137
	ScrollChainId   uint64 = 534352 // 534352 overflows uint16
	LineaChainId    uint64 = 59144
	MantleChainId   uint64 = 5000
	ArbitrumChainId uint64 = 42161
	OptimismChainId uint64 = 10

	ListenerPauseTime  = 3 * time.Second
	ErrorRetryInterval = 1 * time.Second
	AssembleInterval   = 500 * time.Millisecond

	TxDelayAlertThreshHold = 300 // in second

	DefaultBscMirrorZkmeSBTRelayerFee    = "1300000000000000" // 0.0013
	DefaultBscMirrorZkmeSBTAckRelayerFee = "250000000000000"  // 0.00025
)
