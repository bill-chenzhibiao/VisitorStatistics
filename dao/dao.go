package dao

import (
	"VisitorStatistics/entity"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

//数据库配置
const (
	userName = "root"
	password = "root"
	ip = "127.0.0.1"
	port = "3306"
	dbName = "ecp_test"
)

var DB *sql.DB

func InitDB()  {
	path := strings.Join([]string{userName, ":", password, "@tcp(",ip, ":", port, ")/", dbName, "?charset=utf8"}, "")

	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil{
		fmt.Println("opon database fail")
		return
	}
	fmt.Println("connect database success")
}


func SelectById(id int) (*entity.Data_Summary_Visit) {
	var visit entity.Data_Summary_Visit
	err := DB.QueryRow("SELECT * FROM data_summary_visit WHERE id = ?", id).Scan(&visit.Id,&visit.Date_type,&visit.Date_value,&visit.Pv,&visit.Uv,&visit.Group,&visit.Union,&visit.Channel)
	if err != nil{
		fmt.Println("查询出错了")
	}
	return &visit
}

func SelectCountByCondition(data *map[string]interface{}) int{
	total := 0
	queryCountSql := buildQueryCountSql(data)

	err := DB.QueryRow(queryCountSql,(*data)["dateType"]).Scan(&total)
	if err != nil{
		fmt.Println("查询出错了")
		return 0
	}
	return total
}



func SelectAllByCondition(data *map[string]interface{},start int,size int) (*[]entity.Data_Summary_Visit) {
	querySql := buildQuerySql(data)

	rows, err := DB.Query(querySql,(*data)["dateType"],start,size)
	if err != nil{
		fmt.Println("查询出错了")
		return nil
	}
	var visits []entity.Data_Summary_Visit
	//loop
	for rows.Next(){
		var visit entity.Data_Summary_Visit
		err := rows.Scan(&visit.Id,&visit.Date_type,&visit.Date_value,&visit.Pv,&visit.Uv,&visit.Group,&visit.Union,&visit.Channel)
		if err != nil {
			fmt.Println("rows fail")
		}
		visits = append(visits, visit)
	}
	return &visits
}

func InsertOne(entity *entity.Data_Summary_Visit) (bool){
	tx, err := DB.Begin()
	if err != nil{
		fmt.Println("tx fail")
		return false
	}

	stmt, err := tx.Prepare("INSERT INTO data_summary_visit (`date_type`,`date_value`,`pv`,`uv`,`group`,`union`,`channel`) VALUES (?, ?,?, ?,?, ?,?)")
	if err != nil{
		fmt.Println("Prepare fail")
		return false
	}

	_, err = stmt.Exec(entity.Date_type,entity.Date_value,entity.Pv,entity.Uv,entity.Group,entity.Union,entity.Channel)
	if err != nil{
		fmt.Println("Exec fail")
		return false
	}
	tx.Commit()
	return true
}

func InsertAll(data *[]entity.Data_Summary_Visit) (bool){
	//开启事务
	tx, err := DB.Begin()
	if err != nil{
		fmt.Println("tx fail")
		return false
	}
	sql := "INSERT INTO data_summary_visit (`date_type`,`date_value`,`pv`,`uv`,`group`,`union`,`channel`) VALUES "
	// 循环data数组,组合sql语句
	for key, entity := range (*data) {
		if len(*data) - 1 == key {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("('%s',%d,%d,%d,'%s','%s','%s');", entity.Date_type,entity.Date_value,entity.Pv,entity.Uv,entity.Group,entity.Union,entity.Channel)
		} else {
			sql += fmt.Sprintf("('%s',%d,%d,%d,'%s','%s','%s'),", entity.Date_type,entity.Date_value,entity.Pv,entity.Uv,entity.Group,entity.Union,entity.Channel)
		}
	}
	_, err = DB.Exec(sql)
	if err != nil{
		fmt.Println("Exec fail")
		return false
	}
	//将事务提交
	tx.Commit()
	return true
}


func buildQueryCountSql(data *map[string]interface{}) string {
	sql := "SELECT count(*) from data_summary_visit where date_type = ?"
	sql = addOptionalCondition(data, sql)
	return sql
}

func buildQuerySql(data *map[string]interface{}) string{
	sql := "SELECT * from data_summary_visit where date_type = ?"
	sql = addOptionalCondition(data, sql)
	sql += " order by date_value desc limit ?,?"
	return sql
}

func addOptionalCondition(data *map[string]interface{}, sql string) string {
	if group, ok := (*data)["group"].(string); ok && group != "" {
		sql += " and group = " + "'" + group + "'"
	}
	if union, ok := (*data)["union"].(string); ok && union != "" {
		sql += " and union = " + "'" + union + "'"
	}
	if channel, ok := (*data)["channel"].(string); ok && channel != "" {
		sql += " and channel = " + "'" + channel + "'"
	}
	return sql
}