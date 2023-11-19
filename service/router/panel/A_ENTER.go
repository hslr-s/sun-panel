package panel

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	InitItemIcon(routerGroup)
	InitUserConfig(routerGroup)
	InitUsersRouter(routerGroup)
	InitItemIconGroup(routerGroup)
}
