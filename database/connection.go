package database

import (
	//"fmt"
	"bufio"
	"os"
	"system/memotest"
	"system/models"
	//"gorm.io/driver/mysql"
	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"sync"
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
	//defer database.DBConn.Close()

	if err != nil {
		panic("could not connection to database")
	}

	DB = connection
	connection.AutoMigrate(&models.User{}, &models.Data{})
}

var dataLoaded = new(sync.Once)
func GetFullSymbolSet() ([][]*memotest.Symbol, error) {
	ret, err := getFullSymbolSet()
	if (err!=nil) {
		return ret, err
	}
	if( len(ret) > 0) {
		return ret, err
	}
	// Si al primer intento no hay nada, intentará llenar
	// la base de datos por única vez
	dataLoaded.Do(LoadData)
	return getFullSymbolSet()
}
func getFullSymbolSet() ([][]*memotest.Symbol, error) {
	var pairs []models.Data
	var fullSet [][]*memotest.Symbol
	result := DB.Find(&pairs)
	if(result.Error != nil) {
		
		return nil, result.Error
	}
	for _, pair := range pairs {
		imagen  := memotest.Symbol{Text: pair.Imagen,  Pair: pair.Id}
		japones := memotest.Symbol{Text: pair.Japones, Pair: pair.Id}
		pair := []*memotest.Symbol { & imagen, & japones }
		fullSet = append(fullSet, pair)
	}
	return fullSet, nil
}
func LoadData() {
	archivo, error := os.Open("./database/memotest_data.txt")
	defer archivo.Close()
	if error != nil {
		panic("error opening file")
	}

	scanner := bufio.NewScanner(archivo)

	for scanner.Scan() {
		linea := scanner.Text()
		info := strings.Split(linea, ",")

		var data models.Data 
		data.Name = info[0]
		data.Imagen = info[1]
		data.Japones = info[2]

		DB.Create(&data)
	}
}