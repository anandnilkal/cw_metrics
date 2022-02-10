package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type MetricStatCount struct {
	Namespace   string
	MetricName  string
	SampleCount float64
	Dimensions  []types.Dimension
}

// var metrics []MetricStatCount

type CloudwatchClient struct {
	Client *cloudwatch.Client
	Region string
}

func NewCloudwatchClient(ctx context.Context, region string) (*CloudwatchClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		fmt.Println(err, "failed to load SDK configuration, %v", err)
		return nil, err
	}
	return &CloudwatchClient{
		Client: cloudwatch.NewFromConfig(cfg),
		Region: region,
	}, nil
}

func (c *CloudwatchClient) getMetricsStats(ctx context.Context, start time.Time, end time.Time, period int32, metricName, nameSpace string, dimensions []types.Dimension) (*cloudwatch.GetMetricStatisticsOutput, error) {
	return c.Client.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
		MetricName: &metricName,
		Namespace:  &nameSpace,
		StartTime:  &start,
		EndTime:    &end,
		Period:     &period,
		Statistics: types.StatisticSampleCount.Values(),
		Dimensions: dimensions,
	})
}

func (c *CloudwatchClient) getMetrics(ctx context.Context) error {
	paginator := cloudwatch.NewListMetricsPaginator(
		c.Client,
		&cloudwatch.ListMetricsInput{
			// RecentlyActive: types.RecentlyActivePt3h,
		})

	end := time.Now()
	start := time.Now().Add(-24 * time.Hour)
	period := 24 * 60 * 60
	var totalSamples float64
	for {
		metricsListOut, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, value := range metricsListOut.Metrics {
			metricStatistic, err01 := c.getMetricsStats(
				ctx,
				start,
				end,
				int32(period),
				*value.MetricName,
				*value.Namespace,
				value.Dimensions,
			)
			if err01 != nil {
				fmt.Printf("failed with error :%s", err01.Error())
				return err01
			}
			if len(metricStatistic.Datapoints) > 0 {
				sampleCount := *metricStatistic.Datapoints[0].SampleCount
				totalSamples = totalSamples + sampleCount
				var dimensions []string
				for _, dimension := range value.Dimensions {
					dimensions = append(dimensions, strings.Join([]string{*dimension.Name, *dimension.Value}, "/"))
				}
				dimension := strings.Join(dimensions, ",")
				fmt.Printf("%s,%s,%f,%s\n", *value.Namespace, *value.MetricName, sampleCount, dimension)
			}
		}
		if !paginator.HasMorePages() {
			break
		}
	}
	fmt.Printf("%s,%s,%f", "totalSamples", "totalSamples", totalSamples)
	return nil
}
