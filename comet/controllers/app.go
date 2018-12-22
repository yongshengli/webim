package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

var langTypes []string // Languages that are supported.

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
	c.Lang = "" // This field is from i18n.Locale.

	// 1. Get language information from 'Accept-Language'.
	al := c.Ctx.Request.Header.Get("Accept-Language")
	if len(al) > 4 {
		al = al[:5] // Only compare first 5 letters.
		if i18n.IsExist(al) {
			c.Lang = al
		}
	}

	// 2. Default language is English.
	if len(c.Lang) == 0 {
		c.Lang = "en-US"
	}

	// Set template level language option.
	c.Data["Lang"] = c.Lang
}

// AppController handles the welcome screen that allows user to pick a technology and username.
type AppController struct {
	BaseController // Embed to use methods that are implemented in baseController.
}

// Get implemented Get() method for AppController.
func (c *AppController) Get() {
	c.TplName = "welcome.html"
	c.Render()
}

// Join method handles POST requests for AppController.
func (c *AppController) Join() {
	// Get form value.
	uname := c.GetString("uname")
	tech := c.GetString("tech")

	// Check valid.
	if len(uname) == 0 {
		c.Redirect("/", 302)
		return
	}

	switch tech {
	case "longpolling":
		c.Redirect("/lp?uname="+uname, 302)
	case "websocket":
		c.Redirect("/ws?uname="+uname, 302)
	default:
		c.Redirect("/", 302)
	}

	// Usually put return after redirect.
	return
}
