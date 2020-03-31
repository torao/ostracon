package types

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/vrf"
)

// PrivValidator defines the functionality of a local Tendermint validator
// that signs votes and proposals, and never double signs.
type PrivValidator interface {
	// TODO: Extend the interface to return errors too. Issue: https://github.com/tendermint/tendermint/issues/3602
	GetPubKey() crypto.PubKey

	SignVote(chainID string, vote *Vote) error
	SignProposal(chainID string, proposal *Proposal) error

	GenerateVRFProof(message []byte) (vrf.Proof, error)
}

//----------------------------------------
// Misc.

type PrivValidatorsByAddress []PrivValidator

func (pvs PrivValidatorsByAddress) Len() int {
	return len(pvs)
}

func (pvs PrivValidatorsByAddress) Less(i, j int) bool {
	return bytes.Compare(pvs[i].GetPubKey().Address(), pvs[j].GetPubKey().Address()) == -1
}

func (pvs PrivValidatorsByAddress) Swap(i, j int) {
	it := pvs[i]
	pvs[i] = pvs[j]
	pvs[j] = it
}

//----------------------------------------
// MockPV

// MockPV implements PrivValidator without any safety or persistence.
// Only use it for testing.
type MockPV struct {
	PrivKey              crypto.PrivKey
	breakProposalSigning bool
	breakVoteSigning     bool
}

func NewMockPV() MockPV {
	return MockPV{ed25519.GenPrivKey(), false, false}
}

// NewMockPVWithParams allows one to create a MockPV instance, but with finer
// grained control over the operation of the mock validator. This is useful for
// mocking test failures.
func NewMockPVWithParams(privKey crypto.PrivKey, breakProposalSigning, breakVoteSigning bool) MockPV {
	return MockPV{privKey, breakProposalSigning, breakVoteSigning}
}

// Implements PrivValidator.
func (pv MockPV) GetPubKey() crypto.PubKey {
	return pv.PrivKey.PubKey()
}

// Implements PrivValidator.
func (pv MockPV) SignVote(chainID string, vote *Vote) error {
	useChainID := chainID
	if pv.breakVoteSigning {
		useChainID = "incorrect-chain-id"
	}
	signBytes := vote.SignBytes(useChainID)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	vote.Signature = sig
	return nil
}

// Implements PrivValidator.
func (pv MockPV) SignProposal(chainID string, proposal *Proposal) error {
	useChainID := chainID
	if pv.breakProposalSigning {
		useChainID = "incorrect-chain-id"
	}
	signBytes := proposal.SignBytes(useChainID)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	proposal.Signature = sig
	return nil
}

// Implements PrivValidator.
func (pv *MockPV) GenerateVRFProof(message []byte) (vrf.Proof, error) {
	privKey, ok := pv.privKey.(ed25519.PrivKeyEd25519)
	if !ok {
		return nil, NewErrUnsupportedKey("ed25519")
	}
	return vrf.Prove(privKey, message)
}

// String returns a string representation of the MockPV.
func (pv MockPV) String() string {
	addr := pv.GetPubKey().Address()
	return fmt.Sprintf("MockPV{%v}", addr)
}

// XXX: Implement.
func (pv MockPV) DisableChecks() {
	// Currently this does nothing,
	// as MockPV has no safety checks at all.
}

type ErroringMockPV struct {
	MockPV
}

var ErroringMockPVErr = errors.New("erroringMockPV always returns an error")

// Implements PrivValidator.
func (pv *ErroringMockPV) SignVote(chainID string, vote *Vote) error {
	return ErroringMockPVErr
}

// Implements PrivValidator.
func (pv *ErroringMockPV) SignProposal(chainID string, proposal *Proposal) error {
	return ErroringMockPVErr
}

// NewErroringMockPV returns a MockPV that fails on each signing request. Again, for testing only.

func NewErroringMockPV() *ErroringMockPV {
	return &ErroringMockPV{MockPV{ed25519.GenPrivKey(), false, false}}
}
