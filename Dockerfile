FROM golang:1.15-alpine3.12 as bin

ARG cmd

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd/${cmd} cmd/${cmd}
COPY internal internal
COPY pkg pkg

RUN go install ./cmd/${cmd}

FROM alpine:3.12

ARG cmd

COPY --from=bin /go/bin/${cmd} /usr/local/bin/app

ENTRYPOINT ["app"]
