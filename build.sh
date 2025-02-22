#!/bin/bash

# Linux
GOOS=linux GOARCH=amd64 go build -o bin/vanity_bitcoin_linux_amd64 main.go
echo "Собран бинарник для Linux"

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/vanity_bitcoin.exe main.go
echo "Собран бинарник для Windows"

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/vanity_bitcoin_macos main.go
echo "Собран бинарник для macOS"
