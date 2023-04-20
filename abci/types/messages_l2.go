//go:build layer2
// +build layer2

package types

func ToRequestGetAppHash(req RequestGetAppHash) *Request {
	return &Request{
		Value: &Request_GetAppHash{&req},
	}
}

func ToRequestGenerateFraudProof(req RequestGenerateFraudProof) *Request {
	return &Request{
		Value: &Request_GenerateFraudProof{&req},
	}
}

func ToRequestVerifyFraudProof(req RequestVerifyFraudProof) *Request {
	return &Request{
		Value: &Request_VerifyFraudProof{&req},
	}
}

func ToResponseGetAppHash(res ResponseGetAppHash) *Response {
	return &Response{
		Value: &Response_GetAppHash{&res},
	}
}

func ToResponseGenerateFraudProof(res ResponseGenerateFraudProof) *Response {
	return &Response{
		Value: &Response_GenerateFraudProof{&res},
	}
}

func ToResponseVerifyFraudProof(res ResponseVerifyFraudProof) *Response {
	return &Response{
		Value: &Response_VerifyFraudProof{&res},
	}
}
