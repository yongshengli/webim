package controllers

import "comet/common"

var langTypes []string // Languages that are supported.

// AppController handles the welcome screen that allows user to pick a technology and username.
type WebIMController struct {
	BaseController // Embed to use methods that are implemented in baseController.
}

//method for WebIMController.
func (c *WebIMController) Welcome() {
	c.Redirect("/webim", 302)
}

// Get method handles GET requests for WebImController.
func (c *WebIMController) Get() {
	// Safe check.
	c.TplName = "websocket.html"
	c.Data["IsWebSocket"] = true
	c.Data["DeviceId"] = common.Uuid()
	c.Render()
}
