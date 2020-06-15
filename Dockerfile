# Start from the golang base image.
FROM golang:1.13.12 as build
LABEL maintainer="Measurement Lab <support@measurementlab.net>"

WORKDIR /measure-upload

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and
# go.sum files are not changed.
RUN go mod download

# Copy the source from the current directory to the Working Directory inside
# the container.
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build -o measure-upload ./cmd/measure-upload/

# Copy the build output into a minimal alpine image.
FROM alpine:3.12.0

COPY --from=build /measure-upload/measure-upload /measure-upload
# Command to run the executable
CMD ["/measure-upload"]
