.PHONY: build

DAGPOOL_TARGET=./dagpool
DATANODE_TARGET=./datanode
OBJECTSTORE_TARGET=./objectstore
IAMTOOLS_TARGET=./iam-tools
PLATFORM=darwin

build: clean dagpool datanode objectstore iamtools


dagpool:
	CGO_ENABLED=0 GOOS=${PLATFORM} GOARCH=amd64 go build -ldflags "-s -w" -o ${DAGPOOL_TARGET} ./cmd/dagpool

datanode:
	CGO_ENABLED=0 GOOS=${PLATFORM} GOARCH=amd64 go build -ldflags "-s -w" -o ${DATANODE_TARGET} ./cmd/datanode

objectstore:
	CGO_ENABLED=0 GOOS=${PLATFORM} GOARCH=amd64 go build -ldflags "-s -w" -o ${OBJECTSTORE_TARGET} ./cmd/objectstore

iamtools:
	CGO_ENABLED=0 GOOS=${PLATFORM} GOARCH=amd64 go build -ldflags "-s -w" -o ${IAMTOOLS_TARGET} ./cmd/tools/iam-tools

.PHONY: clean
clean:
	-rm -f ${DAGPOOL_TARGET}
	-rm -f ${DATANODE_TARGET}
	-rm -f ${OBJECTSTORE_TARGET}
	-rm -f ${IAMTOOLS_TARGET}

proto:
	protoc --go_out=./dag/proto --go-grpc_out=./dag/proto dag/proto/*.proto --proto_path=./dag/proto
