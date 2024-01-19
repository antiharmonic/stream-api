package main

import (
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/gookit/config/v2"
  "github.com/gookit/config/v2/ini"
  "database/sql"
  _ "github.com/lib/pq"
  "net/url"
)

//type NullString struct {
// sql.NullString
//}

type StreamItem struct {
  Id      int     `json:"id"`
  Title   string  `json:"title"`
  Comment string  `json:"comment"`
  Url     string  `json:"url"`
}

var db *sql.DB

func main() {
  config.AddDriver(ini.Driver)
  err := config.LoadFiles("config.ini") //should make this configurable too
  if err != nil {
    panic(err)
  }
  connstr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
    config.String("database.user"),
    url.QueryEscape(config.String("database.pass")),
    config.String("database.host"),
    config.String("database.port"),
    config.String("database.db_name"),
  )

  db, err = sql.Open("postgres", connstr)
  if err != nil {
    panic(err)
  }
  http.HandleFunc("/hello", HelloHandler)
  http.HandleFunc("/list", ListHandler)
  fmt.Println(http.ListenAndServe("0.0.0.0:" + config.String("server.port"), nil))
}

func HelloHandler(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w, "hello\n")
}

func ListHandler(w http.ResponseWriter, req *http.Request) {
  rows, err := db.Query("select id, title, comment, url from stream order by posted_date")
  if err != nil {
    panic(err)
  }
  defer rows.Close()

  items := make([]StreamItem, 0)
  for rows.Next() {
    var item StreamItem //:= StreamItem{Id: 0, Title: "Tenet - Algorithm", Comment: "The bit at 3:30 always gets me"}
    if err := rows.Scan(&item.Id, &item.Title, &item.Comment, &item.Url); err != nil {
      fmt.Println(err)
      panic(err)
    }
    items = append(items, item)
  }
  json_items, err := json.Marshal(items)
  if err != nil {
    panic(err)
  }
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, string(json_items))
}
