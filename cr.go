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
    "math/rand"
    "strings"
    "strconv"
    "time"
)

const (
    PORT = "8081"
    HTML = "/home/graham/Projects/Go/src/carriagereturn/index.html"
)

var (
    allids []int
)

func main() {

    allids = ids()

    fmt.Println("carriagereturn listening on port", PORT)

    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":" + PORT, nil))
}

type Context struct {
    Entry *Entry
}

func handler(response http.ResponseWriter, request *http.Request) {

    path := strings.Trim(request.URL.Path, "/")
    entryId, converr := strconv.Atoi(path)
    if converr != nil {
        redirectToEntry(response, request)
        return
    }

    entry := LoadEntry(entryId)

    tmpl, terr := template.ParseFiles(HTML)
    if terr != nil {
        log.Fatal(terr)
    }

    context := Context{Entry:entry}
    tmpl.Execute(response, context)
}

// No entryId in URL, pick one at random and redirect to it
func redirectToEntry(response http.ResponseWriter, request *http.Request) {

    // Seed random generator to today, so all requests in same day get same entry
    now := time.Now()
    mid := time.Date(
        now.Year(),
        now.Month(),
        now.Day(),
        12, 0, 0, 0,
        now.Location())
    rand.Seed(mid.Unix())

    choice := rand.Intn(len(allids))
    entryId := allids[choice]
    http.Redirect(response, request, "/" + strconv.Itoa(entryId) + "/", 302)
}

type Entry struct {
    Id int
    Content string
    Author string
    Tags string
}

func LoadEntry(entryId int) *Entry {

    db, dberr := sql.Open("postgres", "user=graham dbname=carriagereturn")
    if dberr != nil {
        log.Fatal("Error connecting. ", dberr)
    }
    defer db.Close()

    rows, qerr := db.Query(`SELECT content, author, tags
                            FROM cr_entry
                            WHERE id = $1`,
                            entryId)
    if qerr != nil {
        log.Fatal("Error reading from cr_entry in db", qerr)
    }

    entry := Entry{Id:entryId}
    rows.Next()
    rows.Scan(&entry.Content, &entry.Author, &entry.Tags)

    return &entry
}

func ids() []int {

    db, dberr := sql.Open("postgres", "user=graham dbname=carriagereturn")
    if dberr != nil {
        log.Fatal("Error connecting. ", dberr)
    }
    defer db.Close()

    rows, _ := db.Query("SELECT id FROM cr_entry")
    ids := make([]int, 0)
    var id int

    for rows.Next() {
        rows.Scan(&id)
        ids = append(ids, id)
    }

    return ids
}
