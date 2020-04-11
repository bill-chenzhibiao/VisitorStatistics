package main

func initializeRouters(){
	engine.Use(logBeforeHandler())

	engine.GET("/",showIndexPage)

	engine.POST("/data/summary/visit/get",querySummaryVisit)
	engine.POST("/data/summary/visit/set",insertSummaryVisit)
}
