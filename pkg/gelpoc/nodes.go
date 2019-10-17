package gelpoc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/gel-app/pkg/mathexp"
	"github.com/grafana/grafana-plugin-sdk-go"
	"gonum.org/v1/gonum/graph/simple"
)

// baseNode includes commmon properties used across DPNodes.
type baseNode struct {
	id    int64
	refID string
}

type rawNode struct {
	RefID string `json:"refId"`
	Query map[string]interface{}
}

func (rn *rawNode) GetDatasourceName() (string, error) {
	rawDs, ok := rn.Query["datasource"]
	if !ok {
		return "", fmt.Errorf("no datasource in query for refId %v", rn.RefID)
	}
	dsName, ok := rawDs.(string)
	if !ok {
		return "", fmt.Errorf("expted datasource identifer to be a string, got %T", rawDs)
	}
	return dsName, nil
}

func (rn *rawNode) GetGELType() (c CommandType, err error) {
	rawType, ok := rn.Query["type"]
	if !ok {
		return c, fmt.Errorf("no gel type in query for refId %v", rn.RefID)
	}
	typeString, ok := rawType.(string)
	if !ok {
		return c, fmt.Errorf("expected gel type to be a string, got %T", rawType)
	}
	return ParseCommandType(typeString)
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
func (b *baseNode) ID() int64 {
	return b.id
}

// RefID returns the refId of the node.
func (b *baseNode) RefID() string {
	return b.refID
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

func buildGELNode(dp *simple.DirectedGraph, tr grafana.TimeRange, rn *rawNode) (*GELNode, error) {

	commandType, err := rn.GetGELType()
	if err != nil {
		return nil, fmt.Errorf("invalid GEL type in '%v'", rn.RefID)
	}

	node := &GELNode{
		baseNode: baseNode{
			id:    dp.NewNode().ID(),
			refID: rn.RefID,
		},
	}

	switch commandType {
	case TypeMath:
		node.GELCommand, err = UnmarshalMathCommand(rn)
	case TypeReduce:
		node.GELCommand, err = UnmarshalReduceCommand(rn)
	case TypeResample:
		node.GELCommand, err = UnmarshalResampleCommand(rn, tr)
	default:
		return nil, fmt.Errorf("gel type '%v' in '%v' not implemented", commandType, rn.RefID)
	}
	if err != nil {
		return nil, err
	}

	return node, nil
}

const (
	defaultIntervalMS = int64(64)
	defaultMaxDP      = int64(5000)
)

// DSNode is a DPNode that holds a datasource request.
type DSNode struct {
	baseNode
	query        json.RawMessage
	datasourceID int64
	orgID        int64
	timeRange    grafana.TimeRange
	intervalMS   int64
	maxDP        int64
	dsAPI        grafana.GrafanaAPIHandler
}

// NodeType returns the data pipeline node type.
func (dn *DSNode) NodeType() NodeType {
	return TypeDatasourceNode
}

func buildDSNode(dp *simple.DirectedGraph, rn *rawNode, tr grafana.TimeRange, dsAPI grafana.GrafanaAPIHandler) (*DSNode, error) {
	encodedQuery, err := json.Marshal(rn.Query)
	if err != nil {
		return nil, err
	}

	dsNode := &DSNode{
		baseNode: baseNode{
			id:    dp.NewNode().ID(),
			refID: rn.RefID,
		},
		query:      json.RawMessage(encodedQuery),
		timeRange:  tr,
		intervalMS: defaultIntervalMS,
		maxDP:      defaultMaxDP,
		dsAPI:      dsAPI,
	}

	rawDsID, ok := rn.Query["datasourceId"]
	if !ok {
		return nil, fmt.Errorf("no datasourceId in gel command for refId %v", rn.RefID)
	}
	floatDsID, ok := rawDsID.(float64)
	if !ok {
		return nil, fmt.Errorf("expected datasourceId to be a float64, got %T for refId %v", rawDsID, rn.RefID)
	}
	dsNode.datasourceID = int64(floatDsID)

	rawOrgID, ok := rn.Query["orgId"]
	if !ok {
		return nil, fmt.Errorf("no orgId in gel command for refId %v", rn.RefID)
	}
	floatOrgID, ok := rawOrgID.(float64)
	if !ok {
		return nil, fmt.Errorf("expected orgId to be a float64, got %T for refId %v", rawOrgID, rn.RefID)
	}
	dsNode.orgID = int64(floatOrgID)

	var floatIntervalMS float64
	if rawIntervalMS := rn.Query["intervalMs"]; ok {
		if floatIntervalMS, ok = rawIntervalMS.(float64); !ok {
			return nil, fmt.Errorf("expected intervalMs to be an float64, got %T for refId %v", rawIntervalMS, rn.RefID)
		}
		dsNode.intervalMS = int64(floatIntervalMS)
	}

	var floatMaxDP float64
	if rawMaxDP := rn.Query["maxDataPoints"]; ok {
		if floatMaxDP, ok = rawMaxDP.(float64); !ok {
			return nil, fmt.Errorf("expected maxDataPoints to be an float64, got %T for refId %v", rawMaxDP, rn.RefID)
		}
		dsNode.maxDP = int64(floatMaxDP)
	}

	return dsNode, nil
}

// Execute runs the node and adds the results to vars. If the node requires
// other nodes they must have already been executed and their results must
// already by in vars.
func (dn *DSNode) Execute(ctx context.Context, vars mathexp.Vars) (mathexp.Results, error) {

	q := []grafana.Query{
		grafana.Query{
			RefID:         dn.refID,
			MaxDataPoints: dn.maxDP,
			Interval:      time.Duration(int64(time.Millisecond) * dn.intervalMS),
			ModelJSON:     dn.query,
		},
	}

	resp, err := dn.dsAPI.QueryDatasource(ctx, dn.orgID, dn.datasourceID, dn.timeRange, q)

	if err != nil {
		return mathexp.Results{}, err
	}

	vals := make([]mathexp.Value, 0)
	for _, res := range resp {
		for _, frame := range res.DataFrames {
			vals = append(vals, mathexp.Series{Frame: frame})
		}
		//vals = append(vals, mathexp.FromGRPC(dsRes.GetSeries()).Values...)
	}

	return mathexp.Results{
		Values: vals,
	}, nil
}
