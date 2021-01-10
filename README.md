<p style="text-align:center;" align="center">
  <img align="center" src="https://raw.githubusercontent.com/layer5io/layer5/master/assets/images/layer5/layer5-tag-white-bg.png" width="45%" /></p>

![GitHub contributors](https://img.shields.io/github/contributors/layer5io/layer5.svg)
![GitHub](https://img.shields.io/github/license/layer5io/layer5.svg) 
[![Docker Pulls](https://img.shields.io/docker/pulls/layer5/learn-layer5.svg)](https://hub.docker.com/r/layer5/learn-layer5)
[![Go Report Card](https://goreportcard.com/badge/github.com/layer5io/learn-layer5)](https://goreportcard.com/report/github.com/layer5io/learn-layer5)
[![GitHub issues by-label](https://img.shields.io/github/issues/layer5io/learn-layer5/help%20wanted.svg)](https://github.com/issues?utf8=✓&q=is%3Aopen+is%3Aissue+archived%3Afalse+org%3Alayer5io+label%3A%22help+wanted%22+")
[![Website](https://img.shields.io/website/https/layer5.io/meshery.svg)](https://layer5.io)
[![Twitter Follow](https://img.shields.io/twitter/follow/layer5.svg?label=Follow&style=social)](https://twitter.com/intent/follow?screen_name=mesheryio)
[![Slack](http://slack.layer5.io/badge.svg)](http://slack.layer5.io)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3564/badge)](https://bestpractices.coreinfrastructure.org/projects/3564)

<p style="clear:both;">

# Learn Layer5 (learn to service mesh)

The Learn Layer5 sample application is to be available for use across all service meshes that Meshery supports and is to used as:

- a learning device (for [service mesh workshops](https://layer5.io/workshops))
- for [Service Mesh Interface conformance](https://docs.google.com/document/d/1HL8Sk7NSLLj-9PRqoHYVIGyU6fZxUQFotrxbmfFtjwc/edit#)

## Application Architecture
The Learn Layer5 application includes three services: `app-a`, `app-b`, and `app-c`. Though they are different services, they are defined using the same app (source code in ./service). Each service is listening on port `9091/tcp`.

### Service

The following are the routes defined by the `service` app and their functionality.

#### POST /call

This route makes the service make requests to another service. Metrics are collected for this route. sample usage given below:


```shell
# Command
curl --location --request POST 'http://service-a:9091/call' \
--header 'Content-Type: application/json' \
--data-raw '{
"url": "http://service-b:9091/call",
"body": "{\r\n\"url\": \"http:\/\/service-c:9091\/echo\",\r\n\"body\": \"\",\r\n\"method\": \"GET\"\r\n}",
"method": "POST",
"headers": {
    "h1":"v1"
}
}'

# No Output
GET /echo HTTP/1.1
Host: service-c:9091
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip
Servicename: Service-B
```

In the above example, we are making a post request from `service-a` with the body:
```json
{
  "url": "http://service-b:9091/call",
  "body": "{\r\n\"url\": \"http:\/\/service-c:9091\/echo\",\r\n\"body\": \"\",\r\n\"method\": \"GET\"\r\n}",
  "method": "POST",
  "headers": {
      "h1":"v1"
  }
}
```
This will make `service-a` to make a `POST` request to `http://service-b:9091/call` with the headers specified above, and the body:
```json
{
  "url": "http://service-c:9091/echo",
  "body":"",
  "method": "GET"
}
```
This inturn will make `service-b` to make a `GET` request to `http://service-c:9091/echo`.

#### GET /metrics

Gets the metrics from `service-a`
```shell
# Command
curl --location --request GET 'http://service-b:9091/metrics'
# Output
{
    "ReqReceived": [
        "Service-A"
    ],
    "RespSucceeded": [
        {
            "URL": "http://service-c:9091/echo",
            "Method": "GET",
            "Headers": null
        }
    ],
    "RespFailed": []
}
```
* In ReqReceived we see list of requests `service-b` received and from whom it received. Here we see `service-A`. Actually each of the service sets a header `ServiceName` which is read  by the service to determine the sender. 
* As `service-b` made a request to `service-c` and the request succeeded, we can see the details in the list of successful responses (RespSucceeded).

#### DELETE /metrics

Clears the counters in `service`
```shell
# Command
curl --location --request DELETE 'http://34.68.35.174:9091/metrics'
# No Output
```

> Note: metrics are collected only for `/call` and `/echo`.

# Service Mesh Interface (SMI) Conformance
The `learn-layer5` application serves as a sample application to validate the conformance of any service mesh with the SMI specifications. To verify SMI conformance, run Meshery and the Meshery adapter for the specific service mesh you wish to test. Invoke the suite of SMI conformance tests on the specific service mesh you would like to validate.

## Invoking conformance tests

**As a Service Mesh user**
Meshery allows you to schedule tests and invoke them programmatically. Meshery will store these test results and allow you to retrieve them later. 

**As a Service Mesh maker**
Meshery guarantees provenance of these tests and facilitates the public publicing of this suite of tests results into a versioned, public matrix of conformance status (consisting of both supported capabilities and test compliance).
<!--
## Testing Manually

To manually invoke SMI Conformance test for a deployed service mesh, you can apply tests from the `smi-test` directory of this repository. Execute the following command to run the smi-conformace tests:

# To check for smi conformance of a deployed service mesh
Meshery, the service mesh management plane, facililtates the execution and results gathering of these conformance tests. All the tests are written in `smi-test` directory of this repository.
Build and run with the below commands to create a container for running the smi-conformace tests:-

```shell
docker build . -t meshery-smi-conformance:latest
docker run -p 10008:10008 meshery-smi-conformance:latest
```


```shell
kubectl kuttl test  --skip-cluster-delete=true --start-kind=false ./smi-test
```
-->

<br /><br /><p align="center"><i>If you’re using Learn Layer5 or if you like the project, please <a href="https://github.com/layer5io/meshery/stargazers">★</a> star this repository to show your support! 🤩</i></p>
</p>

<p style="clear:both;">
<h2><a name="contributing"></a><a name="community"></a> <a href="http://slack.layer5.io">Community</a> and <a href="https://github.com/layer5io/layer5/blob/master/CONTRIBUTING.md">Contributing</a></h2>
Our projects are community-built and welcome collaboration. 👍 Be sure to see the <a href="https://docs.google.com/document/d/17OPtDE_rdnPQxmk2Kauhm3GwXF1R5dZ3Cj8qZLKdo5E/edit">Layer5 Community Welcome Guide</a> for a tour of resources available to you and jump into our <a href="http://slack.layer5.io">Slack</a>! Contributors are expected to adhere to the <a href="https://github.com/cncf/foundation/blob/master/code-of-conduct.md">CNCF Code of Conduct</a>.

<a href="https://meshery.io/community"><img alt="Layer5 Service Mesh Community" src="img/readme/slack-128.png" style="margin-left:10px;padding-top:5px;" width="110px" align="right" /></a>

<a href="http://slack.layer5.io"><img alt="Layer5 Service Mesh Community" src="img/readme/community.svg" style="margin-right:8px;padding-top:5px;" width="140px" align="left" /></a>

<p>
✔️ <em><strong>Join</strong></em> any or all of the weekly meetings on the <a href="https://calendar.google.com/calendar/b/1?cid=bGF5ZXI1LmlvX2VoMmFhOWRwZjFnNDBlbHZvYzc2MmpucGhzQGdyb3VwLmNhbGVuZGFyLmdvb2dsZS5jb20">community calendar</a>.<br />
✔️ <em><strong>Watch</strong></em> community <a href="https://www.youtube.com/channel/UCFL1af7_wdnhHXL1InzaMvA?sub_confirmation=1">meeting recordings</a>.<br />
✔️ <em><strong>Access</strong></em> the <a href="https://drive.google.com/drive/u/4/folders/0ABH8aabN4WAKUk9PVA">community drive</a>.<br />
</p>
<p align="center">
<i>Not sure where to start?</i> Grab an open issue with the <a href="https://github.com/issues?utf8=✓&q=is%3Aopen+is%3Aissue+archived%3Afalse+org%3Alayer5io+label%3A%22help+wanted%22+">help-wanted label</a>.
</p>

## About Layer5

**Community First**
<p>The <a href="https://layer5.io">Layer5</a> community represents the largest collection of service mesh projects and their maintainers in the world.</p>

**Open Source First**
<p>We build projects to provide learning environments, deployment and operational best practices, performance benchmarks, create documentation, share networking opportunities, and more. Our shared commitment to the open source spirit pushes Layer5 projects forward.</p>

**License**

This repository and site are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).
