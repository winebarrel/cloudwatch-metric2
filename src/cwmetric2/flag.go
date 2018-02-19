package cwmetric2

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var version string

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
	cwm2 = &CloudWatchMetric2{}
	dimensionsStr := ""
	var showVersion bool

	flag.StringVar(&cwm2.Region, "region", "", "region")
	flag.StringVar(&cwm2.Namespace, "namespace", "", "namespace")
	flag.StringVar(&cwm2.Metric, "metric", "", "metric")
	flag.StringVar(&dimensionsStr, "dimensions", "", "dimensions")
	flag.StringVar(&cwm2.Statistics, "statistics", "", "statistics")
	flag.Int64Var(&cwm2.Period, "period", 60, "period")
	flag.Int64Var(&cwm2.Delay, "delay", 0, "delay")
	flag.BoolVar(&cwm2.FailIfZero, "fail-if-zero", false, "fail-if-zero")
	flag.BoolVar(&showVersion, "version", false, "version")
	flag.Parse()

	if showVersion {
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

	if dimensionsStr == "" {
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
	dimensionNvs := strings.Split(dimensionsStr, ",")

	for _, nv := range dimensionNvs {
		nameValue := strings.SplitN(nv, "=", 2)

		if len(nameValue) != 2 {
			err = fmt.Errorf("invalid dimensions: %s", dimensionsStr)
			return
		}

		cwm2.Dimensions[nameValue[0]] = nameValue[1]
	}

	return
}
