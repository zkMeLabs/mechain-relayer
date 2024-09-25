package executor

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"sync"
	"time"

	sdkclient "github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zkMeLabs/mechain-relayer/contract/zkmecrosschainupgradeable"
	"github.com/zkMeLabs/mechain-relayer/logging"
)

type GreenfieldClient struct {
	sdkclient.IClient
	ethClient            *ethclient.Client
	zkmeCrossChainClient *zkmecrosschainupgradeable.IZKMECrossChainUpgradeable
	Height               int64
}

type GnfdCompositeClients struct {
	clients []*GreenfieldClient
}

func getEthRPCAddress(rpcAddr string) string {
	u, err := url.Parse(rpcAddr)
	if err != nil {
		return rpcAddr
	}
	host, portStr, _ := net.SplitHostPort(u.Host)
	if portStr == "" {
		return fmt.Sprintf("%s://%s:8545%s", u.Scheme, host, u.Path)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return rpcAddr
	}
	if port != 8545 {
		return fmt.Sprintf("%s://%s:8545%s", u.Scheme, host, u.Path)
	}
	return rpcAddr
}

func NewGnfdCompositClients(rpcAddrs []string, chainId string, account *types.Account, useWebsocket bool, srcZkmeSBTContractAddr string) GnfdCompositeClients {
	clients := make([]*GreenfieldClient, 0)
	for i := 0; i < len(rpcAddrs); i++ {
		sdkClient, err := sdkclient.New(chainId, rpcAddrs[i], sdkclient.Option{DefaultAccount: account, UseWebSocketConn: useWebsocket})
		if err != nil {
			logging.Logger.Errorf("rpc node %s is not available", rpcAddrs[i])
			continue
		}

		ethClient, err := ethclient.Dial(getEthRPCAddress(rpcAddrs[0]))
		if err != nil {
			panic("new eth client error")
		}
		zkmeCrossChainClient, err := zkmecrosschainupgradeable.NewIZKMECrossChainUpgradeable(
			common.HexToAddress(srcZkmeSBTContractAddr),
			ethClient)
		if err != nil {
			panic("new zkmeCrossChain client error")
		}
		clients = append(clients, &GreenfieldClient{
			IClient:              sdkClient,
			ethClient:            ethClient,
			zkmeCrossChainClient: zkmeCrossChainClient,
		})
		if len(clients) == 0 {
			panic("no Greenfield client available")
		}
	}
	return GnfdCompositeClients{
		clients: clients,
	}
}

func (gc *GnfdCompositeClients) GetClient() *GreenfieldClient {
	wg := new(sync.WaitGroup)
	wg.Add(len(gc.clients))
	clientCh := make(chan *GreenfieldClient)
	waitCh := make(chan struct{})
	go func() {
		for _, c := range gc.clients {
			go getClientBlockHeight(clientCh, wg, c)
		}
		wg.Wait()
		close(waitCh)
	}()
	var maxHeight int64
	maxHeightClient := gc.clients[0]
	for {
		select {
		case c := <-clientCh:
			if c.Height > maxHeight {
				maxHeight = c.Height
				maxHeightClient = c
			}
		case <-waitCh:
			return maxHeightClient
		}
	}
}

func getClientBlockHeight(clientChan chan *GreenfieldClient, wg *sync.WaitGroup, client *GreenfieldClient) {
	defer wg.Done()
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	status, err := client.GetStatus(ctxWithTimeout)
	if err != nil {
		return
	}
	latestHeight := status.SyncInfo.LatestBlockHeight
	client.Height = latestHeight
	clientChan <- client
}
