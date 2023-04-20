//go:build !layer2
// +build !layer2

package server

import "github.com/Finschia/ostracon/abci/types"

func layer2(app types.Application, req *types.Request, responses chan<- *types.Response) {
	responses <- types.ToResponseException("Unknown request")
}
