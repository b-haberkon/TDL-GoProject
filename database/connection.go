package database

import (
	"system/models"

	//"gorm.io/driver/mysql"
	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func conn_to_mysql() (db *gorm.DB, err error) {
	Usuario := "root"
	Contrasenia := ""
	Nombre := "sistema"
	return gorm.Open(mysql.Open(Usuario+":"+Contrasenia+"@/"+Nombre), &gorm.Config{})
}

func conn_to_sqlite() (db *gorm.DB, err error) {
	return gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
}

func Connect() {
	connection, err := conn_to_sqlite()

	if err != nil {
		panic("could not connection to database")
	}

	DB = connection
	connection.AutoMigrate(&models.User{})
}
