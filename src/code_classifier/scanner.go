package scanner

import(
  "fmt"
  "os"
  "path/filepath"
  "io/ioutil"
  "strings"
  "log"
  "github.com/jbrukh/bayesian"
)


type Scanner struct {
  classifier * bayesian.Classifier
  save_file string
}

func InitFromFile(path string, classes ... bayesian.Class ) ( * Scanner, error) {
  log.Printf("Loading %s", path)
  classifier, err := bayesian.NewClassifierFromFile(path)
  if err != nil {
    if os.IsNotExist(err) {
      classifier, err := bayesian.NewClassifier(classes)
      if err != nil {
        log.Fatal(err)
      }
    }
    log.Fatal(err)
  }
  return &Scanner{classifier, path}
}

func (scanner * Scanner) Scan(path string, lang string) {
  
  wf := func (path string, info os.FileInfo, err error) error {
    if ! info.IsDir() {
      scanner.Classify(path, lang)
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
func (scanner * Scanner) Classify(path string, lang string) {
  fmt.Printf("Scanning %s...\n", filepath.Base(path))
  contents, err := ioutil.ReadFile(path);
  if err == nil {
    scanner.classifier.Learn(strings.Split(" ",string(contents)), Go)
  } else {
    log.Fatal(err)
  }
}