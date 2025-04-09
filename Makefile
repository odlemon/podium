.PHONY: build run test clean

build:
	go build -o bin/podium cmd/server/main.go

run: build
	./bin/podium

test:
	go test ./...

clean:
	go clean
	-rm -rf bin 2>/dev/null || true
	-if [ -d "bin" ]; then rmdir /s /q bin; fi 2>/dev/null || true