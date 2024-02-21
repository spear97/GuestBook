# Use the golang:1.15 image as a builder stage
FROM golang:1.15 as builder

# Run commands to get necessary Go packages
RUN go get github.com/codegangsta/negroni
RUN go get github.com/gorilla/mux 
RUN go get github.com/xyproto/simpleredis/v2

# Copy the main.go file into the container
COPY main.go .

# Build the Go application
RUN go build main.go

# Start a new stage from the ubuntu:18.04 image
FROM ubuntu:18.04

# Copy the compiled binary from the builder stage to /app/guestbook in the new image
COPY --from=builder /go//main /app/guestbook

# Copy the HTML, JavaScript, CSS, and jQuery files to /app/public/ in the new image
COPY public/index.html /app/public/index.html
COPY public/script.js /app/public/script.js
COPY public/style.css /app/public/style.css
COPY public/jquery.min.js /app/public/jquery.min.js

# Set the working directory to /app
WORKDIR /app

# Command to run the application when the container starts
CMD ["./guestbook"]

# Expose port 3000 to allow outside connections
EXPOSE 3000
