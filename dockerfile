# Use the official Golang base image bookworm
FROM golang:bookworm

# Set the working directory in the container
WORKDIR /app

# Clone the Git repository
RUN git clone https://github.com/NullpointerW/rbnb-pow-mint.git

# Change the working directory to the project folder
WORKDIR /app/rbnb-pow-mint

# Build the project
RUN go build

# run the built binary with docker run -it
CMD ["./rbnb"]
