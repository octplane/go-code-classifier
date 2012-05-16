package catalog

import (
	"encoding/gob"
	"hash/adler32"
	"os"
	"sort"
)

/* a Catalog is just a simple list of uint32 which are in fact
a Hash32 of the file
*/

type Catalog struct {
	Filename string
	Files    Int32Slice
}
type serializableCatalog struct {
	files Int32Slice
}

type Int32Slice []uint32

func (p Int32Slice) Len() int           { return len(p) }
func (p Int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Int32Slice) Less(i, j int) bool { return p[i] < p[j] }

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
func (c *Catalog) Append(content []byte) {
	crc := adler32.Checksum([]byte(content))
	c.Files = append(c.Files, crc)
}
func (c *Catalog) Include(content []byte) (ret bool) {
	crc := adler32.Checksum([]byte(content))
	sort.Sort(&c.Files)
	exists := sort.Search(len(c.Files), func(i int) bool { return c.Files[i] >= crc })

	return exists != len(c.Files)
	return false

}
