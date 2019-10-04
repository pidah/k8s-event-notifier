FROM golang:1.13-alpine as builder

RUN apk --update add gcc libc-dev

ADD . /go/src/github.com/pidah/k8s-event-notifier

WORKDIR /go/src/github.com/pidah/k8s-event-notifier

RUN go build -buildmode=pie -ldflags "-linkmode external -extldflags -static -w" -o k8s-event-notifier

FROM alpine

RUN apk --update add ca-certificates

COPY --from=builder /go/src/github.com/pidah/k8s-event-notifier /

# Create a group and user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Tell docker that all future commands should run as the appuser user
USER appuser

CMD ["/k8s-event-notifier","--logtostderr","-v=4","2>&1"]
