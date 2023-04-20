//go:build layer2
// +build layer2

package kvstore

import (
	"github.com/Finschia/ostracon/abci/types"
)

func (app *PersistentKVStoreApplication) GetAppHash(
	req types.RequestGetAppHash) types.ResponseGetAppHash {
	return types.ResponseGetAppHash{}
}

func (app *PersistentKVStoreApplication) GenerateFraudProof(
	req types.RequestGenerateFraudProof) types.ResponseGenerateFraudProof {
	return types.ResponseGenerateFraudProof{}
}

func (app *PersistentKVStoreApplication) VerifyFraudProof(
	req types.RequestVerifyFraudProof) types.ResponseVerifyFraudProof {
	return types.ResponseVerifyFraudProof{}
}
