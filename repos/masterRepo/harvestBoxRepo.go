package masterrepo

import (
	"HarvestBox/config"
	"HarvestBox/models"
	"HarvestBox/utils"
	"fmt"
	"log"
	"time"
)

type HarvestRepositoryInterface interface {
	AddNewFeedback(input *models.Feedback) bool
	GetAllFeedback() ([]models.Feedback, bool)
	AddNewUser(input *models.User) bool
	FindOrCreateUser(name, email string) (*models.User, error)
}

type HarvestRepository struct{}

func (r *HarvestRepository) AddNewFeedback(input *models.Feedback) bool {
	myDb, err := config.DbConnection()
	if err != nil {
		log.Println("Database Connection Error:", err)
		return false
	}
	defer myDb.Close()

	queryInsert := `
        INSERT INTO Feedback 
		(user_id, name, age, occupation, mail_id, would_recommend, suggestion, likes, created_at)
        VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING feedback_id;
    `
	var id int
	err = myDb.QueryRow(queryInsert, input.UserID, input.Name, input.Age, input.Occupation, input.MailID, input.WouldRecommend, input.Suggestion, input.Likes, time.Now()).Scan(&id)
	if err != nil {
		log.Println("Error adding new feedback:", err)
		return false
	}

	fmt.Println("New feedback added with ID:", id)
	return true
}

func (r *HarvestRepository) GetAllFeedback() ([]models.Feedback, bool) {
	myDb, err := config.DbConnection()
	if err != nil {
		log.Println("Database Connection Error:", err)
		return nil, false
	}
	defer myDb.Close()

	query := "SELECT feedback_id, user_id, name, age, occupation, mail_id, would_recommend, suggestion, likes, created_at FROM Feedback"
	rows, err := myDb.Query(query)
	if err != nil {
		log.Println("Error retrieving feedback:", err)
		return nil, false
	}
	defer rows.Close()

	var feedbacks []models.Feedback
	for rows.Next() {
		var feedback models.Feedback
		err := rows.Scan(&feedback.FeedbackID,
			&feedback.UserID,
			&feedback.Name,
			&feedback.Age,
			&feedback.Occupation,
			&feedback.MailID,
			&feedback.WouldRecommend,
			&feedback.Suggestion,
			&feedback.Likes,
			&feedback.CreatedAt)
		if err != nil {
			log.Println("Error scanning feedback row:", err)
			return nil, false
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, false
	}

	return feedbacks, true
}

// func (r *HarvestRepository) AdminLogin(email, password string) bool {
// 	myDb, err := config.DbConnection()
// 	if err != nil {
// 		log.Println("Database Connection Error:", err)
// 		return false
// 	}
// 	defer myDb.Close()

// 	var dbPassword string
// 	query := "SELECT password FROM users WHERE email = $1"
// 	err = myDb.QueryRow(query, email).Scan(&dbPassword)
// 	if err != nil {
// 		log.Println("Error retrieving admin credentials:", err)
// 		return false
// 	}

// 	if dbPassword != password {
// 		log.Println("Invalid login credentials")
// 		return false
// 	}

//		return true
//	}
func (r *HarvestRepository) AddNewUser(input *models.User) bool {
	myDb, err := config.DbConnection()
	if err != nil {
		log.Println("Database Connection Error:", err)
		return false
	}
	defer myDb.Close()

	queryInsert := `
        INSERT INTO Users (name, email, password, is_admin, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id;
    `
	var id int
	err = myDb.QueryRow(queryInsert, input.Name, input.Email, input.Password, false, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		log.Println("Error adding new user:", err)
		return false
	}

	fmt.Println("New user added with ID:", id)
	return true
}

func (r *HarvestRepository) FindOrCreateUser(name, email string) (*models.User, error) {
	myDb, err := config.DbConnection()
	if err != nil {
		log.Println("Database Connection Error:", err)
		return nil, err
	}
	defer myDb.Close()

	user := &models.User{}
	querySelect := `SELECT user_id, name, email FROM Users WHERE email = $1`
	err = myDb.QueryRow(querySelect, email).Scan(&user.UserID, &user.Name, &user.Email)
	if err == nil {
		return user, nil
	}

	// User does not exist, create a new user
	queryInsert := `
        INSERT INTO Users (name, email, password, is_admin, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id;
    `
	err = myDb.QueryRow(queryInsert, name, email, "", false, time.Now(), time.Now()).Scan(&user.UserID)
	if err != nil {
		log.Println("Error creating new user:", err)
		return nil, err
	}

	user.Name = name
	user.Email = email
	return user, nil
}

func (r *HarvestRepository) AdminLogin(email, password string) (string, bool) {
	myDb, err := config.DbConnection()
	if err != nil {
		log.Println("Database Connection Error:", err)
		return "", false
	}
	defer myDb.Close()

	var hashedPassword string
	var userID int
	query := "SELECT password, user_id FROM Users WHERE email = $1 AND is_admin = TRUE"
	err = myDb.QueryRow(query, email).Scan(&hashedPassword, &userID)
	if err != nil {
		log.Println("Error retrieving admin credentials:", err)
		return "", false
	}

	// Generate a token for the admin user
	token, err := utils.GenerateToken(userID, email, "admin")
	if err != nil {
		log.Println("Error generating token:", err)
		return "", false
	}

	return token, true
}
