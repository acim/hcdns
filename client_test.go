package hcdns_test

import (
	"context"
	"os"
	"testing"
	"time"

	"go.acim.net/hcdns"
)

func TestZones(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	c := hcdns.NewClient(os.Getenv("TOKEN"))

	zs, err := c.Zones(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(zs) == 0 {
		t.Error("got 0 zones; want more than zero")
	}
}

func TestZone(t *testing.T) { //nolint:cyclop
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	client := hcdns.NewClient(os.Getenv("TOKEN"))
	zoneName := "hcdns-test.com"

	ctx := context.Background()

	zone, err := client.CreateZone(ctx, zoneName)
	if err != nil {
		t.Fatal(err)
	}

	if zone.Name != zoneName {
		t.Errorf("got name %s; want %s", zone.Name, zoneName)
	}

	if zone.TTL != 86400 {
		t.Errorf("got ttl %d; want %d", zone.TTL, 86400)
	}

	zone, err = client.Zone(ctx, zone.ID)
	if err != nil {
		t.Errorf("zone: %v", err)
	}

	if zone.Name != zoneName {
		t.Errorf("got name %s; want %s", zone.Name, zoneName)
	}

	zone, err = client.ZoneByName(ctx, zone.Name)
	if err != nil {
		t.Errorf("zone by name: %v", err)
	}

	if zone.Name != zoneName {
		t.Errorf("got name %s; want %s", zone.Name, zoneName)
	}

	if err := zone.UpdateDefaultTTL(ctx, time.Hour); err != nil {
		t.Errorf("update zone: %v", err)
	}

	if zone.TTL != 3600 {
		t.Errorf("got ttl %d; want %d", zone.TTL, 3600)
	}

	if err := zone.Delete(ctx); err != nil {
		t.Fatal(err)
	}
}
