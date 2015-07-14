# go-contrib-mobile
A middleware for gin-gonic that determines whether a user is using a mobile, tablet or normal device.

# Examples
```go
package hello

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/floresj/go-contrib-mobile"
)


func init() {
	r := gin.Default()
	
  // Set up Mobile Resolver
	r.Use(mobile.Resolver())

	r.GET("/", func(c *gin.Context){
	  // Get handle to Device
		d := mobile.GetDevice(c)
		
		switch {
		// Hey I'm a desktop!... or laptop but not a mobile or tablet!
		case d.Normal():
			c.HTML(http.StatusOK, "index.tmpl", nil)
		// Hey I'm a mobile device!
		case d.Mobile():
			c.HTML(http.StatusOK, "index.mobile.tmpl", nil)
		// Woa I'm a tablet!
		case d.Tablet():
			c.HTML(http.StatusOK, "index.tablet.tmpl", nil)
		}
	})
```