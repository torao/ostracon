//go:build layer2
// +build layer2

package abcicli

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/Finschia/ostracon/abci/types"
)

func (cli *grpcClient) GetAppHashAsync(params types.RequestGetAppHash) *ReqRes {
	req := types.ToRequestGetAppHash(params)
	res, err := cli.client.GetAppHash(context.Background(), req.GetGetAppHash(), grpc.WaitForReady(true))
	if err != nil {
		cli.StopForError(err)
	}
	return cli.finishAsyncCall(req, &types.Response{Value: &types.Response_GetAppHash{GetAppHash: res}}, func(*types.Response) {})
}

func (cli *grpcClient) GenerateFraudProofAsync(params types.RequestGenerateFraudProof) *ReqRes {
	req := types.ToRequestGenerateFraudProof(params)
	res, err := cli.client.GenerateFraudProof(context.Background(), req.GetGenerateFraudProof(), grpc.WaitForReady(true))
	if err != nil {
		cli.StopForError(err)
	}
	return cli.finishAsyncCall(req, &types.Response{Value: &types.Response_GenerateFraudProof{GenerateFraudProof: res}}, func(*types.Response) {})
}

func (cli *grpcClient) VerifyFraudProofAsync(params types.RequestVerifyFraudProof) *ReqRes {
	req := types.ToRequestVerifyFraudProof(params)
	res, err := cli.client.VerifyFraudProof(context.Background(), req.GetVerifyFraudProof(), grpc.WaitForReady(true))
	if err != nil {
		cli.StopForError(err)
	}
	return cli.finishAsyncCall(req, &types.Response{Value: &types.Response_VerifyFraudProof{VerifyFraudProof: res}}, func(*types.Response) {})
}

func (cli *grpcClient) GetAppHashSync(
	params types.RequestGetAppHash) (*types.ResponseGetAppHash, error) {
	reqres := cli.GetAppHashAsync(params)
	return cli.finishSyncCall(reqres).GetGetAppHash(), cli.Error()
}

func (cli *grpcClient) GenerateFraudProofSync(
	params types.RequestGenerateFraudProof) (*types.ResponseGenerateFraudProof, error) {
	reqres := cli.GenerateFraudProofAsync(params)
	return cli.finishSyncCall(reqres).GetGenerateFraudProof(), cli.Error()
}

func (cli *grpcClient) VerifyFraudProofSync(
	params types.RequestVerifyFraudProof) (*types.ResponseVerifyFraudProof, error) {
	reqres := cli.VerifyFraudProofAsync(params)
	return cli.finishSyncCall(reqres).GetVerifyFraudProof(), cli.Error()
}
