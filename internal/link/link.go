package link

import (
	"strings"
	"time"
)

type Schema string

const (
	SchemaHTTP  Schema = "http"
	SchemaHTTPS Schema = "https"
)

func (s Schema) String() string {
	return string(s)
}

type ID = string

type Link struct {
	ID        ID
	Schema    Schema
	Domain    string
	URI       string
	CreatedAt time.Time
	ExpiredAt *time.Time
}

func (l *Link) Path() string {
	parts := strings.Split(l.URI, "?")

	if parts[0] == "" {
		return "/"
	}

	return parts[0]
}

func (l *Link) RawQuery() string {
	parts := strings.Split(l.URI, "?")

	if len(parts) <= 1 {
		return ""
	}

	return parts[1]
}

func (l *Link) RawURL() string {
	return l.Schema.String() + "://" + l.Domain + l.URI
}
