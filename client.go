package hcdns

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const baseURL = "https://dns.hetzner.com/api/v1"

var ErrStatus = errors.New("http status")

type Client struct {
	inner *http.Client
	bu    string
	token string
}

func NewClient(token string, opts ...Option) *Client {
	client := &Client{
		inner: http.DefaultClient,
		bu:    baseURL,
		token: token,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) Zones(ctx context.Context) ([]Zone, error) {
	return c.ZonesByKeyword(ctx, "")
}

func (c *Client) ZonesByKeyword(ctx context.Context, keyword string) ([]Zone, error) {
	var (
		zones []Zone
		query url.Values
		page  uint = 1
	)

	if keyword != "" {
		query = make(url.Values, 2) //nolint:gomnd
		query.Set("search_name", keyword)
	} else {
		query = make(url.Values, 1)
	}

	for {
		query.Set("page", strconv.FormatUint(uint64(page), 10))

		root, err := c.do(ctx, http.MethodGet, "zones", http.NoBody, query)
		if err != nil {
			return nil, fmt.Errorf("request: %w", err)
		}

		zones = append(zones, root.Zones...)

		if root.Meta.Pagination.NextPage == page {
			break
		}

		page = root.Meta.Pagination.NextPage
	}

	for i := range zones {
		z := &zones[i]
		z.c = c
	}

	return zones, nil
}

func (c *Client) Zone(ctx context.Context, id string) (*Zone, error) {
	root, err := c.do(ctx, http.MethodGet, fmt.Sprintf("zones/%s", id), http.NoBody, nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	root.Zone.c = c

	return &root.Zone, nil
}

func (c *Client) ZoneByName(ctx context.Context, n string) (*Zone, error) {
	q := make(url.Values, 1)
	q.Set("name", n)

	root, err := c.do(ctx, http.MethodGet, "zones", http.NoBody, q)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	root.Zones[0].c = c

	return &root.Zones[0], nil
}

func (c *Client) CreateZone(ctx context.Context, name string) (*Zone, error) {
	return c.CreateZoneWithDefaultTTL(ctx, name, 0)
}

func (c *Client) CreateZoneWithDefaultTTL(ctx context.Context, name string, ttl time.Duration) (*Zone, error) {
	payload := zoneReq{
		Name: name,
		TTL:  uint64(ttl.Seconds()),
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	root, err := c.do(ctx, http.MethodPost, "zones", bytes.NewBuffer(json), nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	root.Zone.c = c

	if err := root.Zone.createNSRecords(ctx); err != nil {
		return nil, fmt.Errorf("create ns records: %w", err)
	}

	return &root.Zone, nil
}

func (c *Client) do(ctx context.Context, method, path string, body io.Reader, queryParams url.Values) (*root, error) {
	u := fmt.Sprintf("%s/%s", c.bu, path)

	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Auth-API-Token", c.token)

	if body != http.NoBody {
		req.Header.Set("Content-Type", "application/json")
	}

	if queryParams != nil {
		req.URL.RawQuery = queryParams.Encode()
	}

	res, err := c.inner.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer func() {
		err = errors.Join(err, res.Body.Close())
	}()

	var root root

	if err := json.NewDecoder(res.Body).Decode(&root); err != nil {
		return nil, errors.Join(err, fmt.Errorf("%w: %s", ErrStatus, res.Status))
	}

	if root.Error.Code >= http.StatusBadRequest {
		return nil, fmt.Errorf("%w: %s", ErrStatus, root.Error.Message)
	}

	return &root, nil
}

func WithClient(hc *http.Client) Option {
	return func(c *Client) {
		c.inner = hc
	}
}

func WithBaseURL(u url.URL) Option {
	return func(c *Client) {
		c.bu = u.String()
	}
}

type Option func(*Client)
