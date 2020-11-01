package libraries

import (
	"errors"
	"kwanjai/config"
	"log"
	"os"

	"cloud.google.com/go/firestore"
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

// FirestoreFind by collection and document ID.
func FirestoreFind(collecttion string, id string) (*firestore.DocumentSnapshot, error) {
	if collecttion == "" || id == "" {
		// create blank document
		blank := new(firestore.DocumentSnapshot)
		blank.Ref = new(firestore.DocumentRef)
		blank.Ref.Parent = new(firestore.CollectionRef)
		return blank, errors.New("invalid document reference")
	}
	firestoreClient, err := FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	document, err := firestoreClient.Collection(collecttion).Doc(id).Get(config.Context)
	return document, err
}

// FirestoreDelete by collection and document ID.
func FirestoreDelete(collecttion string, id string) (*firestore.WriteResult, error) {
	if collecttion == "" || id == "" {
		// create blank result
		blank := new(firestore.WriteResult)
		return blank, errors.New("invalid document reference")
	}
	firestoreClient, err := FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	result, err := firestoreClient.Collection(collecttion).Doc(id).Delete(config.Context)
	return result, err
}

// FirestoreSearch by collection and condition
func FirestoreSearch(collecttion string, field string, condition string, property interface{}) ([]*firestore.DocumentSnapshot, error) {
	firestoreClient, err := FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	search := firestoreClient.Collection(collecttion).Where(field, condition, property).Documents(config.Context)
	documents, err := search.GetAll()
	return documents, err
}

// FirestoreCreatedOrSet by collection, id.
func FirestoreCreatedOrSet(collecttion string, id string, data interface{}) (*firestore.WriteResult, error) {
	if collecttion == "" || id == "" {
		// create blank result
		blank := new(firestore.WriteResult)
		return blank, errors.New("invalid document reference")
	}
	firestoreClient, err := FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	result, err := firestoreClient.Collection(collecttion).Doc(id).Set(config.Context, data)
	return result, err
}

// FirestoreUpdateField by collection, id, and field.
func FirestoreUpdateField(collecttion string, id string, field string, property interface{}) (*firestore.WriteResult, error) {
	if collecttion == "" || id == "" || field == "" {
		// create blank result
		blank := new(firestore.WriteResult)
		return blank, errors.New("invalid document reference")
	}
	firestoreClient, err := FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	result, err := firestoreClient.Collection(collecttion).Doc(id).Update(config.Context, []firestore.Update{
		{
			Path:  field,
			Value: property,
		},
	})
	return result, err
}

// FirestoreDeleteField by collection, id, and field.
func FirestoreDeleteField(collecttion string, id string, field string) (*firestore.WriteResult, error) {
	firestoreClient, err := FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	result, err := firestoreClient.Collection(collecttion).Doc(id).Update(config.Context, []firestore.Update{
		{
			Path:  field,
			Value: firestore.Delete,
		},
	})
	return result, err
}
