package database

import (
	"system/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	Usuario := "root"
	Contrasenia := ""
	Nombre := "sistema"

	connection, err := gorm.Open(mysql.Open(Usuario+":"+Contrasenia+"@/"+Nombre), &gorm.Config{})

	if err != nil {
		panic("could not connection to database")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
}
