package cronsched_test

import (
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/cronsched"
)

func TestTimeVal_Set(t *testing.T) {
	is := is.New(t)
	flag := "16:10"
	val := &cronsched.TimeVal{}
	err := val.Set(flag)
	is.NoErr(err)
	is.Equal(16, val.Hour())
	is.Equal(10, val.Minute())
	is.True(val.IsSet)
}
