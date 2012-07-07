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
    PORT = "8082"
    ROOT = "/usr/local/carriagereturn/"
    HTML = ROOT + "index.html"
    ATOM = ROOT + "index.atom"
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

    if path == "feed" {
        atomFeed(response)
        return
    }

    // If no ID in path redirect to todays entry
    entryId, converr := strconv.Atoi(path)
    if converr != nil {
        http.Redirect(response, request, "/" + strconv.Itoa(todaysId()) + "/", 302)
        return
    }

    entry := LoadEntry(entryId)
    outputTemplate(HTML, entry, response)

}

// Write out a template with given entry. Response is finished after this runs.
func outputTemplate(tmplFilename string, entry *Entry, response http.ResponseWriter) {

    tmpl, terr := template.ParseFiles(tmplFilename)
    if terr != nil {
        log.Fatal(terr)
    }

    context := Context{Entry:entry}
    tmpl.Execute(response, context)
}

// ATOM 1.0 feed
func atomFeed(response http.ResponseWriter) {
    response.Header().Add("Content-Type", "application/atom+xml")
    entry := LoadEntry(todaysId())
    outputTemplate(ATOM, entry, response)
}

// Id of the entry for today
func todaysId() int {

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

    return entryId
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
