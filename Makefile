PROJECT = $(shell basename $(CURDIR))

deps:
	go mod download

apply:
	go build -o build/$(PROJECT) && ./build/$(PROJECT)
