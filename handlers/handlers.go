package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/harshgupta9473/fi/dto"
	"github.com/harshgupta9473/fi/services"
	"net/http"
	"strconv"
)

type Handlers struct {
	ProductService services.ProductServiceIntf
	UserService    services.UserServiceIntf
}

func NewHandler(productService services.ProductServiceIntf, userService services.UserServiceIntf) *Handlers {
	return &Handlers{
		ProductService: productService,
		UserService:    userService,
	}
}

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user *dto.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {

	}
	if user.Username == "" || user.Password == "" {

	}
	err = h.UserService.LoginUser(user)
	if err != nil {

	}

}

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {

	}
	if user.Username == "" || user.Password == "" {

	}
	err = h.UserService.CreateUserAccount(&user)
	if err != nil {

	}

}

func (h *Handlers) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {

	}
	id, err := h.ProductService.AddProduct(&product)
	if err != nil {

	}
	// return id if success
}

func (h *Handlers) GetAllProducts(w http.ResponseWriter, r *http.Request) {

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page := 1
	limit := 10
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
	}

	products, err := h.ProductService.GetALLProducts(page, limit)
	if err != nil {

	}
}

func (h *Handlers) UpdateProductQuantitty(w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {

	}
	var quantity struct {
		Quantity int64 `json:"quantity"`
	}
	err = json.NewDecoder(r.Body).Decode(&quantity)
	if err != nil {

	}
	product, err, _ := h.ProductService.UpdateProduct(id, quantity.Quantity)
	if err != nil {

	}
	if product == nil {

	}

}
