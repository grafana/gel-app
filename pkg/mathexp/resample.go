package mathexp

import (
	"fmt"

	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
)

var aliasToDuration = map[string]time.Duration{
	"D":   86400 * time.Second,
	"W":   604800 * time.Second,
	"MS":  2629800 * time.Second,
	"Y":   31557600 * time.Second,
	"H":   time.Hour,
	"T":   time.Minute,
	"min": time.Minute,
	"S":   time.Second,
	"L":   time.Millisecond,
	"ms":  time.Millisecond,
	"U":   time.Microsecond,
	"us":  time.Microsecond,
	"N":   time.Nanosecond,
}

var re *regexp.Regexp

func parseRule(rule string) (time.Duration, error) {
	if re == nil {
		// run only once
		aliases := make([]string, 0)
		for k := range aliasToDuration {
			aliases = append(aliases, k)
		}
		expr := strings.Join(aliases, "|")
		re = regexp.MustCompile(fmt.Sprintf(`^(\d*)(%v)$`, expr))
	}
	match := re.FindStringSubmatch(rule)

	errDuration, _ := time.ParseDuration("0s")
	if len(match) == 0 {
		// What should I return instead of s?
		return errDuration, fmt.Errorf("resample rule %v not implemented", rule)
	}
	var multiplier int64
	if match[1] != "" {
		valueInt64, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			// Different message for ErrSyntax and ErrRange
			return errDuration, fmt.Errorf("string %v cannot be converted to integer", match[1])
		}
		multiplier = valueInt64
	} else {
		multiplier = 1
	}
	return time.Duration(multiplier) * aliasToDuration[match[2]], nil
}

// Resample turns the Series into a Number based on the given reduction function
func (s Series) Resample(rule string, downsampler string, upsampler string, tr *datasource.TimeRange) (Series, error) {
	interval, err := parseRule(rule)
	if err != nil {
		return s, fmt.Errorf(`failed to parse "rule" field "%v": %v`, rule, err)
	}

	from := time.Unix(0, tr.FromEpochMs*int64(time.Millisecond))
	to := time.Unix(0, tr.ToEpochMs*int64(time.Millisecond))

	newSeriesLength := int(float64(to.Sub(from).Nanoseconds()) / float64(interval.Nanoseconds()))
	if newSeriesLength <= 0 {
		return s, fmt.Errorf("The series cannot be sampled further; the time range is shorter than the interval")
	}
	resampled := NewSeries(s.GetName(), s.GetLabels(), newSeriesLength+1)
	bookmark := 0
	var lastSeen *float64
	idx := 0
	t := from
	for !t.After(to) && idx <= newSeriesLength {
		vals := make([]float64, 0)
		sIdx := bookmark
		for {
			if sIdx == s.Len() {
				break
			}
			st, v := s.GetPoint(sIdx)
			if st.After(t) {
				break
			}
			bookmark++
			sIdx++
			lastSeen = v
			vals = append(vals, *v)
		}
		var value *float64
		var err error
		if len(vals) == 0 { // upsampling
			value, err = s.upsample(upsampler, sIdx, lastSeen)
		} else { // downsampling
			value, err = s.downsample(downsampler, vals)
		}
		if err != nil {
			return resampled, err
		}

		tv := t // his is required otherwise all points keep the latest timestamp; anything better?
		resampled.SetPoint(idx, &tv, value)
		t = t.Add(interval)
		idx++
	}
	return resampled, nil
}

func (s Series) upsample(method string, sIdx int, lastSeen *float64) (*float64, error) {
	switch method {
	case "pad":
		if lastSeen != nil {
			return lastSeen, nil
		} else {
			return nil, nil
		}
	case "backfilling":
		if sIdx == s.Len() { // no vals left
			return nil, nil
		} else {
			_, v := s.GetPoint(sIdx)
			return v, nil
		}
	case "fillna":
		return nil, nil
	default:
		return nil, fmt.Errorf("Upsampling method %q not implemented", method)
	}
}

func (s Series) downsample(method string, vals []float64) (*float64, error) {
	fVec := dataframe.NewField("", dataframe.FieldTypeNumber, vals).Vector
	switch method {
	case "sum":
		return Sum(fVec), nil
	case "mean":
		return Avg(fVec), nil
	case "min":
		return Min(fVec), nil
	case "max":
		return Max(fVec), nil
	default:
		return nil, fmt.Errorf("Upsampling method %q not implemented", method)
	}
}
