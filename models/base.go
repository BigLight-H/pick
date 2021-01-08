package models

import (
	"github.com/beego/beego/client/orm"
	"github.com/beego/beego/core/config"
)

//初始化数据库
func Init() {
	dbhost, _ := config.String("db_host")
	dbport, _ := config.String("db_port")
	dbuser, _ := config.String("db_user")
	dbname, _ := config.String("db_name")
	dbpwd,  _ := config.String("db_password")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpwd + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8&loc=Asia%2FShanghai"
	_ = orm.RegisterDataBase("default", "mysql", dsn)
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
	prifix, _ := config.String("db_prifix")
	return prifix + str
}