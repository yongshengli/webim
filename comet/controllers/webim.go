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

// AppController handles the welcome screen that allows user to pick a technology and username.
type WebIMController struct {
	BaseController // Embed to use methods that are implemented in baseController.
}

// Get implemented Get() method for AppController.
func (c *WebIMController) Welcome() {
	c.TplName = "welcome.html"
	c.Render()
}

// Join method handles POST requests for AppController.
func (c *WebIMController) Join() {
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
// Get method handles GET requests for WebSocketController.
func (c *WebIMController) Get() {
	// Safe check.
	uname := c.GetString("uname")
	if len(uname) == 0 {
		c.Redirect("/", 302)
		return
	}

	c.TplName = "websocket.html"
	c.Data["IsWebSocket"] = true
	c.Data["UserName"] = uname
	c.Render()
}

