package cwmetric2

import (
	"fmt"
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

func describeLoadBalancer(alb *elbv2.ELBV2, name string) (out *elbv2.DescribeLoadBalancersOutput, err error) {
	params := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{aws.String(name)},
	}

	out, err = alb.DescribeLoadBalancers(params)

	return
}

func describeTargetGroups(alb *elbv2.ELBV2, lbArn *string) (out *elbv2.DescribeTargetGroupsOutput, err error) {
	params := &elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: lbArn,
	}

	out, err = alb.DescribeTargetGroups(params)

	return
}

func (cwm2 *CloudWatchMetric2) buildDimensions() (dimensions []*cloudwatch.Dimension, err error) {
	alb := elbv2.New(session.New(), aws.NewConfig().WithRegion(cwm2.Region))
	dimensions = []*cloudwatch.Dimension{}

	for name, value := range cwm2.Dimensions {
		if cwm2.Namespace == "AWS/ApplicationELB" && (name == "LoadBalancerName" || name == "LoadBalancerNameWithTG") {
			var albout *elbv2.DescribeLoadBalancersOutput
			albout, err = describeLoadBalancer(alb, value)

			if err != nil {
				return
			}

			albId := strings.SplitN(*albout.LoadBalancers[0].LoadBalancerArn, "/", 2)[1]
			albDim := &cloudwatch.Dimension{Name: aws.String("LoadBalancer"), Value: aws.String(albId)}
			dimensions = append(dimensions, albDim)

			if name == "LoadBalancerNameWithTG" {
				var tgout *elbv2.DescribeTargetGroupsOutput
				tgout, err = describeTargetGroups(alb, albout.LoadBalancers[0].LoadBalancerArn)

				if len(tgout.TargetGroups) < 1 {
					err = fmt.Errorf("cannot find TargetGroup")
				}

				if err != nil {
					return
				}

				tgId := strings.SplitN(*tgout.TargetGroups[0].TargetGroupArn, ":", 6)[5]
				tgDim := &cloudwatch.Dimension{Name: aws.String("TargetGroup"), Value: aws.String(tgId)}
				dimensions = append(dimensions, tgDim)
			}
		} else {
			dimension := &cloudwatch.Dimension{Name: aws.String(name), Value: aws.String(value)}
			dimensions = append(dimensions, dimension)
		}
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
