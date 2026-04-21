.PHONY: build run test clean

BINARY=recipe-site
RECIPES_DIR?=..

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY) -addr=:8080 -recipes=$(RECIPES_DIR)

test:
	go test ./...

clean:
	rm -f $(BINARY)
