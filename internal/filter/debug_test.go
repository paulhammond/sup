package filter

import (
	"bytes"
	"fmt"
)

func newMockDebug() mockDebug {
	return mockDebug{&bytes.Buffer{}}
}

type mockDebug struct {
	*bytes.Buffer
}

func (d *mockDebug) debugFunc(format string, a ...any) {
	fmt.Fprintf(d.Buffer, format+"\n", a...)
}
