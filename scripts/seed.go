package main

import (
	"fmt"
	"log"

	"flypro-assessment-ayo/internal/config"
	"flypro-assessment-ayo/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create users
	users := []models.User{
		{Email: "john.doe@example.com", Name: "John Doe"},
		{Email: "jane.smith@example.com", Name: "Jane Smith"},
		{Email: "bob.wilson@example.com", Name: "Bob Wilson"},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			log.Printf("User %s might already exist: %v", users[i].Email, err)
			// Try to fetch existing user
			db.Where("email = ?", users[i].Email).First(&users[i])
		}
		fmt.Printf("Created/Found user: %s (ID: %d)\n", users[i].Email, users[i].ID)
	}

	// Create expenses for user 1
	expenses := []models.Expense{
		{
			UserID:      users[0].ID,
			Amount:      150.50,
			Currency:    "USD",
			Category:    models.CategoryTravel,
			Description: "Flight ticket to New York",
			Status:      models.ExpenseStatusPending,
		},
		{
			UserID:      users[0].ID,
			Amount:      75.00,
			Currency:    "USD",
			Category:    models.CategoryMeals,
			Description: "Dinner with client",
			Status:      models.ExpenseStatusApproved,
		},
		{
			UserID:      users[0].ID,
			Amount:      45.25,
			Currency:    "EUR",
			Category:    models.CategoryOfficeSupplies,
			Description: "Office supplies",
			Status:      models.ExpenseStatusPending,
		},
		{
			UserID:      users[1].ID,
			Amount:      200.00,
			Currency:    "USD",
			Category:    models.CategoryTravel,
			Description: "Hotel accommodation",
			Status:      models.ExpenseStatusPending,
		},
		{
			UserID:      users[1].ID,
			Amount:      50.00,
			Currency:    "USD",
			Category:    models.CategoryMeals,
			Description: "Lunch meeting",
			Status:      models.ExpenseStatusApproved,
		},
	}

	for i := range expenses {
		if err := db.Create(&expenses[i]).Error; err != nil {
			log.Printf("Failed to create expense: %v", err)
			continue
		}
		fmt.Printf("Created expense: %s - %.2f %s (ID: %d)\n",
			expenses[i].Description, expenses[i].Amount, expenses[i].Currency, expenses[i].ID)
	}

	// Create expense reports
	reports := []models.ExpenseReport{
		{
			UserID: users[0].ID,
			Title:  "Q1 2024 Travel Expenses",
			Status: models.ReportStatusDraft,
			Total:  270.75,
		},
		{
			UserID: users[1].ID,
			Title:  "March 2024 Business Trip",
			Status: models.ReportStatusSubmitted,
			Total:  250.00,
		},
	}

	for i := range reports {
		if err := db.Create(&reports[i]).Error; err != nil {
			log.Printf("Failed to create report: %v", err)
			continue
		}
		fmt.Printf("Created report: %s (ID: %d)\n", reports[i].Title, reports[i].ID)

		// Add expenses to reports
		if i == 0 && len(expenses) >= 3 {
			// Add first 3 expenses to first report
			if err := db.Model(&reports[i]).Association("Expenses").Append(&expenses[0], &expenses[1], &expenses[2]); err != nil {
				log.Printf("Failed to add expenses to report: %v", err)
			}
		} else if i == 1 && len(expenses) >= 5 {
			// Add last 2 expenses to second report
			if err := db.Model(&reports[i]).Association("Expenses").Append(&expenses[3], &expenses[4]); err != nil {
				log.Printf("Failed to add expenses to report: %v", err)
			}
		}
	}

	fmt.Println("\nDatabase seeded successfully!")
}
