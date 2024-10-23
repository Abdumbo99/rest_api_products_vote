# Use the official Golang image as the base image
FROM golang:1.21.3-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code
COPY ./ ./ 

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o votes_system ./cmd/api/main.go 

RUN chmod a+x votes_system

# Command to run the executable
CMD ["./votes_system"]

# Expose the port 
EXPOSE 8080
