package timex

import (
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestLocalDiff(t *testing.T) {
	t.Log(LocalDiff())
}

func TestPbTimestamp(t *testing.T) {
	var now = time.Now()
	var ts = timestamppb.New(now)
	var ano = ts.AsTime()
	t.Log(now)
	t.Log(ano)
	t.Log(now.Equal(ano))
}
