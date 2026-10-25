package main

import (
	"errors"
	"fmt"
	"log"
	"myapi/db"
	"myapi/models"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type categorySeed struct {
	ID     uint
	Name   string
	Slug   string
	Status string
}

type productSeed struct {
	ID                uint
	Name              string
	Slug              string
	Quantity          int
	Price             float64
	Description       string
	Sku               string
	PublicationStatus string
	Image             string
	CategoryID        uint
}

var categories = []categorySeed{
	{ID: 22, Name: "Pastries", Slug: "pastries", Status: "active"},
	{ID: 23, Name: "Breads", Slug: "breads", Status: "active"},
	{ID: 24, Name: "Cakes", Slug: "cakes", Status: "active"},
	{ID: 25, Name: "Cookies", Slug: "cookies", Status: "active"},
	{ID: 26, Name: "Beverages", Slug: "beverages", Status: "active"},
}

var products = []productSeed{
	{
		ID: 35, Name: "Butter Croissant", Slug: "butter-croissant", Quantity: 24,
		Price: 35000, Description: "Flaky buttery croissant baked fresh every morning",
		Sku: "PAS-CRO-001", PublicationStatus: "active",
		Image: "https://gpggjd7cx2cizeyj.public.blob.vercel-storage.com/products/f50c96ff9bd73f1f-RwHKiS5q2PiOaAR3o6isqifPCC1tfS.jpg", CategoryID: 22,
	},
	{
		ID: 36, Name: "Almond Croissant", Slug: "almond-croissant", Quantity: 18,
		Price: 42000, Description: "Buttery croissant filled with almond cream and toasted almonds",
		Sku: "PAS-ALM-001", PublicationStatus: "active",
		Image: "https://gpggjd7cx2cizeyj.public.blob.vercel-storage.com/products/771ec0ac65033f38-iDOgf8KwX7k4038cGjPMXpRFRAHbwg.jpg", CategoryID: 22,
	},
	{
		ID: 37, Name: "Cinnamon Morning Roll", Slug: "cinnamon-morning-roll", Quantity: 20,
		Price: 28000, Description: "Soft cinnamon roll finished with a light vanilla glaze",
		Sku: "PAS-CIN-002", PublicationStatus: "active",
		Image: "https://gpggjd7cx2cizeyj.public.blob.vercel-storage.com/products/95db4863d11bef5c-cefMtdmw8MfXlby3MYYpBAdYKUXw8U.jpg", CategoryID: 22,
	},
	{
		ID: 38, Name: "Country Sourdough", Slug: "country-sourdough", Quantity: 12,
		Price: 48000, Description: "Crackly sourdough loaf with a tender open crumb",
		Sku: "BRD-SOU-001", PublicationStatus: "active",
		Image: "https://gpggjd7cx2cizeyj.public.blob.vercel-storage.com/products/fa933cd8d12feef4-2H8UTRp8IxXKUBW1Nb29kn1oNQzRNm.jpg", CategoryID: 23,
	},
	{
		ID: 39, Name: "Strawberry Celebration Cake", Slug: "strawberry-celebration-cake", Quantity: 6,
		Price: 185000, Description: "Vanilla sponge layered with cream and fresh strawberry jam",
		Sku: "CAK-STR-001", PublicationStatus: "active",
		Image: "https://gpggjd7cx2cizeyj.public.blob.vercel-storage.com/products/fa731725575f7973-MRCtBCzES50NSsBcnGnIL7hrVoR0vu.jpg", CategoryID: 24,
	},
	{
		ID: 40, Name: "Chocolate Chip Cookie", Slug: "chocolate-chip-cookie", Quantity: 36,
		Price: 18000, Description: "Golden cookie with dark chocolate chunks and crisp edges",
		Sku: "CKE-CHO-001", PublicationStatus: "active",
		Image: "https://gpggjd7cx2cizeyj.public.blob.vercel-storage.com/products/caccf4da75703bb5-Uqw6ADgcmNb04BvLebJV4GRhpXEit8.jpg", CategoryID: 25,
	},
	{
		ID: 41, Name: "Iced Vanilla Latte", Slug: "iced-vanilla-latte", Quantity: 30,
		Price: 32000, Description: "Chilled espresso with milk and house vanilla syrup",
		Sku: "BEV-VAN-001", PublicationStatus: "active",
		Image: "https://gpggjd7cx2cizeyj.public.blob.vercel-storage.com/products/9072fd741bd311df-qEJ1FwJekGLW6mEgAWQLpiG5FSgVKO.jpg", CategoryID: 26,
	},
}

func main() {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		log.Fatal("ADMIN_PASSWORD is required")
	}

	db.Init()
	defer db.Close()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash admin password: %v", err)
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := upsertAdmin(tx, string(passwordHash)); err != nil {
			return err
		}

		if err := upsertCategories(tx); err != nil {
			return err
		}

		return upsertProducts(tx)
	})
	if err != nil {
		log.Fatalf("migrate data: %v", err)
	}

	if err := resetSequences(); err != nil {
		log.Fatalf("reset sequences: %v", err)
	}

	log.Println("Admin, categories, and products migrated successfully")
}

func upsertAdmin(tx *gorm.DB, passwordHash string) error {
	var user models.User
	err := tx.Where("username = ?", "admin").First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tx.Create(&models.User{
			ID:       1,
			Username: "admin",
			Password: passwordHash,
			Email:    "admin@example.com",
		}).Error
	}
	if err != nil {
		return err
	}

	user.Password = passwordHash
	user.Email = "admin@example.com"
	return tx.Save(&user).Error
}

func upsertCategories(tx *gorm.DB) error {
	for _, seed := range categories {
		var category models.Category
		err := tx.First(&category, seed.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			category = models.Category{
				ID:     seed.ID,
				Name:   seed.Name,
				Slug:   seed.Slug,
				Status: seed.Status,
			}
			if err := tx.Create(&category).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}

		category.Name = seed.Name
		category.Slug = seed.Slug
		category.Status = seed.Status
		if err := tx.Save(&category).Error; err != nil {
			return err
		}
	}
	return nil
}

func upsertProducts(tx *gorm.DB) error {
	for _, seed := range products {
		var product models.Product
		err := tx.First(&product, seed.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			product = models.Product{ID: seed.ID}
		} else if err != nil {
			return err
		}

		product.Name = seed.Name
		product.Slug = seed.Slug
		product.Quantity = seed.Quantity
		product.Price = seed.Price
		product.Description = seed.Description
		product.Sku = seed.Sku
		product.PublicationStatus = seed.PublicationStatus
		product.Image = seed.Image

		if product.ID == seed.ID {
			if product.CreatedAt.IsZero() {
				if err := tx.Create(&product).Error; err != nil {
					return err
				}
			} else if err := tx.Save(&product).Error; err != nil {
				return err
			}
		} else if err := tx.Save(&product).Error; err != nil {
			return err
		}

		var category models.Category
		if err := tx.First(&category, seed.CategoryID).Error; err != nil {
			return fmt.Errorf("find category %d for product %d: %w", seed.CategoryID, seed.ID, err)
		}
		if err := tx.Model(&product).Association("Categories").Replace([]models.Category{category}); err != nil {
			return err
		}
	}
	return nil
}

func resetSequences() error {
	for _, table := range []string{"users", "categories", "products"} {
		query := fmt.Sprintf(
			"SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE((SELECT MAX(id) FROM %s), 1), true)",
			table,
			table,
		)
		if err := db.DB.Exec(query).Error; err != nil {
			return err
		}
	}
	return nil
}
