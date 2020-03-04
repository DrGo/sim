package stata

import (
	"fmt"
	"testing"
)

const (
	testDir = "/Users/salah/Dropbox/code/go/src/github.com/drgo/simStudy/stata/testing"
)

func TestRunStata(t *testing.T) {
	output, err := RunStataDo(testDir, "do.do")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	t.Logf("out: %s", output)
	dict := GetKeyValuePairs(output)
	t.Logf("dict: %v", dict)
}
