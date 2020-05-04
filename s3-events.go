package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	sp "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
	"os"
	"time"
)

type Config struct {
	Schema       string
	CollectorUri string
}

type SnowplowEvent struct {
	Source string
	Time   time.Time
	Bucket string
	Key    string
}

func ToSDJ(event SnowplowEvent) map[string]interface{} {
	var eventInterface map[string]interface{}
	eventrec, _ := json.Marshal(event)
	json.Unmarshal(eventrec, &eventInterface)
	return eventInterface
}

// Loads configuration from environment in lambda
func LoadConfig() Config {
	return Config{Schema: os.Getenv("EVENT_SCHEMA"), CollectorUri: os.Getenv("COLLECTOR_URI")}
}

// Pushes events to Snowplow Pipeline
func Push(cfg Config, tracker *sp.Tracker, event SnowplowEvent) {
	data := ToSDJ(event)
	sdj := sp.InitSelfDescribingJson(cfg.Schema, data)
	tracker.TrackSelfDescribingEvent(sp.SelfDescribingEvent{
		Event: sdj,
	})
}

// Handles S3events with specific configuration loaded from environment
func Handle(cfg Config) func(context.Context, events.S3Event) {
	return func(ctx context.Context, s3Event events.S3Event) {
		emitter := sp.InitEmitter(sp.RequireCollectorUri(cfg.CollectorUri))
		tracker := sp.InitTracker(sp.RequireEmitter(emitter))
		for _, record := range s3Event.Records {
			s3 := record.S3
			Push(cfg, tracker, SnowplowEvent{record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key})
		}
	}
}

func main() {
	lambda.Start(Handle(LoadConfig()))
}
