package controllers

import (
	"net/http"
	"tusk-bwa/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB //konek ke db
}

func (u *UserController) Login(ctx *gin.Context){
	//tangkap data dari flutter yg key nya harus sama dengan di model User
	user := models.User{}
	errBindJson := ctx.ShouldBindJSON(&user) 
	// cek error
	if errBindJson != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : errBindJson.Error()})
		return;
	}

	password := user.Password

	errEmail:= u.DB.Where("email = ? ", user.Email).Take(&user).Error

	if errEmail != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : errEmail.Error()})
		return;
	}
	// compare password yang dikirim dengan yg ada di db
	errHash := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if errHash != nil{
		ctx.JSON(http.StatusUnauthorized, gin.H{"error" : "Email or password is wrong"})
		return;
	}

	// ctx.JSON(http.StatusOK, map[string]interface{}{
    //     "role":       user.Role,
    //     "name":       user.Name,
    //     "email":      user.Email,
    //     "created_at": user.CreatedAt,
    //     "updated_at": user.UpdatedAt,
    // })
	ctx.JSON(http.StatusOK, user)
}

func (u *UserController) Create(ctx *gin.Context){
	//tangkap data dari flutter yg key nya harus sama dengan di model User
	user := models.User{}
	errBindJson := ctx.ShouldBindJSON(&user) 
	// cek error
	if errBindJson != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : errBindJson.Error()})
		return;
	}

	emailExist := u.DB.Where("email = ? ", user.Email).Take(&user).RowsAffected != 0

	if emailExist{
		ctx.JSON(http.StatusBadRequest, gin.H{"error" : "Email already exists"})
		return;
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	
	if errHash != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error" : errHash.Error()})
	}

	user.Password = string(hashedPassword)
	user.Role = "Employee"

	// simpan ke tabel users
	errDB := u.DB.Create(&user).Error

	if errDB != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : errDB.Error()})
		return;
	}

	ctx.JSON(http.StatusOK, user)
}

func (u *UserController) Delete(ctx *gin.Context){
	id := ctx.Param("id")
	user := models.User{}

	err := u.DB.First(&user, id).Error

	if err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{"error" : err.Error()})
		return;
	}

	errDB := u.DB.Where("id = ?", id).Delete(&models.User{}).Error

	if errDB != nil{
		ctx.JSON(http.StatusNotFound, gin.H{"error" : errDB.Error()})
		return;
	}

	ctx.JSON(http.StatusOK, gin.H{"message" : "Deleted success"})
}

func (u *UserController) GetEmployee(ctx *gin.Context){
	//tangkap data dari flutter yg key nya harus sama dengan di model User
	users := []models.User{}
	
	errDB := u.DB.Select("id", "name", "email").Where("role = ?", "Employee").Find(&users).Error

	if errDB != nil{
		ctx.JSON(http.StatusNotFound, gin.H{"error" : errDB.Error()})
		return;
	}

	ctx.JSON(http.StatusOK, users)
}