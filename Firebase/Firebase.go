package firebase

import (
	"log"
	"golang.org/x/net/context"
	firebase "firebase.google.com/go"
	//"firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"cloud.google.com/go/storage"
)

var Bucket *storage.BucketHandle
var Attrs *storage.BucketAttrs
var ClientOpts option.ClientOption

// InitializeFirebase - initialise firebase
func InitializeFirebase () {
	ClientOpts = option.WithCredentialsFile("filename")
	config := &firebase.Config {
		DatabaseURL: "DBURL",
  		StorageBucket: "StorageURL",
	}

	app, err := firebase.NewApp(context.Background(), config, ClientOpts)
	if err != nil {
		panic(err)
	}

	client, err := app.Storage(context.Background())
	if err != nil {
			log.Fatalln(err)
	}

	Bucket, err = client.DefaultBucket()
	Attrs, err = Bucket.Attrs(context.Background())
 	if err != nil {
			log.Fatalln(err)
	}
}
