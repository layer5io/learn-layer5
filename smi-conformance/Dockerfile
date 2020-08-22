FROM golang:1.14-alpine3.11 as build-img

RUN apk update && apk add --no-cache git libc-dev gcc pkgconf && mkdir /home/meshery
COPY ${PWD} /go/src/github.com/layer5io/learn-layer5/smi-conformance/
WORKDIR /go/src/github.com/layer5io/learn-layer5/smi-conformance/
# RUN git rev-parse HEAD > /home/meshery/version
# RUN git describe --tags `git rev-list --tags --max-count=1` >> /home/com/version
RUN go mod vendor && go build -a -ldflags "-s -w" -o /home/meshery/smi_conformance main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates && mkdir /home/test-yamls && mkdir /home/test-yamls/traffic-access && mkdir /home/test-yamls/traffic-spec && mkdir /home/test-yamls/traffic-split
COPY --from=build-img /home/meshery/** /home/
COPY --from=build-img /go/src/github.com/layer5io/learn-layer5/smi-conformance/test-gen/test-yamls/traffic-access/** /home/test-yamls/traffic-access/
COPY --from=build-img /go/src/github.com/layer5io/learn-layer5/smi-conformance/test-gen/test-yamls/traffic-split/** /home/test-yamls/traffic-split/
COPY --from=build-img /go/src/github.com/layer5io/learn-layer5/smi-conformance/test-gen/test-yamls/traffic-spec/** /home/test-yamls/traffic-spec/
WORKDIR /home/
EXPOSE 10008
CMD ["sh","-c","./smi_conformance"]