package main

import (
	"log"
	"myapi/db"
	"myapi/models"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	db.Init()

	seedUsers()
	seedCategories()
	seedProducts()
	seedOrder()

	log.Println("Seed data created successfully")
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to generate password hash: %v", err)
	}
	return string(hash)
}

func seedUsers() {
	if db.DB.Find(&models.User{}).RowsAffected > 0 {
		log.Println("Users already exist, skipping")
		return
	}

	users := []models.User{
		{Username: "admin", Password: hashPassword("admin123"), Email: "admin@example.com"},
		{Username: "user1", Password: hashPassword("password1"), Email: "user1@example.com"},
	}

	for _, user := range users {
		if err := db.DB.Create(&user).Error; err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
	}
	log.Println("Users seeded")
}

func seedCategories() {
	if db.DB.Find(&models.Category{}).RowsAffected > 0 {
		log.Println("Categories already exist, skipping")
		return
	}

	categories := []models.Category{
		{Name: "Elektronik"},
		{Name: "Pakaian"},
		{Name: "Makanan"},
		{Name: "Gaming"},
	}

	for _, category := range categories {
		if err := db.DB.Create(&category).Error; err != nil {
			log.Fatalf("Failed to create category: %v", err)
		}
	}
	log.Println("Categories seeded")
}

func seedProducts() {
	if db.DB.Find(&models.Product{}).RowsAffected > 0 {
		log.Println("Products already exist, skipping")
		return
	}

	var categories []models.Category
	db.DB.Find(&categories)
	catMap := make(map[string]models.Category)
	for _, c := range categories {
		catMap[c.Name] = c
	}

	products := []struct {
		Name       string
		Qty        int
		Price      float64
		Desc       string
		Categories []string
	}{
		{"Laptop Gaming", 10, 15000000, "Laptop gaming 16GB RAM RTX4060", []string{"Elektronik", "Gaming"}},
		{"Kemeja Flanel", 50, 150000, "Kemeja flanel kotak-kotak", []string{"Pakaian"}},
		{"Snack Keripik", 100, 10000, "Keripik kentang pedas", []string{"Makanan"}},
		{"Mouse Wireless", 30, 250000, "Mouse gaming wireless ringan", []string{"Elektronik", "Gaming"}},
	}

	for _, p := range products {
		var cats []models.Category
		for _, name := range p.Categories {
			if c, ok := catMap[name]; ok {
				cats = append(cats, c)
			}
		}

		product := models.Product{
			Name:        p.Name,
			Quantity:    p.Qty,
			Price:       p.Price,
			Description: p.Desc,
			Categories:  cats,
		}
		if err := db.DB.Create(&product).Error; err != nil {
			log.Fatalf("Failed to create product: %v", err)
		}
	}
	log.Println("Products seeded")
}

func seedOrder() {
	if db.DB.Find(&models.Order{}).RowsAffected > 0 {
		log.Println("Orders already exist, skipping")
		return
	}

	var user models.User
	if err := db.DB.First(&user).Error; err != nil {
		log.Println("No user found, skipping order seed")
		return
	}

	var products []models.Product
	db.DB.Find(&products)
	if len(products) < 2 {
		log.Println("Not enough products, skipping order seed")
		return
	}

	order := models.Order{
		UserID:     user.ID,
		TotalPrice: products[0].Price*2 + products[1].Price*1,
		Status:     "pending",
		Items: []models.OrderItem{
			{ProductID: products[0].ID, Quantity: 2, Price: products[0].Price},
			{ProductID: products[1].ID, Quantity: 1, Price: products[1].Price},
		},
	}

	if err := db.DB.Create(&order).Error; err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}
	log.Println("Orders seeded")
}
