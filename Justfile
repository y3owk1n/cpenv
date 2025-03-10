build:
    mkdir -p build
    # Build for darwin-arm64
    env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/y3owk1n/cpenv/cmd.Version=local-build" -trimpath -o ./build/cpenv-darwin-arm64 ./main.go

    # Build for darwin-amd64
    env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/y3owk1n/cpenv/cmd.Version=local-build" -trimpath -o ./build/cpenv-darwin-amd64 ./main.go

    # Build for linux-arm64
    env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/y3owk1n/cpenv/cmd.Version=local-build" -trimpath -o ./build/cpenv-linux-arm64 ./main.go

    # Build for linux-amd64
    env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/y3owk1n/cpenv/cmd.Version=local-build" -trimpath -o ./build/cpenv-linux-amd64 ./main.go

    # Build for windows-amd64
    env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/y3owk1n/cpenv/cmd.Version=local-build" -trimpath -o ./build/cpenv-windows64.exe ./main.go

vet:
	go vet ./...

fmt:
	go fmt ./...

test:
	go test -v ./...
