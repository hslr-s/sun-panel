package openness

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"

	"github.com/gin-gonic/gin"
)

type Openness struct {
}

func (p *Openness) ProIsExpired(c *gin.Context) {
	apiReturn.SuccessData(c, gin.H{
		"isExpired": biz.ProAuth.ProIsExpired(),
		"hideProBadge":biz.ProAuth.GetHideProBadgeStatus(),
	})
}
