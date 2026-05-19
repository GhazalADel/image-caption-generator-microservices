FROM golang:1.23-alpine as builder

RUN apk add --no-cache git

WORKDIR /piccap2

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM alpine:3.19

RUN apk add --no-cache ca-certificates curl

WORKDIR /piccap2

COPY --from=builder /piccap2/main .

COPY --from=builder /piccap2/config/.env ./config/.env

EXPOSE 8080

CMD ["./main"]
