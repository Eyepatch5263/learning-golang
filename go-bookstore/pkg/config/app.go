package config
import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" 
	"fmt"
)

var db *gorm.DB

func Connect() {
	dns:= "root:Costa8228@10@tcp(localhost:3306)/bookstore?charset=utf8&parseTime=True&loc=Local"
	d,err:= gorm.Open("mysql",dns)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		panic(err)
	}
	if err = d.DB().Ping(); err != nil {
        fmt.Printf("Failed to ping database: %v\n", err)
        panic(err)
    }
	db=d
}

func GetDB() *gorm.DB {
	return db
}