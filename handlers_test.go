package main

import (
	"VisitorStatistics/dao"
	"VisitorStatistics/entity"
	"encoding/json"
	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var CodeSuccess = 0
var CodeFail = 500
var MsgSuccess = "success"
var MsgError = "error"
var MsgFail = "fail"

func TestGetSummaryJSON(t *testing.T) {
	Convey("test querySummaryVisit",t,func(){
		Convey("test querySummaryVisit success",func(){
			total := 1
			entities := []entity.Data_Summary_Visit{entity.Data_Summary_Visit{
				Id:         1,
				Date_type:  "week",
				Date_value: 201951,
				Pv:         129323,
				Uv:         93234,
				Group:      "9bang",
				Union:      "",
				Channel:    "SnapTub",
			}}
			queryCountPatche := ApplyFunc(dao.SelectCountByCondition, func(_ *map[string]interface{}) int {
				return total
			})
			defer queryCountPatche.Reset()
			queryAllPatche := ApplyFunc(dao.SelectAllByCondition, func(_ *map[string]interface{},_ int,_ int) (*[]entity.Data_Summary_Visit ){
				return &entities
			})
			defer queryAllPatche.Reset()

			r := getRouter(false)
			r.POST("/data/summary/visit/get",querySummaryVisit)

			requestBody := make(map[string]interface{})
			data := make(map[string]interface{})
			requestBody["data"] = data
			data["dateType"] = "week"
			data["dateValue"] = 201951
			data["channe"] = "SnapTube"
			jsons,_ := json.Marshal(requestBody)

			req, _ := http.NewRequest("POST", "/data/summary/visit/get", strings.NewReader(string(jsons)))
			req.Header.Add("Content-Type", "application/json")

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				So(w.Code,ShouldEqual,http.StatusOK)

				p, err := ioutil.ReadAll(w.Body)
				So(err,ShouldBeNil)

				var response map[string]interface{}
				json.Unmarshal(p, &response)
				So(response["code"],ShouldEqual,CodeSuccess)
				So(response["msg"],ShouldEqual,MsgSuccess)

				dataMap:= response["data"].(map[string]interface{})
				So(dataMap["total"],ShouldEqual,total)
				So(dataMap["pageIndex"],ShouldEqual,1)
				So(dataMap["pageSize"],ShouldEqual,10)

				list:= dataMap["item"].([]interface {})
				itemMap := list[0].(map[string]interface{})
				So(itemMap["dateType"],ShouldEqual,entities[0].Date_type)
				So(itemMap["dateValue"],ShouldEqual,entities[0].Date_value)
				So(itemMap["group"],ShouldEqual,entities[0].Group)
				So(itemMap["union"],ShouldEqual,entities[0].Union)
				So(itemMap["channel"],ShouldEqual,entities[0].Channel)
				So(itemMap["PV"],ShouldEqual,entities[0].Pv)
				So(itemMap["UV"],ShouldEqual,entities[0].Uv)

				return true
			})
		})

		Convey("test querySummaryVisit fail",func(){
			r := getRouter(false)
			r.POST("/data/summary/visit/get",querySummaryVisit)

			requestBody := make(map[string]interface{})
			requestBody["data"] = -1
			jsons,_ := json.Marshal(requestBody)

			req, _ := http.NewRequest("POST", "/data/summary/visit/get", strings.NewReader(string(jsons)))
			req.Header.Add("Content-Type", "application/json")

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				So(w.Code,ShouldEqual,http.StatusOK)

				p, err := ioutil.ReadAll(w.Body)
				So(err,ShouldBeNil)

				var response map[string]interface{}
				json.Unmarshal(p, &response)
				So(response["code"],ShouldEqual,CodeFail)
				So(response["msg"],ShouldEqual,MsgError)
				So(response["data"],ShouldBeNil)

				return true
			})
		})
	})
}

func TestInsertSummaryVisit(t *testing.T){
	Convey("test insertSummaryVisit",t,func(){
		Convey("test insert success",func(){
			patches := ApplyFunc(dao.InsertAll, func(_ *[]entity.Data_Summary_Visit) bool {
				return true
			})
			defer patches.Reset()

			r := getRouter(false)
			r.POST("/data/summary/visit/set",insertSummaryVisit)

			requestBodyString := `[
								{
									"source": "testtoby.myshopify.coms",
									"source_type": "shopifys12",
									"type": "daily_report",
									"created_at": 1578233369267,
									"raw_data": {
										"date": "2020-02-10T12:59:04.709Z",
										"sessions": 22,
										"visitors": 459
									}
								}
							]`
			req,_ :=http.NewRequest("POST","/data/summary/visit/set",strings.NewReader(requestBodyString))
			req.Header.Add("Content-Type", "application/json")

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				So(w.Code,ShouldEqual,http.StatusOK)

				p, err := ioutil.ReadAll(w.Body)
				So(err,ShouldBeNil)

				var response map[string]interface{}
				json.Unmarshal(p, &response)
				So(response["code"],ShouldEqual,CodeSuccess)
				So(response["msg"],ShouldEqual,MsgSuccess)
				So(response["data"],ShouldBeNil)

				return true
			})
		})

		Convey("test insert fail",func(){
			patches := ApplyFunc(dao.InsertAll, func(_ *[]entity.Data_Summary_Visit) bool {
				return false
			})
			defer patches.Reset()

			r := getRouter(false)
			r.POST("/data/summary/visit/set",insertSummaryVisit)

			requestBodyString := `[
								{
									"source": "testtoby.myshopify.coms",
									"source_type": "shopifys12",
									"type": "daily_report",
									"created_at": 1578233369267,
									"raw_data": {
										"date": "2020-02-10T12:59:04.709Z",
										"sessions": 22,
										"visitors": 459
									}
								}
							]`
			req,_ :=http.NewRequest("POST","/data/summary/visit/set",strings.NewReader(requestBodyString))
			req.Header.Add("Content-Type", "application/json")

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				So(w.Code,ShouldEqual,http.StatusOK)

				p, err := ioutil.ReadAll(w.Body)
				So(err,ShouldBeNil)

				var response map[string]interface{}
				json.Unmarshal(p, &response)
				So(response["code"],ShouldEqual,CodeSuccess)
				So(response["msg"],ShouldEqual,MsgFail)
				So(response["data"],ShouldBeNil)

				return true
			})
		})
	})
}
