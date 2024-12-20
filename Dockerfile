FROM golang:1.23-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o app

FROM scratch

COPY --from=builder /app/app /usr/local/bin/

ENTRYPOINT ["app"]
