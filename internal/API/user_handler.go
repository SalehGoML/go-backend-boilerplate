package api

import (
	utils "Salehaskarzadeh/internal/Utils"
	"Salehaskarzadeh/internal/storee"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/lib/pq"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore storee.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore storee.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("STEP 1: decoding request")
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	h.logger.Println("STEP 2: validating request")
	err = h.validateRegisterRequest(&req)
	if err != nil {
		h.logger.Printf("ERROR: validation failed: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	h.logger.Println("STEP 3: preparing user struct")
	user := &storee.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: storee.Password{},
		Bio:          req.Bio,
	}

	h.logger.Println("STEP 4: hashing password")
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: hashing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to hash password"})
		return
	}

	h.logger.Println("STEP 5: creating user in database")
	err = h.userStore.CreateUser(user)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			h.logger.Printf("ERROR: database error: %s (%s)", pgErr.Message, pgErr.Code)
			if pgErr.Code == "23505" {
				utils.WriteJSON(w, http.StatusConflict, utils.Envelope{"error": "username or email already exists"})
				return
			}
		} else {
			h.logger.Printf("ERROR: registering user: %v", err)
		}
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	h.logger.Println("STEP 6: user created successfully")
	response := map[string]interface{}{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"bio":       user.Bio,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": response})
}
