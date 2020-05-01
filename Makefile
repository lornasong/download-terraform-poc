PROJECT = $(shell basename $(CURDIR))

# Update the following variables:
TF_PATH = ""
OS = "" # Options: darwin, freebsd, linux, openbsd, solaris, windows
ARCH = "" # Options: amd64, 386, arm

deps:
	go mod download

apply:
	go build -o build/$(PROJECT) && ./build/$(PROJECT) -tfPath=$(TF_PATH) -os=$(OS) -arch=$(ARCH)

destroy:
	go build -o build/$(PROJECT) && ./build/$(PROJECT) -destroy=true -tfPath=$(TF_PATH) -os=$(OS) -arch=$(ARCH)
