FROM golang:1.23 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /message-service cmd/main.go

FROM alpine:3.17
WORKDIR /root/
COPY --from=build /message-service /usr/local/bin/
EXPOSE 8083
CMD ["message-service"]
