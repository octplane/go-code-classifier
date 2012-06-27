package main

import (
    "net/http"
    "html/template"
)

func init() {
    http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  t := template.Must(template.New("foo").ParseFiles("templates/base.html"))
  tc := make(map[string]interface{})
  tc["ProjectName"] = "CopiePrivee"
  tc["Title"] = "Page title"
  if err := t.ExecuteTemplate(w, "base.html", tc); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
