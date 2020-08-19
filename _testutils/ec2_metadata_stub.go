package testutils

import (
	"net/http"
	"net/http/httptest"
)

const instanceIdentityDocument = `{
  "devpayProductCodes" : null,
  "marketplaceProductCodes" : [ "1abc2defghijklm3nopqrs4tu" ], 
  "availabilityZone" : "us-west-2a",
  "privateIp" : "10.158.112.84",
  "version" : "2010-08-31",
  "region" : "us-west-2",
  "instanceId" : "i-1234567890abcdef0",
  "billingProducts" : null,
  "instanceType" : "t2.micro",
  "accountId" : "123456789012",
  "pendingTime" : "2015-11-19T16:32:11Z",
  "imageId" : "ami-5fb8c835",
  "kernelId" : "aki-919dcaf8",
  "ramdiskId" : null,
  "architecture" : "x86_64"
}`

// EC2MetadataResponseStub fakes an EC2 metadata response
func EC2MetadataResponseStub() *httptest.Server {
	var resp string
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/latest/api/token":
			resp = "abcdefg"
		case "/latest/dynamic/instance-identity/document":
			resp = instanceIdentityDocument
		case "/latest/meta-data/instance-id":
			resp = "i-12345"
		default:
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(resp))
	}))
}

// InvalidEC2MetadataResponseStub represents a response from a non-EC2 environment
func InvalidEC2MetadataResponseStub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "RequestError: send request failed", http.StatusBadRequest)
	}))
}
