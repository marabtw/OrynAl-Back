ARG GO_VERSION=1.22

FROM golang:$GO_VERSION-alpine as deps

WORKDIR /app

COPY ./app/go.mod ./app/go.sum ./

RUN go mod download -x

###########################################################

FROM deps as builder-main
WORKDIR /app
COPY app /app
RUN pwd && go build -o /app/main ./cmd/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder-main /app/main /app/main

COPY ./app/.env /app/.env
COPY ./app/config.yml /app/config.yml

EXPOSE 5000

CMD ["/app/main"]
