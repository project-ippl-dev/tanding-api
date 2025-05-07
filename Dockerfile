FROM golang:1.16-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o tanding-api cmd/server/main.go

FROM alpine:3.14
WORKDIR /app
RUN apk add --no-cache tzdata
COPY --from=builder /app/tanding-api .
COPY --from=builder /app/public public
EXPOSE 9000
CMD ["./tanding-api"]
