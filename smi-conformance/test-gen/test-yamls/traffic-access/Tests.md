# Traffic Access

Note that the following tests are executed sequentially in an temporarily spawned namespace.

* *01* : Asserts if the TrafficTarget CRD exists.
* *02* : Deploys the app and asserts if the app has been deployed
* *03* : Custom test is run where we verify if by default all traffic is blocked.
* *04* : Create and assert a TrafficTarget which allows traffic from `service-a` to `service-b`.
* *05* : Custom test is run where we verify if traffic from `service-a` to `service-b` succeeds.
* *06* : Deletes the CRDs created in step *04* and asserts its deletion.
* *07* : Custom test to verify if the traffic from `service-a` to `service-b` fails as the CRD has been deleted.