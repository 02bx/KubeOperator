package middleware

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/gin-gonic/gin"
	"strconv"
)

func PagerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		num := ctx.Query(constant.PageNumQueryKey)
		limit := ctx.Query(constant.PageNumQueryKey)
		limitInt, err := strconv.Atoi(limit)
		if err != nil || limitInt < 0 {
			ctx.Set("page", false)
			ctx.Next()
		}
		numInt, err := strconv.Atoi(num)
		if err != nil || numInt < 0 {
			ctx.Set("page", false)
			ctx.Next()
		}
		ctx.Set("page", true)
		ctx.Set(constant.PageNumQueryKey, numInt)
		ctx.Set(constant.PageSizeQueryKey, limitInt)
		ctx.Next()
	}
}
