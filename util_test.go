package blockfrost

import (
	"net/http"
	"testing"
)

func boolPtr(b bool) *bool { return &b }

func TestFormatParams(t *testing.T) {
	tests := []struct {
		query APIQueryParams
		want  string
	}{
		{APIQueryParams{Count: 5}, "count=5"},
		{APIQueryParams{Page: 10}, "page=10"},
		{APIQueryParams{}, ""},
		{APIQueryParams{Order: "asc"}, "order=asc"},
		{APIQueryParams{Count: 5, Page: 10}, "count=5&page=10"},
		{APIQueryParams{Count: 5, Page: 10, Order: "desc"}, "count=5&order=desc&page=10"},
		{APIQueryParams{From: "8929261"}, "from=8929261"},
		{APIQueryParams{To: "9999269:10"}, "to=9999269%3A10"},
		{APIQueryParams{OrderBy: "amount", Retired: boolPtr(true), Expired: boolPtr(true)}, ""},
	}
	req, err := http.NewRequest(http.MethodGet, "/go", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		tt, want := tt, tt.want
		v := req.URL.Query()
		t.Run("", func(t *testing.T) {
			got := formatParams(v, tt.query).Encode()
			if got != want {
				t.Fatalf("expected %s got %s", want, got)
			}
		})
	}
}

func TestFormatDrepsParams(t *testing.T) {
	tests := []struct {
		query APIQueryParams
		want  string
	}{
		{APIQueryParams{}, ""},
		{APIQueryParams{OrderBy: "amount"}, "order_by=amount"},
		{APIQueryParams{OrderBy: "invalid"}, ""},
		{APIQueryParams{Retired: boolPtr(true)}, "retired=true"},
		{APIQueryParams{Expired: boolPtr(false)}, "expired=false"},
		{APIQueryParams{OrderBy: "amount", Retired: boolPtr(false), Expired: boolPtr(true)}, "expired=true&order_by=amount&retired=false"},
	}
	req, err := http.NewRequest(http.MethodGet, "/go", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		tt, want := tt, tt.want
		v := req.URL.Query()
		t.Run("", func(t *testing.T) {
			got := formatDrepsParams(v, tt.query).Encode()
			if got != want {
				t.Fatalf("expected %s got %s", want, got)
			}
		})
	}
}
