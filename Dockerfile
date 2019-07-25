FROM alpine:3.10.0

# Add system certificates as Go application needs to authorize connections over HTTPS.
RUN apk add --no-cache curl ca-certificates

ADD bin/crawler-http-server /crawler

ENTRYPOINT [ "/crawler" ]
