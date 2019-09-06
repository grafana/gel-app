package gelpoc

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/gel-app/pkg/data"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/components/null"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/tsdb"
	"gopkg.in/macaron.v1"
)

var (
	pluginID = "test-plugin"
	dsRefID  = "GA"
)

func TestService_Scenario(t *testing.T) {
	tsdb.RegisterTsdbQueryEndpoint(pluginID, func(ds *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
		return newTestQueryEndpoint(
			dsRefID,
			pt(1567345500283, 0.0),
			pt(1567346100283, 10.0),
		), nil
	})

	var (
		rr = routing.NewRouteRegister()
		cs = newTestCacheService(pluginID)
	)

	svc := &GelService{
		RouterRegister:  rr,
		DatasourceCache: cs,
	}

	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}

	ga := simplejson.New()
	ga.Set("datasource", "gdev-testdata")
	ga.Set("refId", dsRefID)
	ga.Set("datasourceId", 123)

	gb := simplejson.New()
	gb.Set("datasource", gelDataSourceName)
	gb.Set("refId", "GB")
	gb.Set("type", "reduce")
	gb.Set("reducer", "mean")
	gb.Set("expression", "$GA")

	gc := simplejson.New()
	gc.Set("datasource", gelDataSourceName)
	gc.Set("refId", "GC")
	gc.Set("type", "math")
	gc.Set("expression", "$GA + $GB")

	gd := simplejson.New()
	gd.Set("datasource", gelDataSourceName)
	gd.Set("refId", "GD")
	gd.Set("type", "reduce")
	gd.Set("reducer", "sum")
	gd.Set("expression", "$GC")

	tr := tsdb.NewTimeRange("", "")

	p, err := buildPipeline([]*simplejson.Json{ga, gb, gc, gd}, tr, cs)
	if err != nil {
		t.Fatal(err)
	}

	ctx := models.ReqContext{}

	mctx := macaron.Context{}
	mctx.Req = macaron.Request{httptest.NewRequest("GET", "/", nil)}

	ctx.Context = &mctx

	vars, err := p.Execute(&ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(vars) != 4 {
		t.Errorf("unexpected number of variables: %v", len(vars))
	}

	for _, val := range vars["GA"].Values {
		df := val.AsDataFrame()
		ts := fieldsAsTime(df.Fields)
		fs := fieldsAsFloat64(df.Fields)

		if len(ts) != len(fs) {
			t.Error("vectors should be of same length")
		}
		if want := 0.0; *fs[0] != want {
			t.Errorf("got = %v; want = %v", *fs[0], want)
		}
		if want := 10.0; *fs[1] != want {
			t.Errorf("got = %v; want = %v", *fs[1], want)
		}
	}

	for _, val := range vars["GB"].Values {
		df := val.AsDataFrame()
		fs := fieldsAsFloat64(df.Fields)

		if len(fs) != 1 {
			t.Error("mean should only be one number")
		}
		if want := 5.0; *fs[0] != want {
			t.Errorf("got = %v; want = %v", *fs[0], want)
		}
	}

	for _, val := range vars["GC"].Values {
		df := val.AsDataFrame()
		ts := fieldsAsTime(df.Fields)
		fs := fieldsAsFloat64(df.Fields)

		if len(ts) != len(fs) {
			t.Error("vectors should be of same length")
		}
		if want := 5.0; *fs[0] != want {
			t.Errorf("got = %v; want = %v", *fs[0], want)
		}
		if want := 15.0; *fs[1] != want {
			t.Errorf("got = %v; want = %v", *fs[1], want)
		}
	}

	for _, val := range vars["GD"].Values {
		df := val.AsDataFrame()
		fs := fieldsAsFloat64(df.Fields)

		if len(fs) != 1 {
			t.Error("sum should only be one number")
		}
		if want := 20.0; *fs[0] != want {
			t.Errorf("got = %v; want = %v", *fs[0], want)
		}
	}
}

// Helpers

func pt(ts, val float64) tsdb.TimePoint {
	return tsdb.TimePoint{null.FloatFrom(val), null.FloatFrom(ts)}
}

func fieldsAsTime(fs data.Fields) []*time.Time {
	var res []*time.Time
	for _, field := range fs {
		if v, ok := field.Vector.(*data.TimeVector); ok {
			for _, value := range *v {
				res = append(res, value)
			}
		}
	}
	return res
}

func fieldsAsFloat64(fs data.Fields) []*float64 {
	var res []*float64
	for _, field := range fs {
		if v, ok := field.Vector.(*data.Float64Vector); ok {
			for _, value := range *v {
				res = append(res, value)
			}
		}
	}
	return res
}

// Mocks

type testCacheService struct {
	GetDatasourceFn func(datasourceID int64, user *models.SignedInUser, skipCache bool) (*models.DataSource, error)
}

func newTestCacheService(datasourceType string) *testCacheService {
	return &testCacheService{
		GetDatasourceFn: func(datasourceID int64, user *models.SignedInUser, skipCache bool) (*models.DataSource, error) {
			return &models.DataSource{
				Type: pluginID,
			}, nil
		},
	}
}

func (s *testCacheService) GetDatasource(datasourceID int64, user *models.SignedInUser, skipCache bool) (*models.DataSource, error) {
	return s.GetDatasourceFn(datasourceID, user, skipCache)
}

type testQueryEndpoint struct {
	QueryFn func(ctx context.Context, ds *models.DataSource, query *tsdb.TsdbQuery) (*tsdb.Response, error)
}

func newTestQueryEndpoint(refID string, points ...tsdb.TimePoint) *testQueryEndpoint {
	return &testQueryEndpoint{
		QueryFn: func(ctx context.Context, ds *models.DataSource, query *tsdb.TsdbQuery) (*tsdb.Response, error) {
			return &tsdb.Response{
				Results: map[string]*tsdb.QueryResult{
					refID: &tsdb.QueryResult{
						RefId: refID,
						Series: tsdb.TimeSeriesSlice{
							&tsdb.TimeSeries{
								Name:   refID + "-series",
								Points: tsdb.TimeSeriesPoints(points),
							},
						},
					},
				},
			}, nil
		},
	}
}

func (e *testQueryEndpoint) Query(ctx context.Context, ds *models.DataSource, query *tsdb.TsdbQuery) (*tsdb.Response, error) {
	return e.QueryFn(ctx, ds, query)
}

func TestReqGoNumPlay(t *testing.T) {
	s := &GelService{}
	aGelBlob := []byte(`
	{
		"datasource": "-- GEL --",
		"type": "math",
		"refId": "GA",
		"expression": "1 + $GB"
	}
	`)
	bGelBlob := []byte(`
	{
		"datasource": "-- GEL --",
		"type": "math",
		"refId": "GB",
		"expression": "1 + $GC"
	}
	`)
	datasourceBlob := []byte(`
	{
		"refId": "GC",
		"datasource": "gdev-testdata",
		"scenarioId": "predictable_pulse",
		"stringInput": "",
		"pulseWave": {
			"timeStep": 60,
			"onCount": 3,
			"onValue": 2,
			"offCount": 3,
			"offValue": 1
		},
		"datasourceId": 4
	}
	`)
	dGelBlob := []byte(`
	{
		"datasource": "-- GEL --",
		"type": "math",
		"refId": "GD",
		"expression": "1 + $GE"
	}
	`)
	eGelBlob := []byte(`
	{
		"datasource": "-- GEL --",
		"type": "math",
		"refId": "GE",
		"expression": "1"
	}
	`)

	aGelSJSON, err := simplejson.NewJson(aGelBlob)
	if err != nil {
		t.FailNow()
	}
	bGelSJSON, err := simplejson.NewJson(bGelBlob)
	if err != nil {
		t.FailNow()
	}
	dGelSJSON, err := simplejson.NewJson(dGelBlob)
	if err != nil {
		t.FailNow()
	}
	eGelSJSON, err := simplejson.NewJson(eGelBlob)
	if err != nil {
		t.FailNow()
	}
	dsJSON, err := simplejson.NewJson(datasourceBlob)
	if err != nil {
		t.FailNow()
	}
	req := GelAppReq{
		ReqOptions{
			Targets: []*simplejson.Json{
				aGelSJSON,
				bGelSJSON,
				dsJSON,
				dGelSJSON, // d and e are connect to each other, but not the rest of the graph
				eGelSJSON, //
				// P.O.S !@#$ simplejson seems to ignore struct tags
				// simplejson.NewFromAny(
				// 	&struct {
				// 		Datasource string `json:"datasource"`
				// 		Type       string `json:"type"`
				// 		RefID      string `json:"refId"`
				// 		Expr       string `json:"expression"`
				// 	}{
				// 		gelDataSourceName,
				// 		"math",
				// 		"$GB",
				// 		"1 + $GA",
				// 	},
				// ),
			},
		},
	}

	nodes, err := buildPipeline(
		req.Options.Targets,
		tsdb.NewTimeRange(req.Options.Range.Raw.From, req.Options.Range.Raw.To),
		s.DatasourceCache,
	)

	spew.Dump(err)
	spew.Dump(nodes)
}
