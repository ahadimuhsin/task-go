package config

import (
	"fmt"
	"tusk-bwa/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = ""
	dbName   = "tusk"
)

func DatabaseConnection() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}

func CreateOwnerAccount(db *gorm.DB){
	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password"), bcrypt.DefaultCost)

	//variabel owner
	owner := models.User{
		Role: "Admin",
		Name: "Owner",
		Password: string(hashedPassword),
		Email: "admin@mail.com",
	}

	if db.Where("email = ?", owner.Email).First(&owner).RowsAffected == 0 {
		db.Create(&owner)
	}else{
		fmt.Println("Owner exist")
	}
}
