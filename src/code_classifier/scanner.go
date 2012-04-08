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
}

var (
  VALID_LANGUAGES = map[string] *Syntax {
  "go" : &Syntax{"go",    regexp.MustCompile("go"), bayesian.Class("go")},
  "ruby" : &Syntax{"ruby",  regexp.MustCompile("rb"), bayesian.Class("ruby")}}
  BAYESIAN_CLASSES [] bayesian.Class
  EXTENSIONS_TO_SYNTAX map[*regexp.Regexp] * Syntax
  CLASS_TO_SYNTAX map[bayesian.Class] * Syntax
)

// Prepare the other data structures used by the scanner
func init() {
  EXTENSIONS_TO_SYNTAX := make(map[*regexp.Regexp]*Syntax)
  CLASS_TO_SYNTAX := make(map[bayesian.Class]*Syntax)

  for _, syntax := range(VALID_LANGUAGES) {
    fmt.Printf("Adding %s language.\n", syntax.name)
    EXTENSIONS_TO_SYNTAX[syntax.extensions] = syntax
    CLASS_TO_SYNTAX[syntax.class] = syntax
  }
  BAYESIAN_CLASSES := make([] bayesian.Class, 0)
  for re, _ := range CLASS_TO_SYNTAX {
    BAYESIAN_CLASSES = append(BAYESIAN_CLASSES, re)
  }
}

// This is the only constructor you should use
func InitFromFile(path string) ( * Scanner) {
  log.Printf("Loading %s", path)
  classifier, err := bayesian.NewClassifierFromFile(path)
  if err != nil {
    if os.IsNotExist(err) {
      classifier = bayesian.NewClassifier(BAYESIAN_CLASSES ...)
    }
  }

  return &Scanner{classifier, path}
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
  for re, syntax:= range(EXTENSIONS_TO_SYNTAX) {
    if re.MatchString(filepath.Ext(path)) {
      lang = EXTENSIONS_TO_SYNTAX[re].class
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
