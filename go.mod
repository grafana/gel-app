module github.com/grafana/gel-app

replace github.com/grafana/grafana => github.com/grafana/grafana v0.0.0-20190906082001-e7dd025bf379

replace google.golang.org/grpc => google.golang.org/grpc v1.11.1

replace github.com/hashicorp/go-hclog => github.com/hashicorp/go-hclog v0.0.0-20180402200405-69ff559dc25f

replace github.com/hashicorp/go-plugin => github.com/hashicorp/go-plugin v0.0.0-20180331002553-e8d22c780116

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20180514194645-7bb2a897381c

go 1.12

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/google/go-cmp v0.3.1
	github.com/grafana/grafana v0.0.0-00010101000000-000000000000
	github.com/grafana/grafana-plugin-model v0.0.0-20190906002829-a7284adfb995
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/go-plugin v0.0.0-20190220160451-3f118e8ee104
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20190909003024-a7b16738d86b
	gonum.org/v1/gonum v0.0.0-20190908220844-1d8f8b2ee4ce
	google.golang.org/grpc v1.11.1
)
