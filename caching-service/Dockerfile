FROM golang:1.23 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /caching-service cmd/main.go

FROM alpine:3.17
WORKDIR /root/
COPY --from=build /caching-service /usr/local/bin/
EXPOSE 8084
CMD ["caching-service"]
