//go:build grpc

package protos

// This file is only an entrypoint for `go generate`
// When triggered all protocol buffer definitions in the same folder as this file will be compiled with protoc

//go:generate make all
