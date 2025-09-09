package helpers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStringFromContext(ctx *gin.Context, key string, defaultValue int64) (int64, error) {
	raw, ok := ctx.Get(key)
	if !ok {
		return defaultValue, errors.New("something went wrong") // TODO!
	}

	valStr, ok := raw.(string)
	if !ok {
		return defaultValue, errors.New("something went wrong") // TODO!
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return defaultValue, err
	}

	return val, nil
}
