package gelpoc

import (
	"context"
	"fmt"
	"strings"

	"github.com/grafana/gel-app/pkg/mathexp"
	"github.com/grafana/grafana-plugin-model/go/datasource"
)

// Command is an interface for all GEL commands.
type Command interface {
	NeedsVars() []string
	Execute(c context.Context, vars mathexp.Vars) (mathexp.Results, error)
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

// UnmarshalMathCommand creates a MathCommand from Grafana's frontend query.
func UnmarshalMathCommand(rn *rawNode) (*MathCommand, error) {
	rawExpr, ok := rn.Query["expression"]
	if !ok {
		return nil, fmt.Errorf("no expression in gel command for refId %v", rn.RefID)
	}
	exprString, ok := rawExpr.(string)
	if !ok {
		return nil, fmt.Errorf("expected expression to be a string, got %T", rawExpr)
	}

	gm, err := NewMathCommand(exprString)
	if err != nil {
		return nil, fmt.Errorf("invalid math command type in '%v': %v", rn.RefID, err)
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
func (gm *MathCommand) Execute(ctx context.Context, vars mathexp.Vars) (mathexp.Results, error) {
	return gm.Expression.Execute(vars)
}

// ReduceCommand is a GEL command for reduction of a timeseries such as a min, mean, or max.
type ReduceCommand struct {
	Reducer     string
	VarToReduce string
}

// NewReduceCommand creates a new ReduceCMD.
func NewReduceCommand(reducer, varToReduce string) *ReduceCommand {
	// TODO: validate reducer here, before execution
	return &ReduceCommand{
		Reducer:     reducer,
		VarToReduce: varToReduce,
	}
}

// UnmarshalReduceCommand creates a MathCMD from Grafana's frontend query.
func UnmarshalReduceCommand(rn *rawNode) (*ReduceCommand, error) {
	rawVar, ok := rn.Query["expression"]
	if !ok {
		return nil, fmt.Errorf("no variable to reduce in gel command for refId %v", rn.RefID)
	}
	varToReduce, ok := rawVar.(string)
	if !ok {
		return nil, fmt.Errorf("expected variable to be a string, got %T for refId %v", rawVar, rn.RefID)
	}
	varToReduce = strings.TrimPrefix(varToReduce, "$")

	rawReducer, ok := rn.Query["reducer"]
	if !ok {
		return nil, fmt.Errorf("no reducer specified in gel command for refId %v", rn.RefID)
	}
	redFunc, ok := rawReducer.(string)
	if !ok {
		return nil, fmt.Errorf("expected reducer to be a string, got %T for refId %v", rawReducer, rn.RefID)
	}

	return NewReduceCommand(redFunc, varToReduce), nil
}

// NeedsVars returns the variable names (refIds) that are dependencies
// to execute the command and allows the command to fulfill the Command interface.
func (gr *ReduceCommand) NeedsVars() []string {
	return []string{gr.VarToReduce}
}

// Execute runs the command and returns the results or an error if the command
// failed to execute.
func (gr *ReduceCommand) Execute(ctx context.Context, vars mathexp.Vars) (mathexp.Results, error) {
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

// ResampleCommand is a GEL command for resampling of a timeseries
type ResampleCommand struct {
	Rule          string
	VarToResample string
	Downsampler   string
	Upsampler     string
	TimeRange     *datasource.TimeRange
}

// NewResampleCommand creates a new ResampleCMD.
func NewResampleCommand(rule, varToResample string, downsampler string, upsampler string, tr *datasource.TimeRange) *ResampleCommand {
	// TODO: validate reducer here, before execution
	return &ResampleCommand{
		Rule:          rule,
		VarToResample: varToResample,
		Downsampler:   downsampler,
		Upsampler:     upsampler,
		TimeRange:     tr,
	}
}

// UnmarshalResampleCommand creates a ResampleCMD from Grafana's frontend query.
func UnmarshalResampleCommand(rn *rawNode, tr *datasource.TimeRange) (*ResampleCommand, error) {
	rawVar, ok := rn.Query["expression"]
	if !ok {
		return nil, fmt.Errorf("no variable to resample in gel command for refId %v", rn.RefID)
	}
	varToReduce, ok := rawVar.(string)
	if !ok {
		return nil, fmt.Errorf("expected variable to be a string, got %T for refId %v", rawVar, rn.RefID)
	}
	varToReduce = strings.TrimPrefix(varToReduce, "$")
	varToResample := varToReduce

	rawRule, ok := rn.Query["rule"]
	if !ok {
		return nil, fmt.Errorf("no rule specified in gel command for refId %v", rn.RefID)
	}
	rule, ok := rawRule.(string)
	if !ok {
		return nil, fmt.Errorf("expected reducer to be a string, got %T for refId %v", rawRule, rn.RefID)
	}

	rawDownsampler, ok := rn.Query["downsampler"]
	if !ok {
		return nil, fmt.Errorf("no downsampler specified in gel command for refId %v", rn.RefID)
	}
	downsampler, ok := rawDownsampler.(string)
	if !ok {
		return nil, fmt.Errorf("expected downsampler to be a string, got %T for refId %v", downsampler, rn.RefID)
	}

	rawUpsampler, ok := rn.Query["upsampler"]
	if !ok {
		return nil, fmt.Errorf("no downsampler specified in gel command for refId %v", rn.RefID)
	}
	upsampler, ok := rawUpsampler.(string)
	if !ok {
		return nil, fmt.Errorf("expected downsampler to be a string, got %T for refId %v", upsampler, rn.RefID)
	}
	return NewResampleCommand(rule, varToResample, downsampler, upsampler, tr), nil
}

// NeedsVars returns the variable names (refIds) that are dependencies
// to execute the command and allows the command to fulfill the Command interface.
func (gr *ResampleCommand) NeedsVars() []string {
	return []string{gr.VarToResample}
}

// Execute runs the command and returns the results or an error if the command
// failed to execute.
func (gr *ResampleCommand) Execute(ctx context.Context, vars mathexp.Vars) (mathexp.Results, error) {
	newRes := mathexp.Results{}
	for _, val := range vars[gr.VarToResample].Values {
		series, ok := val.(mathexp.Series)
		if !ok {
			return newRes, fmt.Errorf("can only resample type series, got type %v", val.Type())
		}
		num, err := series.Resample(gr.Rule, gr.Downsampler, gr.Upsampler, gr.TimeRange)
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
	// TypeResample is the CMDType for a GEL resampling function.
	TypeResample
)

func (gt CommandType) String() string {
	switch gt {
	case TypeMath:
		return "math"
	case TypeReduce:
		return "reduce"
	case TypeResample:
		return "resample"
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
	case "resample":
		return TypeResample, nil
	default:
		return TypeUnknown, fmt.Errorf("'%v' is not a GEL Type", s)
	}
}
