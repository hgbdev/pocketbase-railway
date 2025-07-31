package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

func main() {
	app := pocketbase.New()

	// Register the OTP request hook for the users collection
	app.OnRecordAuthRequestOTP("users").Add(func(e *core.RecordAuthRequestOTPEvent) error {
		// If no user record exists, create one
		if e.Record == nil {
			// Extract email from request body
			requestData := apis.RequestInfo(e.HttpContext).Data
			emailValue, exists := requestData["email"]
			if !exists {
				return apis.NewBadRequestError("Email is required", nil)
			}

			email, ok := emailValue.(string)
			if !ok {
				return apis.NewBadRequestError("Email must be a string", nil)
			}

			// Create new user record
			collection, err := app.Dao().FindCollectionByNameOrId("users")
			if err != nil {
				return err
			}

			record := collection.NewRecord()
			record.SetEmail(email)
			
			// Generate a random password (user will use OTP to login)
			randomPassword := security.RandomString(30)
			record.SetPassword(randomPassword)

			// Save the new user
			if err := app.Dao().SaveRecord(record); err != nil {
				return err
			}

			// Set the record for the OTP process
			e.Record = record
		}

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}