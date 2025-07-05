FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/first-debug/lk-auth/cmd/schema-fetcher@latest

COPY . .

RUN /go/bin/schema-fetcher --url https://raw.githubusercontent.com/first-debug/lk-graphql-schemas/master/schemas/user-provider/schema.graphql --output api/graphql/schema.graphql

RUN go generate ./...

RUN CGO_ENABLE=0 go build -ldflags="-w -s" -o /user-provider ./cmd/main.go

FROM alpine:latest

COPY --from=builder /user-provider /user-provider

COPY config/config_local.yml /config/config_local.yml
COPY .env /.env

WORKDIR /

EXPOSE 8080

CMD ["/user-provider"]
