# syntax=docker/dockerfile:1
##
## Build
##

FROM golang:1.17-buster AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /httpserver ./cmd/appserver


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