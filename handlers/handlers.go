package handlers

import (
	"context"
	"encoding/json"
	"github.com/harshgupta9473/fi/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/fi/dto"
	"github.com/harshgupta9473/fi/services"
	"github.com/harshgupta9473/fi/utils"
)

type Handlers struct {
	ProductService services.ProductServiceIntf
	UserService    services.UserServiceIntf
	TimeOut        time.Duration
}

func NewHandler(productService services.ProductServiceIntf, userService services.UserServiceIntf) *Handlers {
	return &Handlers{
		ProductService: productService,
		UserService:    userService,
		TimeOut:        5,
	}
}

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.TimeOut*time.Second)
	defer cancel()

	var user *dto.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}
	if user.Username == "" || user.Password == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Username and password are required"})
		return
	}
	err = h.UserService.LoginUser(ctx, user)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	token, err := middleware.CreateJWTToken(user.Username)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"access_token": token})
}

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.TimeOut*time.Second)
	defer cancel()
	var user dto.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		return
	}
	if user.Username == "" || user.Password == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Username and password are required"})
		return
	}
	err = h.UserService.CreateUserAccount(ctx, &user)
	if err != nil {
		utils.WriteJSON(w, 409, map[string]string{"error": err.Error()})
		return
	}

	token, err := middleware.CreateJWTToken(user.Username)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"access_token": token})
}

func (h *Handlers) AddProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.TimeOut*time.Second)
	defer cancel()
	var product dto.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		return
	}
	id, err := h.ProductService.AddProduct(ctx, &product)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
    "product_id": id,
})
}

func (h *Handlers) GetAllProducts(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), h.TimeOut*time.Second)
	defer cancel()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page := 1
	limit := 10
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid page number"})
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid limit"})
			return
		}
	}

	products, err := h.ProductService.GetALLProducts(ctx, page, limit)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handlers) UpdateProductQuantitty(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.TimeOut*time.Second)
	defer cancel()

	idstr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}
	var quantity struct {
		Quantity int64 `json:"quantity"`
	}
	err = json.NewDecoder(r.Body).Decode(&quantity)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		return
	}
	product, err, _ := h.ProductService.UpdateProduct(ctx, id, quantity.Quantity)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if product == nil {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"success": "Updated product quantity",
			"data": "Product not able to fetch",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
    "quantity": product.Quantity,
   })

}
