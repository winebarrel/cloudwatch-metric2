package cwmetric2

import (
	"flag"
	"fmt"
	"strings"
)

type CloudWatchMetric2 struct {
	Region     string
	Namespace  string
	Metric     string
	Dimensions map[string]string
	Statistics string
}

func ParseFlag() (cwm2 *CloudWatchMetric2, err error) {
	cwm2 = &CloudWatchMetric2{}
	dimensionsStr := ""

	flag.StringVar(&cwm2.Region, "region", "", "region")
	flag.StringVar(&cwm2.Namespace, "namespace", "", "namespace")
	flag.StringVar(&cwm2.Metric, "metric", "", "metric")
	flag.StringVar(&dimensionsStr, "dimensions", "", "dimensions")
	flag.StringVar(&cwm2.Statistics, "statistics", "", "statistics")
	flag.Parse()

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
