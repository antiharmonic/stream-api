package main

import (
  "fmt"
  "net/http"
  "encoding/json"
)

type StreamItem struct {
  Id      int     `json:"id"`
  Title   string  `json:"title"`
  Comment string  `json:"comment"`
}

func main() {
  http.HandleFunc("/api/hello", HelloHandler)
  http.HandleFunc("/api/list", ListHandler)
  fmt.Println(http.ListenAndServe("0.0.0.0:5689", nil))
}

func HelloHandler(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w, "hello\n")
}

func ListHandler(w http.ResponseWriter, req *http.Request) {
  item := StreamItem{Id: 0, Title: "Tenet - Algorithm", Comment: "The bit at 3:30 always gets me"}
  items := []StreamItem{item}
  json_items, err := json.Marshal(items)
  if err != nil {
    panic(err)
  }
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, string(json_items))
}
