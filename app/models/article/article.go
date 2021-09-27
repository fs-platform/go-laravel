package article

import (
	"go_blog/pkg/route"
	"strconv"
)



func (a Article) Link(name string) string {
	return route.RouteName2URL("articles.show", "id", strconv.FormatInt(int64(a.ID), 10))
}