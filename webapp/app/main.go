package main

import (
    "net/http"
    "html/template"
    "cclassifier"
)

func init() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/guess", guess_handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  t := template.Must(template.New("foo").ParseFiles("templates/base.html"))
  tc := make(map[string]interface{})
  tc["ProjectName"] = "CopiePrivee"
  tc["Title"] = "Page title"
  tc["Code"] = ""
  if err := t.ExecuteTemplate(w, "base.html", tc); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

type Result struct {
  Score float64
  Lang string
}

func guess_handler(w http.ResponseWriter, r *http.Request) {
  
  r.ParseForm()
  code := r.Form.Get("code")
  my_scanner := scanner.InitFromFile(".", "plop")

  scores, langs := my_scanner.Guess(code)
  sc := make([]Result, len(scores), len(scores))
  for index, score := range scores {
    r := Result{Score: score, Lang: langs[index]}
    sc[index] = r
  }


  t := template.Must(template.New("foo").ParseFiles("templates/base.html"))
  tc := make(map[string]interface{})
  tc["ProjectName"] = "CopiePrivee"
  tc["Title"] = "Page title"
  tc["Res"] = sc
  tc["Code"] = code
  if err := t.ExecuteTemplate(w, "base.html", tc); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
