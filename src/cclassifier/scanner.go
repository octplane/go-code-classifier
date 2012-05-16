package scanner

import (
	"cclassifier/catalog"
	"fmt"
	"github.com/jbrukh/bayesian"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	Invalid bayesian.Class = "invalid"
)

type Syntax struct {
	name       string
	extensions *regexp.Regexp
	class      bayesian.Class
}

type Scanner struct {
	classifier *bayesian.Classifier
	catalog    *catalog.Catalog
	base_fname string
}

var (
	VALID_LANGUAGES = map[string]*Syntax{
		"go":    &Syntax{"go", regexp.MustCompile("go"), bayesian.Class("go")},
		"ruby":  &Syntax{"ruby", regexp.MustCompile("rb"), bayesian.Class("ruby")},
		"shell": &Syntax{"shell", regexp.MustCompile("sh"), bayesian.Class("shell")}}

	BAYESIAN_CLASSES     []bayesian.Class
	EXTENSIONS_TO_SYNTAX map[*regexp.Regexp]*Syntax
	CLASS_TO_SYNTAX      map[bayesian.Class]*Syntax
)

// Prepare the other data structures used by the scanner
func init() {
	EXTENSIONS_TO_SYNTAX = make(map[*regexp.Regexp]*Syntax)
	CLASS_TO_SYNTAX = make(map[bayesian.Class]*Syntax)

	for _, syntax := range VALID_LANGUAGES {
		fmt.Printf("Adding %s language.\n", syntax.name)
		EXTENSIONS_TO_SYNTAX[syntax.extensions] = syntax
		CLASS_TO_SYNTAX[syntax.class] = syntax
	}
	BAYESIAN_CLASSES = make([]bayesian.Class, 0)
	for re, _ := range CLASS_TO_SYNTAX {
		BAYESIAN_CLASSES = append(BAYESIAN_CLASSES, re)
	}
}

/* A constructor for a Scanner object. Creates or load a Scanner.
Bayesian data file is written to dir_name + base_name + ".bay"
Scanner data file is written to dir_name + base_name + ".sca"
*/
func InitFromFile(dir_name string, base_name string) *Scanner {
	ret := &Scanner{base_fname: filepath.Join(dir_name, base_name)}
	ret.LoadOrCreate()
	return ret
}

func (scanner *Scanner) LoadOrCreate() {
	classifier, err := bayesian.NewClassifierFromFile(scanner.BayesianFile())

	if err != nil {
		if os.IsNotExist(err) {
			classifier = bayesian.NewClassifier(BAYESIAN_CLASSES...)
		}
	}

	cat, err := catalog.NewCatalogFromFile(scanner.CatalogFile())

	if err != nil {
		if os.IsNotExist(err) {
			cat = &catalog.Catalog{scanner.CatalogFile(), make([]uint32, 0)}
		}
	}
	scanner.classifier = classifier
	scanner.catalog = cat
}

func (scanner *Scanner) BayesianFile() string {
	return scanner.base_fname + ".bay"
}
func (scanner *Scanner) CatalogFile() string {
	return scanner.base_fname + ".sca"
}

// Scan a file or folder according to a provided Class
func (scanner *Scanner) Scan(path string) {

	wf := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			scanner.Classify(path)
		}
		return nil
	}
	filepath.Walk(path, wf)
}
func (scanner *Scanner) Snapshot() {
	fmt.Printf("Scanner knows %d documents.\n", scanner.classifier.Learned())
	log.Printf("Saving %s", scanner.BayesianFile())

	err := scanner.classifier.WriteToFile(scanner.BayesianFile())
	if err != nil {
		log.Fatal(err)
	}
	err = scanner.catalog.Write()
	if err != nil {
		log.Fatal(err)
	}

}
func (scanner *Scanner) Classify(path string) {
	lang := Invalid
	for re, syntax := range EXTENSIONS_TO_SYNTAX {
		if re.MatchString(filepath.Ext(path)) {
			lang = EXTENSIONS_TO_SYNTAX[re].class
			fmt.Printf("Found a %s file: %s\n", syntax.name, path)
		}
	}

	if lang != Invalid {
		contents, err := ioutil.ReadFile(path)
		if err == nil {
			if scanner.catalog.Include(contents) {
				log.Printf("We have alreay %s in our catalog\n", path)
			} else {
				scanner.catalog.Append(contents)
				scanner.classifier.Learn(strings.Split(" ", string(contents)), lang)
			}
		} else {
			log.Fatal(err)
		}
	}
}
