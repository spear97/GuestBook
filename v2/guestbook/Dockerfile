# Use golang:1.15 as builder image
FROM golang:1.15 as builder

# Install necessary Go packages
RUN go get github.com/codegangsta/negroni
RUN go get github.com/xyproto/simpleredis/v2
RUN go get github.com/gorilla/mux

# Copy the main.go file into the container
COPY main.go .

# Build the Go application
RUN go build main.go

# Start a new stage with Ubuntu 18.04 as base image
FROM ubuntu:18.04

# Copy the built executable from the builder stage to the new image
COPY --from=builder /go/main /app/guestbook

# Add static files to the application directory
ADD public/index.html /app/public/index.html
ADD public/script.js /app/public/script.js
ADD public/style.css /app/public/style.css
ADD public/jquery.min.js /app/public/jquery.min.js

# Set the working directory inside the container
WORKDIR /app

# Define the command to run the application when the container starts
CMD ["./guestbook"]

# Expose port 3000 for the application
EXPOSE 3000
