module github.com/grafana/gel-app

replace (
	github.com/apache/arrow/go/arrow => github.com/apache/arrow/go/arrow v0.0.0-20190716210558-5f564424c71c
	github.com/grafana/grafana-plugin-sdk-go => github.com/grafana/grafana-plugin-sdk-go v0.5.1-0.20200114092430-dcbac1da7f41
)

go 1.13

require (
	github.com/google/go-cmp v0.3.1
	github.com/grafana/grafana-plugin-sdk-go v0.5.0
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/go-plugin v1.0.1
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20190923162816-aa69164e4478
	gonum.org/v1/gonum v0.0.0-20190923210107-40fa6a493b3d
	google.golang.org/grpc v1.24.0
)
