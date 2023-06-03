package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Customer struct {
	ID        int    `json:"id"`
	Name      string `json:"name" binding:"required"`
	Role      string `json:"role" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Contacted bool   `json:"contacted"`
}

var db = make(map[string]interface{})
var customersSlice []Customer

func findCustomerByID(id int) (*Customer, error) {
	for _, customer := range customersSlice {
		if customer.ID == id {
			return &customer, nil
		}
	}

	return nil, fmt.Errorf("Customer with ID %d not found", id)
}

func index(c *gin.Context) {
	c.File("./static/index.html")
}

func postman(c *gin.Context) {
	c.File("./crm-backend.postman_collection.json")
}

func customerIndex(c *gin.Context) {
	customers := db["customers"].([]Customer)

	c.JSON(http.StatusOK, gin.H{
		"data": customers,
	})
}

func customerShow(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID!",
		})

		return
	}

	customer, err := findCustomerByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Customer not found!",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"customer": customer,
		},
	})
}

func customerStore(c *gin.Context) {
	var customer Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request!",
			"err":     err.Error(),
		})

		return
	}

	customers := db["customers"].([]Customer)
	customer.ID = len(customers) + 1
	customers = append(customers, customer)
	db["customers"] = customers

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"customer": customer,
		},
	})
}

func customerUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID!",
		})
		return
	}

	customer, err := findCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Customer not found!",
		})
		return
	}

	// Bind the new data to the existing customer object
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request!",
			"err":     err.Error(),
		})
		return
	}

	// Update the customer in the customersSlice
	for i, cust := range customersSlice {
		if cust.ID == customer.ID {
			customersSlice[i] = *customer
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": customer,
	})
}

func customerDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID!",
		})
	}

	customer, err := findCustomerByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Customer not found!",
		})

		return
	}

	customers := db["customers"].([]Customer)
	for i, c := range customers {
		if c.ID == customer.ID {
			customers = append(customers[:i], customers[i+1:]...)
			break
		}
	}

	db["customers"] = customers

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer deleted successfully!",
	})
}

func initializeCustomers() {
	customersSlice = []Customer{
		{
			ID:        1,
			Name:      "John Doe",
			Role:      "CEO",
			Email:     "john.doe@example.com",
			Phone:     "1234567890",
			Contacted: false,
		},
		{
			ID:        2,
			Name:      "Jane Doe",
			Role:      "CTO",
			Email:     "jane.doe@example.com",
			Phone:     "1234567890",
			Contacted: false,
		},
		{
			ID:        3,
			Name:      "John Smith",
			Role:      "CFO",
			Email:     "jon.smith@example.com",
			Phone:     "1234567890",
			Contacted: false,
		},
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/", index)
	r.GET("/postman", postman)
	r.GET("/customers", customerIndex)
	r.GET("/customers/:id", customerShow)
	r.POST("/customers", customerStore)
	r.PUT("/customers/:id", customerUpdate)
	r.DELETE("/customers/:id", customerDelete)

	return r
}

func main() {
	initializeCustomers()
	db["customers"] = customersSlice

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := setupRouter()
	r.Run(":8080")
}
