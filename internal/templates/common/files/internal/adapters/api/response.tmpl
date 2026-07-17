package api

import (
	"net/url"
	"strconv"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	// NextCursor is the opaque token a client passes back as ?cursor= to fetch
	// the next page of a list. It is omitted when there are no further rows.
	NextCursor string           `json:"next_cursor,omitempty"`
	Error      ErrorMessageType `json:"error,omitempty"`
}

// PageParams reads the pagination inputs from a request's query string. A
// missing or non-numeric limit yields 0, which the service reads as "use the
// default page size" — so a client that omits ?limit still gets a bounded page.
func PageParams(query url.Values) (limit int, cursor string) {
	if raw := query.Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			limit = n
		}
	}

	return limit, query.Get("cursor")
}
