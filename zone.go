package hcdns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Zone struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	TTL           int      `json:"ttl"`
	Registrar     string   `json:"registrar"`
	LegacyDNSHost string   `json:"legacy_dns_host"`
	LegacyNS      []string `json:"legacy_ns"`
	NS            []string `json:"ns"`
	Created       Time     `json:"created"`
	Verified      Time     `json:"verified"`
	Modified      Time     `json:"modified"`
	Project       string   `json:"project"`
	Owner         string   `json:"owner"`
	Permission    string   `json:"permission"`
	ZoneType      struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Prices      any    `json:"prices"`
	} `json:"zone_type"`
	Status          string `json:"status"`
	Paused          bool   `json:"paused"`
	IsSecondaryDNS  bool   `json:"is_secondary_dns"`
	TxtVerification struct {
		Name  string `json:"name"`
		Token string `json:"token"`
	} `json:"txt_verification"`
	RecordsCount int     `json:"records_count"`
	c            *Client `json:"-"`
}

func (z *Zone) UpdateDefaultTTL(ctx context.Context, ttl time.Duration) error {
	payload := zoneReq{
		Name: z.Name,
		TTL:  uint64(ttl.Seconds()),
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	root, err := z.c.do(ctx, http.MethodPut, "zones/"+z.ID, bytes.NewBuffer(json), nil)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}

	z.TTL = root.Zone.TTL

	return nil
}

func (z *Zone) Delete(ctx context.Context) error {
	_, err := z.c.do(ctx, http.MethodDelete, "zones/"+z.ID, http.NoBody, nil)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}

func (z *Zone) Records(ctx context.Context) ([]Record, error) {
	var (
		records []Record
		query        = make(url.Values, 2) //nolint:mnd
		page    uint = 1
	)

	query.Set("zone_id", z.ID)

	for {
		query.Set("page", strconv.FormatUint(uint64(page), 10))

		root, err := z.c.do(ctx, http.MethodGet, "records", http.NoBody, query)
		if err != nil {
			return nil, fmt.Errorf("request: %w", err)
		}

		records = append(records, root.Records...)

		if root.Meta.Pagination.NextPage == page {
			break
		}

		page = root.Meta.Pagination.NextPage
	}

	for i := range records {
		r := &records[i]
		r.c = z.c
		r.zoneID = z.ID
	}

	return records, nil
}

func (z *Zone) Record(ctx context.Context, id string) (*Record, error) {
	root, err := z.c.do(ctx, http.MethodGet, "records/"+id, http.NoBody, nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	root.Record.c = z.c
	root.Record.zoneID = z.ID

	return &root.Record, nil
}

func (z *Zone) CreateRecord(ctx context.Context, _type RecordType, name, value string) (*Record, error) {
	return z.CreateRecordWithTTL(ctx, _type, name, value, 0)
}

func (z *Zone) CreateRecordWithTTL(ctx context.Context, _type RecordType, name, value string,
	ttl time.Duration,
) (*Record, error) {
	payload := recordReq{
		Type:   _type,
		Name:   name,
		Value:  value,
		TTL:    uint64(ttl.Seconds()),
		ZoneID: z.ID,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	root, err := z.c.do(ctx, http.MethodPost, "records", bytes.NewBuffer(json), nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	root.Record.c = z.c
	root.Record.zoneID = z.ID

	return &root.Record, nil
}

func (z *Zone) createNSRecords(ctx context.Context) error {
	for _, ns := range []string{"helium.ns.hetzner.de.", "hydrogen.ns.hetzner.com.", "oxygen.ns.hetzner.com."} {
		if _, err := z.CreateRecord(ctx, NS, "@", ns); err != nil {
			return fmt.Errorf("create NS record %s: %w", ns, err)
		}
	}

	return nil
}

func (z *Zone) PrimaryServers(ctx context.Context) ([]PrimaryServer, error) {
	query := make(url.Values, 1)

	query.Set("zone_id", z.ID)

	root, err := z.c.do(ctx, http.MethodGet, "primary_servers", http.NoBody, query)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return root.PrimaryServers, nil
}

type zoneReq struct {
	Name string `json:"name"`
	TTL  uint64 `json:"ttl,omitempty"`
}
