FROM docker.io/golang:1.27-alpine3.24 AS build

RUN apk update && \
    apk upgrade && \
    apk add --no-cache git libcap

WORKDIR /go/src/app
ENV GO111MODULE=on

# Download dependencies only
COPY go.mod go.sum .
RUN go mod download

# Copy and build the webhook example.go file separately to build dependencies for caching
COPY assets/example.go example.go
RUN CGO_ENABLED=0 go build -trimpath -o /go/bin/app -ldflags '-s -w -extldflags "-static"' . && \
    rm example.go /go/bin/app

# Copy and build the main source code
COPY src src
RUN CGO_ENABLED=0 go build -trimpath -o /go/bin/app -ldflags '-s -w -extldflags "-static"' ./src/

FROM alpine:3.24

RUN apk update && \
    apk upgrade --no-cache

COPY --from=build /go/bin/app /

COPY --from=build /usr/sbin/setcap /usr/sbin/setcap
COPY --from=build /usr/lib/libcap.so.2* /usr/lib/

RUN /usr/sbin/setcap cap_net_bind_service=+ep /app

USER nobody:nogroup
ENTRYPOINT ["/app"]
