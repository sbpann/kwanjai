package libraries

import (
	"context"
	"gin-sandbox/config"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// FirebaseApp initialize firebase by credential.json
func FirebaseApp() *firebase.App {
	ctx := context.Background()
	sa := option.WithCredentialsFile(config.BaseDirectory + "/.secret/credential.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Println(err)
	}
	return app
}
