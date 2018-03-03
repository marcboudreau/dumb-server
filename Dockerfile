FROM golang:1.10-alpine

LABEL maintainer="marc.a.boudreau@gmail.com"

WORKDIR /go/src/app

COPY . .

RUN go get -v -d ./...
RUN go install -v ./...

FROM alpine:latest

LABEL maintainer="marc.a.boudreau@gmail.com"

RUN addgroup -g 1001 dumbsrvr && adduser -S -u 1001 dumbsrvr -g dumbsrvr dumbsrvr

COPY --chown=dumbsrvr:dumbsrvr --from=0 /go/bin/app /usr/bin/app

USER dumbsrvr

EXPOSE 7979

CMD [ "app" ]