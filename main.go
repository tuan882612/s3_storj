package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/tuan882612/apiutils"
)

func main() {
	// Load env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get env variables
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	endpoint := os.Getenv("ENDPOINT")

	// Create a session using the aws config
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-west-1"), // set your region
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	router := chi.NewRouter()
	router.Get("/buckets", func(w http.ResponseWriter, r *http.Request) {
		buckets, err := svc.ListBuckets(nil)
		if err != nil {
			apiutils.HandleHttpErrors(w, err)
			return
		}

		name := []string{}
		for _, bucket := range buckets.Buckets {
			name = append(name, *bucket.Name)
		}

		resp := apiutils.NewRes(http.StatusOK, "", name)
		resp.SendRes(w)
	})
	router.Get("/buckets/{bucketName}", func(w http.ResponseWriter, r *http.Request) {
		bucketName := chi.URLParam(r, "bucketName")
		if bucketName == "" {
			resp := apiutils.NewRes(http.StatusBadRequest, "Bucket name is required", nil)
			resp.SendRes(w)
			return
		}

		// List objects in bucket
		objects, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
		if err != nil {
			apiutils.HandleHttpErrors(w, err)
			return
		}

		// Get object names
		name := []string{}
		for _, object := range objects.Contents {
			name = append(name, *object.Key)
		}

		resp := apiutils.NewRes(http.StatusOK, "", name)
		resp.SendRes(w)
	})

	log.Default().Println("Listening on port 8080")
	http.ListenAndServe(":8080", router)
}
