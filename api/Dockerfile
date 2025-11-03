# Build stage
FROM golang:1.25.3 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd

# Run stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/server .
EXPOSE `
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
