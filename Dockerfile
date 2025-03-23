FROM golang:1.20 AS builder

WORKDIR /app



# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Install gleece tool
RUN go get github.com/gopher-fleece/gleece
RUN go install github.com/gopher-fleece/gleece

# Copy source code
COPY . .

# Build using gleece
RUN gleece

# Final stage
FROM golang:1.20-alpine

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bin/app .

# The EXPOSE instruction is mainly documentation
# Railway will set and use the PORT environment variable
EXPOSE ${PORT:-8080}

# Run the app
CMD ["./app"]