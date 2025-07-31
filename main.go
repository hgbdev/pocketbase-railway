package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

func main() {
	app := pocketbase.New()

	// Register hook for record auth with OTP (before request processing)
	app.OnRecordAuthWithOTPRequest().BindFunc(func(e *core.RecordAuthWithOTPRequestEvent) error {
		// Check if this is for the users collection
		if e.Collection.Name != "users" {
			return e.Next()
		}

		// If no user record exists for the email, create one
		if e.Record == nil {
			// Extract email from the request context
			email := e.RequestEvent.Request.PostFormValue("email")
			if email == "" {
				email = e.RequestEvent.Request.FormValue("email")
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