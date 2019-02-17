package cronsched_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/cheekybits/is"

	"github.com/burrbd/cronsched"
)

func TestRun(t *testing.T) {
	is := is.New(t)
	exp := `1:30 tomorrow - /bin/run_me_daily
16:45 today - /bin/run_me_hourly
16:10 today - /bin/run_me_every_minute
19:00 today - /bin/run_me_sixty_times
`
	stdin := strings.NewReader(`30 1 /bin/run_me_daily
45 * /bin/run_me_hourly
* * /bin/run_me_every_minute
* 19 /bin/run_me_sixty_times`)
	baseTime, _ := time.Parse(time.RFC822, "02 Jan 06 16:10 UTC")
	var stdout bytes.Buffer
	is.NoErr(cronsched.Run(stdin, &stdout, baseTime))
	is.Equal(exp, stdout.String())
}

func TestRun_EveryMinute(t *testing.T) {
	is := is.New(t)
	exp := `16:10 today - /bin/run_me_every_minute
`
	stdin := strings.NewReader(`
* * /bin/run_me_every_minute
`)
	baseTime, _ := time.Parse(time.ANSIC, "Mon Jan 2 16:10:00 2006")
	var stdout bytes.Buffer
	is.NoErr(cronsched.Run(stdin, &stdout, baseTime))
	is.Equal(exp, stdout.String())
}

func TestRun_EveryMinuteInHourDuringHour(t *testing.T) {
	is := is.New(t)
	exp := `16:10 today - /bin/run_me_sixty_times
`
	stdin := strings.NewReader(`
* 16 /bin/run_me_sixty_times
`)
	baseTime, _ := time.Parse(time.ANSIC, "Mon Jan 2 16:10:00 2006")
	var stdout bytes.Buffer
	is.NoErr(cronsched.Run(stdin, &stdout, baseTime))
	is.Equal(exp, stdout.String())
}

func TestRun_EveryMinuteInHourAfterHour(t *testing.T) {
	is := is.New(t)
	exp := `16:00 tomorrow - /bin/run_me_sixty_times
`
	stdin := strings.NewReader(`
* 16 /bin/run_me_sixty_times
`)
	baseTime, _ := time.Parse(time.ANSIC, "Mon Jan 2 17:10:00 2006")
	var stdout bytes.Buffer
	is.NoErr(cronsched.Run(stdin, &stdout, baseTime))
	is.Equal(exp, stdout.String())
}

func TestRun_EveryMinuteInHourBeforeHour(t *testing.T) {
	is := is.New(t)
	exp := `16:00 today - /bin/run_me_sixty_times
`
	stdin := strings.NewReader(`
* 16 /bin/run_me_sixty_times
`)
	baseTime, _ := time.Parse(time.ANSIC, "Mon Jan 2 15:10:00 2006")
	var stdout bytes.Buffer
	is.NoErr(cronsched.Run(stdin, &stdout, baseTime))
	is.Equal(exp, stdout.String())
}

func TestRun_EveryHourBeforeStartMinute(t *testing.T) {
	is := is.New(t)
	exp := `17:20 today - /bin/run_me_hourly
`
	stdin := strings.NewReader(`
20 * /bin/run_me_hourly
`)
	baseTime, _ := time.Parse(time.ANSIC, "Mon Jan 2 17:10:00 2006")
	var stdout bytes.Buffer
	is.NoErr(cronsched.Run(stdin, &stdout, baseTime))
	is.Equal(exp, stdout.String())
}

func TestRun_EveryHourPastStartMinuteNextDay(t *testing.T) {
	is := is.New(t)
	exp := `0:05 tomorrow - /bin/run_me_hourly
`
	stdin := strings.NewReader(`
05 * /bin/run_me_hourly
`)
	baseTime, _ := time.Parse(time.ANSIC, "Mon Jan 2 23:10:00 2006")
	var stdout bytes.Buffer
	is.NoErr(cronsched.Run(stdin, &stdout, baseTime))
	is.Equal(exp, stdout.String())
}

func TestRun_GivenInvalidHour(t *testing.T) {
	is := is.New(t)
	stdin := strings.NewReader("30 FOO /bin/run_me_daily")
	baseTime, _ := time.Parse("15:04", "16:10")
	var stdout bytes.Buffer
	is.Err(cronsched.Run(stdin, &stdout, baseTime))
}

func TestRun_GivenInvalidMinute(t *testing.T) {
	is := is.New(t)
	stdin := strings.NewReader("BAR 10 /bin/run_me_daily")
	baseTime, _ := time.Parse("15:04", "16:10")
	var stdout bytes.Buffer
	is.Err(cronsched.Run(stdin, &stdout, baseTime))
}

func TestRun_GivenInvalidLine(t *testing.T) {
	is := is.New(t)
	stdin := strings.NewReader("FOO BAR")
	baseTime, _ := time.Parse("15:04", "16:10")
	var stdout bytes.Buffer
	is.Err(cronsched.Run(stdin, &stdout, baseTime))
}
