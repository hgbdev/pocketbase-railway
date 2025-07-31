package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

func main() {
	app := pocketbase.New()

	// Register hook for OTP request - fires only for "users" collection
	app.OnRecordRequestOTPRequest("users").BindFunc(func(e *core.RecordCreateOTPRequestEvent) error {
		// If no user record exists for the email, create one
		if e.Record == nil {
			// Extract email from the request
			email := e.Request.PostFormValue("email")
			if email == "" {
				email = e.Request.FormValue("email")
			}
			if email == "" {
				return e.Next() // Let the original handler deal with validation
			}

			// Create new user record
			record := core.NewRecord(e.Collection)
			record.SetEmail(email)
			
			// Generate a random password (user will use OTP to login)
			randomPassword := security.RandomString(30)
			record.SetPassword(randomPassword)

			// Save the new user
			if err := e.App.Save(record); err != nil {
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