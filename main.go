package main
//GO111MODULE=on go get -u github.com/cosmtrek/air
import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleFunc)
	http.ListenAndServe(":3000", nil)
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	switch r.URL.Path {
	case "/about":
		fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 Aron "+
			"<a href=\"Aron.Yao@feisu.com\">251109226@qq.com</a>")
		break
	case "/":
		fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
			"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
		break
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w,"<h1>请求页面未找到 :(</h1>")
	}

}
