package application

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	a := "/bind a   b  "
	t.Logf("%+v", strings.Fields(a))
}
