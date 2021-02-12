FROM golang:1.15-alpine as builder

COPY . /src
WORKDIR /src

RUN apk add --no-cache make && make build


FROM alpine:3.12

WORKDIR /srv

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /src/build/http-server /srv/http-server

RUN apk --update add tzdata ca-certificates && \
  cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
  echo "Europe/Moscow" > /etc/timezone && \
  date && \
  apk del tzdata

EXPOSE 3000

ENTRYPOINT ["/srv/http-server"]
