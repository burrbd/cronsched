package cronsched

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TimeVal struct {
	IsSet bool
	time.Time
}

func (arg TimeVal) String() string {
	return fmt.Sprintf("%d:%d", arg.Time.Hour(), arg.Time.Minute())
}

func (arg *TimeVal) Set(flag string) error {
	hourMinute := strings.Split(flag, ":")
	if len(hourMinute) != 2 {
		return newDecodeTimeErr(flag)
	}
	hour, err := strconv.Atoi(hourMinute[0])
	if err != nil {
		return newDecodeTimeErr(flag)
	}
	minute, err := strconv.Atoi(hourMinute[1])
	if err != nil {
		return newDecodeTimeErr(flag)
	}
	year, month, day := time.Now().Date()
	arg.Time = time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
	arg.IsSet = true
	return nil
}

func newDecodeTimeErr(timeArg string) error {
	return fmt.Errorf("invalid time arg: %s", timeArg)
}
