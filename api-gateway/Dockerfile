FROM golang:1.23 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway cmd/main.go

FROM alpine:3.17
WORKDIR /root/
COPY --from=build /api-gateway /usr/local/bin/
EXPOSE 8080
CMD ["api-gateway"]
