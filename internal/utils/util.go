package utils

import (
	"encoding/json"
	. "fmt"
	. "golang.org/x/exp/constraints"
	"net/http"
	. "time"
)

func Abs[T Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func FromRemoteJson[T any](url string) (*T, error) {
	c := &http.Client{Timeout: 1 * Second}

	res, err := c.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var ret T
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func FormatRelativeToday(target Time) string {
	now := Now()

	nowDate := now.Truncate(24 * Hour)
	targetDate := target.Truncate(24 * Hour)

	daysApart := int(targetDate.Sub(nowDate).Hours() / 24.0)
	nowYear, nowWeek := nowDate.ISOWeek()
	targetYear, targetWeek := targetDate.ISOWeek()

	var weeksApart int
	switch targetYear - nowYear {
	case -1:
		weeksApart = targetWeek - 52 - nowWeek
	case 0:
		weeksApart = targetWeek - nowWeek
	case 1:
		weeksApart = targetWeek + 52 - nowWeek
	default:
		panic("date results should not have been more than a year apart!")
	}

	weekday := targetDate.Weekday()

	switch daysApart {
	case -1:
		return "yesterday"
	case 0:
		return "today"
	case 1:
		return "tomorrow"
	default:
		if targetDate.Before(nowDate) {
			return Sprintf("%d days ago", Abs(daysApart))
		} else {
			switch weeksApart {
			case 0:
				return Sprintf("this %s", weekday)
			case 1:
				return Sprintf("next %s", weekday)
			default:
				return Sprintf("%s in %d weeks", weekday, Abs(weeksApart))
			}
		}
	}
}
