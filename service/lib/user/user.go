package file

import (
	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	c.SetCookie("cloud_tk", "", 0, "/source/", "", false, true)

}
