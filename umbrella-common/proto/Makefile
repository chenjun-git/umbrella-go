all: pb

pb: *.proto
	protoc --proto_path=$(GOPATH)/src:. --go_out=plugins=grpc:. *.proto

.PHONY: pb
