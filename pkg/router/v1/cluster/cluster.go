package cluster

import (
	"errors"
	"github.com/gin-gonic/gin"
	"ko3-gin/pkg/constant"
	clusterModel "ko3-gin/pkg/model/cluster"
	commonModel "ko3-gin/pkg/model/common"
	"ko3-gin/pkg/router/v1/common"
	clusterService "ko3-gin/pkg/service/cluster"
	"net/http"
)

func List(ctx *gin.Context) {
	models, err := clusterService.List()
	items := make([]Cluster, 0)
	for _, model := range models {
		items = append(items, FromModel(model))
	}
	if err != nil {
		_ = ctx.Error(err)
	}
	ctx.JSON(http.StatusOK, ListResponse{Items: items})
}

func Page(ctx *gin.Context) {
	page := ctx.GetBool("page")
	if page {
		pageNum := ctx.GetInt(constant.PageNumQueryKey)
		pageSize := ctx.GetInt(constant.PageSizeQueryKey)
		models, total, err := clusterService.Page(pageNum, pageSize)
		if err != nil {
			_ = ctx.Error(err)
		}
		var resp = common.PageResponse{
			Items: []interface{}{},
			Total: total,
		}
		for _, model := range models {
			resp.Items = append(resp.Items, model)
		}
		ctx.JSON(http.StatusOK, resp)
	} else {
		_ = ctx.Error(common.InvalidPageParam)
	}
}

var invalidClusterName = errors.New("invalid cluster name")

func Get(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		_ = ctx.Error(invalidClusterName)
	}
	model, err := clusterService.Get(name)
	if err != nil {
		_ = ctx.Error(err)
	}
	ctx.JSON(http.StatusOK, GetResponse{Item: FromModel(model)})

}

func Create(ctx *gin.Context) {
	var req CreateRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	model := clusterModel.Cluster{
		BaseModel: commonModel.BaseModel{
			Name: req.Name,
		},
	}
	err = clusterService.Save(&model)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, CreateResponse{Item: FromModel(model)})
}

func Update(ctx *gin.Context) {
	var req UpdateRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	model := clusterModel.Cluster{
		BaseModel: commonModel.BaseModel{
			Name: req.Name,
		},
	}
	err = clusterService.Save(&model)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, UpdateResponse{Item: FromModel(model)})

}

func Delete(ctx *gin.Context) {
	var req DeleteRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		_ = ctx.Error(err)
	}
	err = clusterService.Delete(req.Name)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, DeleteResponse{})
}

func Batch(ctx *gin.Context) {
	var req BatchRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	models := make([]clusterModel.Cluster, 0)
	for _, item := range req.Items {
		models = append(models, ToModel(item))
	}
	models, err = clusterService.Batch(req.Operation, models)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	var resp BatchResponse

	for _, model := range models {
		resp.Items = append(resp.Items, FromModel(model))
	}
	ctx.JSON(http.StatusOK, resp)
}
