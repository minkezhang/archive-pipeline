package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/golang/protobuf/proto"
	"github.com/minkezhang/archive-pipeline/import/googlevoice"
)

func main() {
	const root = "data/Calls"
	files := map[string]io.Reader{}
	if err := fs.WalkDir(os.DirFS(root), ".", func(p string, d fs.DirEntry, err error) error {
		p = path.Join(root, p)
		files[p], err = os.Open(p)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gv := googlevoice.New(files)
	r, err := gv.Import()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.WriteFile("/tmp/record.textproto", []byte(proto.MarshalTextString(r)), 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
