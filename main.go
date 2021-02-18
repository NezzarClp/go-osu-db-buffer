package main

import (
	"fmt"
	"github.com/NezzarClp/go-osu-db-buffer/v2/internal/buffer"
)

type Collection struct {
	Name string
	Maps []string
}

type CollectionData struct {
	Path string
	readBuf buffer.Buffer
	ver int
	Collections []Collection
}

func (colData *CollectionData) Load() error {
	b := buffer.Buffer { Path: "./collection.db" }
	err := b.Load()

	if err != nil {
		return err
	}

	colData.ver, _ = b.ReadInt()
	numCol, _ := b.ReadInt()

	for i := 0; i < numCol; i++ {
		name, _ := b.ReadString()
		numMaps, _ := b.ReadInt()

		var m []string

		fmt.Println(name)

		for j := 0; j < numMaps; j++ {
			md5, _ := b.ReadString()

			m = append(m, md5)
		}

		col := Collection { Name: name, Maps: m }

		colData.Collections = append(colData.Collections, col)
	}

	colData.readBuf = b

	return nil
}

func (colData *CollectionData) Save(path string) error {
	b := buffer.Buffer { Path: path }
	b.Load()

	b.WriteInt(colData.ver)

	numCol := len(colData.Collections)

	b.WriteInt(numCol)

	for _, col := range colData.Collections {
		b.WriteString(col.Name)
		b.WriteInt(len(col.Maps))

		for _, md5 := range col.Maps {
			b.WriteString(md5)
		}
	}

	return nil
}