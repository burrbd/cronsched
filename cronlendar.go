package cronlendar

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
		if err = task.printSchedule(w, baseTime); err != nil {
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
	return task{minute: minute, hour: hour, cmd: fields[2]}, nil
}

type task struct {
	minute, hour int
	cmd          string
}

func (t task) printSchedule(w io.Writer, baseTime time.Time) error {
	next := t.nextRun(baseTime)
	_, err := fmt.Fprintf(w, "%s %s - %s\n",
		strings.TrimPrefix(next.Format(timeLayout), "0"),
		day(baseTime, next),
		t.cmd)
	return err
}

func (t task) nextRun(baseTime time.Time) time.Time {
	hour, minute := t.hour, t.minute
	switch {
	case t.minute == star && t.hour == star:
		hour, minute = baseTime.Hour(), baseTime.Minute()
	case t.hour == star:
		hour = baseTime.Hour()
	case t.minute == star:
		minute = 0
	}
	baseHour, baseMinute, _ := baseTime.Clock()
	delta := time.Duration(hour-baseHour)*time.Hour +
		time.Duration(minute-baseMinute)*time.Minute
	if delta < 0 {
		delta = 24*time.Hour + delta
	}
	return baseTime.Add(delta)
}

func day(baseTime, next time.Time) string {
	if baseTime.Day() != next.Day() {
		return "tomorrow"
	}
	return "today"
}
