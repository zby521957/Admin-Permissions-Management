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
	protected.GET("/test", func(c *gin.Context) {
		c.String(200, "路由正常")
	})
	{
		// 用户管理
		protected.GET("/users", controller.GetUsers)              // 列表
		protected.GET("/users/:id", controller.GetUser)           // 详情
		protected.PUT("/users/:id", controller.UpdateUser)        // 更新
		protected.DELETE("/users/:id", controller.DeleteUser)     // 删除
		protected.POST("/users/:id/roles", controller.AssignRole) // 分配角色

		// 角色管理
		protected.GET("/roles", controller.GetRoles)
		protected.GET("/roles/:id", controller.GetRole)
		protected.POST("/roles", controller.CreateRole)
		protected.PUT("/roles/:id", controller.UpdateRole)
		protected.DELETE("/roles/:id", controller.DeleteRole)
		protected.POST("/roles/:id/permissions", controller.AssignPermissions)

		// 权限管理
		protected.POST("/permissions", controller.CreatePermission)
		protected.GET("/permissions", controller.GetPermissions)
		protected.GET("/permissions/:id", controller.GetPermission)
		protected.PUT("/permissions/:id", controller.UpdatePermission)
		protected.DELETE("/permissions/:id", controller.DeletePermission)

		// 后续角色管理路由放这里
	}

	return r
}
