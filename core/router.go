package core

import (
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type Route struct {
	Node   string
	Method map[string]Action
	Regex  bool
}

type Router struct {
	ListPath []Route
}

type bySlash []Route

func (l bySlash) Len() int {
	return len(l)
}

func (l bySlash) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l bySlash) Less(i, j int) bool {
	if strings.Count(l[i].Node, "/") < strings.Count(l[j].Node, "/") {
		return false
	}

	if strings.Count(l[i].Node, "/") > strings.Count(l[j].Node, "/") {
		return true
	}
	if l[i].Node > l[j].Node {
		return true
	}
	return false
}

func (r *Router) init() {
	r.ListPath = make([]Route, 0, 2)
}

func (r *Router) AddRoute(path string, method string, ac Action) {
	//append
	r.ListPath = append(r.ListPath, Route{})
	r.ListPath[len(r.ListPath)-1].Node = path
	r.ListPath[len(r.ListPath)-1].Method = make(map[string]Action)
	r.ListPath[len(r.ListPath)-1].Method[strings.ToUpper(method)] = ac
	sort.Sort(bySlash(r.ListPath))
	log.Printf("%+v\n", Routes)
}

var Routes Router
var re = regexp.MustCompile(`{[0-9a-zA-Z]+}`)
var substitution = "[0-9A-Za-z]+"

func init() {
	Routes.init()
	log.Printf("%+v\n", Routes)
}

func Routing(w http.ResponseWriter, r *http.Request) {
	log.Println("----Routing-----")
	// log.Printf("%+v\n", *r)
	// log.Printf("%+v\n", Routes)
	// log.Printf("%+v\n", Routes.ListPath[0])
	// log.Printf("%+v\n", Routes.ListPath[0].Node)
	for i := 0; i < len(Routes.ListPath); i++ {
		// For REST
		if strings.Count(Routes.ListPath[i].Node, "{") > 0 {
			path := "^" + re.ReplaceAllString(Routes.ListPath[i].Node, substitution) + "$"
			//log.Println(path)
			//log.Println(r.URL.Path)
			reg, _ := regexp.Compile(path)
			if reg.MatchString(r.URL.Path) {
				if val, ok := Routes.ListPath[i].Method[r.Method]; ok {
					val(w, r)
					continue
				}
			}
		}
		// constant
		if Routes.ListPath[i].Node == r.URL.Path {
			if val, ok := Routes.ListPath[i].Method[r.Method]; ok {
				val(w, r)
			}
		}
	}
	w.WriteHeader(404)
}
