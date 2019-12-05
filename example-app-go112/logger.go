package exampleservice

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

type loggingClient struct {
	lg *logging.Logger
}

func (l *loggingClient) doLog(msg string, severity logging.Severity) {
	if l.lg == nil {
		return
	}
	l.lg.Log(logging.Entry{
		Severity: severity,
		Payload:  msg,
	})

	l.lg.StandardLogger(severity).Println(msg)
}

func (l *loggingClient) Info(msg string) {
	l.doLog(msg, logging.Info)
}

func (l *loggingClient) Infof(formatString string, toFormat ...interface{}) {
	msg := fmt.Sprintf(formatString, toFormat...)
	l.doLog(msg, logging.Info)
}

func (l *loggingClient) Debug(msg string) {
	l.doLog(msg, logging.Debug)
}

func (l *loggingClient) Debugf(formatString string, toFormat ...interface{}) {
	msg := fmt.Sprintf(formatString, toFormat...)
	l.doLog(msg, logging.Debug)
}

func (l *loggingClient) Error(msg string) {
	l.doLog(msg, logging.Error)
}

func (l *loggingClient) Errorf(formatString string, toFormat ...interface{}) {
	msg := fmt.Sprintf(formatString, toFormat...)
	l.doLog(msg, logging.Error)
}

func (l *loggingClient) Warning(msg string) {
	l.doLog(msg, logging.Error)
}

func (l *loggingClient) Warningf(formatString string, toFormat ...interface{}) {
	msg := fmt.Sprintf(formatString, toFormat...)
	l.doLog(msg, logging.Error)
}

func initLoggerClient(logName string, projectID string, versionID string, instanceID string, serviceID string) (*loggingClient, error) {

	ctx := context.Background()
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	resource := &monitoredres.MonitoredResource{
		Labels: map[string]string{
			"module_id":   serviceID,
			"project_id":  projectID,
			"version_id":  versionID,
			"instance_id": instanceID,
		},
		Type: "gae_app",
	}

	lg := client.Logger(
		logName,
		logging.CommonLabels(map[string]string{
			"commonChildParent": "commonLabelChildValue",
		}),
		logging.CommonResource(resource),
	)

	return &loggingClient{
		lg: lg,
	}, err
}
