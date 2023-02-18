package hcdns_test

import (
	"context"
	"os"
	"testing"
	"time"

	"go.acim.net/hcdns"
)

func TestRecords(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	client := hcdns.NewClient(os.Getenv("TOKEN"))
	zoneName := "hcdns-test.io"
	ctx := context.Background()

	zone, err := client.CreateZone(ctx, zoneName)
	if err != nil {
		t.Fatal(err)
	}

	records, err := zone.Records(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 4 {
		t.Errorf("got %d records; want 4 records", len(records))
	}

	if err := zone.Delete(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestRecord(t *testing.T) { //nolint:cyclop
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	c := hcdns.NewClient(os.Getenv("TOKEN"))
	zoneName := "hcdns-test.co"
	ctx := context.Background()

	zone, err := c.CreateZoneWithDefaultTTL(ctx, zoneName, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	record, err := zone.CreateRecordWithTTL(ctx, hcdns.A, "foo", "127.0.0.1", time.Minute)
	if err != nil {
		t.Error(err)
	}

	if record.Type != hcdns.A {
		t.Errorf("got type %s; want A", record.Type)
	}

	if record.Name != "foo" {
		t.Errorf("got name %s; want foo", record.Name)
	}

	if record.Value != "127.0.0.1" {
		t.Errorf("got value %s; want 127.0.0.1", record.Value)
	}

	if record.TTL != 60 {
		t.Errorf("got ttl %d; want 60", record.TTL)
	}

	if record.ZoneID != zone.ID {
		t.Errorf("got zone id %s; want %s", record.ZoneID, zone.ID)
	}

	_, err = zone.Record(ctx, record.ID)
	if err != nil {
		t.Error(err)
	}

	if err := zone.Delete(ctx); err != nil {
		t.Fatal(err)
	}
}
