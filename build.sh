#! /bin/bash
env GOOS=darwin GOARCH=amd64 go build -o build/ip_addr_server ./src/main.go
