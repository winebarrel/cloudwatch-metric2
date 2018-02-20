package cwmetric2

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var version string

var region = flag.String("region", "", "region")
var namespace = flag.String("namespace", "", "namespace")
var metric = flag.String("metric", "", "metric")
var dimensionsStr = flag.String("dimensions", "", "dimensions")
var statistics = flag.String("statistics", "", "statistics")
var period = flag.Int64("period", 60, "period")
var delay = flag.Int64("delay", 0, "delay")
var failIfZero = flag.Bool("fail-if-zero", false, "fail-if-zero")
var showVersion = flag.Bool("version", false, "version")

func init() {
	flag.StringVar(region, "r", "", "region")
	flag.StringVar(namespace, "n", "", "namespace")
	flag.StringVar(metric, "m", "", "metric")
	flag.StringVar(dimensionsStr, "d", "", "dimensions")
	flag.StringVar(statistics, "s", "", "statistics")
	flag.Int64Var(period, "p", 60, "period")
	flag.Int64Var(delay, "l", 0, "delay")
	flag.BoolVar(failIfZero, "f", false, "fail-if-zero")
	flag.BoolVar(showVersion, "v", false, "version")
}

type CloudWatchMetric2 struct {
	Region     string
	Namespace  string
	Metric     string
	Dimensions map[string]string
	Statistics string
	Period     int64
	Delay      int64
	FailIfZero bool
}

func ParseFlag() (cwm2 *CloudWatchMetric2, err error) {
	flag.Parse()

	cwm2 = &CloudWatchMetric2{}
	cwm2.Region = *region
	cwm2.Namespace = *namespace
	cwm2.Metric = *metric
	cwm2.Statistics = *statistics
	cwm2.Period = *period
	cwm2.Delay = *delay
	cwm2.FailIfZero = *failIfZero

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if cwm2.Region == "" {
		err = fmt.Errorf("'-region' is required")
		return
	}

	if cwm2.Namespace == "" {
		err = fmt.Errorf("'-namespace' is required")
		return
	}

	if cwm2.Metric == "" {
		err = fmt.Errorf("'-metric' is required")
		return
	}

	if *dimensionsStr == "" {
		err = fmt.Errorf("'-dimensions' is required")
		return
	}

	if cwm2.Statistics == "" {
		err = fmt.Errorf("'-statistics' is required")
		return
	}

	if cwm2.Period < 1 {
		err = fmt.Errorf("invalid period")
		return
	}

	if cwm2.Delay < 0 {
		err = fmt.Errorf("invalid delay")
		return
	}

	cwm2.Dimensions = map[string]string{}
	dimensionNvs := strings.Split(*dimensionsStr, ",")

	for _, nv := range dimensionNvs {
		nameValue := strings.SplitN(nv, "=", 2)

		if len(nameValue) != 2 {
			err = fmt.Errorf("invalid dimensions: %s", *dimensionsStr)
			return
		}

		cwm2.Dimensions[nameValue[0]] = nameValue[1]
	}

	fmt.Println(cwm2)

	return
}
