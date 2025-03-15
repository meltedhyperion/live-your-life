# Use Golang 1.24 as the base image for both build and runtime
FROM golang:1.24 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -o bin/globetrotter-server ./cmd/*.go

# Use the same base image for the final stage
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=build /app/bin/globetrotter-server /app/bin/globetrotter-server

# Copy the .env file into the container
COPY .env .env

# Expose the API port specified in the .env file
EXPOSE ${API_PORT}

# Run the application
CMD ["./bin/globetrotter-server"]
