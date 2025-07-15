package router

import (
	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "blog/docs" // 导入生成的docs包
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// 初始化服务
	userService := service.NewUserService()
	articleService := service.NewArticleService()
	fileService := service.NewFileService()
	testService := service.NewTestService()

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	articleHandler := handler.NewArticleHandler(articleService)
	fileHandler := handler.NewFileHandler(fileService)
	testHandler := handler.NewTestHandler(testService)
	apifoxHandler := handler.NewApifoxHandler()

	// Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Apifox导入页面
	r.GET("/apifox", apifoxHandler.GetApifoxQuickImport)

	// 公开路由
	api := r.Group("/api/v1")
	{
		// 用户管理
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		// 文章管理
		api.GET("/articles", articleHandler.GetArticles)
		api.GET("/articles/:id", articleHandler.GetArticle)

		// 测试管理
		api.POST("/test", testHandler.CreateTest)
		api.DELETE("/test/:id", testHandler.DeleteTest)
		api.PUT("/test/:id", testHandler.UpdateTest)
		api.GET("/tests", testHandler.GetTests)
		api.GET("/test/:id", testHandler.GetTest)

		// Apifox导入
		api.GET("/apifox/import", apifoxHandler.GetApifoxImportInfo)
	}

	// 需要认证的路由
	auth := api.Group("")
	auth.Use(middleware.Auth(userService))
	{
		// 用户管理
		auth.GET("/profile", userHandler.GetProfile)
		auth.POST("/logout", userHandler.Logout)
		auth.PUT("/profile", userHandler.UpdateProfile)

		// 文章管理
		auth.POST("/articles", articleHandler.CreateArticle)
		auth.PUT("/articles/:id", articleHandler.UpdateArticle)
		auth.DELETE("/articles/:id", articleHandler.DeleteArticle)

		// 文件管理
		auth.POST("/upload", fileHandler.Upload)
	}

	return r
}
