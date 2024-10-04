package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// User struct with GORM model
type User struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Email  string `json:"email" gorm:"unique"`
	Salary int    `json:"salary"`
}

var db *gorm.DB

// Initialize the database connection
func initDB() {
	var err error
	url := string("postgres://username:kris@localhost:5432/kris")
	dsn := os.Getenv(url) // e.g. "postgres://username:password@localhost:5432/dbname"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Unable to Auto migrate the schema : ", err)
		return
	}
}

// Create a new user
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Get all users
func getUsers(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Get a single user by ID
func getUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update a user
func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete a user
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// download csv
func downloadCSV(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//创建csv文件
	file, err := os.Create("goUsers.csv")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//defer file.Close()

	writer := csv.NewWriter(file)
	//defer writer.Flush()

	//写入csv 头部文件
	//writer.Write([]string{"ID", "Name", "Email", "Salary"})
	if err := writer.Write([]string{"ID", "Name", "Email", "Salary"}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing header: " + err.Error()})
		return
	}

	//写入数据
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

	// Flush 数据到文件
	writer.Flush()
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error flushing writer: " + err.Error()})
		return
	}

	// 关闭文件以确保所有数据都写入文件
	err = file.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file.Close(): " + err.Error()})
		return
	}

	c.File("goUsers.csv")
}

func main() {
	// Initialize DB connection
	initDB()

	// Create a Gin router
	r := gin.Default()

	// Define API routes
	r.POST("/users/createUser", createUser)
	r.GET("/users/getUsers", getUsers)
	r.GET("/users/getUser/:id", getUser)
	r.PUT("/users/updateUser/:id", updateUser)
	r.DELETE("/users/deleteUser/:id", deleteUser)

	r.GET("/users/downloadCSV", downloadCSV)

	// Start the server
	log.Println("Server running on port 8080")
	r.Run(":8080")
}
