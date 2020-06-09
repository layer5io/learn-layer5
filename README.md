
<p style="text-align:center;" align="center"><a href="https://layer5.io/meshery"><img align="center" style="margin-bottom:20px;" src="https://raw.githubusercontent.com/layer5io/layer5/master/assets/images/meshery/meshery-logo-tag-light-text-side.png"  width="70%" /></a><br /><br /></p>



[![Docker Pulls](https://img.shields.io/docker/pulls/layer5/learn-layer5.svg)](https://hub.docker.com/r/layer5/learn-layer5)
[![Go Report Card](https://goreportcard.com/badge/github.com/layer5io/learn-layer5)](https://goreportcard.com/report/github.com/layer5io/learn-layer5)
[![GitHub issues by-label](https://img.shields.io/github/issues/layer5io/learn-layer5/help%20wanted.svg)](https://github.com/layer5io/learn-layer5/issues?q=is%3Aopen+is%3Aissue+label%3A"help+wanted")
[![Website](https://img.shields.io/website/https/layer5.io/meshery.svg)](https://layer5.io/meshery/)
[![Twitter Follow](https://img.shields.io/twitter/follow/layer5.svg?label=Follow&style=social)](https://twitter.com/intent/follow?screen_name=mesheryio)
[![Slack](http://slack.layer5.io/badge.svg)](http://slack.layer5.io)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3564/badge)](https://bestpractices.coreinfrastructure.org/projects/3564)


# learn-layer5


## Service

The following are the routes defined by the `service` app and their functionality.

##### POST /call

This is the route whose metrics will be collected by the app. This route can be used to make the service call any other web service.

Simple POST request
```shell
# Command
curl --location --request POST 'http://localhost:9091/call' \
--data-raw ''
# No Output
```

`service` makes a POST request to `"http://httpbin.org/post"`.
```shell
# Command
curl --location --request POST 'http://localhost:9091/call' \
--header 'Content-Type: application/json' \
--data-raw '{
"host": "http://httpbin.org/post",
"body": "{\n\t\"hello\": \"bye\"\n}"
}'
# Output
{
  "args": {}, 
  "data": "{\n\t\"hello\": \"bye\"\n}", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Content-Length": "19", 
    "Content-Type": "application/json", 
    "Host": "httpbin.org", 
    "User-Agent": "Go-http-client/1.1", 
  }, 
  "json": {
    "hello": "bye"
  }, 
  "origin": "...", 
  "url": "http://httpbin.org/post"
}
```

`service` makes a get request (as body is not provided) to `http://httpbin.org/get`.
```shell
# Command
curl --location --request POST 'http://localhost:9091/call' \
--header 'Content-Type: application/json' \
--data-raw '{
"host": "http://httpbin.org/get",
}'
# Output
{
  "args": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Host": "httpbin.org", 
    "User-Agent": "Go-http-client/1.1", 
  }, 
  "origin": "...", 
  "url": "http://httpbin.org/get"
}
```

##### GET /metrics

Gets the metrics from `service`
```shell
# Command
curl --location --request GET 'localhost:9091/metrics' \
--header 'Content-Type: application/json' \
--data-raw '{
"hello": "bye"
}'
# Output
{
    "requestsReceived": "19", # Total requests service recieved
    "responsesFailed": "3",   # The responses of the requests the service made that failed
    "responsesSucceeded": "7" # The responses of the requests the service made that succeeded
}
```

##### DELETE /metrics

Clears the counters in `service`
```shell
# Command
curl --location --request DELETE 'localhost:9091/metrics' \
--header 'Content-Type: application/json' \
--data-raw '{
	"hello": "bye"
}'
# No Output
```