//go:build layer2
// +build layer2

package abcicli

import (
	"github.com/Finschia/ostracon/abci/types"
)

func (app *localClient) GetAppHashAsync(req types.RequestGetAppHash) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.GetAppHash(req)
	return app.callback(
		types.ToRequestGetAppHash(req),
		types.ToResponseGetAppHash(res),
	)
}

func (app *localClient) GenerateFraudProofAsync(req types.RequestGenerateFraudProof) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.GenerateFraudProof(req)
	return app.callback(
		types.ToRequestGenerateFraudProof(req),
		types.ToResponseGenerateFraudProof(res),
	)
}

func (app *localClient) VerifyFraudProofAsync(req types.RequestVerifyFraudProof) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.VerifyFraudProof(req)
	return app.callback(
		types.ToRequestVerifyFraudProof(req),
		types.ToResponseVerifyFraudProof(res),
	)
}

func (app *localClient) GetAppHashSync(
	req types.RequestGetAppHash) (*types.ResponseGetAppHash, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.GetAppHash(req)
	return &res, nil
}

func (app *localClient) GenerateFraudProofSync(
	req types.RequestGenerateFraudProof) (*types.ResponseGenerateFraudProof, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.GenerateFraudProof(req)
	return &res, nil
}

func (app *localClient) VerifyFraudProofSync(
	req types.RequestVerifyFraudProof) (*types.ResponseVerifyFraudProof, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.VerifyFraudProof(req)
	return &res, nil
}
