//go:build !layer2
// +build !layer2

package abcicli

import (
	ocabci "github.com/Finschia/ostracon/abci/types"
	"github.com/Finschia/ostracon/libs/service"
	"github.com/tendermint/tendermint/abci/types"
)

// Client defines an interface for an ABCI client.
// All `Async` methods return a `ReqRes` object.
// All `Sync` methods return the appropriate protobuf ResponseXxx struct and an error.
// Note these are client errors, eg. ABCI socket connectivity issues.
// Application-related errors are reflected in response via ABCI error codes and logs.
type Client interface {
	service.Service

	SetGlobalCallback(GlobalCallback)
	GetGlobalCallback() GlobalCallback
	Error() error

	FlushAsync(ResponseCallback) *ReqRes
	EchoAsync(string, ResponseCallback) *ReqRes
	InfoAsync(types.RequestInfo, ResponseCallback) *ReqRes
	SetOptionAsync(types.RequestSetOption, ResponseCallback) *ReqRes
	DeliverTxAsync(types.RequestDeliverTx, ResponseCallback) *ReqRes
	CheckTxAsync(types.RequestCheckTx, ResponseCallback) *ReqRes
	QueryAsync(types.RequestQuery, ResponseCallback) *ReqRes
	CommitAsync(ResponseCallback) *ReqRes
	InitChainAsync(types.RequestInitChain, ResponseCallback) *ReqRes
	BeginBlockAsync(ocabci.RequestBeginBlock, ResponseCallback) *ReqRes
	EndBlockAsync(types.RequestEndBlock, ResponseCallback) *ReqRes
	BeginRecheckTxAsync(ocabci.RequestBeginRecheckTx, ResponseCallback) *ReqRes
	EndRecheckTxAsync(ocabci.RequestEndRecheckTx, ResponseCallback) *ReqRes
	ListSnapshotsAsync(types.RequestListSnapshots, ResponseCallback) *ReqRes
	OfferSnapshotAsync(types.RequestOfferSnapshot, ResponseCallback) *ReqRes
	LoadSnapshotChunkAsync(types.RequestLoadSnapshotChunk, ResponseCallback) *ReqRes
	ApplySnapshotChunkAsync(types.RequestApplySnapshotChunk, ResponseCallback) *ReqRes

	FlushSync() (*types.ResponseFlush, error)
	EchoSync(string) (*types.ResponseEcho, error)
	InfoSync(types.RequestInfo) (*types.ResponseInfo, error)
	SetOptionSync(types.RequestSetOption) (*types.ResponseSetOption, error)
	DeliverTxSync(types.RequestDeliverTx) (*types.ResponseDeliverTx, error)
	CheckTxSync(types.RequestCheckTx) (*ocabci.ResponseCheckTx, error)
	QuerySync(types.RequestQuery) (*types.ResponseQuery, error)
	CommitSync() (*types.ResponseCommit, error)
	InitChainSync(types.RequestInitChain) (*types.ResponseInitChain, error)
	BeginBlockSync(ocabci.RequestBeginBlock) (*types.ResponseBeginBlock, error)
	EndBlockSync(types.RequestEndBlock) (*types.ResponseEndBlock, error)
	BeginRecheckTxSync(ocabci.RequestBeginRecheckTx) (*ocabci.ResponseBeginRecheckTx, error)
	EndRecheckTxSync(ocabci.RequestEndRecheckTx) (*ocabci.ResponseEndRecheckTx, error)
	ListSnapshotsSync(types.RequestListSnapshots) (*types.ResponseListSnapshots, error)
	OfferSnapshotSync(types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error)
	LoadSnapshotChunkSync(types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error)
	ApplySnapshotChunkSync(types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error)
}
