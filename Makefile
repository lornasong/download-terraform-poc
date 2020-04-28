PROJECT = $(shell basename $(CURDIR))

TF_PATH ?= ""
deps:
	go mod download

apply:
	go build -o build/$(PROJECT) && ./build/$(PROJECT) -tfPath=$(TF_PATH)
