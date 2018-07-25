FROM golang:1.10
RUN go get upspin.io/...
RUN go get github.com/minio/minio-go
RUN CGO_ENABLED=0 go build -o go/bin/upspinserver github.com/alinz/upspin-go/cmd/upspinserver-do/main.go

FROM alpine:latest
RUN apk update && apk add ca-certificates shadow libcap && rm -rf /var/cache/apk/*

COPY --from=0 go/bin/upspinserver /upspinserver
COPY ./run.sh /run.sh

CMD ["sh", "./run.sh"]



