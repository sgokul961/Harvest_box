package master

import (
	masterrepo "HarvestBox/repos/masterRepo"
	"HarvestBox/utils"
	middleware "HarvestBox/utils/midleware"
	"encoding/json"
	"log"
	"net/http"
)

type AdminHandler struct{}

func (h *AdminHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Panic recovered:", rec)
			utils.HandleError(w, nil, "Internal server error", http.StatusInternalServerError)
		}
	}()

	loginReq := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		utils.HandleError(w, err, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	adminRepo := &masterrepo.HarvestRepository{}

	token, success := adminRepo.AdminLogin(loginReq.Email, loginReq.Password)
	if !success {
		utils.HandleError(w, nil, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	// Respond with the token
	response := map[string]string{"token": token}
	utils.SuccessResponse(w, "Admin logged in successfully", response, http.StatusOK)
}
func (h *AdminHandler) GetFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			utils.HandleError(w, r.(error), "Panic recovered", http.StatusInternalServerError)
		}
	}()

	// Extract claims from request context
	claims, ok := r.Context().Value(middleware.TokenClaimsKey).(*utils.CustomClaims)
	if !ok {
		utils.HandleError(w, nil, "Unable to retrieve claims", http.StatusInternalServerError)
		return
	}

	// Authorization check - only allow access if the user is an admin
	if claims.Role != "admin" { // Adjust this check as needed
		utils.HandleError(w, nil, "Unauthorized", http.StatusForbidden)
		return
	}

	// Access the repository and get feedback
	userRepo := masterrepo.HarvestRepositoryInterface(&masterrepo.HarvestRepository{})
	feedbacks, success := userRepo.GetAllFeedback()
	if !success {
		utils.HandleError(w, nil, "Failed to retrieve feedback", http.StatusInternalServerError)
		return
	}

	// Return success response with feedback data
	utils.SuccessResponse(w, "Feedback retrieved successfully", feedbacks, http.StatusOK)
}
