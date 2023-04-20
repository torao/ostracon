//go:build !layer2
// +build !layer2

package types

import (
	"github.com/tendermint/tendermint/abci/types"
)

//go:generate mockery --case underscore --name Application

type CheckTxCallback func(ResponseCheckTx)

// Application is an interface that enables any finite, deterministic state machine
// to be driven by a blockchain-based replication engine via the ABCI.
// All methods take a RequestXxx argument and return a ResponseXxx argument,
// except CheckTx/DeliverTx, which take `tx []byte`, and `Commit`, which takes nothing.
type Application interface {
	// Info/Query Connection
	Info(types.RequestInfo) types.ResponseInfo                // Return application info
	SetOption(types.RequestSetOption) types.ResponseSetOption // Set application option
	Query(types.RequestQuery) types.ResponseQuery             // Query for state

	// Mempool Connection
	CheckTxSync(types.RequestCheckTx) ResponseCheckTx            // Validate a tx for the mempool
	CheckTxAsync(types.RequestCheckTx, CheckTxCallback)          // Asynchronously validate a tx for the mempool
	BeginRecheckTx(RequestBeginRecheckTx) ResponseBeginRecheckTx // Signals the beginning of rechecking
	EndRecheckTx(RequestEndRecheckTx) ResponseEndRecheckTx       // Signals the end of rechecking

	// Consensus Connection
	InitChain(types.RequestInitChain) types.ResponseInitChain // Initialize blockchain w validators/other info from OstraconCore
	BeginBlock(RequestBeginBlock) types.ResponseBeginBlock    // Signals the beginning of a block
	DeliverTx(types.RequestDeliverTx) types.ResponseDeliverTx // Deliver a tx for full processing
	EndBlock(types.RequestEndBlock) types.ResponseEndBlock    // Signals the end of a block, returns changes to the validator set
	Commit() types.ResponseCommit                             // Commit the state and return the application Merkle root hash

	// State Sync Connection
	ListSnapshots(types.RequestListSnapshots) types.ResponseListSnapshots                // List available snapshots
	OfferSnapshot(types.RequestOfferSnapshot) types.ResponseOfferSnapshot                // Offer a snapshot to the application
	LoadSnapshotChunk(types.RequestLoadSnapshotChunk) types.ResponseLoadSnapshotChunk    // Load a snapshot chunk
	ApplySnapshotChunk(types.RequestApplySnapshotChunk) types.ResponseApplySnapshotChunk // Apply a shapshot chunk
}
