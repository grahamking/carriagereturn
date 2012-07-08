package main

import (
    "regexp"
    "net/http"
)

var (
    URLS []Route
)

type Route struct {
    Re *regexp.Regexp
    Target RouteTarget
}
type RouteTarget func (http.ResponseWriter, *http.Request, map[string] string)

func AddRoute(url string, target RouteTarget) {
    if URLS == nil {
        URLS = make([]Route, 0)
    }

    re := regexp.MustCompile(url)
    URLS = append(URLS, Route{re, target})
}

func FindRoute(url string) (Route, map[string] string) {

    var match []string
    args := make(map[string] string)

    for _, route := range URLS {

        match = route.Re.FindStringSubmatch(url)
        if len(match) == 0 {
            continue
        }

        for idx, val := range match[1:] {
            args[route.Re.SubexpNames()[1:][idx]] = val
        }

        return route, args
    }

    return Route{}, nil
}
