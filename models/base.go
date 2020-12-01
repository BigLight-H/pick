package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//初始化数据库
func Init() {
	dbhost := beego.AppConfig.String("db_host")
	dbport := beego.AppConfig.String("db_port")
	dbuser := beego.AppConfig.String("db_user")
	dbname := beego.AppConfig.String("db_name")
	dbpwd  := beego.AppConfig.String("db_password")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpwd + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8&loc=Asia%2FShanghai"
	orm.RegisterDataBase("default", "mysql", dsn)
	orm.RegisterModel(
		new(Book),
		new(Chapter),
		new(Photo),
		new(Links),
		new(BookList),
		new(BookEpisode),
	)
}

//返回带前缀的表名
func TableName(str string) string {
	return beego.AppConfig.String("db_prifix") + str
}