.PHONY: build

DAGPOOL_TARGET=./dagpool
DATANODE_TARGET=./datanode
OBJECTSTORE_TARGET=./objectstore
IAMTOOLS_TARGET=./iam-tools

#VERSION ?= $(shell git describe --tags)
#TAG_DATANODE ?= "filedag/datanode:$(VERSION)"
TAG_DATANODE ?= "filedag/datanode:latest"
TAG_DAGPOOL ?= "filedag/dagpool:latest"
TAG_OBJECTSTORE ?= "filedag/objectstore:latest"

build: clean dagpool datanode objectstore iamtools


dagpool:
	go build -ldflags "-s -w" -o ${DAGPOOL_TARGET} ./cmd/dagpool

datanode:
	go build -ldflags "-s -w" -o ${DATANODE_TARGET} ./cmd/datanode

objectstore:
	go build -ldflags "-s -w" -o ${OBJECTSTORE_TARGET} ./cmd/objectstore

iamtools:
	go build -ldflags "-s -w" -o ${IAMTOOLS_TARGET} ./cmd/tools/iam-tools

docker-datanode:
	docker build -q --no-cache -t $(TAG_DATANODE) . -f Dockerfile.datanode

docker-dagpool:
	docker build -q --no-cache -t $(TAG_DAGPOOL) . -f Dockerfile.dagpool

docker-objectstore:
	docker build -q --no-cache -t $(TAG_OBJECTSTORE) . -f Dockerfile.objectstore

docker: docker-datanode docker-dagpool docker-objectstore
	docker buildx prune -f

.PHONY: clean
clean:
	-rm -f ${DAGPOOL_TARGET}
	-rm -f ${DATANODE_TARGET}
	-rm -f ${OBJECTSTORE_TARGET}
	-rm -f ${IAMTOOLS_TARGET}

proto:
	protoc --go_out=./dag/proto --go-grpc_out=./dag/proto dag/proto/*.proto --proto_path=./dag/proto
