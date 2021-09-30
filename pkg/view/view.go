package view

import (
	"go_blog/pkg/auth"
	"go_blog/pkg/flash"
	"go_blog/pkg/logger"
	"go_blog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

type D map[string]interface{}

func Render(w io.Writer, data D, tplFiles ...string) {
	renderTemplate(w, "app", data, tplFiles...)
}

func AuthRender(w io.Writer, data D, tplFiles ...string) {
	renderTemplate(w, "auth", data, tplFiles...)
}

func renderTemplate(w io.Writer, name string, data D, tplFiles ...string) {
	// 1 设置模板相对路径
	viewDir := "./resources/views"
	// 2. 遍历传参文件列表 Slice，设置正确的路径，支持 dir.filename 语法糖
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + "/" + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}
	// 3. 所有布局模板文件 Slice
	layoutFiles, err := filepath.Glob(viewDir + "/" + "layouts/*.gohtml")
	logger.LogError(err)
	data["flash"] = flash.All()
	// 4. 合并所有文件
	allFiles := append(layoutFiles, tplFiles...)
	// 5 解析所有模板文件
	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"isLogin":       auth.Check,
			"RouteName2URL": route.RouteName2URL,
		}).ParseFiles(allFiles...)
	logger.LogError(err)
	// 6 渲染模板
	tmpl.ExecuteTemplate(w, name, data)
}
