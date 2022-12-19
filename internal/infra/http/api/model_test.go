package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setTimersRequest_Validate(t *testing.T) {
	type fields struct {
		Hours   int
		Minutes int
		Seconds int
		URL     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "valid",
			fields:  fields{URL: "http://valid.url"},
			wantErr: assert.NoError,
		},
		{
			name:    "invalid",
			fields:  fields{URL: "invalid.url"},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SetTimersRequest{
				Hours:   tt.fields.Hours,
				Minutes: tt.fields.Minutes,
				Seconds: tt.fields.Seconds,
				URL:     tt.fields.URL,
			}
			tt.wantErr(t, r.Validate(), "Validate()")
		})
	}
}
