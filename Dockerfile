FROM golang:1.12.7-alpine

WORKDIR /go/src/github.com/otm/cloudwatch-alarm-exporter
COPY . /go/src/github.com/otm/cloudwatch-alarm-exporter

RUN go build -o cloudwatch-alarm-exporter \
    && cp cloudwatch-alarm-exporter /bin/ \
    && mkdir -p /cloudwatch-alarm-exporter /etc/cloudwatch-alarm-exporter/aws \
    && chown -R nobody:nogroup /cloudwatch-alarm-exporter /etc/cloudwatch-alarm-exporter \
    && rm -rf /go

USER       nobody
ENV        AWS_CONFIG_FILE=/etc/cloudwatch-alarm-exporter/aws/config
ENV        AWS_SHARED_CREDENTIALS_FILE=/etc/cloudwatch-alarm-exporter/aws/credentials
EXPOSE     8080
WORKDIR    /cloudwatch-alarm-exporter
ENTRYPOINT [ "/bin/cloudwatch-alarm-exporter" ]
CMD        [ "-port=8080", "-retries=1", "-refresh=10" ]
