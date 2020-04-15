package filesystem

import "testing"

func TestFile(t *testing.T) {
	var fs = createFileSystem("store", 0755)
	t.Logf("root: %s", fs.GetRootPath())
	t.Logf("subpath: %s", fs.Path("1.jpg"))
	t.Logf("filesep: %s", fs.FileSeparator())
	t.Logf("2.png exists: %t", fs.Exists("2.png"))
	err1 := fs.Create("1.png")
	if err1 != nil {
		t.Error(err1)
	}
	err2 := fs.Delete("1.png")
	if err2 != nil {
		t.Error(err2)
	}
}
