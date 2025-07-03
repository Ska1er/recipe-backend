FROM golang:1.24-alpine

ENV GOCACHE=/tmp/go
WORKDIR /var/www/app


RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
