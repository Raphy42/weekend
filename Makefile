.PHONY: compile

PROTOC := $(shell which protoc)
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go

CWD := $(shell pwd)
OUT_DIR := $(CWD)/grpc

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install
# it.
$(PROTOC_GEN_GO):
	go install github.com/golang/protobuf/protoc-gen-go@latest

$(OUT_DIR)/api.pb.go: protos/api.proto | $(PROTOC_GEN_GO)
	protoc --go_out=plugins=grpc:. protos/api.proto

clean:
	rm -rf grpc/api

all: $(OUT_DIR)/api.pb.go