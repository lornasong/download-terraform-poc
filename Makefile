PROJECT = $(shell basename $(CURDIR))

# Update the following variables:
TF_PATH = ""

deps:
	go mod download

apply:
	go build -o build/$(PROJECT) && ./build/$(PROJECT) -tfPath=$(TF_PATH)

destroy:
	go build -o build/$(PROJECT) && ./build/$(PROJECT) -destroy=true -tfPath=$(TF_PATH)
