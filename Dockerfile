FROM golang:1.23-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o app

FROM alpine:3.21

RUN apk add --no-cache bash=5.2.37-r0

COPY --from=builder /app/app /usr/local/bin/

COPY ./entrypoint.sh /usr/local/bin/

ENTRYPOINT ["entrypoint.sh"]
