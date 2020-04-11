package main

import (
	"VisitorStatistics/dao"
	"VisitorStatistics/entity"
	"strconv"
	"strings"
	"webengine"
)

func showIndexPage(c *webengine.Context){
	render(c,webengine.H{
		"payload":buildResponseMap(0,"success",nil),
	},"")
}

func querySummaryVisit(c *webengine.Context){
	object := c.GetJsonObject()
	if data,ok := (*object)["data"].(map[string]interface{});ok{
		setDefaultValueForQuary(&data)
		total := dao.SelectCountByCondition(&data)
		start,size := getPageInfo(&data,total)
		result := dao.SelectAllByCondition(&data,start,size)
		items := convertToResponseItem(result)
		pageData := buildPageData(items,total,start,size)
		render(c,webengine.H{
			"payload": buildResponseMap(0,"success",pageData),
		},"")
	}else{
		render(c,webengine.H{
			"payload":buildResponseMap(500,"error",nil),
		},"")
	}
}

func convertToResponseItem(result *[]entity.Data_Summary_Visit) *[]map[string]interface{} {
	responseItems := make( []map[string]interface{},0)
	for _,r := range *result{
		item := make(map[string]interface{})
		item["dateType"] = r.Date_type
		item["dateValue"] = r.Date_value
		item["group"] = r.Group
		item["union"] = r.Union
		item["channel"] = r.Channel
		item["PV"] = r.Pv
		item["UV"] = r.Uv
		responseItems = append(responseItems,item)
	}
	return &responseItems
}

func buildPageData(result *[]map[string]interface{}, total int, start int, size int) *map[string]interface{} {
	pageData := make(map[string]interface{})
	pageData["total"] = total
	pageData["pageIndex"] = (start / size) + 1
	pageData["pageSize"] = size
	pageData["item"] = result
	return &pageData
}

func getPageInfo(data *map[string]interface{}, total int) (int,int){
	size := obj2int((*data)["pageSize"],10)
	start := (obj2int((*data)["pageIndex"],1) - 1) * size
	if start >= total{
		start = (total / size) * size
	}
	return start,size
}

func insertSummaryVisit(c *webengine.Context){
	array := c.GetJsonArray()
	summaryVisitList := convert2SummayVisit(array)
	msg := "success"
	if len(*summaryVisitList) != 0{
		isSuccess := dao.InsertAll(summaryVisitList)
		if !isSuccess{
			msg = "fail"
		}
	}
	render(c,webengine.H{
		"payload":buildResponseMap(0,msg,nil),
	},"")
}

func convert2SummayVisit(array *[]map[string]interface{}) *[]entity.Data_Summary_Visit {
	var inserts []entity.Data_Summary_Visit

	for _, item := range (*array){
		if rowMap,ok := item["raw_data"].(map[string]interface{});ok{
			var dateType int = 0
			if dateString,ok := rowMap["date"].(string);ok{
				if date,err:=strconv.Atoi(strings.ReplaceAll(dateString[:10],"-",""));err == nil{
					dateType = date
				}
			}

			visit := entity.Data_Summary_Visit{
				Date_type:  "day",
				Date_value: dateType,
				Pv:         obj2int(rowMap["visitors"], 0),
				Uv:         0,
			}
			inserts = append(inserts,visit)
		}
	}

	return &inserts
}

func setDefaultValueForQuary(data *map[string]interface{}) {
	if (*data)["dateType"] == nil{
		(*data)["dateType"] = "day"
	}
	if (*data)["pageSize"] == nil{
		(*data)["pageSize"] = 10
	}
	if (*data)["pageIndex"] == nil{
		(*data)["pageIndex"] = 1
	}
}

func buildResponseMap(code int,msg string,data interface{}) interface{}{
	response := make(map[string]interface{})
	response["code"] = code
	response["msg"] = msg
	response["data"] = data
	return response
}
