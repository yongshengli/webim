package controllers

var langTypes []string // Languages that are supported.

// AppController handles the welcome screen that allows user to pick a technology and username.
type WebIMController struct {
	BaseController // Embed to use methods that are implemented in baseController.
}

//method for WebIMController.
func (c *WebIMController) Welcome() {
	c.TplName = "welcome.html"
	c.Render()
}

// Get method handles GET requests for WebImController.
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

