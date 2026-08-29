default:
    just --list

build-client: 
    cd client && go build -o ./bin/ ./cmd/client/

build-server:
    cd server && cargo build

build: build-client build-server

test-client:
    cd client && go test ./...

test-server: 
    cd server && cargo test

test: build test-client test-server

format-client: 
    cd client && go fmt ./...

format-server:
    cd server && cargo fmt

format: format-client format-server

lint-client:
    cd client && go vet ./...

lint-server:
    cd server && cargo clippy

lint: lint-client lint-server

run-client: build-client
    cd client && ./bin/client

run-server: build-server
    cd server && cargo run

clean-client:
    cd client && rm -rf ./bin

clean-server:
    cd server && cargo clean

clean: clean-client clean-server
