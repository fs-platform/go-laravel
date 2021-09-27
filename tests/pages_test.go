package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHomePage(t *testing.T) {
	baseURL := "http://localhost:3000"
	var (
		resp *http.Response
		err  error
	)
	resp, err = http.Get(baseURL + "/")
	assert.NoError(t, err, "有错误发生，err 不为空")
	assert.Equal(t, 200, resp.StatusCode, "应返回状态码 200")
}
