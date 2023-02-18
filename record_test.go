package hcdns_test

import (
	"context"
	"os"
	"testing"
	"time"

	"go.acim.net/hcdns"
)

func TestUpdateAndDelete(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	c := hcdns.NewClient(os.Getenv("TOKEN"))
	zoneName := "hcdns-test.rs"
	ctx := context.Background()

	zone, err := c.CreateZoneWithDefaultTTL(ctx, zoneName, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	record, err := zone.CreateRecordWithTTL(ctx, hcdns.A, "foo", "127.0.0.1", time.Minute)
	if err != nil {
		t.Error(err)
	}

	if err := record.UpdateValue(ctx, "192.168.0.1"); err != nil {
		t.Error(err)
	}

	if record.Value != "192.168.0.1" {
		t.Errorf("got value %s; want 192.168.0.1", record.Value)
	}

	if err := record.UpdateValueAndTTL(ctx, "192.168.1.1", time.Hour); err != nil {
		t.Error(err)
	}

	if record.Value != "192.168.1.1" {
		t.Errorf("got value %s; want 192.168.1.1", record.Value)
	}

	if record.TTL != 3600 {
		t.Errorf("got ttl %d; want 3600", record.TTL)
	}

	if err := record.Delete(ctx); err != nil {
		t.Error(err)
	}

	records, err := zone.Records(ctx)
	if err != nil {
		t.Error(err)
	}

	if len(records) != 4 {
		t.Errorf("got %d records; want 4 records", len(records))
	}

	if err := zone.Delete(ctx); err != nil {
		t.Fatal(err)
	}
}
