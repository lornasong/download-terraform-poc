PROJECT = $(shell basename $(CURDIR))

# Update the following variables:
TF_PATH = ""
TF_VERSION = "" # Latest version: https://www.terraform.io/downloads.html
OS = "" # Options: darwin, freebsd, linux, openbsd, solaris, windows
ARCH = "" # Options: amd64, 386, arm

deps:
	go mod download

apply:
	go build -o build/$(PROJECT) && ./build/$(PROJECT) -tfPath=$(TF_PATH) -tfv=$(TF_VERSION) -os=$(OS) -arch=$(ARCH)

destroy:
	go build -o build/$(PROJECT) && ./build/$(PROJECT) -destroy=true -tfPath=$(TF_PATH) -tfv=$(TF_VERSION) -os=$(OS) -arch=$(ARCH)
