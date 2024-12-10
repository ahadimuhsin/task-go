package controllers

import (
	// "fmt"
	"fmt"
	"net/http"
	"os"
	"time"
	"tusk-bwa/helpers"
	"tusk-bwa/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB //konek ke db
}

func (t *TaskController) Create(ctx *gin.Context) {
	task := models.Task{}

	errBindJson := ctx.ShouldBindJSON(&task)

	if errBindJson != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errBindJson.Error()})
		return
	}

	currentTime := time.Now()

	// Format waktu ke format Y-m-d H:i:s
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	task.SubmitDate = formattedTime

	errDB := t.DB.Create(&task).Error

	if errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (t *TaskController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	task := models.Task{}

	err := t.DB.First(&task, id).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	errDB := t.DB.Where("id", id).Delete(&models.Task{}).Error

	if errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	if task.Attachment != "" {
		os.Remove("assets/" + task.Attachment)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Deleted success"})
}

func (t *TaskController) Submit(ctx *gin.Context) {

	task := models.Task{}

	id := ctx.Param("id")

	// ambil dari form data
	file, errFile := ctx.FormFile("attachment")

	if errFile != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errFile.Error()})
		return
	}

	// get data by id
	err := t.DB.First(&task, id).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// get current date time
	currentTime := time.Now()
	// Format waktu ke format Y-m-d H:i:s
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	submitDate := formattedTime

	// remove old attachment
	attachment := task.Attachment
	fileInfo, _ := os.Stat("assets/" + attachment)

	// found
	if fileInfo != nil {
		os.Remove("assets/" + attachment)
	}

	// create new attachment
	newAttachment := file.Filename
	errSave := ctx.SaveUploadedFile(file, "assets/"+newAttachment)

	if errSave != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errSave.Error()})
		return
	}

	errDB := t.DB.Where("id = ?", id).Updates(models.Task{
		Status:     "Review",
		SubmitDate: submitDate,
		Attachment: newAttachment,
	}).Error

	if errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Submit success",
	})
}

func (t *TaskController) Reject(ctx *gin.Context) {

	task := models.Task{}

	id := ctx.Param("id")

	reason := ctx.PostForm("reason")

	// get data by id
	err := t.DB.First(&task, id).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// get current date time
	currentTime := time.Now()
	// Format waktu ke format Y-m-d H:i:s
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	rejectDate := formattedTime

	errDB := t.DB.Where("id = ?", id).Updates(models.Task{
		Status:       "Rejected",
		RejectedDate: rejectDate,
		Reason:       reason,
	}).Error

	if errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Reject success",
	})
}

func (t *TaskController) Fix(ctx *gin.Context) {
	task := models.Task{}

	id := ctx.Param("id")

	// get data by id
	err := t.DB.First(&task, id).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	revision := int(task.Revision) + 1

	errUpdate := t.DB.Where("id = ?", id).Updates(models.Task{
		Status:       "Queue",
		Revision:       int8(revision), //convert ke int8
	}).Error

	if errUpdate != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errUpdate.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Fix success",
	})
}

func (t *TaskController) Approve(ctx *gin.Context) {
	task := models.Task{}

	id := ctx.Param("id")

	// get data by id
	err := t.DB.First(&task, id).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// get current date time
	currentTime := time.Now()
	// Format waktu ke format Y-m-d H:i:s
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	approveDate := formattedTime

	errUpdate := t.DB.Model(&task).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        "Approved",
		"approved_date": approveDate,
		"rejected_date": gorm.Expr("NULL"), // Gunakan gorm.Expr untuk mengatur NULL
		"reason": gorm.Expr("NULL"),
	}).Error

	if errUpdate != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errUpdate.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Approved success",
	})
}

func (t *TaskController) FindById(ctx *gin.Context){
	task := models.Task{}

	id := ctx.Param("id")

	err := t.DB.First(&models.Task{}, id).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	errDB := t.DB.Preload("User").Find(&task, id).Error

	if errDB != nil{
		ctx.JSON(http.StatusNotFound, gin.H{"error" : "Data Not Found"})
		return;
	}

	ctx.JSON(http.StatusOK, task)
}

func (t *TaskController) Review(ctx *gin.Context){
	//tipe datanya list
	tasks := []models.Task{}

	orderBy := ctx.Query("order_by")

	// Validate and sanitize the order_by parameter to avoid SQL injection
	validOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !validOrders[orderBy] {
		orderBy = "asc" // Default to "asc" if the parameter is invalid
	}

	errDB := t.DB.Preload("User").
	Where("status = ?", "Review").
	Order(fmt.Sprintf("submit_date %s", orderBy)).
	Limit(5).
	Find(&tasks).
	Error

	if errDB != nil{
		helpers.JSONResponse(ctx, http.StatusNotFound, false, "error", gin.H{"error" : "Data Not Found"})
		return;
	}

	helpers.JSONResponse(ctx, http.StatusOK, true, "success", tasks)
}

func (t *TaskController) ProgressTask(ctx *gin.Context){
	//tipe datanya list
	tasks := []models.Task{}

	userId := ctx.Param("userId")

	orderBy := ctx.Query("order_by")

	// Validate and sanitize the order_by parameter to avoid SQL injection
	validOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !validOrders[orderBy] {
		orderBy = "desc" // Default to "asc" if the parameter is invalid
	}

	errDB := t.DB.Preload("User").
	Where(
		"(status != ? AND user_id = ?) OR (revision > 0 AND user_id = ?)", "Queue", userId, userId).
	Order(fmt.Sprintf("updated_at %s", orderBy)).
	Limit(5).
	Find(&tasks).
	Error

	if errDB != nil{
		helpers.JSONResponse(ctx, http.StatusNotFound, false, "error", gin.H{"error" : "Data Not Found"})
		return;
	}

	helpers.JSONResponse(ctx, http.StatusOK, true, "success", tasks)
}

func (t *TaskController) Statistic(ctx *gin.Context){

	userId := ctx.Param("userId")

	stat := []map[string]interface{}{}

	errDB := t.DB.Model(models.Task{}).
	Select("status, COUNT(status) as total").
	Where("user_id = ?", userId).
	Group("status").
	Find(&stat).
	Error

	if errDB != nil{
		helpers.JSONResponse(ctx, http.StatusNotFound, false, "error", gin.H{"error" : errDB.Error()})
		return;
	}

	helpers.JSONResponse(ctx, http.StatusOK, true, "success", stat)
}

func (t *TaskController) FindByUserAndStatus(ctx *gin.Context){
	//tipe datanya list
	tasks := []models.Task{}

	userId := ctx.Param("userId")
	status := ctx.Param("status")

	orderBy := ctx.Query("order_by")

	// Validate and sanitize the order_by parameter to avoid SQL injection
	validOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !validOrders[orderBy] {
		orderBy = "asc" // Default to "asc" if the parameter is invalid
	}

	errDB := t.DB.Preload("User").
	Where("user_id = ? AND status = ?", userId, status).
	Order(fmt.Sprintf("submit_date %s", orderBy)).
	Find(&tasks).
	Error

	if errDB != nil{
		helpers.JSONResponse(ctx, http.StatusNotFound, false, "error", gin.H{"error" : "Data Not Found"})
		return;
	}

	helpers.JSONResponse(ctx, http.StatusOK, true, "success", tasks)
}


