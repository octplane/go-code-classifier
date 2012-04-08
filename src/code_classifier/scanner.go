package scanner

import(
  "fmt"
  "os"
  "path/filepath"
  "io/ioutil"
  // Store the syntax extension
  "regexp"
  "strings"
  "log"
  "github.com/jbrukh/bayesian"
)

const (
    Invalid bayesian.Class = "invalid"
)

type Syntax struct {
  name string
  extensions *regexp.Regexp
  class bayesian.Class
}

type Scanner struct {
  classifier * bayesian.Classifier
  save_file string
  name_to_syntax  map[string] *Syntax
  extensions_to_syntax  map[*regexp.Regexp] *Syntax
  class_to_syntax map[bayesian.Class] *Syntax
}

func validLanguages() map[string] *Syntax {

  s := map[string] *Syntax {
  "go" : &Syntax{"go",    regexp.MustCompile("go"), bayesian.Class("go")},
  "ruby" : &Syntax{"ruby",  regexp.MustCompile("rb"), bayesian.Class("ruby")}}

  return s
}

// This is the only constructor you should use
func InitFromFile(path string) ( * Scanner) {
  log.Printf("Loading %s", path)
  vl := validLanguages()
  by_ext := make(map[*regexp.Regexp]*Syntax)
  by_class := make(map[bayesian.Class]*Syntax)

  for _, syntax := range(vl) {
    fmt.Printf("Adding %s language.\n", syntax.name)
    by_ext[syntax.extensions] = syntax
    by_class[syntax.class] = syntax
  }

  classifier, err := bayesian.NewClassifierFromFile(path)
  if err != nil {
    if os.IsNotExist(err) {
      classes := make([] bayesian.Class, 0)
      for re, _ := range by_class {
        classes = append(classes, re)
      }
      classifier = bayesian.NewClassifier(classes ...)
    }
  }

  return &Scanner{classifier, path, vl, by_ext, by_class }
}

// Scan a file or folder according to a provided Class
func (scanner * Scanner) Scan(path string) {
  
  wf := func (path string, info os.FileInfo, err error) error {
    if ! info.IsDir() {
      scanner.Classify(path)
    }
    return nil
  }
  filepath.Walk(path, wf)
}
func (scanner * Scanner) Snapshot() {
  fmt.Printf("Learned %d documents.\n", scanner.classifier.Learned())
  log.Printf("Saving %s", scanner.save_file)

  err := scanner.classifier.WriteToFile(scanner.save_file)
  if err != nil {
    log.Fatal(err)
  }
}
func (scanner * Scanner) Classify(path string) {
  fmt.Printf("Scanning %s...\n", filepath.Base(path))
  lang := Invalid
  for re, syntax:= range(scanner.extensions_to_syntax) {
    if re.MatchString(filepath.Ext(path)) {
      lang = scanner.extensions_to_syntax[re].class
      fmt.Printf("Found a %s file: %s\n", syntax.name, path)
    }
  }

  if lang != Invalid {
    contents, err := ioutil.ReadFile(path);
    if err == nil {
      scanner.classifier.Learn(strings.Split(" ",string(contents)), lang)
    } else {
      log.Fatal(err)
    }
  }
}
