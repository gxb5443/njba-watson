.PHONY: build clean run json proto stop
package = server
GOFLAGS ?= $(GOFLAGS:)

all:clean json build

build:
	@go build $(GOFLAGS)
	
clean:
	@go clean $(GOFLAGS) -i ./...

run:build
	./arnoldb $(GOFLAGS) 

json:
	python db/gen-postgres.py

stop: 
	killall server

proto:
	protoc -I ./proto ./proto/arnoldb.proto --go_out=plugins=grpc:./proto
