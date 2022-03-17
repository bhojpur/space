FROM alpine:3.15
RUN apk add --no-cache ca-certificates

ADD spacesvr /usr/local/bin
ADD spacectl /usr/local/bin

RUN addgroup -S bhojpur && \
    adduser -S -G bhojpur bhojpur && \
    mkdir /data && chown bhojpur:bhojpur /data

VOLUME /data

EXPOSE 9851
CMD ["spacesvr", "-d", "/data"]