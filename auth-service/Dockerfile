FROM golang:1.23 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service cmd/main.go

FROM alpine:3.17
WORKDIR /root/
COPY --from=build /auth-service /usr/local/bin/
EXPOSE 8081
CMD ["auth-service"]
