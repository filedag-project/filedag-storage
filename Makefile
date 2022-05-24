.PHONY: build

DAGPOOL_TARGET=./dagpool
DATANODE_TARGET=./datanode
OBJECTSTORE_TARGET=./objectstore

dagpool:
	go build -ldflags "-s -w" -o ${DAGPOOL_TARGET} ./cmd/dagpool

datanode:
	go build -ldflags "-s -w" -o ${DATANODE_TARGET} ./cmd/datanode

objectstore:
	go build -ldflags "-s -w" -o ${OBJECTSTORE_TARGET} ./cmd/objectstore

build: clean dagpool datanode objectstore

.PHONY: clean
clean:
	-rm -f ${DAGPOOL_TARGET}
	-rm -f ${DATANODE_TARGET}
	-rm -f ${OBJECTSTORE_TARGET}

