package main

import (
	"fmt"
	"time"
	"webengine"
)

func logBeforeHandler() webengine.HandlerFunc{
	return func(c *webengine.Context){
		fmt.Printf("request come!method:%s\tpath:%s\ttime:%s",c.Request.Method,c.Request.URL.Path,time.Now())
	}
}
