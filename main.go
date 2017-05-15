package main

import (
  "log"
  "net/http"
  "os"
  "runtime"
  "strings"

  "github.com/gorilla/websocket"
)

var Hive *hive = NewHive()
var Config *config = LoadConfig()
// TODO: These buffer sizes should be configurable.
var Upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
}

func main() {
  log.Printf("[Main] Initializing Philotic Network\n")
  log.Printf("[Main] Version: %v\n", VERSION)
  log.Printf("[Main] Port: %v\n", Config.port)
  log.Printf("[Main] Cores: %v\n", runtime.NumCPU())

  http.HandleFunc("/", ServeNewConnection)

  err := http.ListenAndServe(":" + Config.port, nil); if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}

func ServeNewConnection(w http.ResponseWriter, r *http.Request) {
  auth := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
  accessKey, err := NewAccessKey(auth); if err != nil {
    if os.Getenv("DEBUG") == "true" { log.Println("Access Key error: " + err.Error()) }
    w.Write([]byte(err.Error()))
    return
  }

  connection, err := Upgrader.Upgrade(w, r, nil); if err != nil {
    if os.Getenv("DEBUG") == "true" { log.Println("Connection Upgrader error: " + err.Error()) }
    w.Write([]byte(err.Error()))
    return
  }

  philote := NewPhilote(accessKey, connection)
  Hive.Connect <- philote
}
