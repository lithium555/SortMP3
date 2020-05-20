# FIRST IMAGE - builder
FROM golang:1.13.1 AS builder

ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app

# Copying go.mod / go.sum files into docker first
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY . /app

# 1. Create binary file for package cmd/
RUN cd /app/cmd/ \
        && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/cmd/binary

# SECOND IMAGE
FROM alpine:latest

COPY --from=builder /app/cmd/binary /app/cmd

EXPOSE 8877

CMD /app/cmd