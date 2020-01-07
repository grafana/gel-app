module github.com/grafana/gel-app

replace github.com/apache/arrow/go/arrow => github.com/apache/arrow/go/arrow v0.0.0-20190716210558-5f564424c71c

//replace github.com/grafana/grafana-plugin-sdk-go => ../grafana-plugin-sdk-go

go 1.13

require (
	github.com/fatih/structtag v1.1.0 // indirect
	github.com/google/go-cmp v0.3.1
	github.com/grafana/grafana-plugin-sdk-go v0.4.1-0.20200107153407-ccbf1374e434
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/go-plugin v1.0.1
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mattn/go-runewidth v0.0.6 // indirect
	github.com/mgechev/dots v0.0.0-20190921121421-c36f7dcfbb81 // indirect
	github.com/mgechev/revive v0.0.0-20191017201419-88015ccf8e97 // indirect
	github.com/securego/gosec v0.0.0-20191104154532-b4c76d4234af // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20190923162816-aa69164e4478
	golang.org/x/tools v0.0.0-20191104213129-fda23558a172 // indirect
	gonum.org/v1/gonum v0.0.0-20190923210107-40fa6a493b3d
	google.golang.org/grpc v1.24.0
)
