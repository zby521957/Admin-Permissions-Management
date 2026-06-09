package router

import (
	"rbac-admin/controller"
	"rbac-admin/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// 公开路由组
	public := r.Group("/api/v1")
	{
		public.POST("/register", controller.Register)
		public.POST("/login", controller.Login)
	}

	// 需认证路由组
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		// 后续用户管理、角色管理路由放这里
	}

	return r
}
