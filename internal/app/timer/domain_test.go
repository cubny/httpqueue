package timer

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTimer(t *testing.T) {
	now := time.Now()

	type args struct {
		rawURL  string
		hours   time.Duration
		minutes time.Duration
		seconds time.Duration
	}
	tests := []struct {
		name        string
		args        args
		wantURLRaw  string
		wantFireAt  time.Time
		wantErrFunc assert.ErrorAssertionFunc
	}{
		{
			name: "valid",
			args: args{
				rawURL:  "http://valid.url",
				hours:   0,
				minutes: 0,
				seconds: 0,
			},
			wantURLRaw:  "http://valid.url",
			wantFireAt:  now,
			wantErrFunc: assert.NoError,
		},
		{
			name: "invalid url",
			args: args{
				rawURL: "invalid.url",
			},
			wantErrFunc: assert.Error,
		},
		{
			name: "empty url",
			args: args{
				rawURL: "",
			},
			wantErrFunc: assert.Error,
		},
		{
			name: "timer in the past",
			args: args{
				rawURL: "http://valid.url",
				hours:  -1,
			},
			wantErrFunc: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTimer(tt.args.rawURL, tt.args.hours, tt.args.minutes, tt.args.seconds)
			tt.wantErrFunc(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, fmt.Sprintf("%s/%s", tt.wantURLRaw, got.ID), got.URL.String())
			assert.WithinDuration(t, tt.wantFireAt, got.FireAt, time.Second)
		})
	}
}

func TestNewTimerFromCommand(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		cmd        SetTimerCommand
		wantURLRaw string
		wantFireAt time.Time
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "valid",
			cmd: SetTimerCommand{
				Hours:   0,
				Minutes: 0,
				Seconds: 0,
				URLRaw:  "http://valid.url",
			},
			wantURLRaw: "http://valid.url",
			wantFireAt: now,
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTimerFromCommand(tt.cmd)
			if !tt.wantErr(t, err, fmt.Sprintf("NewTimerFromCommand(%v)", tt.cmd)) {
				return
			}
			assert.Equal(t, fmt.Sprintf("%s/%s", tt.wantURLRaw, got.ID), got.URL.String())
			assert.WithinDuration(t, tt.wantFireAt, got.FireAt, time.Second)
		})
	}
}

func TestTimer_DelayFromNowSeconds(t1 *testing.T) {
	now := time.Now()

	type fields struct {
		ID     string
		URL    url.URL
		FireAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "now",
			fields: fields{
				FireAt: now,
			},
			want: 0,
		},
		{
			name: "in am hour",
			fields: fields{
				FireAt: now.Add(time.Hour),
			},
			want: float64(3600),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Timer{
				ID:     tt.fields.ID,
				URL:    tt.fields.URL,
				FireAt: tt.fields.FireAt,
			}
			assert.InDelta(t1, tt.want, t.DelayFromNowSeconds(), 1)
		})
	}
}

func TestTimer_Validate(t1 *testing.T) {
	type fields struct {
		ID     string
		URL    url.URL
		FireAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "constructor produces valid timer",
			fields: func() fields {
				tm, err := NewTimer("http://valid.url", 0, 0, 0)
				require.NoError(t1, err)
				return fields{
					ID:     tm.ID,
					URL:    tm.URL,
					FireAt: tm.FireAt,
				}
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "empty URL is invalid",
			fields: fields{
				ID:     "1",
				URL:    url.URL{},
				FireAt: time.Now(),
			},
			wantErr: assert.Error,
		},
		{
			name: "empty ID is iinvalid",
			fields: fields{
				ID: "",
				URL: url.URL{
					Scheme: "http://",
					Host:   "valid.url",
				},
				FireAt: time.Now(),
			},
			wantErr: assert.Error,
		},
		{
			name: "zero FireAt is invalid",
			fields: fields{
				ID: "1",
				URL: url.URL{
					Scheme: "http://",
					Host:   "valid.url",
				},
				FireAt: time.Time{},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Timer{
				ID:     tt.fields.ID,
				URL:    tt.fields.URL,
				FireAt: tt.fields.FireAt,
			}
			tt.wantErr(t1, t.Validate(), "Validate()")
		})
	}
}
