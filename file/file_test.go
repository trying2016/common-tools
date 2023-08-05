package file

import "testing"

func TestCopyFile(t *testing.T) {
	err := CopyFile("/Volumes/video/奇异博士2.mkv", "/Users/trying/Downloads/奇异博士2.mkv", func(current, total int64) {
		t.Logf("%.02f", float64(current)/float64(total)*100)
	})
	if err != nil {
		t.Fatal(err)
	}
}
