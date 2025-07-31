package main

import (
	"log"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

func main() {
	app := pocketbase.New()

	// Register the OTP request hook for the users collection
	app.OnRecordAuthRequestOTPRequest("users").BindFunc(func(e *core.RecordAuthRequestOTPRequestEvent) error {
		// If no user record exists, create one
		if e.Record == nil {
			// Extract email from request body
			email := e.HttpContext.FormValue("email")
			if email == "" {
				return core.NewBadRequestError("Email is required", nil)
			}

			// Create new user record
			collection, err := app.FindCollectionByNameOrId("users")
			if err != nil {
				return err
			}

			record := core.NewRecord(collection)
			record.SetEmail(email)
			
			// Generate a random password (user will use OTP to login)
			randomPassword := security.RandomString(30)
			record.SetPassword(randomPassword)

			// Save the new user
			if err := app.Save(record); err != nil {
				return err
			}

			// Set the record for the OTP process
			e.Record = record
		}

		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}