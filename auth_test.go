package sendowl

import (
	"net/url"
	"testing"
)

func TestSigningText(t *testing.T) {
	var tests = []struct {
		signingKeySecret string
		url.Values
		signingText string
	}{
		{
			"sauce",
			url.Values{},
			"&secret=sauce",
		},
		{
			"sauce",
			url.Values{"expires_at": {"1234"}},
			"expires_at=1234&secret=sauce",
		},
		{
			"sauce",
			url.Values{"z": {"1"}, "expires_at": {"1234"}},
			"expires_at=1234&z=1&secret=sauce",
		},
	}
	for i, tt := range tests {
		actual := SigningText(tt.signingKeySecret, tt.Values)
		if actual != tt.signingText {
			t.Errorf("%d. SigningText(%s, %q) = %q, want %q", i, tt.signingKeySecret, tt.Values, actual, tt.signingKeySecret)
		}
	}
}

func TestSigningKey(t *testing.T) {
	var tests = []struct {
		signingKey       string
		signingKeySecret string
		expected         string
	}{
		{
			"barbeque",
			"sauce",
			"barbeque&sauce",
		},
	}
	for i, tt := range tests {
		actual := SigningKey(tt.signingKey, tt.signingKeySecret)
		if actual != tt.expected {
			t.Errorf("%d. SigningKey(%q, %q) = %q, want %q", i, tt.signingKey, tt.signingKeySecret, actual, tt.expected)
		}
	}
}
