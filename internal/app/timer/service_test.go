package timer_test

import (
	"context"
	"fmt"
	"github.com/cubny/httpqueue/internal/app/timer"
	timer2 "github.com/cubny/httpqueue/internal/infra/redis/timer"
	mocks "github.com/cubny/httpqueue/internal/mocks/app/timer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestServiceImp_CreateTimer(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name       string
		repoError  error
		cmd        timer.SetTimerCommand
		wantURLRaw string
		wantFireAt time.Time
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:      "happy",
			repoError: nil,
			cmd: timer.SetTimerCommand{
				URLRaw: "http://valid.url",
			},
			wantFireAt: time.Now(),
			wantURLRaw: "http://valid.url",
			wantErr:    assert.NoError,
		},
		{
			name:      "repo errors",
			repoError: assert.AnError,
			cmd: timer.SetTimerCommand{
				URLRaw: "http://valid.url",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepo(ctrl)
			repo.EXPECT().AddTimer(gomock.Any(), gomock.Any()).Return(tt.repoError)

			s, err := timer.NewService(repo)
			require.NoError(t, err)

			got, err := s.CreateTimer(context.Background(), tt.cmd)
			if !tt.wantErr(t, err, fmt.Sprintf("CreateTimer(ctx, %v)", tt.cmd)) {
				return
			}
			if err != nil {
				return
			}

			assert.Equal(t, fmt.Sprintf("%s/%s", tt.wantURLRaw, got.ID), got.URL.String())
			assert.WithinDuration(t, tt.wantFireAt, got.FireAt, time.Second)
		})
	}
}

func TestServiceImp_GetTimer(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name                   string
		timerID                string
		repoFindError          error
		repoFindTimer          *timer.Timer
		isRepoIsArchivedCalled bool
		repoIsArchivedError    error
		repoIsArchived         bool
		want                   *timer.Timer
		wantErr                error
	}{
		{
			name:                   "found, not archived",
			timerID:                "1",
			repoFindError:          nil,
			repoFindTimer:          &timer.Timer{ID: "1"},
			isRepoIsArchivedCalled: false,
			want:                   &timer.Timer{ID: "1"},
			wantErr:                nil,
		},
		{
			name:                   "repo error on deserialization",
			timerID:                "1",
			repoFindError:          timer2.ErrDeserialization,
			repoFindTimer:          nil,
			isRepoIsArchivedCalled: false,
			want:                   nil,
			wantErr:                timer2.ErrDeserialization,
		},
		{
			name:                   "archived",
			timerID:                "1",
			repoFindError:          nil,
			repoFindTimer:          nil,
			isRepoIsArchivedCalled: true,
			repoIsArchived:         true,
			repoIsArchivedError:    nil,
			want:                   nil,
			wantErr:                timer.ErrTimerArchived,
		},
		{
			name:                   "not found, isArchived errors",
			timerID:                "1",
			repoFindError:          nil,
			repoFindTimer:          nil,
			isRepoIsArchivedCalled: true,
			repoIsArchived:         false,
			repoIsArchivedError:    assert.AnError,
			want:                   nil,
			wantErr:                assert.AnError,
		},
		{
			name:                   "not found, not archived",
			timerID:                "1",
			repoFindError:          nil,
			repoFindTimer:          nil,
			isRepoIsArchivedCalled: true,
			repoIsArchived:         false,
			repoIsArchivedError:    nil,
			want:                   nil,
			wantErr:                timer.ErrTimerNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepo(ctrl)
			repo.EXPECT().Find(gomock.Any(), tt.timerID).Return(tt.repoFindTimer, tt.repoFindError)
			if tt.isRepoIsArchivedCalled {
				repo.EXPECT().IsArchived(gomock.Any(), tt.timerID).Return(tt.repoIsArchived, tt.repoIsArchivedError)
			}

			s, err := timer.NewService(repo)
			require.NoError(t, err)

			got, err := s.GetTimer(context.Background(), tt.timerID)
			if err != nil {
				assert.EqualError(t, tt.wantErr, err.Error())
				return
			}

			assert.Equal(t, tt.want, got)

		})
	}
}
