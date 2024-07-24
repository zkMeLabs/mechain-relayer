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

	ListenerPauseTime  = 3 * time.Second
	ErrorRetryInterval = 1 * time.Second
	AssembleInterval   = 500 * time.Millisecond

	TxDelayAlertThreshHold = 300 // in second

	DefaultBscMirrorZkmeSBTRelayerFee    = "1300000000000000" // 0.0013
	DefaultBscMirrorZkmeSBTAckRelayerFee = "250000000000000"  // 0.00025
)
