FROM golang:1.23-alpine as builder

RUN apk add --no-cache git

WORKDIR /piccap

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM alpine:3.19

RUN apk add --no-cache ca-certificates curl

WORKDIR /piccap

COPY --from=builder /piccap/main .

COPY --from=builder /piccap/config/.env ./config/.env
COPY --from=builder /piccap/config/config.yml ./config/config.yml

EXPOSE 8080

CMD ["./main"]
