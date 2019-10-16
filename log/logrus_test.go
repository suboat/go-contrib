package log

import (
	"testing"
)

func Test_File(t *testing.T) {
	var (
		l = NewLogFile("./test.log")
	)
	l.Infof("ha ha")
	l.Errorf("la la")
}
