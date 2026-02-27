package aws


import (
	"context"
	"time"
    "fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)


func GetCPUUtilization(cfg aws.Config, instanceID string) (float64, error){
client := cloudwatch.NewFromConfig(cfg)

endTime := time.Now()
startTime := endTime.Add(-5 * time.Minute)

input := &cloudwatch.GetMetricStatisticsInput{
	Namespace:  aws.String("AWS/EC2"),
	MetricName: aws.String("CPUUtilization"),	
	Dimensions: []cloudwatchTypes.Dimension{
{
	Name: aws.String("InstanceId"),
	Value: aws.String(instanceID),
},
},


StartTime: aws.Time(startTime),
EndTime: aws.Time(endTime),
Period: aws.Int32(300),
Statistics: []cloudwatchTypes.Statistic{
	cloudwatchTypes.StatisticAverage,
},

}

	output, err := client.GetMetricStatistics(context.TODO(), input)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch CPU metrics: %w", err)
	} 
	if len(output.Datapoints) == 0 {
		return 0, fmt.Errorf("no CPU data available")
	}

	return *output.Datapoints[0].Average, nil

}