# go-contrib-mobile
[![wercker status](https://app.wercker.com/status/f51162a8f18f296de9285aa36a8e984c/s "wercker status")](https://app.wercker.com/project/bykey/f51162a8f18f296de9285aa36a8e984c)

A middleware for gin-gonic that determines whether a user is using a mobile, tablet or normal device. This middleware was inspired by the [spring-mobile](https://github.com/spring-projects/spring-mobile) project by the good folks from the spring framework.

# Examples
```go
package main

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/floresj/go-contrib-mobile"
)


func main() {
    r := gin.Default()
    fmt.Println("Main!")
  // Set up Mobile Resolver
    r.Use(mobile.Resolver())

    r.GET("/", func(c *gin.Context){
      // Get handle to Device
        d := mobile.GetDevice(c)

        switch {
        // Hey I'm a desktop!... or laptop but not a mobile or tablet!
        case d.Normal():
            c.JSON(http.StatusOK, "Hello I'm a normal device")
        // Hey I'm a mobile device!
        case d.Mobile():
            c.JSON(http.StatusOK, "Hello I'm a mobile device")
        // Woa I'm a tablet!
        case d.Tablet():
            c.JSON(http.StatusOK, "Hello I'm a tablet device")
        }
    })

    r.Run(":9000")
}
```
