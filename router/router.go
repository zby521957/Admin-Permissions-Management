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
		//protected.GET("/users", controller.GetUsers)              // 列表
		//protected.GET("/users/:id", controller.GetUser)           // 详情
		//protected.PUT("/users/:id", controller.UpdateUser)        // 更新
		//protected.DELETE("/users/:id", controller.DeleteUser)     // 删除
		//protected.POST("/users/:id/roles", controller.AssignRole) // 分配角色
		//
		//// 角色管理
		//protected.GET("/roles", controller.GetRoles)
		//protected.GET("/roles/:id", controller.GetRole)
		//protected.POST("/roles", controller.CreateRole)
		//protected.PUT("/roles/:id", controller.UpdateRole)
		//protected.DELETE("/roles/:id", controller.DeleteRole)
		//protected.POST("/roles/:id/permissions", controller.AssignPermissions)
		//
		//// 权限管理
		//protected.POST("/permissions", controller.CreatePermission)
		//protected.GET("/permissions", controller.GetPermissions)
		//protected.GET("/permissions/:id", controller.GetPermission)
		//protected.PUT("/permissions/:id", controller.UpdatePermission)
		//protected.DELETE("/permissions/:id", controller.DeletePermission)

		// 用户管理（需要 user:list、user:update、user:delete 等权限）
		protected.GET("/users", middleware.RequirePermission("user:list"), controller.GetUsers)
		protected.GET("/users/:id", middleware.RequirePermission("user:list"), controller.GetUser)
		protected.PUT("/users/:id", middleware.RequirePermission("user:update"), controller.UpdateUser)
		protected.DELETE("/users/:id", middleware.RequirePermission("user:delete"), controller.DeleteUser)
		protected.POST("/users/:id/roles", middleware.RequirePermission("user:update"), controller.AssignRole)

		// 角色管理
		protected.GET("/roles", middleware.RequirePermission("role:list"), controller.GetRoles)
		protected.GET("/roles/:id", middleware.RequirePermission("role:list"), controller.GetRole)
		protected.POST("/roles", middleware.RequirePermission("role:create"), controller.CreateRole)
		protected.PUT("/roles/:id", middleware.RequirePermission("role:update"), controller.UpdateRole)
		protected.DELETE("/roles/:id", middleware.RequirePermission("role:delete"), controller.DeleteRole)
		protected.POST("/roles/:id/permissions", middleware.RequirePermission("role:update"), controller.AssignPermissions)

		// 权限管理
		protected.GET("/permissions", middleware.RequirePermission("permission:list"), controller.GetPermissions)
		protected.GET("/permissions/:id", middleware.RequirePermission("permission:list"), controller.GetPermission)
		protected.POST("/permissions", middleware.RequirePermission("permission:create"), controller.CreatePermission)
		protected.PUT("/permissions/:id", middleware.RequirePermission("permission:update"), controller.UpdatePermission)
		protected.DELETE("/permissions/:id", middleware.RequirePermission("permission:delete"), controller.DeletePermission)
	}

	return r
}
