//go:build layer2
// +build layer2

package abcicli

import (
	"github.com/Finschia/ostracon/abci/types"
)

func (cli *socketClient) GetAppHashAsync(req types.RequestGetAppHash) *ReqRes {
	return cli.queueRequest(types.ToRequestGetAppHash(req), func(*types.Response) {})
}

func (cli *socketClient) GenerateFraudProofAsync(
	req types.RequestGenerateFraudProof,
) *ReqRes {
	return cli.queueRequest(types.ToRequestGenerateFraudProof(req), func(*types.Response) {})
}

func (cli *socketClient) GetAppHashSync(
	req types.RequestGetAppHash) (*types.ResponseGetAppHash, error) {
	reqres := cli.queueRequest(types.ToRequestGetAppHash(req), func(*types.Response) {})
	if _, err := cli.FlushSync(); err != nil {
		return nil, err
	}
	return reqres.Response.GetGetAppHash(), cli.Error()
}

func (cli *socketClient) GenerateFraudProofSync(
	req types.RequestGenerateFraudProof) (*types.ResponseGenerateFraudProof, error) {
	reqres := cli.queueRequest(types.ToRequestGenerateFraudProof(req), func(*types.Response) {})
	if _, err := cli.FlushSync(); err != nil {
		return nil, err
	}
	return reqres.Response.GetGenerateFraudProof(), cli.Error()
}

func (cli *socketClient) VerifyFraudProofSync(
	req types.RequestVerifyFraudProof) (*types.ResponseVerifyFraudProof, error) {
	reqres := cli.queueRequest(types.ToRequestVerifyFraudProof(req), func(*types.Response) {})
	if _, err := cli.FlushSync(); err != nil {
		return nil, err
	}
	return reqres.Response.GetVerifyFraudProof(), cli.Error()
}
