package handler

import (
	"testing"
)

func prepareMockRedis() {

}

func TestURLHandler(t *testing.T) {
	tests := []struct {
		name      string
		callType  string
		prepare   func()
		expectErr bool
	}{
		{
			name:      "GET: Redirect to original URL",
			callType:  "GET",
			expectErr: true,
		},
		{
			name:      "POST: Submit a long url for shortening",
			callType:  "POST",
			expectErr: false,
		},
	}

}
