package libraries

import (
	"kwanjai/config"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// FirebaseApp initialize firebase by credential.json.
func FirebaseApp() *firebase.App {
	var err error
	var app *firebase.App
	if os.Getenv("GIN_MODE") == "release" {
		conf := &firebase.Config{ProjectID: config.FirebaseProjectID}
		app, err = firebase.NewApp(config.Context, conf)
	} else {
		sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
		app, err = firebase.NewApp(config.Context, nil, sa)
	}
	if err != nil {
		log.Println(err)
	}
	return app
}
