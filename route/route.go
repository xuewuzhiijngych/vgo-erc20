package route

import (
	"github.com/gin-gonic/gin"
	"vgo/controller"
	"vgo/controller/Info"
	Chain "vgo/controller/chain"
)

// CollectRoute 注册路由
func CollectRoute(app *gin.Engine) *gin.Engine {
	//AntiShake := middle.RateLimiter(1, time.Second) // 防抖

	app.GET("/ping", Info.Ping)

	app.GET("/user/info", controller.UserInfo)
	app.GET("/info", Info.Index)
	app.Any("/info/detail", Info.Detail)

	app.POST("/chain/trans/query", Chain.TransactionQuery)
	app.POST("/chain/balance/bnb", Chain.BalanceBnb)
	app.POST("/chain/balance/other", Chain.BalanceOther)
	app.POST("/chain/transfer/submit", Chain.Transfer)

	return app
}
