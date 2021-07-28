package needle

import (
	"os/exec"
	"regexp"
	"testing"
)

func TestNeedle(t *testing.T) {
	out, err := exec.Command("needle", "-version").CombinedOutput()
	if err != nil {
		t.Fatalf(`Cannot find needle tool. Please install EMBOSS tool suite and place binaries in your PATH.`)
	}

	re := regexp.MustCompile(`([\d\.]+)`)
	ver := re.Find(out)

	t.Logf(`Found needle version: %s`, ver)
}