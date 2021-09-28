package article

import (
	"go_blog/pkg/route"
	"strconv"
)

func (article Article) Link(name string) string {
	return route.RouteName2URL(name, "id", strconv.FormatInt(int64(article.ID), 10))
}