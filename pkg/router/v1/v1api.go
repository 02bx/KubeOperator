package v1

import (
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/credential"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/host"
	"github.com/gin-gonic/gin"
)

func V1(root *gin.RouterGroup) *gin.RouterGroup {
	v1Api := root.Group("v1")
	{
		v1HostApi := v1Api.Group("/hosts")
		{
			v1HostApi.GET("/", host.List)
			v1HostApi.POST("/", host.Create)
			v1HostApi.GET("/:name/", host.Get)
			v1HostApi.PATCH("/:name/", host.Update)
			v1HostApi.DELETE("/:name/", host.Delete)
			v1HostApi.POST("/batch/", host.Batch)
			v1HostApi.PATCH("/:name/sync/", host.Sync)
		}
		v1ClusterApi := v1Api.Group("/clusters")
		{
			v1ClusterApi.GET("/", cluster.List)
			v1ClusterApi.POST("/", cluster.Create)
			v1ClusterApi.GET("/:name/", cluster.Get)
			v1ClusterApi.PATCH("/:name/", cluster.Update)
			v1ClusterApi.DELETE("/:name/", cluster.Delete)
			v1ClusterApi.POST("/batch/", cluster.Batch)
			v1ClusterApi.GET("/:name/status/", cluster.Status)
			v1ClusterApi.POST("/initial/:name/", cluster.Init)
		}
		v1CredentialApi := v1Api.Group("/credentials")
		{
			v1CredentialApi.GET("/", credential.List)
			v1CredentialApi.POST("/", credential.Create)
			v1CredentialApi.GET("/:name/", credential.Get)
			v1CredentialApi.PATCH("/:name/", credential.Update)
			v1CredentialApi.DELETE("/:name/", credential.Delete)
			v1CredentialApi.POST("/batch/", credential.Batch)
		}
	}
	return v1Api
}
