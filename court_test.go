package courttime

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCourt_properties(t *testing.T) {
	t.Parallel()
	checkFn := func(
		start time.Time,
		end time.Time,
		gameDuration time.Duration,
		slackTime time.Duration,
	) bool {
		court, err := BuildCourt("testCourt", start, end, gameDuration, slackTime)
		if err != nil {
			return assert.ErrorIs(t, err, ErrInvalidInput)
		}
		assert.NotNil(t, court)
		// slots should be sorted by start
		for i := range len(court.Slots) - 1 {
			assert.Less(t, court.Slots[i].Start, court.Slots[i+1].Start)
		}
		// slots should not overlap
		for i := range len(court.Slots) - 1 {
			assert.LessOrEqual(t, court.Slots[i].End, court.Slots[i+1].Start)
		}
		// slots should have a slackTime break between end of slot[i] and start of slot [i+1]
		for i := range len(court.Slots) - 1 {
			assert.Equal(t, slackTime, court.Slots[i+1].Start.Sub(court.Slots[i].End))
		}
		// slots should never go past end (slot.Start and slot.End must be before end)
		assert.False(
			t,
			court.Slots[len(court.Slots)-1].End.After(end),
			"%v should never end be after %v",
			court.Slots[len(court.Slots)-1],
			end,
		)
		// slots should have start < end
		for _, slot := range court.Slots {
			assert.Less(t, slot.Start, slot.End)
		}
		// slots should have end - start = gameDuration
		for _, slot := range court.Slots {
			assert.Equal(t, slot.End.Sub(slot.Start), gameDuration)
		}
		return true
	}
	qCfg := &quick.Config{
		Values: func(v []reflect.Value, r *rand.Rand) {
			start := time.Unix(r.Int63(), 0)
			end := calcBoundedTime(start, 8*time.Hour, r)
			v[0] = reflect.ValueOf(start)
			v[1] = reflect.ValueOf(end)
			v[2] = reflect.ValueOf(calcBoundedDuration(30*time.Minute, r))
			v[3] = reflect.ValueOf(calcBoundedDuration(30*time.Minute, r))
		},
	}
	require.NoError(t, quick.Check(checkFn, qCfg))
}

func TestBuildCourt(t *testing.T) {
	t.Parallel()

	startTime := time.Date(2026, time.April, 1, 1, 0, 0, 0, time.UTC)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		start        time.Time
		end          time.Time
		gameDuration time.Duration
		slackTime    time.Duration
		want         *Court
		wantErr      error
	}{
		{
			name:         "negative game duration",
			gameDuration: -1 * time.Minute,
			wantErr:      ErrInvalidGameDuration,
		},
		{
			name:         "zero game duration",
			gameDuration: 0,
			wantErr:      ErrInvalidGameDuration,
		},
		{
			name:         "negative slack time",
			gameDuration: 10 * time.Minute,
			slackTime:    -1 * time.Minute,
			wantErr:      ErrInvalidSlack,
		},
		{
			name:         "No way to schedule a game",
			gameDuration: 10 * time.Minute,
			start:        startTime,
			end:          startTime.Add(5 * time.Minute),
			wantErr:      ErrNoGames,
		},
		{
			name:         "Simple happy 1 hours no slack",
			start:        startTime,
			end:          startTime.Add(1 * time.Hour),
			gameDuration: 20 * time.Minute,
			want: &Court{
				Name: "testCourt",
				Slots: []Slot{
					{
						Start: startTime,
						End:   startTime.Add(20 * time.Minute),
					},
					{
						Start: startTime.Add(20 * time.Minute),
						End:   startTime.Add(40 * time.Minute),
					},
					{
						Start: startTime.Add(40 * time.Minute),
						End:   startTime.Add(60 * time.Minute),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := BuildCourt("testCourt", tt.start, tt.end, tt.gameDuration, tt.slackTime)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got, "want: %v\ngot:  %v", tt.want, got)
		})
	}
}

func calcBoundedTime(start time.Time, max time.Duration, r *rand.Rand) time.Time {
	return start.Add(calcBoundedDuration(max, r))
}

func calcBoundedDuration(max time.Duration, r *rand.Rand) time.Duration {
	return time.Duration(r.Int63n(2*int64(max)) - int64(max))
}

func TestGetCap(t *testing.T) {
	t.Parallel()
	checkFn := func(
		start time.Time,
		r *rand.Rand,
	) bool {
		maxBooking := 8 * time.Hour
		end := calcBoundedTime(start, maxBooking, r)
		assert.Less(t, end.Sub(start), 8*time.Hour)
		return true
	}
	qCfg := &quick.Config{
		Values: func(v []reflect.Value, r *rand.Rand) {
			v[0] = reflect.ValueOf(time.Unix(r.Int63(), 0))
			v[1] = reflect.ValueOf(r)
		},
	}
	assert.NoError(t, quick.Check(checkFn, qCfg))
}
