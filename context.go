package Rrpc

import (
	"Rrpc/render"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	engineer   *Engine
	StatusCode int
}

func (c *Context) HTML(status int, html string) error {
	//状态是200 默认不设置的话 如果调用了 write这个方法 实际上默认返回状态 200
	return c.Render(status, &render.HTML{Data: html, IsTemplate: false})
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

func (c *Context) Template(name string, data any) error {
	//状态是200 默认不设置的话 如果调用了 write这个方法 实际上默认返回状态 200
	return c.Render(http.StatusOK, &render.HTML{
		Data:       data,
		IsTemplate: true,
		Template:   c.engineer.HTMLRender.Template,
		Name:       name,
	})
}

func (c *Context) JSON(status int, data any) error {
	//状态是200 默认不设置的话 如果调用了 write这个方法 实际上默认返回状态 200
	return c.Render(status, &render.JSON{Data: data})
}
func (c *Context) XML(status int, data any) error {
	//状态是200 默认不设置的话 如果调用了 write这个方法 实际上默认返回状态 200
	return c.Render(status, &render.XML{
		Data: data,
	})
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

func (c *Context) Redirect(status int, location string) {
	if (status < http.StatusMultipleChoices || status > http.StatusPermanentRedirect) && status != http.StatusCreated {
		panic(fmt.Sprintf("Cannot redirect with status code %d", status))
	}
	http.Redirect(c.W, c.R, location, status)
}

func (c *Context) String(status int, format string, values ...any) error {
	return c.Render(status, &render.String{Format: format, Data: values})
}

func (c *Context) Render(statusCode int, r render.Render) error {
	//如果设置了statusCode，对header的修改就不生效了
	err := r.Render(c.W, statusCode)
	c.StatusCode = statusCode
	//多次调用 WriteHeader 就会产生这样的警告 superfluous response.WriteHeader
	return err
}
