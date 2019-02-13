package cronlendar_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/burrbd/cronlendar"
	"github.com/cheekybits/is"
)

const timeLayout = "15:04"

func TestRun(t *testing.T) {
	is := is.New(t)

	exp := `1:30 tomorrow - /bin/run_me_daily
16:45 today - /bin/run_me_hourly
16:10 today - /bin/run_me_every_minute
19:00 today - /bin/run_me_sixty_times`

	stdin := strings.NewReader(`30 1 /bin/run_me_daily
45 * /bin/run_me_hourly
* * /bin/run_me_every_minute
* 19 /bin/run_me_sixty_times
`)

	tm, err := time.Parse(timeLayout, "16:10")
	is.NoErr(err)

	var stdout bytes.Buffer
	cronlendar.Run(tm, stdin, &stdout)

	is.Equal(exp, stdout.String())
}