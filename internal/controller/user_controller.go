package controller

import (
	"awesomeGo/internal/config"
	"awesomeGo/internal/model"
	"awesomeGo/internal/service"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUser Create a new user
func CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUsers Get all users
func GetUsers(c *gin.Context) {
	users, err := service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser Get a single user by ID
func GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update a user
func UpdateUser(c *gin.Context) {
	//id := c.Param("id")
	idStr := c.Param("id")                      // 获取字符串类型的id
	id, err := strconv.ParseUint(idStr, 10, 32) // 转换为 uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var user model.User
	user, err = service.GetUserByID(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = uint(id)

	if err := service.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := service.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Download CSV
func DownloadCSV(c *gin.Context) {
	users, err := service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create CSV file
	file, err := os.Create("goUsers.csv")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	writer := csv.NewWriter(file)

	// Write CSV header
	if err := writer.Write([]string{"ID", "Name", "Email", "Salary"}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing header: " + err.Error()})
		return
	}

	// Write data
	for _, user := range users {
		record := []string{
			fmt.Sprintf("%d", user.ID),
			user.Name,
			user.Email,
			fmt.Sprintf("%d", user.Salary),
		}
		if err := writer.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing record: " + err.Error()})
			return
		}
	}
	writer.Flush()
	err2 := file.Close()
	if err2 != nil {
		return
	}
	c.File("goUsers.csv")
}

// ImportCSV UploadFive
func ImportCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	src, err := file.Open()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reader := csv.NewReader(src)
	// 读取 CSV 文件头部
	if _, err := reader.Read(); err != nil { // 假设 CSV 的第一行是头部
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var users []model.User
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading record"})
			return
		}
		// 假设 CSV 列顺序为 ID, Name, Email, Salary
		salary, _ := strconv.Atoi(record[3]) // 将字符串转为整数
		user := model.User{
			Name:   record[1],
			Email:  record[2],
			Salary: salary,
		}
		users = append(users, user)
	}
	// 将用户数据批量写入数据库
	if err := config.DB.Create(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "ImportCSV success")

}
