package main

import (
	"VisitorStatistics/dao"
	"net/http"
	we "webengine"
)

var engine *we.Engine

func main(){
	dao.InitDB()

	engine = we.Default()

	initializeRouters()

	engine.Run(":8090")

}

func render(c *we.Context, data we.H, templateName string) {

	switch c.Request.Header.Get("Content-Type") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	default:
		c.JSON(http.StatusOK, data["payload"])
	}
}