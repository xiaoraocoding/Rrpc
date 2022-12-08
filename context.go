package Rrpc

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type Context struct {
	W        http.ResponseWriter
	R        *http.Request
	engineer *Engine
}

func (c *Context) HTML(status int, html string) {
	c.W.WriteHeader(status)
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := c.W.Write([]byte(html))
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) HTMLTemplate(name string, funcMap template.FuncMap, data any, fileName ...string) {
	t := template.New(name)
	t.Funcs(funcMap)
	t, err := t.ParseFiles(fileName...)
	if err != nil {
		log.Println(err)
		return
	}
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(c.W, data)
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) HTMLTemplateGlob(name string, funcMap template.FuncMap, pattern string, data any) {
	t := template.New(name)
	t.Funcs(funcMap)
	t, err := t.ParseGlob(pattern)
	if err != nil {
		log.Println(err)
		return
	}
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(c.W, data)
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) Template(name string, data any) {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	template := c.engineer.HTMLRender.Template
	err := template.ExecuteTemplate(c.W, name, data)
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) JSON(status int, data any) error {
	c.W.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.W.WriteHeader(status)
	rsp, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.W.Write(rsp)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) XML(status int, data any) error {
	c.W.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.W.WriteHeader(status)
	res, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.W.Write(res)
	if err != nil {
		return err
	}
	return nil
}

// 下载文件
func (c *Context) File(filePath string) {
	http.ServeFile(c.W, c.R, filePath)
}

// 修改下载文件的名称
func (c *Context) FileAttachment(filepath, filename string) {
	if isASCII(filename) {
		c.W.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	} else {
		c.W.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(filename))
	}
	http.ServeFile(c.W, c.R, filepath)
}

// 从服务器指定文件中下载文件
func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.R.URL.Path = old
	}(c.R.URL.Path)

	c.R.URL.Path = filepath

	http.FileServer(fs).ServeHTTP(c.W, c.R)
}
