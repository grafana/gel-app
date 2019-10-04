module github.com/grafana/gel-app

replace google.golang.org/grpc => google.golang.org/grpc v1.11.1

replace github.com/hashicorp/go-hclog => github.com/hashicorp/go-hclog v0.0.0-20180402200405-69ff559dc25f

replace github.com/hashicorp/go-plugin => github.com/hashicorp/go-plugin v0.0.0-20180331002553-e8d22c780116

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20180514194645-7bb2a897381c

go 1.12

require (
	github.com/apache/arrow/go/arrow v0.0.0-20191004105443-1e2cf1f95df0
	github.com/google/go-cmp v0.3.1
	github.com/grafana/grafana-plugin-model v0.0.0-20190925141336-5d93412845bc
	github.com/grafana/grafana-plugin-sdk-go v0.0.0-20191004114449-2aa3c1124792
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/go-plugin v1.0.1
	github.com/mattetti/filebuffer v1.0.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20190909003024-a7b16738d86b
	gonum.org/v1/gonum v0.0.0-20190923210107-40fa6a493b3d
	google.golang.org/grpc v1.24.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)
