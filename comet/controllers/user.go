package controllers

import (
	"strings"
	"webim/comet/common"
	"webim/comet/models"

	"github.com/astaxie/beego/logs"
)

type UserController struct {
	BaseController
}

func (c *UserController) Register() {

	c.Ctx.Request.ParseForm()
	if strings.ToUpper(c.Ctx.Request.Method) == "POST" {
		username := strings.Trim(c.Ctx.Request.FormValue("uname"), " ")
		password := strings.Trim(c.Ctx.Request.FormValue("passwd"), " ")
		if username == "" || password == "" {
			c.Data["msg"] = "用户名和密码不能为空"
		} else {
			res := models.InsertUser(&models.User{Username: username, Password: common.Md5(password)})
			if res.Error != nil {
				c.Data["msg"] = "注册失败请稍后重试"
				logs.Error("注册失败:", res.Error.Error())
			} else {
				c.Data["msg"] = "注册成功"
			}
		}
	}
	c.TplName = "register.html"
	c.Render()
}
