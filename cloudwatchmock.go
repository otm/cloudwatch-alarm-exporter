package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type cloudwatchMock struct{}

// DescribeAlarms implements the AlarmDescriber interface
func (c *cloudwatchMock) DescribeAlarms(*cloudwatch.DescribeAlarmsInput) (*cloudwatch.DescribeAlarmsOutput, error) {
	return &cloudwatch.DescribeAlarmsOutput{
		MetricAlarms: []*cloudwatch.MetricAlarm{
			{
				AlarmArn:         aws.String("arn:cloudwatch:foo"),
				AlarmName:        aws.String("test1"),
				AlarmDescription: aws.String("test description one"),
				Threshold:        aws.Float64(3),
				StateReason:      aws.String("Foo is to much"),
				StateValue:       aws.String(cloudwatch.StateValueOk),
				MetricName:       aws.String("CloudWatchMetricName"),
				Namespace:        aws.String("ASG"),
			},
			{
				AlarmArn:         aws.String("arn:cloudwatch:foo2"),
				AlarmName:        aws.String("test2"),
				AlarmDescription: aws.String("test description two"),
				Threshold:        aws.Float64(10.3),
				StateReason:      aws.String("Foo2 is to much"),
				StateValue:       aws.String(cloudwatch.StateValueAlarm),
				MetricName:       aws.String("CloudWatchMetricName2"),
				Namespace:        aws.String("ASG"),
			},
		},
	}, nil
}
