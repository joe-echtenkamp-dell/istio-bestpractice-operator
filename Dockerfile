# syntax=docker/dockerfile:1

FROM golang:alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o ./unittest.out

##
## Deploy
##
FROM alpine
RUN apk update 
WORKDIR /

COPY --from=build /app/unittest.out /unittest.out

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot 
USER nonroot

ENTRYPOINT ["/unittest.out"]
