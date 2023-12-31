# Start with the golang base image
FROM golang:alpine as builder

#ENV GO111MODULE=on
ENV config=docker

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container 
WORKDIR /app

# Copy go mod and sum files 
COPY ./ /app

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# # Copy the source from the current directory to the working Directory inside the container 
COPY . .

WORKDIR /app/cmd

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest
LABEL language="golang"
RUN apk --no-cache add ca-certificates

WORKDIR /root/
ENV config=docker
RUN addgroup -g 1001 -S gorunner
RUN adduser -S elf -u 1001

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/cmd/main .
COPY --from=builder /app/config/config-docker.yaml .
copy /project_map.json ./project_map.json

# Expose port 8080 to the outside world
EXPOSE 5000

#Command to run the executable
CMD ["./main"]