package hcdns

import (
	"fmt"
	"strings"
	"time"
)

type root struct {
	Zones          []Zone          `json:"zones"`
	Zone           Zone            `json:"zone"`
	Records        []Record        `json:"records"`
	Record         Record          `json:"record"`
	PrimaryServers []PrimaryServer `json:"primary_servers"`
	PrimaryServer  PrimaryServer   `json:"primary_server"`
	Meta           struct {
		Pagination pagination `json:"pagination"`
	} `json:"meta"`
	Error err `json:"error"`
}

type pagination struct {
	Page         uint `json:"page"`
	PerPage      uint `json:"per_page"`
	PreviousPage uint `json:"previous_page"`
	NextPage     uint `json:"next_page"`
	LastPage     uint `json:"last_page"`
	TotalEntries uint `json:"total_entries"`
}

type err struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Time struct {
	inner *time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	value := strings.ReplaceAll(string(b), `"`, "")

	if value == "" {
		return nil
	}

	tt, err := time.Parse("2006-01-02 15:04:05.999 -0700 MST", value)
	if err != nil {
		return fmt.Errorf("parse time: %w", err)
	}

	t.inner = &tt

	return nil
}

func (t Time) String() string {
	return t.inner.Format(time.RFC3339)
}

type PrimaryServer struct {
	Port     int       `json:"port"`
	ID       string    `json:"id"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	ZoneID   string    `json:"zone_id"`
	Address  string    `json:"address"`
}
