package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/golang/protobuf/proto"
	"github.com/minkezhang/archive-pipeline/import/googlevoice/importer"
	"github.com/minkezhang/archive-pipeline/import/interfaces"
	"github.com/minkezhang/archive-pipeline/import/record"
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

	importers := []interfaces.I{
		importer.New(files),
	}

	r := record.New(nil, nil)
	for _, i := range importers {
		if rec, err := i.Import(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			r.Merge(rec)
		}
	}

	if err := os.WriteFile("/tmp/record.textproto", []byte(proto.MarshalTextString(r.Export())), 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
