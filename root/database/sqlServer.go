package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

type SQLServer struct {
	UserName     string
	Password     string
	Host         string
	DatabaseName string
}

func (m *SQLServer) Connection() *gorm.DB {
	db, err := gorm.Open("mssql", "sqlserver://"+m.UserName+":"+m.Password+"@"+m.Host+"?database="+m.DatabaseName)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func (m *SQLServer) Parse() DatabaseConnector {
	//m.UserName = ethereal.GetCnf("DATABASE.SQLSERVER.USER").(string)
	//m.Host = ethereal.GetCnf("DATABASE.SQLSERVER.HOST").(string)
	//m.Password = ethereal.GetCnf("DATABASE.SQLSERVER.PASSWORD").(string)
	//m.DatabaseName = ethereal.GetCnf("DATABASE.SQLSERVER.NAME").(string)
	return m
}
