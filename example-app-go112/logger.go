package exampleservice

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/logging"
	stdLog "github.com/sirupsen/logrus"
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

}

func (l *loggingClient) Info(msg string) {
	stdLog.Info(msg)
	l.doLog(msg, logging.Info)
}

func (l *loggingClient) Infof(formatString string, toFormat ...interface{}) {
	stdLog.Infof(formatString, toFormat...)
	l.doLog(fmt.Sprintf(formatString, toFormat...), logging.Info)
}

func (l *loggingClient) Debug(msg string) {
	stdLog.Debug(msg)
	l.doLog(msg, logging.Debug)
}

func (l *loggingClient) Debugf(formatString string, toFormat ...interface{}) {
	stdLog.Debugf(formatString, toFormat...)
	l.doLog(fmt.Sprintf(formatString, toFormat...), logging.Debug)
}

func (l *loggingClient) Error(msg string) {
	stdLog.Error(msg)
	l.doLog(msg, logging.Error)
}

func (l *loggingClient) Errorf(formatString string, toFormat ...interface{}) {
	stdLog.Errorf(formatString, toFormat...)
	l.doLog(fmt.Sprintf(formatString, toFormat...), logging.Error)
}

func (l *loggingClient) Warning(msg string) {
	stdLog.Warning(msg)
	l.doLog(msg, logging.Error)
}

func (l *loggingClient) Warningf(formatString string, toFormat ...interface{}) {
	stdLog.Warningf(formatString, toFormat...)
	l.doLog(fmt.Sprintf(formatString, toFormat...), logging.Error)
}

func initLoggerClient(logName string, projectID string, versionID string, instanceID string, serviceID string) *loggingClient {

	ctx := context.Background()
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		panic(fmt.Sprintf("Error while logger init: %s", err.Error()))
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

	stdLog.SetFormatter(&stdLog.JSONFormatter{})
	stdLog.SetOutput(os.Stdout)

	return &loggingClient{lg: lg}
}
