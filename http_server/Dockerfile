FROM golang:1.16-alpine AS build
WORKDIR /httpserver
COPY . .
RUN apk update && go mod download && GOOS=linux GOARCH=amd64 go build -o bin/linux/httpserver

FROM alpine
WORKDIR /httpserver
COPY --from=build /httpserver /httpserver
EXPOSE 800
ENTRYPOINT ["/httpserver"]