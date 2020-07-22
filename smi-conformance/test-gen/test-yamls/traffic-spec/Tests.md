# Traffic Spec

Note that the following tests are executed sequentially in an temporarily spawned namespace.

* *01* : Asserts if the HttpRoute group CRD exists.
* *02* : Deploys the app and asserts if the app has been deployed
* *03* : Configures HTTPRouteGroup such that traffic from `service-a` to only `service-b:PORT/metrics` (all HTTP methods) is allowed and the rest is blocked.
* *04* : Custom test which verifies if the above configuration works as intended.
* *05* : Configures the above created HTTPRouteGroup such that traffic from `service-a` to `service-b:PORT/*` (only GET HTTP Method) is allowed and the rest is blocked.
* *06* : Custom test which verifies if the above configuration works as intended.

> We aren't validating TCPRouteGroup CRD here as it has been used in the traffic access tests, so if it is not conformant even TrafficAccess will not be.