# Traffic Access

Note that the following tests are executed sequentially in an temporarily spawned namespace.

* *01* : Asserts if the TrafficSplit CRD exists.
* *02* : Deploys the app and asserts that it is deployed.
* *03* : Custom test which verifies that if in default scenario the traffic to `app-svc` is split randomly between `app-b` and `app-c`.
* *04* : Configure a TrafficSplit CRD such that all traffic to `app-svc` is sent to only `app-b` and none to `app-c`.
* *05* : Custom test which verifies the above scenario.
* *06* : Configure a TrafficSplit CRD such that all traffic to `app-svc` is sent to only `app-c` and none to `app-b`.
* *07* : Custom test which verifies the above scenario.
* *08* : Configure a TrafficSplit CRD such that all traffic to `app-svc` is split between the two such that `app-b` gets more traffic (75%) than `app-c` (25%).
* *09* : Custom test which verifies the above scenario.
* *10* : Configure a TrafficSplit CRD such that all traffic to `app-svc` is split between the two such that `app-b` gets more traffic (25%) than `app-c` (75%).
* *11* : Custom test which verifies the above scenario.