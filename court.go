package courttime

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrInvalidGameDuration = fmt.Errorf("invalid game duration: %w", ErrInvalidInput)
	ErrInvalidSlack        = fmt.Errorf("invalid slack time: %w", ErrInvalidInput)
	ErrStartAfterEnd       = fmt.Errorf("start after end: %w", ErrInvalidInput)
	ErrNoGames             = fmt.Errorf("no games possible: %w", ErrInvalidInput)
)

// Court represents a Court
type Court struct {
	Name  string
	Slots []Slot
}

// BuildCourt creates a Court instance with calculated game "slots"
func BuildCourt(
	name string,
	start time.Time,
	end time.Time,
	gameDuration time.Duration,
	slackTime time.Duration) (*Court, error) {
	if err := validate(start, end, gameDuration, slackTime); err != nil {
		return nil, err
	}
	court := &Court{
		Name:  name,
		Slots: buildSlots(start, end, gameDuration, slackTime),
	}
	return court, nil
}

func validate(start time.Time, end time.Time, gameDuration time.Duration, slackTime time.Duration) error {
	if gameDuration == 0 || gameDuration.Abs() != gameDuration {
		return fmt.Errorf("%v: %w", gameDuration, ErrInvalidGameDuration)
	}
	if slackTime.Abs() != slackTime {
		return fmt.Errorf("%v: %w", slackTime, ErrInvalidSlack)
	}
	if start.After(end) {
		return fmt.Errorf("start {%v} end {%v}: %w", start, end, ErrStartAfterEnd)
	}
	if start.Add(gameDuration).After(end) {
		return ErrNoGames
	}
	return nil
}

func buildSlots(start time.Time, end time.Time, gameDuration time.Duration, slackTime time.Duration) []Slot {
	out := []Slot{
		{
			Start: start,
			End:   start.Add(gameDuration),
		},
	}

	currStart := start.Add(gameDuration + slackTime)
	for !currStart.After(end) {
		gameEnd := currStart.Add(gameDuration)
		if gameEnd.After(end) {
			return out
		}
		game := Slot{
			Start: currStart,
			End:   gameEnd,
		}
		out = append(out, game)
		currStart = currStart.Add(gameDuration + slackTime)
	}

	return out
}

type Slot struct {
	// Start is the start of the game
	Start time.Time
	// End is the end of the game
	End time.Time
}
