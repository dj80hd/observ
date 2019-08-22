package app

import (
	"testing"
)

func TestJob(t *testing.T) {
	ok := Job{}.Run()
	if !ok {
		t.Error("default job did not succeed")
	}
}
