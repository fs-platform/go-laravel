package flash

import (
	"encoding/gob"
	"go_blog/pkg/session"
)

type flashes map[string]interface{}

var flashKey string = "_flash"
//session 中存入 map 类型必需使用 gob
func init() {
	gob.Register(flashes{})
}

func addFlash(key string, message string) {
	flashes := flashes{}
	flashes[key] = message
	session.Put(flashKey, flashes)
	session.Save()
}

func Info(message string) {
	addFlash("info", message)
}

func Success(message string) {
	addFlash("success", message)
}

func Danger(message string) {
	addFlash("danger", message)
}

func Warning(message string) {
	addFlash("warning", message)
}

func All() flashes {
	flash := session.Get(flashKey)
	flashMessage, err := flash.(flashes)
	if !err {
		return nil
	}
	defer session.Forget(flashKey)
	return flashMessage
}
