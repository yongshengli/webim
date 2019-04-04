package controllers

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"strings"
	"encoding/json"
)

func init() {
	// Initialize language type list.
	langTypes = strings.Split(beego.AppConfig.String("lang_types"), "|")
	// Load locale files according to language types.
	for _, lang := range langTypes {
		beego.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			beego.Error("Fail to set message file:", err)
			return
		}
	}
}
// baseController represents base router for all other app routers.
// It implemented some methods for the same implementation;
// thus, it will be embedded into other routers.
type BaseController struct {
	beego.Controller // Embed struct that has stub implementation of the interface.
	i18n.Locale      // For i18n usage when process data and render template.
}

// Prepare implemented Prepare() method for baseController.
// It's used for language option check and setting.
func (c *BaseController) Prepare() {
	// Reset language option.
	var params = make(map[string]interface{})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &params); err != nil {
		c.error("请求body必须是json格式"+err.Error())
	}
	c.Data["params"] = params
}

type Response struct {
	Code int
	Msg string
	Data map[string]interface{}
}

func (c *BaseController) success(data map[string]interface{}) {
	c.Data["json"] = &Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
	c.ServeJSON()
}
func (c *BaseController) error(msg string) {
	c.Data["json"] = &Response{
		Code: 1,
		Msg:  msg,
		Data: nil,
	}
	c.ServeJSON()
}