package libraries

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// FirebaseApp initialize firebase by credential.json
func FirebaseApp() *firebase.App {
	mainPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	ctx := context.Background()
	sa := option.WithCredentialsFile(mainPath + "/.secret/credential.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Println(err)
	}
	return app
}
