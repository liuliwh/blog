# syntax=docker/dockerfile:1
##
## Build
##

FROM golang:1.17-buster AS build
WORKDIR /app
COPY ./go.mod .
COPY ./httpserver ./httpserver
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /httpserver ./httpserver 


##
## Deploy
##
FROM gcr.io/distroless/base-debian10
ENV VERSION="v1.0"
WORKDIR /
COPY --from=build /httpserver /httpserver

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/httpserver"]