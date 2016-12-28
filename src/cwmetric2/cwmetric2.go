package cwmetric2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"sort"
	"strings"
	"time"
)

type Datapoints []*cloudwatch.Datapoint

func (dps Datapoints) Len() int {
	return len(dps)
}

func (dps Datapoints) Less(i, j int) bool {
	// reverse sort
	return dps[i].Timestamp.Unix() > dps[j].Timestamp.Unix()
}

func (dps Datapoints) Swap(i, j int) {
	dps[i], dps[j] = dps[j], dps[i]
}

func getValue(dp *cloudwatch.Datapoint) float64 {
	if dp.Average != nil {
		return *dp.Average
	} else if dp.Maximum != nil {
		return *dp.Maximum
	} else if dp.Minimum != nil {
		return *dp.Minimum
	} else if dp.SampleCount != nil {
		return *dp.SampleCount
	} else if dp.Sum != nil {
		return *dp.Sum
	} else {
		return 0.0
	}
}

func (cwm2 *CloudWatchMetric2) buildDimensions() (dimensions []*cloudwatch.Dimension, err error) {
	alb := elbv2.New(session.New(), aws.NewConfig().WithRegion(cwm2.Region))
	dimensions = []*cloudwatch.Dimension{}
	var dimension *cloudwatch.Dimension

	for name, value := range cwm2.Dimensions {
		if cwm2.Namespace == "AWS/ApplicationELB" && name == "LoadBalancerName" {
			params := &elbv2.DescribeLoadBalancersInput{
				Names: []*string{aws.String(value)},
			}

			var out *elbv2.DescribeLoadBalancersOutput
			out, err = alb.DescribeLoadBalancers(params)

			if err != nil {
				return
			}

			albName := strings.SplitN(*out.LoadBalancers[0].LoadBalancerArn, "/", 2)[1]
			dimension = &cloudwatch.Dimension{Name: aws.String("LoadBalancer"), Value: aws.String(albName)}
		} else {
			dimension = &cloudwatch.Dimension{Name: aws.String(name), Value: aws.String(value)}
		}

		dimensions = append(dimensions, dimension)
	}

	return
}

func (cwm2 *CloudWatchMetric2) getMetricStatistics0(svc *cloudwatch.CloudWatch) (value float64, err error) {
	dimensions, err := cwm2.buildDimensions()

	if err != nil {
		return
	}

	now := time.Now()
	_5MinAgo := now.Add(-5 * time.Minute)

	params := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String(cwm2.Namespace),
		MetricName: aws.String(cwm2.Metric),
		Dimensions: dimensions,
		StartTime:  aws.Time(_5MinAgo),
		EndTime:    aws.Time(now),
		Period:     aws.Int64(60),
		Statistics: []*string{aws.String(cwm2.Statistics)},
	}

	out, err := svc.GetMetricStatistics(params)

	if err != nil {
		return
	}

	datapoints := out.Datapoints
	sort.Sort(Datapoints(datapoints))

	if len(datapoints) > 0 {
		dp := datapoints[0]
		value = getValue(dp)
	}

	return
}

func (cwm2 *CloudWatchMetric2) GetMetricStatistics() (value float64, err error) {
	svc := cloudwatch.New(session.New(), aws.NewConfig().WithRegion(cwm2.Region))
	value, err = cwm2.getMetricStatistics0(svc)
	return
}
