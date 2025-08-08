FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" cmd/main.go -o spark

FROM alpine:3.22 AS prod

WORKDIR /app

COPY --from=builder /spark .
COPY --from=builder /config.yml .

CMD [ "/app/spark -config config.yml" ]
