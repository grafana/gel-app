package gelpoc

import (
	"context"
	"fmt"

	"github.com/grafana/gel-app/pkg/mathexp"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"gonum.org/v1/gonum/graph/simple"
)

// baseNode includes commmon properties used across DPNodes.
type baseNode struct {
	id    int64
	refID string
}

// String returns a string representation of the node. In particular for
// %v formating in error messages.
func (b *baseNode) String() string {
	return b.refID
}

// GELNode is a DPNode that holds a GEL command.
type GELNode struct {
	baseNode
	GELType    CommandType
	GELCommand Command
}

// ID returns the id of the node so it can fulfill the gonum's graph Node interface.
func (gn *GELNode) ID() int64 {
	return gn.id
}

// RefID returns the reference ID of the of the node that comes from
// the request.
func (gn *GELNode) RefID() string {
	return gn.refID
}

// NodeType returns the data pipeline node type.
func (gn *GELNode) NodeType() NodeType {
	return TypeGELNode
}

// Execute runs the node and adds the results to vars. If the node requires
// other nodes they must have already been executed and their results must
// already by in vars.
func (gn *GELNode) Execute(ctx context.Context, vars mathexp.Vars) (mathexp.Results, error) {
	return gn.GELCommand.Execute(ctx, vars)
}

func buildGELNode(refID string, dp *simple.DirectedGraph, target *simplejson.Json) (*GELNode, error) {
	commandType, err := ParseCommandType(target.Get("type").MustString())
	if err != nil {
		return nil, fmt.Errorf("invalid GEL type in '%v'", refID)
	}

	node := &GELNode{
		baseNode: baseNode{
			id:    dp.NewNode().ID(),
			refID: refID,
		},
	}

	switch commandType {
	case TypeMath:
		node.GELCommand, err = UnmarshalMathCommand(target)
		if err != nil {
			return nil, err
		}
	case TypeReduce:
		node.GELCommand = UnmarshalReduceCommand(target)
	default:
		return nil, fmt.Errorf("gel type '%v' in '%v' not implemented", commandType, refID)
	}

	return node, nil
}

// DSNode is a DPNode that holds a datasource request.
type DSNode struct {
	baseNode
	query     *simplejson.Json
	timeRange *datasource.TimeRange
	dsAPI     datasource.GrafanaAPI
}

// ID returns the id of the node so it can fulfill the gonum's graph Node interface.
func (dn *DSNode) ID() int64 {
	return dn.id
}

// RefID returns the reference ID of the of the node that comes from
// the request.
func (dn *DSNode) RefID() string {
	return dn.refID
}

// NodeType returns the data pipeline node type.
func (dn *DSNode) NodeType() NodeType {
	return TypeDatasourceNode
}

// Execute runs the node and adds the results to vars. If the node requires
// other nodes they must have already been executed and their results must
// already by in vars.
func (dn *DSNode) Execute(ctx context.Context, vars mathexp.Vars) (mathexp.Results, error) {
	datasourceID, err := dn.query.Get("datasourceId").Int64()
	if err != nil {
		return mathexp.Results{}, fmt.Errorf("query missing datasourceId")
	}
	orgID, err := dn.query.Get("orgId").Int64()
	if err != nil {
		return mathexp.Results{}, fmt.Errorf("query missing orgId")
	}

	// dn.query TO datasource.QueryDatasourceRequest
	qBytes, err := dn.query.Encode()
	if err != nil {
		return mathexp.Results{}, fmt.Errorf("failed to marshal query model: %v", err)
	}

	queries := []*datasource.Query{
		&datasource.Query{
			RefId:         dn.refID,
			IntervalMs:    dn.query.Get("intervalMs").MustInt64(1000),
			MaxDataPoints: dn.query.Get("maxDataPoints").MustInt64(5000),
			ModelJson:     string(qBytes),
		},
	}

	qd := &datasource.QueryDatasourceRequest{
		TimeRange:    dn.timeRange,
		Queries:      queries,
		DatasourceId: datasourceID,
		OrgId:        orgID,
	}

	resp, err := dn.dsAPI.QueryDatasource(ctx, qd)
	if err != nil {
		return mathexp.Results{}, err
	}
	vals := make([]mathexp.Value, 0, len(resp.Results))
	vals = append(vals, mathexp.FromGRPC(resp.Results[0].GetSeries()).Values...)

	// 	vals := make([]mathexp.Value, 0, len(resp.Results))
	// 	for _, tsdbRes := range resp.Results {
	// 		vals = append(vals, mathexp.FromTSDB(tsdbRes.Series).Values...)
	// 	}

	// 	return mathexp.Results{
	// 		Values: vals,
	// 	}, nil

	//_ = someRes

	// ds, err := dn.dsAPI.GetDatasource(datasourceID, c.SignedInUser, c.SkipCache)
	// if err != nil {
	// 	return mathexp.Results{}, fmt.Errorf("unable to load datasource: %v", err)
	// }

	return mathexp.Results{
		Values: vals,
	}, nil

	//return mathexp.Results{}, nil

	//return dn.execute(ctx, ds, vars)
}

// func (dn *DSNode) execute(ctx context.Context, ds *models.DataSource, vars mathexp.Vars) (mathexp.Results, error) {
// 	request := &tsdb.TsdbQuery{
// 		TimeRange: dn.timeRange,
// 		Queries: []*tsdb.Query{
// 			&tsdb.Query{
// 				RefId:      dn.query.Get("refId").MustString(),
// 				IntervalMs: dn.query.Get("intervalMs").MustInt64(1000),
// 				Model:      dn.query,
// 				DataSource: ds,
// 			},
// 		},
// 	}

// 	resp, err := tsdb.HandleRequest(ctx, ds, request)
// 	if err != nil {
// 		return mathexp.Results{}, fmt.Errorf("metric request error: %v", err)
// 	}

// 	for _, res := range resp.Results {
// 		if res.Error != nil {
// 			return mathexp.Results{}, fmt.Errorf("%v : %v", res.ErrorString, res.Error.Error())
// 		}
// 	}

// 	vals := make([]mathexp.Value, 0, len(resp.Results))
// 	for _, tsdbRes := range resp.Results {
// 		vals = append(vals, mathexp.FromTSDB(tsdbRes.Series).Values...)
// 	}

// 	return mathexp.Results{
// 		Values: vals,
// 	}, nil
// }
