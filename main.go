package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/google/uuid"
	"go.opencensus.io/trace"
)

var (
	spannerMinOpened uint64
)

func main() {
	projectID, err := GetProjectID()
	if err != nil {
		panic(err)
	}
	spannerDatabase := os.Getenv("SPANNER_DATABASE")
	fmt.Printf("Env SPANNER_DATABASE:%s\n", spannerDatabase)

	spannerMinOpenedParam := os.Getenv("SPANNER_MIN_OPENED")
	fmt.Printf("Env spannerMinOpened:%s\n", spannerMinOpenedParam)
	v, err := strconv.Atoi(spannerMinOpenedParam)
	if err != nil {
		panic(err)
	}
	spannerMinOpened = uint64(v)

	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()}) // defaultでは10,000回に1回のサンプリングになっているが、リクエストが少ないと出てこないので、とりあえず全部出す

	ctx := context.Background()
	sc := CreateClient(ctx, spannerDatabase, spannerMinOpened)
	ts := TweetStore{
		sc: sc,
	}

	for i := 0; i < 3600; i++ {
		ctx := context.Background()
		id := uuid.New().String()
		if err := ts.Insert(ctx, id); err != nil {
			log.Printf("failed tweet insert, err = %+v", err)
		} else {
			log.Printf("tweet insert id = %s", id)
		}
		time.Sleep(3 * time.Minute)
	}
}

func startSpan(ctx context.Context, name string) (context.Context, *trace.Span) {
	return trace.StartSpan(ctx, fmt.Sprintf("/little_spanner/min-%v/%s", spannerMinOpened, name))
}
