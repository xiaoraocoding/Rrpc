package Rrpc

import (
	"fmt"
	"log"
	"net/http"
)

const ANY = "ANY"

type HandlerFunc func(ctx *Context)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}
type router struct {
	routerGroup []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	g := &routerGroup{
		groupName:        name,
		handlerMap:       make(map[string]map[string]HandlerFunc),
		handlerMethodMap: make(map[string][]string),
		treeNode:         &treeNode{name: "/", children: make([]*treeNode, 0)},
	}
	r.routerGroup = append(r.routerGroup, g)
	return g
}

func (r *routerGroup) Any(name string, handlerFunc HandlerFunc) {
	r.handle(name, "ANY", handlerFunc)
}

func (r *routerGroup) Handle(name string, method string, handlerFunc HandlerFunc) {
	//method有效性做校验
	r.handle(name, method, handlerFunc)
}

func (r *routerGroup) Get(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodGet, handlerFunc)
}
func (r *routerGroup) Post(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPost, handlerFunc)
}
func (r *routerGroup) Delete(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodDelete, handlerFunc)
}
func (r *routerGroup) Put(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPut, handlerFunc)
}
func (r *routerGroup) Patch(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPatch, handlerFunc)
}
func (r *routerGroup) Options(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodOptions, handlerFunc)
}
func (r *routerGroup) Head(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodHead, handlerFunc)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	for _, group := range e.routerGroup {
		routerName := SubStringLast(r.RequestURI, "/"+group.groupName)
		// get/1
		node := group.treeNode.Get(routerName)
		if node != nil {
			//路由匹配上了
			ctx := &Context{
				W: w,
				R: r,
			}
			handle, ok := group.handlerMap[node.routerName][ANY]
			if ok {
				handle(ctx)
				return
			}
			handle, ok = group.handlerMap[node.routerName][method]
			if ok {
				handle(ctx)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "%s %s not allowed \n", r.RequestURI, method)
			return
		}
		//for name, methodHandle := range group.handleFuncMap {
		//	url := "/" + group.name + name
		//	if r.RequestURI == url {
		//
		//	}
		//}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s  not found \n", r.RequestURI)
}

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router{},
	}
}

type routerGroup struct {
	groupName        string
	handlerMap       map[string]map[string]HandlerFunc
	handlerMethodMap map[string][]string
	treeNode         *treeNode
}

func (r *routerGroup) handle(name string, method string, handlerFunc HandlerFunc) {
	_, ok := r.handlerMap[name]
	if !ok {
		r.handlerMap[name] = make(map[string]HandlerFunc)
	}
	r.handlerMap[name][method] = handlerFunc
	r.handlerMethodMap[method] = append(r.handlerMethodMap[method], name)
	methodMap := make(map[string]HandlerFunc)
	methodMap[method] = handlerFunc

	r.treeNode.Put(name)
}
func (e *Engine) Run() {
	//groups := e.router.groups
	//for _, g := range groups {
	//	for name, handle := range g.handlerMap {
	//		http.HandleFunc("/"+g.groupName+name, handle)
	//	}
	//}
	http.Handle("/", e)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
