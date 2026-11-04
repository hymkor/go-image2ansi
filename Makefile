GOOPT:=-ldflags "-s -w"

all :
	go fmt ./...
	go build $(GOOPT) ./cmd/image2ansi
