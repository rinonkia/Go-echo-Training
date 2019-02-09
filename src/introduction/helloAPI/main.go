package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// レイアウト適用済みのテンプレートを保存するmap
var templates map[string]*template.Template

// Template はHTMLテンプレートを利用するためのRenderer Interface
type Template struct {
}

// RenderメソッドはHTMLテンプレートにデータを埋め込んだ結果をWriterに書き込みます。
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return templates[name].ExecuteTemplate(w, "layout.html", data)
}

func main() {
	e := echo.New()

	t := &Template{}
	e.Renderer = t

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/public/css/", "./public/css/")
	e.Static("/public/js/", "./public/js/")
	e.Static("/public/img/", "./public/img/")

	// 書くルーティングに対するハンドラを設定
	e.GET("/", HandleIndexGet)
	e.GET("/hello", HandleHelloGet)
	e.POST("/hello", HandleHelloPost)
	e.GET("/hello_form", HandleHelloFormGet)
	e.GET("/api/hello", HandleAPIHelloGet)
	e.POST("/api/hello", HandleAPIHelloPost)

	// サーバーを開始
	e.Logger.Fatal(e.Start(":3000"))
}

func init() {
	loadTemplates()
}

func loadTemplates() {
	var baseTemplate = "templates/layout.html"
	templates = make(map[string]*template.Template)
	templates["hello"] = template.Must(
		template.ParseFiles(baseTemplate, "./templates/hello.html"))
	templates["hello_form"] = template.Must(
		template.ParseFiles(baseTemplate, "templates/hello_form.html"))
}

func HandleIndexGet(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "world")
}

func HandleHelloGet(c echo.Context) error {
	greetingto := c.QueryParam("greetingto")
	return c.Render(http.StatusOK, "hello", greetingto)
}

func HandleHelloPost(c echo.Context) error {
	greetingto := c.FormValue("greetingto")
	return c.Render(http.StatusOK, "hello_form", greetingto)
}

func HandleHelloFormGet(c echo.Context) error {
	return c.Render(http.StatusOK, "hello_form", nil)
}

func HandleAPIHelloGet(c echo.Context) error {
	greetingto := c.QueryParam("greetingto")
	return c.JSON(http.StatusOK, map[string]interface{}{"hello": greetingto})
}

type HelloParam struct {
	GreetingTo string `json:"greetingto"`
}

func HandleAPIHelloPost(c echo.Context) error {
	param := new(HelloParam)
	if err := c.Bind(param); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"hello": param.GreetingTo})
}
