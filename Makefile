BINARY := bin/iupload
PKG := ./cmd/api

.PHONY: build run test tidy clean

build:
	@mkdir -p $(dir $(BINARY))
	go build -o $(BINARY) $(PKG)

run:
	go run $(PKG)

test:
	go test ./...

tidy:
	go mod tidy

clean:
	rm -f $(BINARY)
