package cmd

import (
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"social-todo-list/common"
	"social-todo-list/middleware"
	ginitem "social-todo-list/module/item/transport/gin"
	"social-todo-list/module/upload"
	userstorage "social-todo-list/module/user/storage"
	ginuser "social-todo-list/module/user/transport/gin"
	ginuserlikeitem "social-todo-list/module/userlikeitem/transport/gin"
	"social-todo-list/plugin/sdkgorm"
	"social-todo-list/plugin/simple"
	"social-todo-list/plugin/tokenprovider/jwt"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("social-todo-list"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main.mysql", common.PluginDBMain)),
		goservice.WithInitRunnable(jwt.NewJWTProvider(common.PluginJWT)),
		goservice.WithInitRunnable(simple.NewSimplePlugin("simple")),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			engine.Use(middleware.Recover())

			// Example for Simple Plugin
			type CanGetValue interface {
				GetValue() string
			}
			log.Println(service.MustGet("simple").(CanGetValue).GetValue())
			/////////

			db := service.MustGet(common.PluginDBMain).(*gorm.DB)

			authStore := userstorage.NewSQLStore(db)

			middlewareAuth := middleware.RequiredAuth(authStore, service)

			v1 := engine.Group("/v1")
			{
				v1.PUT("/upload", upload.Upload(service))

				v1.POST("/register", ginuser.Register(service))
				v1.POST("/login", ginuser.Login(service))
				v1.GET("/profile", middlewareAuth, ginuser.Profile())

				items := v1.Group("/items", middlewareAuth)
				{
					items.POST("", ginitem.CreateItem(service))
					items.GET("", ginitem.ListItem(service))
					items.GET("/:id", ginitem.GetItem(service))
					items.PATCH("/:id", ginitem.UpdateItem(service))
					items.DELETE("/:id", ginitem.DeleteItem(service))

					items.POST("/:id/like", ginuserlikeitem.LikeItem(service))
					items.DELETE("/:id/unlike", ginuserlikeitem.UnlikeItem(service))
					items.GET("/:id/liked-users", ginuserlikeitem.ListUserLiked(service))
				}
			}

			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		})

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)
	rootCmd.AddCommand(cronDemoCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
