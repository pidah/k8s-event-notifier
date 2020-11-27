FROM golang:1.15-alpine as builder

RUN apk --update add gcc libc-dev

ADD . /go/src/github.com/pidah/ethereum-data-fetcher

WORKDIR /go/src/github.com/pidah/ethereum-data-fetcher

RUN go build -buildmode=pie -ldflags "-linkmode external -extldflags -static -w" -o ethereum-data-fetcher

FROM alpine

RUN apk --update add ca-certificates

COPY --from=builder /go/src/github.com/pidah/ethereum-data-fetcher /

# Create a group and user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Tell docker that all future commands should run as the appuser user
USER appuser

CMD ["/ethereum-data-fetcher","--logtostderr","-v=4","2>&1"]
