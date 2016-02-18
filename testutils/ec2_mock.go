package testutils

import (
	"net/http"
	"net/http/httptest"
)

// EC2MetadataResponseStub fakes an EC2 metadata response
func EC2MetadataResponseStub() *httptest.Server {
	var resp string
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/latest/meta-data/instance-id":
			resp = "i-12345"
		case "/latest/meta-data/placement/availability-zone":
			resp = "us-west-2a"
		default:
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(resp))
	}))
}

// InvalidEC2MetadataResponseStub represents a response from a non-EC2 environment
// see https://github.com/aws/aws-sdk-go/blob/368825ea31d6fde9a070070e1f8e2762f72140ca/aws/ec2metadata/api_test.go#L86
func InvalidEC2MetadataResponseStub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "RequestError: send request failed", http.StatusBadRequest)
	}))
}
