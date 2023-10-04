# Use the official Golang image as the base image
FROM golang:1.21

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Build the Go app
RUN go build -o notify .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./notify"]