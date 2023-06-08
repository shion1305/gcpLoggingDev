package main

import (
	"cloud.google.com/go/logging/apiv2"
	logPb "cloud.google.com/go/logging/apiv2/loggingpb"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	ltype "google.golang.org/genproto/googleapis/logging/type"
	"google.golang.org/protobuf/types/known/structpb"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env")
	}
	logger := NewLogger()
	logger.createLogEntry()
	logger.queryLogEntryWithLogging()
}

type Logger struct {
	client *logging.Client
}

func NewLogger() Logger {
	data, err := os.ReadFile(os.Getenv("GCLOUD_CREDENTIAL_JSON"))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	client, err := logging.NewClient(ctx, option.WithCredentialsJSON(data))
	if err != nil {
		panic(err)
	}
	return Logger{
		client: client,
	}
}

func (l Logger) createLogEntry() {
	// create log entry
	payload, err := structpb.NewStruct(map[string]interface{}{
		"message": "test",
	})
	if err != nil {
		panic(err)
	}
	logEntry := logPb.WriteLogEntriesRequest{
		Entries: []*logPb.LogEntry{
			{
				LogName: "projects/" + os.Getenv("GCP_PROJECT_ID") + "/logs/" + os.Getenv("GCP_LOG_NAME"),
				Resource: &monitoredres.MonitoredResource{
					Type: "global",
				},
				Payload: &logPb.LogEntry_JsonPayload{
					JsonPayload: payload,
				},
				Severity: ltype.LogSeverity_INFO,
				Labels: map[string]string{
					"LEVEL": "Warning",
				},
			},
		},
	}
	ctx := context.Background()
	result, err := l.client.WriteLogEntries(ctx, &logEntry)
	if err != nil {
		panic(err)
	}
	fmt.Println("logging completed", result.String())
}

func (l Logger) queryLogEntryWithLogging() {
	ctx := context.Background()
	iter := l.client.ListLogEntries(ctx, &logPb.ListLogEntriesRequest{
		Filter:   "timestamp>=\"2023-06-08T00:00:00Z\" AND jsonPayload.message=\"test\"",
		PageSize: 30,
		ResourceNames: []string{
			"projects/" + os.Getenv("GCP_PROJECT_ID"),
		},
	})
	for {
		entry, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(entry)
	}
}
