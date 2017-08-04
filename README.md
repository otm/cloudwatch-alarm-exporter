# CloudWatch Alarm Exporter
This exporter exposes the state of the CloudWatch alarms. It is also possible to utilize it to directly send alerts to prometheus/alertmanager however it is recommended to use metrics instead. The biggest reason for this is that it's possible to make queries on alarms and also get a bit of historic data.

## Building
Building the project is straight forward by using `go get`

 go get github.com/otm/cloudwatch-alarm-exporter
