package ghost

import (
	"testing"
)

func TestConvert(t *testing.T) {
	g := NewGhost()
	msg, err := g.Convert("./testing/image_test.pdf", "./testing/out")
	if err != nil {
		t.Fatalf("error=%+v\n%s", err, msg)
	}
	g.ParseOutMsg(msg)
}

func TestZip(t *testing.T) {
	g := NewGhost()
	err := g.ZipDirByPath("./testing/out", "./testing/out/x.zip")
	if err != nil {
		t.Fatalf("error=%+v\n", err)
	}
}
