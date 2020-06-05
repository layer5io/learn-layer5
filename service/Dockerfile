FROM golang:1.13.5 as bd
WORKDIR /github.com/layer5io/sample-app-service
ADD . .
RUN GOPROXY=direct GOSUMDB=off go build -a -o /main .
CMD ["/main"]