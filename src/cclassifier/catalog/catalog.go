package catalog

import (
  "encoding/gob"
//  "hash/crc32"
//  "math"
  "os"
)

type Catalog struct {
  filename string
  files []uint32
}
type serializableCatalog struct {
  files []uint32
}


func NewCatalogFromFile(path string) (cat *Catalog, err error) {
  file, err := os.Open(path)
  if err != nil {
      return nil, err
  }
  dec := gob.NewDecoder(file)
  w := new(serializableCatalog)
  err = dec.Decode(w)

  return &Catalog{path, w.files}, err
}
func (c *Catalog) Write() (err error) {
	file, err := os.OpenFile(c.filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(&serializableCatalog{c.files})
    return
}


