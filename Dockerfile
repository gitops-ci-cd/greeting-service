# This stage only downloads and caches Go modules.
FROM golang:1.23 AS deps
WORKDIR /usr/src/app

# Copy the Go Modules manifests first, then download dependencies.
COPY go.* ./
RUN go mod download

# If you have vendoring:
# RUN go mod vendor

# DO NOT copy the rest of the source here for development.
# Rely on the volume mount instead.

# This stage compiles the binary.
FROM golang:1.23 AS builder
WORKDIR /usr/src/app

# Copy dependencies from the deps stage.
COPY --from=deps /go/pkg /go/pkg

# Now copy the rest of the source code.
COPY . ./

# Create a non-root user and group to run the application.
RUN groupadd -r nonroot && useradd --no-log-init -r -g nonroot nonroot

# Compile the binary. If the code has compilation issues,
# you can comment out or skip this step.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -o /go/bin/main ./cmd/server/main.go

# Use a minimal scratch image for the final container.
FROM scratch
WORKDIR /usr/src/app

# Import user and group info
COPY --from=builder /etc/passwd /etc/group /etc/

# Import CA certificates (if your app makes HTTPS calls).
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary from the builder stage.
COPY --from=builder /go/bin/main ./

# Use the nonroot user to run the application
USER nonroot:nonroot

EXPOSE 50051

# Define the entrypoint for the Docker image.
ENTRYPOINT ["/usr/src/app/main"]
