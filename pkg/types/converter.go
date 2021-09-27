package types

import (
	"go_blog/pkg/logger"
	"strconv"
)

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		logger.LogError(err)
	}
	return i
}