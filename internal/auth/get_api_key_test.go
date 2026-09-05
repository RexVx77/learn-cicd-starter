package auth

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetApiKey(t *testing.T) {
	header1 := http.Header{}
	header1.Set("Authorization", "ApiKey myapikey")
	header2 := http.Header{}
	header2.Set("Authorization", "ApiKey ")
	header3 := http.Header{}
	header3.Set("Authorization", "ApiKey")
	header4 := http.Header{}
	header4.Set("Authorization", "Api myapikey")
	header5 := http.Header{}
	tests := map[string]struct {
		input        http.Header
		want_err     error
		want_api_key string
	}{
		"happy path":           {input: header1, want_err: nil, want_api_key: "myapikey"},
		"no key with space":    {input: header2, want_err: nil, want_api_key: ""},
		"no key with no space": {input: header3, want_err: ErrMalformedAuthHeader, want_api_key: ""},
		"incorrect prefix for Authorization header": {input: header4, want_err: ErrMalformedAuthHeader, want_api_key: ""},
		"no Authorization":                          {input: header5, want_err: ErrNoAuthHeaderIncluded, want_api_key: ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got_key, got_err := GetAPIKey(tc.input)
			err_diff := cmp.Diff(tc.want_err, got_err, cmpopts.EquateErrors())
			key_diff := cmp.Diff(tc.want_api_key, got_key)
			if err_diff != "" || key_diff != "" {
				t.Fatalf("Error diff: %v\nKey diff: %v", err_diff, key_diff)
			}
		})
	}
}
