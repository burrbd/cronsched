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
	today      = "today"
	tomorrow   = "tomorrow"
)

func Run(tm time.Time, r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		task := decodeTask(scanner.Text())
		if task == nil {
			continue
		}
		task.printSchedule(w, tm)
	}
}

type Task struct {
	Minute string
	Hour   string
	Cmd    string
}

func (t *Task) printSchedule(w io.Writer, now time.Time) {
	next := t.next(now)
	fmt.Fprintf(w, "%s %s - %s\n",
		strings.TrimPrefix(next.Format(timeLayout), "0"),
		t.day(now, next),
		t.Cmd)
}

func (t *Task) next(now time.Time) time.Time {
	switch {
	case t.Minute == "*" && t.Hour == "*":
		t.Hour = fmt.Sprintf("%d", now.Hour())
		t.Minute = fmt.Sprintf("%d", now.Minute())
	case t.Hour == "*":
		t.Hour = fmt.Sprintf("%d", now.Hour())
	case t.Minute == "*":
		t.Minute = "00"
	}
	hour, _ := strconv.Atoi(t.Hour)
	minute, _ := strconv.Atoi(t.Minute)
	year, month, day := now.Date()
	next := time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
	if next.Before(now) {
		return next.Add(24 * time.Hour)
	}
	return next
}

func (t *Task) day(now, next time.Time) string {
	d := 24 * time.Hour
	if next.Truncate(d).Equal(now.Truncate(d)) {
		return today
	}
	return tomorrow
}

func decodeTask(s string) *Task {
	fields := strings.Split(s, " ")
	if len(fields) != 3 {
		return nil
	}
	return &Task{Minute: fields[0], Hour: fields[1], Cmd: fields[2]}
}
