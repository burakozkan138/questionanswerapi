# Start from golang base image
FROM golang:alpine as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Configure Go
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# Setup folders
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN swag init -d cmd/,internal/api -g main.go --parseDependency --parseInternal

# Build the Go app
RUN go build -o ./build ./cmd/main.go

# Start a new stage from scratch
FROM alpine:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=builder ./app/build /build
COPY --from=builder ./app/config /config
COPY --from=builder ./app/docs /docs

# Command to run the executable
CMD ["/build"]
