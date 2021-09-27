package article

import (
	"go_blog/pkg/model"
	"go_blog/pkg/types"
)

// Article 文章模型
type Article struct {
	ID    int
	Title string
	Body  string
}

func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToInt(idstr)
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}
