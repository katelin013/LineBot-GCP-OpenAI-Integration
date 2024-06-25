package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/storage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func main() {
	log.Println("Starting application...")

	// Setup Google Cloud Storage and Secret Manager clients
	ctx := context.Background()
	log.Println("Creating Google Cloud Storage client...")
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	log.Println("Google Cloud Storage client created.")

	// Download the credential file
	bucketName := "linegpt-427007.appspot.com"
	objectName := "linegpt-427007-aaee6b8a5f53.json"
	localPath := "/tmp/linegpt-427007-aaee6b8a5f53.json"
	log.Printf("Downloading credential file from bucket %s, object %s...", bucketName, objectName)

	err = downloadFile(ctx, storageClient, bucketName, objectName, localPath)
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}
	log.Println("Credential file downloaded.")

	// Set GOOGLE_APPLICATION_CREDENTIALS environment variable
	log.Println("Setting GOOGLE_APPLICATION_CREDENTIALS environment variable...")
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", localPath)
	log.Println("GOOGLE_APPLICATION_CREDENTIALS environment variable set.")

	// Setup Secret Manager client
	log.Println("Creating Secret Manager client...")
	secretClient, err := secretmanager.NewClient(ctx, option.WithCredentialsFile(localPath))
	if err != nil {
		log.Fatalf("Failed to setup secretmanager client: %v", err)
	}
	log.Println("Secret Manager client created.")

	// Get secrets from Secret Manager
	log.Println("Accessing secrets from Secret Manager...")
	lineChannelSecret := accessSecretVersion(secretClient, "projects/linegpt-427007/secrets/LINE_CHANNEL_SECRET/versions/latest")
	lineChannelAccessToken := accessSecretVersion(secretClient, "projects/linegpt-427007/secrets/LINE_CHANNEL_ACCESS_TOKEN/versions/latest")
	openaiAPIKey := accessSecretVersion(secretClient, "projects/linegpt-427007/secrets/OPENAI_API_KEY/versions/latest")
	log.Println("Secrets accessed.")

	// Initialize LINE Bot SDK
	log.Println("Initializing LINE Bot SDK...")
	bot, err := linebot.New(lineChannelSecret, lineChannelAccessToken)
	if err != nil {
		log.Fatalf("Failed to initialize LINE Bot SDK: %v", err)
	}
	log.Println("LINE Bot SDK initialized.")

	// Initialize OpenAI client
	log.Println("Initializing OpenAI client...")
	openaiClient := openai.NewClient(openaiAPIKey)
	log.Println("OpenAI client initialized.")

	// Root handler for testing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
		log.Println("Handled root request.")
	})

	// Webhook callback handler
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		log.Println("Received callback request.")
		signature := req.Header.Get("X-Line-Signature")

		// Verify the request signature
		if !verifySignature(req, lineChannelSecret, signature) {
			log.Println("Invalid signature.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				log.Println("Invalid signature.")
				w.WriteHeader(http.StatusBadRequest)
			} else {
				log.Println("Error parsing request:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					log.Printf("Received message from LINE: %s", message.Text)
					log.Println("Creating OpenAI ChatCompletion request with model GPT-4...")
					resp, err := openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
						Model: openai.GPT4o,
						Messages: []openai.ChatCompletionMessage{
							{Role: "system", Content: "You are a helpful assistant."},
							{Role: "user", Content: message.Text},
						},
					})
					if err != nil {
						log.Println("ChatCompletion error:", err)
						continue
					}
					log.Println("ChatCompletion request successful.")
					log.Printf("OpenAI response: %s", resp.Choices[0].Message.Content)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(resp.Choices[0].Message.Content)).Do(); err != nil {
						log.Printf("ReplyMessage error: %v", err)
					} else {
						log.Println("ReplyMessage sent successfully.")
					}
				}
			}
		}
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func downloadFile(ctx context.Context, client *storage.Client, bucketName, objectName, localPath string) error {
	log.Printf("Starting file download: bucket=%s, object=%s, localPath=%s", bucketName, objectName, localPath)
	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)
	r, err := object.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("object.NewReader: %v", err)
	}
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	if err := ioutil.WriteFile(localPath, data, 0644); err != nil {
		return fmt.Errorf("ioutil.WriteFile: %v", err)
	}

	log.Println("File downloaded successfully.")
	return nil
}

func accessSecretVersion(client *secretmanager.Client, name string) string {
	log.Printf("Accessing secret version: %s", name)
	ctx := context.Background()
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Fatalf("Failed to access secret version: %v", err)
	}

	log.Printf("Secret version accessed: %s", name)
	return string(result.Payload.Data)
}

func verifySignature(req *http.Request, secret, signature string) bool {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		return false
	}
	req.Body = ioutil.NopCloser(strings.NewReader(string(body))) // Reset request body

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	log.Printf("Expected Signature: %s", expectedSignature)
	log.Printf("Received Signature: %s", signature)

	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
