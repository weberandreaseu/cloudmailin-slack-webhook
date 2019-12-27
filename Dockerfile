FROM golang:1-alpine as build

WORKDIR /go/src/app

RUN apk add --virtual buildstuff git --no-cache

COPY . .

RUN go get && \
    go build -i

FROM golang:1-alpine

COPY --from=build /go/src/app/app /usr/local/bin/app

CMD ["/usr/local/bin/app"]