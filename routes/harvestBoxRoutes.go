package routes

import (
	"HarvestBox/controllers/master"
	middleware "HarvestBox/utils/midleware"
	"net/http"

	"github.com/gorilla/mux"
)

func HarvestBoxRoutes(router *mux.Router) {
	HarvestBoxRouter := router.PathPrefix("/master").Subrouter()
	registerUserRoutes(HarvestBoxRouter)
	registerAdminRoutes(HarvestBoxRouter)
}

func registerUserRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()
	user := master.UserHandler{}

	//userRouter.Handle("/addUser", http.HandlerFunc(user.AddNewUserHandler)).Methods("POST")
	userRouter.Handle("/addFeedback", http.HandlerFunc(user.AddFeedbackHandler)).Methods("POST")
}

func registerAdminRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/admin").Subrouter()
	admin := master.AdminHandler{}

	// Public route - no token required
	adminRouter.Handle("/login", http.HandlerFunc(admin.LoginHandler)).Methods("POST")

	// Apply middleware only to protected routes
	protectedRoutes := adminRouter.PathPrefix("/").Subrouter()

	protectedRoutes.Use(middleware.ValidateTokenMiddleware)

	// Protected routes - token required
	protectedRoutes.Handle("/feedbacks", http.HandlerFunc(admin.GetFeedbackHandler)).Methods("GET")
}
