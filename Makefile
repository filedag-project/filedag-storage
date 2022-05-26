.PHONY: build

DAGPOOL_TARGET=./dagpool
DATANODE_TARGET=./datanode
OBJECTSTORE_TARGET=./objectstore

build: clean dagpool datanode objectstore

dagpool:
	go build -ldflags "-s -w" -o ${DAGPOOL_TARGET} ./cmd/dagpool

datanode:
	go build -ldflags "-s -w" -o ${DATANODE_TARGET} ./cmd/datanode

objectstore:
	go build -ldflags "-s -w" -o ${OBJECTSTORE_TARGET} ./cmd/objectstore

.PHONY: clean
clean:
	-rm -f ${DAGPOOL_TARGET}
	-rm -f ${DATANODE_TARGET}
	-rm -f ${OBJECTSTORE_TARGET}

proto:
	protoc --go_out=./dag/proto --go-grpc_out=./dag/proto dag/proto/*.proto --proto_path=./dag/proto