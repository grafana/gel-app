# GEL POC

App for GEL datasource.

As of Sept 9th 2019 - the backend GEL code in `pkg` is no longer runnable. This is because the code was moved from the grafana-enterprise repo. We aim to have it runnable again as a plugin by Sept 13th.

The `pkg/gel` folder in particular was meant for enterprise so needs to be rewroked. The `pkg/data`
and `pkg/mathexp/...` packages can run tests.

Go modules init will be done once the `pkg/gel` folder is sane again.