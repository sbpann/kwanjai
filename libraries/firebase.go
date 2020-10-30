package libraries

import (
	"kwanjai/config"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// FirebaseApp initialize firebase by credential.json.
func FirebaseApp() *firebase.App {
	sa := option.WithCredentialsFile(config.BaseDirectory + "/.secret/credential.json")
	app, err := firebase.NewApp(config.Context, nil, sa)
	if err != nil {
		log.Println(err)
	}
	return app
}
