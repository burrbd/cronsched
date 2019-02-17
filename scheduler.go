package cronsched

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	timeLayout = "15:04"
	star       = -1
)

func Run(r io.Reader, w io.Writer, baseTime time.Time) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		task, err := decodeTask(line)
		if err != nil {
			return err
		}
		timeTillNext := task.timeTillRun(baseTime)
		if err = task.print(w, baseTime, timeTillNext); err != nil {
			return err
		}
	}
	return nil
}

func decodeTask(s string) (task, error) {
	fields := strings.Split(s, " ")
	if len(fields) != 3 {
		return task{}, fmt.Errorf("invalid task: '%s'", s)
	}
	for i := range fields[0:2] {
		if fields[i] == "*" {
			fields[i] = "-1"
		}
	}
	minute, err := strconv.Atoi(fields[0])
	if err != nil {
		return task{}, err
	}
	hour, err := strconv.Atoi(fields[1])
	if err != nil {
		return task{}, err
	}
	return task{
		minute:      minute,
		hour:        hour,
		cmd:         fields[2],
		timeTillRun: chooseTimeTillRun(minute, hour)}, nil
}

type task struct {
	minute, hour int
	cmd          string
	timeTillRun  timeTillRunFunc
}

func (t task) print(w io.Writer, baseTime time.Time, wait time.Duration) error {
	next := baseTime.Add(wait)
	timeStamp := strings.TrimPrefix(next.Format(timeLayout), "0")
	day := day(baseTime, next)

	_, err := fmt.Fprintf(w, "%s %s - %s\n", timeStamp, day, t.cmd)
	return err
}

func day(baseTime, next time.Time) string {
	if baseTime.Day() != next.Day() {
		return "tomorrow"
	}
	return "today"
}

type timeTillRunFunc func(time.Time) time.Duration

func chooseTimeTillRun(minute, hour int) timeTillRunFunc {
	switch {
	case minute == star && hour == star:
		return func(_ time.Time) time.Duration { return 0 }
	case minute != star && hour == star:
		return everyHourCalc(minute)
	case minute == star && hour != star:
		return everyMinuteForOneHourCalc(hour)
	default:
		return everyDayCalc(minute, hour)
	}
}

func everyHourCalc(minute int) timeTillRunFunc {
	return func(baseTime time.Time) time.Duration {
		_, baseMinute, _ := baseTime.Clock()
		d := time.Duration(minute-baseMinute) * time.Minute
		if d < 0 {
			d += 1 * time.Hour
		}
		return d
	}
}

func everyMinuteForOneHourCalc(hour int) timeTillRunFunc {
	return func(baseTime time.Time) time.Duration {
		baseHour, baseMinute, _ := baseTime.Clock()
		if baseHour == hour {
			return 0
		}
		d := time.Duration(hour-baseHour)*time.Hour -
			time.Duration(baseMinute)*time.Minute
		if d < 0 {
			d += 24 * time.Hour
		}
		return d
	}
}

func everyDayCalc(minute, hour int) timeTillRunFunc {
	return func(baseTime time.Time) time.Duration {
		baseHour, baseMinute, _ := baseTime.Clock()
		d := time.Duration(hour-baseHour)*time.Hour +
			time.Duration(minute-baseMinute)*time.Minute
		if d < 0 {
			d += 24 * time.Hour
		}
		return d
	}
}
