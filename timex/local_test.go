package timex

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
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
