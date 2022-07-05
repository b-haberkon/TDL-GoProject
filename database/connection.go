package database

import (
	"fmt"
	"bufio"
	"os"
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
	defer database.DBConn.Close()

	if err != nil {
		panic("could not connection to database")
	}

	DB = connection
	connection.AutoMigrate(&models.User{}, &models.Data{})

	CargarDatos()
}

func CargarDatos() {
	archivo, error := os.Open("./../memotest_data.txt")
	defer archivo.Close()
	if error != nil {
		panic("error opening file")
	}

	scanner := bufio.NewScanner(archivo)

	for scanner.Scan(){
		linea := scanner.Text()
		info := strings.Split(linea, ",")

		var data Data 
		data.Name = info[0]
		data.Imagen = info[1]
		data.Japones = info[2]

		DB.Create(&data)
	}
}