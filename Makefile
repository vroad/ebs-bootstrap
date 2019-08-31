
SHELL := /bin/bash
.DEFAULT_GOAL: ebs-bootstrap-linux-amd64

ebs-bootstrap-linux-amd64: vendor/ *.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ebs-bootstrap-linux-amd64 -ldflags '-s'
