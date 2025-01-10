FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /usr/src/app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Install bee tool
RUN go install github.com/beego/bee/v2@latest

# Configure Git to consider the directory as safe
RUN git config --global --add safe.directory /usr/src/app


# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable with the -buildvcs=false flag
CMD ["bee", "run", "-buildvcs=false"]