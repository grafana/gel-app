package gelpoc

import (
	"fmt"

	"github.com/grafana/gel-app/pkg/mathexp"

	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/datasources"
	"github.com/grafana/grafana/pkg/tsdb"
)

// NodeType is the type of a DPNode. Currently either a GEL or datasource query.
type NodeType int

const (
	// TypeGELNode is a DPNode NodeType for GEL commands.
	TypeGELNode NodeType = iota
	// TypeDatasourceNode is a DPNode NodeType for datasource queries.
	TypeDatasourceNode
)

// Node is a node in a Data Pipeline. Node is either a GEL command or a datasource query.
type Node interface {
	ID() int64 // ID() allows the gonum graph node interface to be fulfilled
	NodeType() NodeType
	RefID() string
	Execute(c *models.ReqContext, vars mathexp.Vars) (mathexp.Results, error)
	String() string
}

// DataPipeline is an ordered set of nodes returned from DPGraph processing.
type DataPipeline []Node

// Execute runs all the command/datasource requests in the pipeline return a
// map of the refId of the of each command
func (dp *DataPipeline) Execute(c *models.ReqContext) (mathexp.Vars, error) {
	vars := make(mathexp.Vars)
	for _, node := range *dp {
		res, err := node.Execute(c, vars)
		if err != nil {
			return nil, err
		}

		vars[node.RefID()] = res
	}
	return vars, nil
}

const gelDataSourceName = "-- GEL --"

// buildPipeline builds a graph of the nodes, and returns the nodes in an
// executable order
func buildPipeline(targets []*simplejson.Json, tr *tsdb.TimeRange, cache datasources.CacheService) (DataPipeline, error) {
	graph, err := buildDependencyGraph(targets, tr, cache)
	if err != nil {
		return nil, err
	}

	nodes, err := buildExecutionOrder(graph)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// buildDependencyGraph returns a dependency graph for a set of target.
func buildDependencyGraph(targets []*simplejson.Json, tr *tsdb.TimeRange, cache datasources.CacheService) (*simple.DirectedGraph, error) {
	graph, err := buildGraph(targets, tr, cache)
	if err != nil {
		return nil, err
	}

	registry := buildNodeRegistry(graph)

	if err := buildGraphEdges(graph, registry); err != nil {
		return nil, err
	}

	return graph, nil
}

// buildExecutionOrder returns a sequence of nodes ordered by dependency.
func buildExecutionOrder(graph *simple.DirectedGraph) ([]Node, error) {
	sortedNodes, err := topo.Sort(graph)
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, len(sortedNodes))
	for i, v := range sortedNodes {
		nodes[i] = v.(Node)
	}

	return nodes, nil
}

// buildNodeRegistry returns a lookup table for reference IDs to respective node.
func buildNodeRegistry(g *simple.DirectedGraph) map[string]Node {
	res := make(map[string]Node)

	nodeIt := g.Nodes()

	for nodeIt.Next() {
		if dpNode, ok := nodeIt.Node().(Node); ok {
			res[dpNode.RefID()] = dpNode
		}
	}

	return res
}

// buildGraph creates a new graph populated with nodes for every target.
func buildGraph(targets []*simplejson.Json, tr *tsdb.TimeRange, cache datasources.CacheService) (*simple.DirectedGraph, error) {
	dp := simple.NewDirectedGraph()

	for _, target := range targets {
		datasource := target.Get("datasource").MustString()
		refID := target.Get("refId").MustString()

		switch datasource {
		case gelDataSourceName:
			node, err := buildGELNode(refID, dp, target)
			if err != nil {
				return nil, err
			}
			dp.AddNode(node)
		default: // If it's not a GEL target, it's a data source.
			dsNode := &DSNode{
				baseNode: baseNode{
					id:    dp.NewNode().ID(),
					refID: refID,
				},
				query:     target,
				timeRange: tr,
				dsCache:   cache,
			}
			dp.AddNode(dsNode)
		}
	}
	return dp, nil
}

// buildGraphEdges generates graph edges based on each node's dependencies.
func buildGraphEdges(dp *simple.DirectedGraph, registry map[string]Node) error {
	nodeIt := dp.Nodes()

	for nodeIt.Next() {
		node := nodeIt.Node().(Node)

		if node.NodeType() != TypeGELNode {
			// datasource node, nothing to do for now. Although if we want GEL results to be
			// used as datasource query params some day this will need change
			continue
		}

		gelNode := node.(*GELNode)

		for _, neededVar := range gelNode.GELCommand.NeedsVars() {
			neededNode, ok := registry[neededVar]
			if !ok {
				return fmt.Errorf("unable to find dependent node '%v'", neededVar)
			}

			if neededNode.ID() == gelNode.ID() {
				return fmt.Errorf("can not add self referencing node for var '%v' ", neededVar)
			}

			edge := dp.NewEdge(neededNode, gelNode)

			dp.SetEdge(edge)
		}
	}
	return nil
}
