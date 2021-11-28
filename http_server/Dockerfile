FROM golang:1.16-alpine AS build
WORKDIR /builder
COPY . .
RUN apk update && go mod download && GOOS=linux GOARCH=amd64 go build -o bin/linux/httpserver

FROM alpine
WORKDIR /bin/linux/httpserver
COPY --from=build builder/bin/linux/httpserver .
EXPOSE 80
ENTRYPOINT ["./httpserver"]