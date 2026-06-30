package router

import (
	"rbac-admin/controller"
	"rbac-admin/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化并配置所有路由
// 路由分两组：
//   - public:   无需认证的公开接口（注册、登录）
//   - protected: 需 JWT 认证 + 权限校验的保护接口
// 权限命名规范：<资源>:<操作>，如 user:list、role:create
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil) // 禁用代理信任

	// === 公开路由组（无需认证） ===
	public := r.Group("/api/v1")
	{
		public.POST("/register", controller.Register) // 用户注册
		public.POST("/login", controller.Login)       // 用户登录，返回 JWT Token
	}

	// === 需认证路由组（JWT 验证 + 权限校验） ===
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuth()) // 应用 JWT 认证中间件
	// 连通性测试路由
	protected.GET("/test", func(c *gin.Context) {
		c.String(200, "路由正常")
	})
	{
		// --- 用户管理 ---
		protected.GET("/users", middleware.RequirePermission("user:list"), controller.GetUsers)
		protected.GET("/users/:id", middleware.RequirePermission("user:list"), controller.GetUser)
		protected.PUT("/users/:id", middleware.RequirePermission("user:update"), controller.UpdateUser)
		protected.DELETE("/users/:id", middleware.RequirePermission("user:delete"), controller.DeleteUser)
		protected.POST("/users/:id/roles", middleware.RequirePermission("user:update"), controller.AssignRole)

		// --- 角色管理 ---
		protected.GET("/roles", middleware.RequirePermission("role:list"), controller.GetRoles)
		protected.GET("/roles/:id", middleware.RequirePermission("role:list"), controller.GetRole)
		protected.POST("/roles", middleware.RequirePermission("role:create"), controller.CreateRole)
		protected.PUT("/roles/:id", middleware.RequirePermission("role:update"), controller.UpdateRole)
		protected.DELETE("/roles/:id", middleware.RequirePermission("role:delete"), controller.DeleteRole)
		protected.POST("/roles/:id/permissions", middleware.RequirePermission("role:update"), controller.AssignPermissions)

		// --- 权限管理 ---
		protected.GET("/permissions", middleware.RequirePermission("permission:list"), controller.GetPermissions)
		protected.GET("/permissions/:id", middleware.RequirePermission("permission:list"), controller.GetPermission)
		protected.POST("/permissions", middleware.RequirePermission("permission:create"), controller.CreatePermission)
		protected.PUT("/permissions/:id", middleware.RequirePermission("permission:update"), controller.UpdatePermission)
		protected.DELETE("/permissions/:id", middleware.RequirePermission("permission:delete"), controller.DeletePermission)
	}

	return r
}