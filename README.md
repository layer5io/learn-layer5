# learn-layer5

# To check for smi conformance of a deployed service mesh
We use kuttl to check for SMI conformance. All the tests are writtten in smi-test directory of this repository.
Execute the following command to run the smi-conformace tests:-

```shell
kubectl kuttl test  --skip-cluster-delete=true --start-kind=false ./smi-test
```

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