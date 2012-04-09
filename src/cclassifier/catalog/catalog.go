package catalog

import (
  "encoding/gob"
  "hash/adler32"
  "os"
)

/* a Catalog is just a simple list of uint32 which are in fact
 a Hash32 of the file
*/

type Catalog struct {
  Filename string
  Files []uint32
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
	file, err := os.OpenFile(c.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(&serializableCatalog{c.Files})
    return
}
func (c *Catalog) Include(content []byte) (ret bool) {
  adler32.Checksum([]byte(content))
  return false

}

