FROM alpine:latest as alpine

RUN apk --no-cache add tzdata zip ca-certificates

WORKDIR /usr/share/zoneinfo

# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -q -r -0 /zoneinfo.zip .

FROM golang:1.15-alpine as builder

ARG SHA1VER
ARG APP_VER

COPY ./ /go/src

WORKDIR /go/src

RUN CGO_ENABLED=0 go build -ldflags="-w -s -X main.sha1ver=${SHA1VER} -X main.buildTime=`date +'%Y-%m-%d_%T'` -X main.version=${APP_VER}" -o /go/bin/portiapp ./cmd/app/...

FROM scratch

ENV TWITTER_LANGUAGE=es \
    TWITTER_HASHTAG=#example \
    TWITTER_SAMPLE=false

LABEL maintainer="labs@portico.net.ar"

# the timezone data:
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /

# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/bin/portiapp /go/bin/portiapp

ENTRYPOINT ["/go/bin/portiapp"]