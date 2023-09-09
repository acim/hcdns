package hcdns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Record struct {
	ID       string     `json:"id"`
	Type     RecordType `json:"type"`
	Name     string     `json:"name"`
	Value    string     `json:"value"`
	TTL      int        `json:"ttl,omitempty"`
	ZoneID   string     `json:"zone_id"`
	Created  Time       `json:"created"`
	Modified Time       `json:"modified"`
	c        *Client    `json:"-"`
	zoneID   string     `json:"-"`
}

func (r *Record) UpdateValue(ctx context.Context, value string) error {
	return r.UpdateValueAndTTL(ctx, value, 0)
}

func (r *Record) UpdateValueAndTTL(ctx context.Context, value string, ttl time.Duration) error {
	payload := recordReq{
		Type:   r.Type,
		Name:   r.Name,
		Value:  value,
		TTL:    uint64(ttl.Seconds()),
		ZoneID: r.zoneID,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	root, err := r.c.do(ctx, http.MethodPut, fmt.Sprintf("records/%s", r.ID), bytes.NewBuffer(json), nil)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}

	r.Value = root.Record.Value
	r.TTL = root.Record.TTL

	return nil
}

func (r *Record) Delete(ctx context.Context) error {
	_, err := r.c.do(ctx, http.MethodDelete, fmt.Sprintf("records/%s", r.ID), http.NoBody, nil)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}

type RecordType string

const (
	A     RecordType = "A"
	AAAA  RecordType = "AAAA"
	NS    RecordType = "NS"
	MX    RecordType = "MX"
	CNAME RecordType = "CNAME"
	RP    RecordType = "RP"
	TXT   RecordType = "TXT"
	SOA   RecordType = "SOA"
	HINFO RecordType = "HINFO"
	SRV   RecordType = "SRV"
	DANE  RecordType = "DANE"
	TLSA  RecordType = "TLSA"
	DS    RecordType = "DS"
	CAA   RecordType = "CAA"
)

type recordReq struct {
	Type   RecordType `json:"type"`
	Name   string     `json:"name"`
	Value  string     `json:"value"`
	TTL    uint64     `json:"ttl,omitempty"`
	ZoneID string     `json:"zone_id"`
}
