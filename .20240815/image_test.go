//go:build ignore

package ledgrid

import (
	"testing"
)

var (
	bmlFileList = []string{
		"data/test1x4.bml",
		"data/test3x4.bml",
		"data/test3x5.bml",
		"data/test3x8.bml",
	}
)

func TestOpenBlinkenFile(t *testing.T) {
	var blinkenFile *BlinkenFile

	for _, fileName := range bmlFileList {
		t.Logf("File  : %s", fileName)
		blinkenFile = ReadBlinkenFile(fileName)
		frame := blinkenFile.Frames[0]
		// t.Logf("Rows  : %+v", frame.Rows)
		t.Logf("Values: %+v", frame.Values)
	}

}
