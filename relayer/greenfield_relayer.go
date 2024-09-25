package relayer

import (
	"github.com/zkMeLabs/mechain-relayer/assembler"
	"github.com/zkMeLabs/mechain-relayer/executor"
	"github.com/zkMeLabs/mechain-relayer/listener"
	"github.com/zkMeLabs/mechain-relayer/vote"
)

type GreenfieldRelayer struct {
	Listener            *listener.GreenfieldListener
	GreenfieldExecutor  *executor.GreenfieldExecutor
	bscExecutor         *executor.BSCExecutor
	voteProcessor       *vote.GreenfieldVoteProcessor
	greenfieldAssembler *assembler.GreenfieldAssembler
}

func NewGreenfieldRelayer(listener *listener.GreenfieldListener, greenfieldExecutor *executor.GreenfieldExecutor, bscExecutor *executor.BSCExecutor, voteProcessor *vote.GreenfieldVoteProcessor, greenfieldAssembler *assembler.GreenfieldAssembler,
) *GreenfieldRelayer {
	return &GreenfieldRelayer{
		Listener:            listener,
		GreenfieldExecutor:  greenfieldExecutor,
		bscExecutor:         bscExecutor,
		voteProcessor:       voteProcessor,
		greenfieldAssembler: greenfieldAssembler,
	}
}

func (r *GreenfieldRelayer) Start() {
	go r.MonitorEventsLoop()
	go r.SignAndBroadcastLoop()
	go r.CollectVotesLoop()
	go r.AssembleTransactionsLoop()
	go r.UpdateCachedLatestValidatorsLoop()
	go r.PurgeLoop()
}

// MonitorEventsLoop will monitor cross chain events for every block and persist into DB
func (r *GreenfieldRelayer) MonitorEventsLoop() {
	r.Listener.StartLoop()
}

func (r *GreenfieldRelayer) SignAndBroadcastLoop() {
	r.voteProcessor.SignAndBroadcastLoop()
}

func (r *GreenfieldRelayer) CollectVotesLoop() {
	r.voteProcessor.CollectVotesLoop()
}

func (r *GreenfieldRelayer) AssembleTransactionsLoop() {
	r.greenfieldAssembler.AssembleTransactionsLoop()
}

func (r *GreenfieldRelayer) UpdateCachedLatestValidatorsLoop() {
	r.GreenfieldExecutor.UpdateCachedLatestValidatorsLoop() // cache validators queried from greenfield, update it every 1 minute
}

func (r *GreenfieldRelayer) PurgeLoop() {
	r.Listener.PurgeLoop()
}
