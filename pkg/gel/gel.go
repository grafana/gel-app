package gelpoc

import (
	"fmt"
	"strings"

	"github.com/grafana/grafana/pkg/components/simplejson"

	"github.com/grafana/gel-app/pkg/mathexp"
	"github.com/grafana/grafana/pkg/models"
)

// Command is an interface for all GEL commands.
type Command interface {
	NeedsVars() []string
	Execute(c *models.ReqContext, vars mathexp.Vars) (mathexp.Results, error)
}

// MathCommand is a GEL commad for a GEL math expression such as "1 + $GA / 2"
type MathCommand struct {
	RawExpression string
	Expression    *mathexp.Expr
}

// NewMathCommand creates a new MathCommand. It will return an error
// if there is an error parsing expr.
func NewMathCommand(expr string) (*MathCommand, error) {
	parsedExpr, err := mathexp.New(expr)
	if err != nil {
		return nil, err
	}
	return &MathCommand{
		RawExpression: expr,
		Expression:    parsedExpr,
	}, nil
}

// UnmarshalMathCommand creates a MathCommand from Grafana's frontend target.
func UnmarshalMathCommand(target *simplejson.Json) (*MathCommand, error) {
	refID := target.Get("refId").MustString()
	exprString := target.Get("expression").MustString()
	gm, err := NewMathCommand(exprString)
	if err != nil {
		return nil, fmt.Errorf("invalid math command type in '%v': %v", refID, err)
	}
	return gm, nil
}

// NeedsVars returns the variable names (refIds) that are dependencies
// to execute the command and allows the command to fulfill the Command interface.
func (gm *MathCommand) NeedsVars() []string {
	return gm.Expression.VarNames
}

// Execute runs the command and returns the results or an error if the command
// failed to execute.
func (gm *MathCommand) Execute(c *models.ReqContext, vars mathexp.Vars) (mathexp.Results, error) {
	return gm.Expression.Execute(vars)
}

// ReduceCommand is a GEL command for reduction of a timeseries such as a min, mean, or max.
type ReduceCommand struct {
	Reducer     string
	VarToReduce string
}

// NewReduceCommand creates a new ReduceCMD.
func NewReduceCommand(reducer, varToReduce string) *ReduceCommand {
	return &ReduceCommand{
		Reducer:     reducer,
		VarToReduce: varToReduce,
	}
}

// UnmarshalReduceCommand creates a MathCMD from Grafana's frontend target.
func UnmarshalReduceCommand(target *simplejson.Json) *ReduceCommand {
	varToReduce := target.Get("expression").MustString()
	varToReduce = strings.TrimPrefix(varToReduce, "$")
	redFunc := target.Get("reducer").MustString()
	return NewReduceCommand(redFunc, varToReduce)
}

// NeedsVars returns the variable names (refIds) that are dependencies
// to execute the command and allows the command to fulfill the Command interface.
func (gr *ReduceCommand) NeedsVars() []string {
	return []string{gr.VarToReduce}
}

// Execute runs the command and returns the results or an error if the command
// failed to execute.
func (gr *ReduceCommand) Execute(c *models.ReqContext, vars mathexp.Vars) (mathexp.Results, error) {
	newRes := mathexp.Results{}
	for _, val := range vars[gr.VarToReduce].Values {
		series, ok := val.(mathexp.Series)
		if !ok {
			return newRes, fmt.Errorf("can only reduce type series, got type %v", val.Type())
		}
		num, err := series.Reduce(gr.Reducer)
		if err != nil {
			return newRes, err
		}
		newRes.Values = append(newRes.Values, num)
	}
	return newRes, nil
}

// CommandType is the type of GelCommand.
type CommandType int

const (
	// TypeUnknown is the CMDType for an unrecognized GEL type.
	TypeUnknown CommandType = iota
	// TypeMath is the CMDType for a GEL math expression.
	TypeMath
	// TypeReduce is the CMDType for a GEL reduction function.
	TypeReduce
)

func (gt CommandType) String() string {
	switch gt {
	case TypeMath:
		return "math"
	case TypeReduce:
		return "reduce"
	default:
		return "unknown"
	}
}

// ParseCommandType returns a CommandType from its string representation.
func ParseCommandType(s string) (CommandType, error) {
	switch s {
	case "math":
		return TypeMath, nil
	case "reduce":
		return TypeReduce, nil
	default:
		return TypeUnknown, fmt.Errorf("'%v' is not a GEL Type", s)
	}
}
