package handler

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRejectsBadUsernames(t *testing.T) {
	for _, u := range []string{"", "ab", "../etc", "a%2Fb", "name with space"} {
		rec := httptest.NewRecorder()
		Handler(rec, httptest.NewRequest("GET", "/api/stats?user="+url.QueryEscape(u), nil))
		if rec.Code != 400 {
			t.Errorf("user %q: got %d, want 400", u, rec.Code)
		}
	}
}
