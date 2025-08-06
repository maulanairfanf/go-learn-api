package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"myapi/db"
	"myapi/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetProducts handles retrieving all products along with their categories
func GetProducts(w http.ResponseWriter, r *http.Request) {
        log.Println("Ini pesan log biasa", )

    var products []models.Product

    if err := db.DB.Preload("Categories").Find(&products).Error; err != nil {
        ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    SuccessResponse(w, products)
}

// GetProduct handles retrieving a single product by ID along with its categories
func GetProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        ErrorResponse(w, http.StatusBadRequest, "Invalid product ID")
        return
    }

    var product models.Product
    if err := db.DB.Preload("Categories").First(&product, id).Error; err != nil {
        ErrorResponse(w, http.StatusNotFound, "Product not found")
        return
    }

    SuccessResponse(w, product)
}

// CreateProduct handles the creation of a new product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
    var product models.Product

    // Decode the JSON request body into the product struct
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    // Create the product in the database
    if err := db.DB.Create(&product).Error; err != nil {
        ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    SuccessResponse(w, product)
}

// UpdateProduct handles the update of an existing product
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        ErrorResponse(w, http.StatusBadRequest, "Invalid product ID")
        return
    }

    var product models.Product
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    // Find the existing product
    // Untuk menampilkan isi body, Anda perlu membaca body terlebih dahulu.
    // Namun, body sudah dibaca oleh json.NewDecoder di atas, sehingga Anda perlu membacanya sebelum decode.
    // Berikut contoh cara membaca body sebagai string (gunakan sebelum decode):

    bodyBytes, _ := io.ReadAll(r.Body)
    log.Println("Isi body:", string(bodyBytes))
    // r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // agar body bisa dibaca ulang oleh decoder

    // Namun, pada posisi ini body sudah dibaca, jadi log body di sini tidak akan menampilkan isinya.
    // Solusi: log body sebelum json.NewDecoder(r.Body).Decode(&product)
    if err := db.DB.First(&product, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            ErrorResponse(w, http.StatusNotFound, "Product not found")
        } else {
            ErrorResponse(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    // Update the product
    if err := db.DB.Save(&product).Error; err != nil {
        ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    SuccessResponse(w, product)
}

// DeleteProduct handles the deletion of a product by ID
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        ErrorResponse(w, http.StatusBadRequest, "Invalid product ID")
        return
    }

    var product models.Product

    // First, find the product and preload its categories
    if err := db.DB.Preload("Categories").First(&product, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            ErrorResponse(w, http.StatusNotFound, "Product not found")
        } else {
            ErrorResponse(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    // Detach from categories first (clears join table entries)
    if err := db.DB.Model(&product).Association("Categories").Clear(); err != nil {
        ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    // Now delete the product
    if err := db.DB.Delete(&product).Error; err != nil {
        ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    SuccessResponse(w, "Product deleted successfully")
}
