//go:build !layer2
// +build !layer2

package proxy

import (
	"github.com/tendermint/tendermint/abci/types"

	abcicli "github.com/Finschia/ostracon/abci/client"
	ocabci "github.com/Finschia/ostracon/abci/types"
)

type AppConnConsensus interface {
	SetGlobalCallback(abcicli.GlobalCallback)
	Error() error

	InitChainSync(types.RequestInitChain) (*types.ResponseInitChain, error)

	BeginBlockSync(ocabci.RequestBeginBlock) (*types.ResponseBeginBlock, error)
	DeliverTxAsync(types.RequestDeliverTx, abcicli.ResponseCallback) *abcicli.ReqRes
	EndBlockSync(types.RequestEndBlock) (*types.ResponseEndBlock, error)
	CommitSync() (*types.ResponseCommit, error)
}
