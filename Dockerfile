FROM golang:1.10
RUN go get upspin.io/...
RUN go get github.com/minio/minio-go
RUN go get github.com/alinz/upspin-do/...
RUN CGO_ENABLED=0 go install github.com/alinz/upspin-do/cmd/upspinserver-do

FROM alpine:latest
RUN apk update && apk add ca-certificates shadow libcap && rm -rf /var/cache/apk/*

COPY --from=0 go/bin/upspinserver-do /upspinserver
COPY ./run.sh /run.sh

CMD ["sh", "./run.sh"]