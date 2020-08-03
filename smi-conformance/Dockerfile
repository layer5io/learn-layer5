FROM golang:1.14-alpine3.11 as build-img
LABEL maintainer "Layer5.io"

ENV GO111MODULE=off

RUN apk update && apk add --no-cache git libc-dev gcc pkgconf && mkdir /home/meshery
COPY ${PWD} /go/src/github.com/layer5io/learn-layer5/smi-conformance/
WORKDIR /go/src/github.com/layer5io/learn-layer5/smi-conformance/
# RUN git rev-parse HEAD > /home/meshery/version
# RUN git describe --tags `git rev-list --tags --max-count=1` >> /home/com/version

RUN go build -a -ldflags "-s -w" -o /home/meshery/smi_conformance main.go

FROM registry.bookmyshow.org/alpine:latest

RUN apk --no-cache add ca-certificates
COPY --from=build-img /home/meshery/** /home/
WORKDIR /home/
EXPOSE 10008
CMD ["sh","-c","./smi_conformance"]