FROM docker.io/golang:1.25-alpine3.22 AS build

RUN apk update && \
    apk upgrade && \
    apk add --no-cache git libcap

WORKDIR /go/src/app
ENV GO111MODULE=on
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app -ldflags '-s -w -extldflags "-static"' .

FROM alpine:3.22

COPY --from=build /go/bin/app /

COPY --from=build /usr/sbin/setcap /usr/sbin/setcap
COPY --from=build /usr/lib/libcap.so.2* /usr/lib/

RUN /usr/sbin/setcap cap_net_bind_service=+ep /app

USER nobody:nogroup
ENTRYPOINT ["/app"]
