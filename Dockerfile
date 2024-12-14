# Start from the official Go image to create a build artifact.
FROM golang:1.23 AS builder

WORKDIR /usr/src/app

# Copy the Go Modules manifests and download the dependencies.
# This is done before copying the code to leverage Docker cache layers.
COPY go.* ./
RUN go mod download

# Copy the source code from the current directory to the working directory inside the container.
COPY . ./

# Create a non-root user and group to run the application.
RUN groupadd -r nonroot && useradd --no-log-init -r -g nonroot nonroot

# Build the binary with full module support and without Cgo.
# Compile the binary statically including all dependencies.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -o /go/bin/main ./cmd/server/main.go

# Start from a scratch image, which is an empty container.
FROM scratch

WORKDIR /usr/src/app

# Import the the user and group information
COPY --from=builder /etc/passwd /etc/group /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
# This is important for applications that make external HTTPS calls.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled binary from the first stage.
COPY --from=builder /go/bin/main ./

# Use the nonroot user to run the application
USER nonroot:nonroot

EXPOSE 50051

# Define the entry point for the docker image.
# This is the command that will be run when the container starts.
ENTRYPOINT ["/usr/src/app/main"]
