package datemath_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/timberio/go-datemath"
)

// much are based on tests from Elasticsearch to ensure we handle dates in a compatible manner
// https://github.com/elastic/elasticsearch/blob/2d3f3cd61ef4b218082609928d6ffc9d20c30ba4/server/src/test/java/org/elasticsearch/common/time/JavaDateMathParserTests.java#L35
func TestParseAndEvaluate(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err error

		now      string
		location *time.Location
	}{

		// basic dates
		{
			in:  "2014",
			out: "2014-01-01T00:00:00.000Z",
		},
		{
			in:  "2014-05",
			out: "2014-05-01T00:00:00.000Z",
		},
		{
			in:  "2014-05-30",
			out: "2014-05-30T00:00:00.000Z",
		},
		{
			in:  "2014-05-30T20",
			out: "2014-05-30T20:00:00.000Z",
		},
		{
			in:  "2014-05-30T20:21",
			out: "2014-05-30T20:21:00.000Z",
		},
		{
			in:  "2014-05-30T20:21:35",
			out: "2014-05-30T20:21:35.000Z",
		},
		{
			in:  "2014-05-30T20:21:35.123",
			out: "2014-05-30T20:21:35.123Z",
		},

		// basic math
		{
			in:  "2014-11-18||+y",
			out: "2015-11-18T00:00:00.000Z",
		},
		{
			in:  "2014-11-18||-2y",
			out: "2012-11-18T00:00:00.000Z",
		},
		{
			in:  "2014-11-18||+3M",
			out: "2015-02-18T00:00:00.000Z",
		},
		{
			in:  "2014-11-18||-M",
			out: "2014-10-18T00:00:00.000Z",
		},
		{
			in:  "2014-11-18||+1w",
			out: "2014-11-25T00:00:00.000Z",
		},
		{
			in:  "2014-11-18||-3w",
			out: "2014-10-28T00:00:00.000Z",
		},
		{
			in:  "2014-11-18||+22d",
			out: "2014-12-10T00:00:00.000Z",
		},
		{
			in:  "2014-11-18||-423d",
			out: "2013-09-21T00:00:00.000Z",
		},
		{
			in:  "2014-11-18T14||+13h",
			out: "2014-11-19T03:00:00.000Z",
		},
		{
			in:  "2014-11-18T14||-1h",
			out: "2014-11-18T13:00:00.000Z",
		},
		{
			in:  "2014-11-18T14||+13H",
			out: "2014-11-19T03:00:00.000Z",
		},
		{
			in:  "2014-11-18T14||-1H",
			out: "2014-11-18T13:00:00.000Z",
		},
		{
			in:  "2014-11-18T14:27||+10240m",
			out: "2014-11-25T17:07:00.000Z",
		},
		{
			in:  "2014-11-18T14:27||-10m",
			out: "2014-11-18T14:17:00.000Z",
		},
		{
			in:  "2014-11-18T14:27:32||+60s",
			out: "2014-11-18T14:28:32.000Z",
		},
		{
			in:  "2014-11-18T14:27:32||-3600s",
			out: "2014-11-18T13:27:32.000Z",
		},
		{
			in:  "2014-11-19T14:27:32||/w",
			out: "2014-11-17T00:00:00.000Z",
		},
		{
			in:  "2014-11-01T14:27:32||/w",
			out: "2014-10-27T00:00:00.000Z",
		},

		// lenient math
		{
			in:  "2014-05-30T20:21||",
			out: "2014-05-30T20:21:00.000Z",
		},

		// multiple adjustments
		{
			in:  "2014-11-18||+1M-1M",
			out: "2014-11-18T00:00:00.000Z",
		},
		{
			in:  "2014-11-18||+1M-1m",
			out: "2014-12-17T23:59:00.000Z",
		},
		{
			in:  "2014-11-18||-1m+1M",
			out: "2014-12-17T23:59:00.000Z",
		},
		{
			in:  "2014-11-18||+1M/M",
			out: "2014-12-01T00:00:00.000Z",
		},
		{
			in:  "2014-11-18||+1M/M+1h",
			out: "2014-12-01T01:00:00.000Z",
		},

		// now
		{
			now: "2014-11-18T14:27:32.000Z",

			in:  "now",
			out: "2014-11-18T14:27:32.000Z",
		},
		{
			now: "2014-11-18T14:27:32.000Z",

			in:  "now+M",
			out: "2014-12-18T14:27:32.000Z",
		},
		{
			now: "2014-11-18T14:27:32.000Z",

			in:  "now/m",
			out: "2014-11-18T14:27:00.000Z",
		},

		// epoch times
		{
			in:  "04:52:20",
			out: "1970-01-01T04:52:20.000Z",
		},

		// timestamps
		{
			in:  "1418248078000",
			out: "2014-12-10T21:47:58.000Z",
		},
		{
			in:  "32484216259000",
			out: "2999-05-20T17:24:19.000Z",
		},
		{
			in:  "253382837059000",
			out: "9999-05-20T17:24:19.000Z",
		},
		{
			in:  "1418248078000||/m",
			out: "2014-12-10T21:47:00.000Z",
		},

		// timezones
		{
			in:  "2014-05-30T20:21+03:00",
			out: "2014-05-30T17:21:00.000Z",
		},
		{
			in:  "2014-05-30T20:21Z",
			out: "2014-05-30T20:21:00.000Z",
		},
		{
			location: time.FixedZone("custom", 2*60*60),

			in:  "2014-05-30T20:21",
			out: "2014-05-30T18:21:00.000Z",
		},
		{
			location: time.FixedZone("custom", 2*60*60),

			in:  "2014-05-30T20:21+03:00",
			out: "2014-05-30T17:21:00.000Z",
		},

		// errors
		{
			in:  "2014-05-35T20:21:21Z",
			out: "2014-05-30T20:21:21Z",

			err: fmt.Errorf(`day 35 out of bounds for month 5 at character 11 starting with "5"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var err error

			now := time.Now()
			if tt.now != "" {
				now, err = time.Parse(time.RFC3339Nano, tt.now)
				if err != nil {
					t.Fatal(err)
				}
			}

			location := time.UTC
			if tt.location != nil {
				location = tt.location
			}

			out, err := datemath.ParseAndEvaluate(tt.in, datemath.WithNow(now), datemath.WithLocation(location))
			switch {
			case err == nil && tt.err != nil:
				t.Errorf("ParseAndEvaluate(%+v) returned no error, expected error %q", tt.in, tt.err)
				return
			case err != nil && tt.err == nil:
				t.Errorf("ParseAndEvaluate(%+v) returned error %q, expected no error", tt.in, err)
				return
			case err != nil && tt.err != nil:
				if tt.err.Error() != err.Error() {
					t.Errorf("ParseAndEvaluate(%+v) returned error %q, expected error %q", tt.in, err, tt.err)
					return
				}
				return
			}

			expected, err := time.Parse(time.RFC3339Nano, tt.out)
			if err != nil {
				t.Fatal(err)
			}

			if !out.Equal(expected) {
				t.Errorf("ParseAndEvaluate(%q) returned %s, expected %s", tt.in, out, expected)
			}
		})
	}
}

var benchmarkParseResult datemath.Expression // used to avoid compiler optimizations

func Benchmark_Parse(b *testing.B) {
	var (
		expr datemath.Expression
		err  error
	)

	bench := func(s string) func(*testing.B) {
		return func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				expr, err = datemath.Parse(s)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	}

	b.Run("now", bench("now"))
	b.Run("rfc3339", bench("2014-05-30T20:21Z"))
	b.Run("now with one operation", bench("now+10240m"))
	b.Run("fixed with one operation", bench("2014-05-30T20:21Z||+10240m"))

	benchmarkParseResult = expr
}

var benchmarkExpressionTime time.Time // used to avoid compiler optimizations
func BenchmarkExpression_Time(b *testing.B) {
	var t time.Time

	bench := func(s string) func(*testing.B) {
		expr, err := datemath.Parse(s)
		if err != nil {
			panic(err)
		}
		return func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				t = expr.Time()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	}

	b.Run("now", bench("now"))
	b.Run("rfc3339", bench("2014-05-30T20:21Z"))
	b.Run("now with one operation", bench("now+10240m"))
	b.Run("fixed with one operation", bench("2014-05-30T20:21Z||+10240m"))

	benchmarkExpressionTime = t
}

func ExampleParse() {
	now, _ := time.Parse(time.RFC3339, "2014-05-30T20:21:35.123Z")

	expressions := []string{
		"now-15m",
		"now/w",
		"now+1M",
		"2014-05-31||+1M/w",
	}

	for _, expression := range expressions {
		t, err := datemath.Parse(expression)
		if err != nil {
			panic(err)
		}
		fmt.Println(t.Time(datemath.WithNow(now)))
	}

	//Output:
	//2014-05-30 20:06:35.123 +0000 UTC
	//2014-05-26 00:00:00 +0000 UTC
	//2014-06-30 20:21:35.123 +0000 UTC
	//2014-06-30 00:00:00 +0000 UTC
}
