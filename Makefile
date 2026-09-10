.PHONY: build test vet fmt lint run clean

build:
	go build -o bin/engageops ./cmd/engageops

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -l .

run: build
	./bin/engageops

clean:
	rm -rf bin
