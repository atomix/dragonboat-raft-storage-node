export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ifdef VERSION
RAFT_STORAGE_NODE_VERSION := $(VERSION)
else
RAFT_STORAGE_NODE_VERSION := latest
endif

all: build

build: # @HELP build the source code
build:
	GOOS=linux GOARCH=amd64 go build -o build/_output/atomix-raft-storage-node ./cmd/atomix-raft-storage-node
	GOOS=linux GOARCH=amd64 go build -o build/_output/atomix-raft-storage-driver ./cmd/atomix-raft-storage-driver

test: # @HELP run the unit tests and source code validation
test: build license_check linters
	go test github.com/atomix/atomix-raft-storage-dragonboat/...

coverage: # @HELP generate unit test coverage data
coverage: build linters license_check
	go test github.com/atomix/atomix-raft-storage-dragonboat/pkg/... -coverprofile=coverage.out.tmp -covermode=count
	@cat coverage.out.tmp | grep -v ".pb.go" > coverage.out

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	./build/licensing/boilerplate.py -v

proto: # @HELP build Protobuf/gRPC generated types
proto:
	docker run -it -v `pwd`:/go/src/github.com/atomix/atomix-raft-storage-dragonboat \
		-w /go/src/github.com/atomix/atomix-raft-storage-dragonboat \
		--entrypoint build/bin/compile_protos.sh \
		onosproject/protoc-go:stable

images: # @HELP build Dragonboat Docker images
images: build
	docker build . -f build/atomix-raft-storage-node/Dockerfile -t atomix/atomix-raft-storage-node:${RAFT_STORAGE_NODE_VERSION}
	docker build . -f build/atomix-raft-storage-driver/Dockerfile -t atomix/atomix-raft-storage-driver:${RAFT_STORAGE_NODE_VERSION}

kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image atomix/atomix-raft-storage-node:${RAFT_STORAGE_NODE_VERSION}
	kind load docker-image atomix/atomix-raft-storage-driver:${RAFT_STORAGE_NODE_VERSION}

clean: # @HELP clean build files
	@rm -rf vendor build/_output

push: # @HELP push atomix-raft-storage-node Docker image
	docker push atomix/atomix-raft-storage-node:${RAFT_STORAGE_NODE_VERSION}
	docker push atomix/atomix-raft-storage-driver:${RAFT_STORAGE_NODE_VERSION}
