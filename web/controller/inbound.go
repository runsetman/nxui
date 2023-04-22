package controller

import (
	"fmt"
	"strconv"
	"x-ui/database/model"
	"x-ui/logger"
	"x-ui/web/global"
	"x-ui/web/service"
	"x-ui/web/session"

	"github.com/gin-gonic/gin"
)

type InboundController struct {
	inboundService  service.InboundService
	xrayService     service.XrayService
	settingsService service.SettingService
	fileService     service.FileService
}

func NewInboundController(g *gin.RouterGroup) *InboundController {
	a := &InboundController{}
	a.initRouter(g)
	a.startTask()
	return a
}

func (a *InboundController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/inbound")

	g.POST("/list", a.getInbounds)
	g.POST("/add", a.addInbound)
	g.POST("/write", a.writeConfig)
	g.POST("/del/:id", a.delInbound)
	g.POST("/update/:id", a.updateInbound)
}

func (a *InboundController) startTask() {
	webServer := global.GetWebServer()
	c := webServer.GetCron()
	c.AddFunc("@every 10s", func() {
		if a.xrayService.IsNeedRestartAndSetFalse() {
			err := a.xrayService.RestartXray(false)
			if err != nil {
				logger.Error("restart xray failed:", err)
			}
		}
	})
}

func (a *InboundController) getInbounds(c *gin.Context) {
	user := session.GetLoginUser(c)
	inbounds, err := a.inboundService.GetInbounds(user.Id)
	if err != nil {
		jsonMsg(c, "获取", err)
		return
	}

	// for _, inbound := range inbounds {
	// 	fmt.Println(inbound)
	// }

	jsonObj(c, inbounds, nil)
}

func (a *InboundController) writeConfig(c *gin.Context) {
	configLink := &model.LinkConfig{}
	err := c.ShouldBind(configLink)
	if err != nil {
		jsonMsg(c, "invalid body", err)
		return
	}

	hostIp, err := a.settingsService.GetHostIp()
	if err != nil {
		jsonMsg(c, "invalid host ip", err)
		return
	}

	// fmt.Println(configLink.Link)

	err = a.fileService.WriteToFile(configLink.Link, hostIp)
	if err != nil {
		jsonMsg(c, "wrtie to file", err)
		return
	}

	jsonMsg(c, "Config writed", nil)
}

func (a *InboundController) addInbound(c *gin.Context) {
	inbound := &model.Inbound{}
	err := c.ShouldBind(inbound)
	if err != nil {
		jsonMsg(c, "添加", err)
		return
	}

	// fmt.Println(inbound)
	// jsonMsg(c, "添加", errors.New("intentional"))
	// return

	user := session.GetLoginUser(c)
	inbound.UserId = user.Id
	inbound.Enable = true
	inbound.Tag = fmt.Sprintf("inbound-%v", inbound.Port)
	err = a.inboundService.AddInbound(inbound)
	jsonMsg(c, "添加", err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *InboundController) delInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "删除", err)
		return
	}
	err = a.inboundService.DelInbound(id)
	jsonMsg(c, "删除", err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *InboundController) updateInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "修改", err)
		return
	}
	inbound := &model.Inbound{
		Id: id,
	}
	err = c.ShouldBind(inbound)
	if err != nil {
		jsonMsg(c, "修改", err)
		return
	}
	err = a.inboundService.UpdateInbound(inbound)
	jsonMsg(c, "修改", err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}
