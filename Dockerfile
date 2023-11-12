FROM golang:1.20.10 AS build-stage

# Set working dir
WORKDIR /app

# Copy go module and lock
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy source files
COPY *.go ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /match-headers-auth

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /match-headers-auth /match-headers-auth

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/match-headers-auth" ]