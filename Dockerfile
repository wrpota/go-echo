FROM golang:alpine as builder
WORKDIR /app

COPY . .

RUN go mod download && \
    go build -o ./app ./main.go

FROM alpine:latest as prod

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=builder /app/app .
COPY ./configs ./configs

CMD ["./app"]
