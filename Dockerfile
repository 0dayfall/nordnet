# Use an official Golang runtime as a parent image
FROM golang:1.23

# Set the working directory in the container to /app
WORKDIR /app

# Add the current directory contents into the container at /app
ADD . /app

# Build the Go app
RUN go build -o main .

# Make port 80 available to the world outside this container
EXPOSE 80

# Run the binary program produced by `go build`
CMD ["./main"]
