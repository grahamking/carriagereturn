/*
Copyright 2012 Graham King <graham@gkg.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

http://www.gnu.org/licenses/agpl.html
*/
package main

import (
    "fmt"
    "log"
    "net/http"
    "text/template"
    _ "pq"
    "database/sql"
)

const (
    PORT = "8081"
    HTML = "/home/graham/Projects/Go/src/carriagereturn/index.html"
)

func main() {
    fmt.Println("carriagereturn listening on port", PORT)
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":" + PORT, nil))
}

type Context struct {
    Entry *Entry
}

func handler(response http.ResponseWriter, request *http.Request) {

    tmpl, err := template.ParseFiles(HTML)
    if err != nil {
        log.Fatal(err)
    }

    entry := LoadEntry()

    context := Context{Entry:entry}
    tmpl.Execute(response, context)
}

type Entry struct {
    Id int
    Content string
    Author string
    Tags string
}

func LoadEntry() *Entry {

    db, dberr := sql.Open("postgres", "user=graham dbname=carriagereturn")
    if dberr != nil {
        log.Fatal("Error connecting. ", dberr)
    }
    defer db.Close()

    rows, qerr := db.Query(`SELECT id, content, author, tags
                            FROM cr_entry
                            LIMIT 1`)
    if qerr != nil {
        log.Fatal("Error reading from cr_entry in db", qerr)
    }

    entry := Entry{}
    rows.Next()
    rows.Scan(&entry.Id, &entry.Content, &entry.Author, &entry.Tags)

    return &entry
}
