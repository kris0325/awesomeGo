package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

/*
*

	gorm实现join 查詢
*/
// Owner struct with GORM model
type Owner struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Pets      []Pet  `gorm:"foreignKey:OwnerID"`
}

type Pet struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	OwnerID uint   `json:"owner_id"`
}

// Pagination struct to handle pagination parameters
type Pagination struct {
	Page      int    `json:"page" binding:"required"`
	Size      int    `json:"size" binding:"required"`
	OwnerName string `json:"ownerName"` // Add ownerName to the struct
	PetName   string `json:"petName"`   // Add ownerName to the struct
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
	/**
	在您提供的Golang代码中，`InitDB`函数的主要作用是初始化与PostgreSQL数据库的连接，并自动迁移数据库模式。下面是对这段代码中相关部分的详细解释：

	### 代码解析

	### 1. 数据库连接初始化

	- **`gorm.Open`**：使用GORM库的`Open`方法打开与指定数据库的连接。这里使用的是PostgreSQL驱动。
	- **`dsn`（Data Source Name）**：这是一个连接字符串，包含了数据库的地址、用户名、密码和数据库名称。在您的代码中，连接字符串是硬编码的，但通常可以通过环境变量来配置，以提高安全性和灵活性。
	- **错误处理**：如果连接失败，程序会记录错误并终止执行。

	### 2. 自动迁移（Auto Migrate）

	- **`DB.AutoMigrate(&model.User{})`**：这个调用会根据`User`模型的结构自动创建或更新数据库中的表结构。
	  - 如果表不存在，`AutoMigrate`会创建它。
	  - 如果表已经存在，但模型结构发生了变化（例如添加了新字段），它会尝试更新表结构以匹配模型。
	- **错误处理**：同样，如果自动迁移过程中出现错误，程序会记录错误并终止执行。

	### 作用总结

	- **初始化数据库连接**：确保应用程序能够与PostgreSQL数据库进行交互。
	- **自动迁移模式**：简化了数据库管理，确保数据库表结构与应用程序中的模型保持同步，从而减少手动管理数据库模式的工作量。

	### 注意事项

	- 在生产环境中，直接使用`AutoMigrate`可能不是最佳实践，因为它可能导致数据丢失或不一致。通常建议在开发阶段使用，而在生产环境中使用更严格的迁移策略（例如手动编写SQL迁移脚本）。
	- 确保在使用环境变量时正确设置，以避免硬编码敏感信息（如用户名和密码）。

	*/
}

// Search for owners with pagination and filtering based on GORM
func searchOwners(c *gin.Context) {
	var pagination Pagination
	if err := c.ShouldBindJSON(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get query parameters and convert them to lowercase for case-insensitive comparison
	// Convert ownerName to lowercase for case-insensitive comparison
	ownerName := strings.ToLower(pagination.OwnerName)
	petName := strings.ToLower(pagination.PetName)

	var owners []Owner
	var total int64

	query := db.Model(&Owner{}).
		Preload("Pets").
		Joins("LEFT JOIN pets ON pets.owner_id = owners.id")

	if ownerName != "" {
		query = query.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?", "%"+ownerName+"%", "%"+ownerName+"%")
	}
	if petName != "" {
		query = query.Where("LOWER(pets.name) LIKE ?", "%"+petName+"%")
	}

	query.Count(&total)

	query = query.Order("last_name ASC, first_name DESC").
		Offset((pagination.Page - 1) * pagination.Size).
		Limit(pagination.Size)

	if err := query.Find(&owners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"page":  pagination.Page,
		"size":  pagination.Size,
		"rows":  owners,
	})
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

	//searchOwners with pagination
	r.POST("/owners/search", searchOwners)

	// Start the server
	log.Println("Server running on port 8090")
	r.Run(":8090")
}
