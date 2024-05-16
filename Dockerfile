FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 go build -v -ldflags="-s -w" -o nero ./cmd/nero

FROM alpine:3.19
WORKDIR /data
COPY --from=build /app/nero /app/nero
ENTRYPOINT ["/app/nero"]
