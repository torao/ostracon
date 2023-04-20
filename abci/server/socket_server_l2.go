//go:build layer2
// +build layer2

package server

import (
	"github.com/Finschia/ostracon/abci/types"
)

func layer2(app types.Application, req *types.Request, responses chan<- *types.Response) {
	switch r := req.Value.(type) {
	case *types.Request_GetAppHash:
		res := app.GetAppHash(*r.GetAppHash)
		responses <- types.ToResponseGetAppHash(res)
	case *types.Request_GenerateFraudProof:
		res := app.GenerateFraudProof(*r.GenerateFraudProof)
		responses <- types.ToResponseGenerateFraudProof(res)
	case *types.Request_VerifyFraudProof:
		res := app.VerifyFraudProof(*r.VerifyFraudProof)
		responses <- types.ToResponseVerifyFraudProof(res)
	default:
		responses <- types.ToResponseException("Unknown request")
	}
}
