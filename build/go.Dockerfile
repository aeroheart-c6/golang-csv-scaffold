FROM golang:1.20-bookworm

RUN apt-get update

ENV GOPRIVATE=code.in.spdigital.sg
ENV GO111MODULE=on

RUN go install golang.org/x/tools/cmd/goimports@v0.2.0\
 && go install github.com/vektah/dataloaden@v0.3.0
