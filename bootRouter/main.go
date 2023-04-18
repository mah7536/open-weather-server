package bootRouter

import (
	"alarm-system/router/api"
	"alarm-system/router/info"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var AppEngine *gin.Engine

// cors的設置
func SetCors(App *gin.Engine) {
	App.Use(
		cors.New(cors.Config{
			AllowHeaders:     []string{"Origin", "Content-Type", "Access-Control-Allow-Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
		}),
	)
}

func StartServer() {

	AppEngine := gin.Default()
	AppEngine.Use(middleware)
	SetCors(AppEngine)

	// pprof
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	// 註冊api的router
	api.RegisterRouter(AppEngine)

	info.RegisterRouter(AppEngine)

}
