package main

import (
	"context"
	"fmt"
	"os"
	"time"

	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	CALENDAR_ID = "e6stidpp5gkq4avsen1geudihk@group.calendar.google.com"
)

func main() {
	ctx := context.Background()
	cal, err := calendar.NewService(ctx, option.WithCredentialsFile("./service_key.json"))
	if err != nil {
		fmt.Printf("Could not create calendar API: %v", err)
		os.Exit(1)
	}

	evsCall := cal.Events.List(CALENDAR_ID).SingleEvents(true).TimeMax(time.Now().Format(time.RFC3339)).MaxResults(2500)
	res, err := evsCall.Do()
	if err != nil {
		fmt.Printf("Could not fetch calendar events: %v", err)
		os.Exit(1)
	}

	var pastEvents []*calendar.Event

	for _, e := range res.Items {
		if e.Summary == "" || e.End == nil || e.End.DateTime == "" {
			continue
		}
		end, err := time.Parse(time.RFC3339, e.End.DateTime)
		if err != nil {
			fmt.Printf("Could not parse date time of %s: %v", e.Summary, e.End.DateTime)
			continue
		}
		if end.Before(time.Now()) {
			pastEvents = append(pastEvents, e)
		}
	}

	for _, e := range pastEvents {
		if err := cal.Events.Delete(CALENDAR_ID, e.Id).Do(); err != nil {
			fmt.Printf("Could not delete past event: %s - %s: %v\n", e.Summary, e.End.DateTime, err)
			continue
		}
		fmt.Printf("Deleted past event: %s\n", e.Summary)
	}
}
