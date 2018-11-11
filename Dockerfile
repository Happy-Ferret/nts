FROM golang:1.11.2-alpine3.7

RUN apk add --no-cache git

RUN mkdir -p /go/src/github.com/krqa/nts/
WORKDIR /go/src/github.com/krqa/nts/

COPY *.go ./
COPY db/ ./db/
COPY routing/ ./routing/

ENV CGO_ENABLED=0
RUN go get -d -v ./...
RUN go build -ldflags="-s -w"


FROM alpine:3.7

RUN apk add --no-cache ca-certificates

COPY --from=0 /go/src/github.com/krqa/nts/nts /home/nts

EXPOSE 9000

ENTRYPOINT ["/home/nts"]
