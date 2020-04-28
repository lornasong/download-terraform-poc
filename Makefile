PROJECT = $(shell basename $(CURDIR))

TF_VERSION = ""
OS = ""
ARCH := ""

deps:
	go mod download

apply:
	go build -o build/$(PROJECT) && ./build/$(PROJECT) -tfPath=$(TF_PATH) -tfv=$(TF_VERSION) -os=$(OS) -arch=$(ARCH)

destroy:
	go build -o build/$(PROJECT) && ./build/$(PROJECT) -destroy=true -tfPath=$(TF_PATH) -tfv=$(TF_VERSION) -os=$(OS) -arch=$(ARCH)
