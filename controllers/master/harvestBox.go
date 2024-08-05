package master

import (
	"HarvestBox/models"
	masterrepo "HarvestBox/repos/masterRepo"
	"HarvestBox/utils"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
}

func (h *UserHandler) AddFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			utils.HandleError(w, r.(error), "Panic recovered", http.StatusInternalServerError)
		}
	}()

	feedbackReq := models.Feedback{}
	err := json.NewDecoder(r.Body).Decode(&feedbackReq)
	if err != nil {
		utils.HandleError(w, err, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	userRepo := masterrepo.HarvestRepositoryInterface(&masterrepo.HarvestRepository{})

	user, err := userRepo.FindOrCreateUser(feedbackReq.Name, feedbackReq.Email)
	if err != nil {
		utils.HandleError(w, err, "Failed to find or create user", http.StatusInternalServerError)
		return
	}

	feedbackReq.UserID = user.UserID

	if !userRepo.AddNewFeedback(&feedbackReq) {
		utils.HandleError(w, nil, "Failed to add feedback", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Successfully added feedback", true, http.StatusOK)
}

func (h *UserHandler) AddNewUserHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			utils.HandleError(w, r.(error), "Panic recovered", http.StatusInternalServerError)
		}
	}()

	userReq := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		utils.HandleError(w, err, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:     userReq.Name,
		Email:    userReq.Email,
		Password: userReq.Password,
	}

	userRepo := masterrepo.HarvestRepositoryInterface(&masterrepo.HarvestRepository{})
	if !userRepo.AddNewUser(&user) {
		utils.HandleError(w, nil, "Failed to add user", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Successfully added user", true, http.StatusOK)
}
