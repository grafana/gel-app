module github.com/grafana/gel-app

replace github.com/apache/arrow/go/arrow => github.com/apache/arrow/go/arrow v0.0.0-20190716210558-5f564424c71c

go 1.12

require (
	github.com/google/go-cmp v0.3.1
	github.com/grafana/grafana-plugin-sdk-go v0.0.0-20191024130641-6756418f682c
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/go-plugin v1.0.1
	github.com/kr/pretty v0.1.0 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20190909003024-a7b16738d86b
	gonum.org/v1/gonum v0.0.0-20190923210107-40fa6a493b3d
	google.golang.org/grpc v1.24.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)
