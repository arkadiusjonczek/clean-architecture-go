FROM golang:1.24.4-alpine as builder

WORKDIR /app

COPY . ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/server

FROM alpine:3.22.1

WORKDIR /

COPY --from=builder /app/server /

ENV HTTP_ADDR :8080
EXPOSE 8080

CMD ["/server"]
