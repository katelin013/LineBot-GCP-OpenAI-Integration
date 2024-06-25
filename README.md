# LineBot-GCP-OpenAI-Integration

## Repository Description
This repository provides a comprehensive guide and implementation for integrating a LINE Bot with Google Cloud Platform (GCP) services and OpenAI API. It covers the steps required to create and configure a LINE Bot, set up necessary GCP services, and deploy the application on Google App Engine.

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Setup Guide](#setup-guide)
3. [File Descriptions](#file-descriptions)
4. [PowerShell Commands](#powershell-commands)
5. [GCP Services Used](#gcp-services-used)
6. [Testing and Verification](#testing-and-verification)

## Prerequisites
- LINE Developer Account
- Google Cloud Platform Account
- OpenAI API Key
- PowerShell installed on your local machine

## Setup Guide

### Step 1: Create a LINE Bot
1. **Enter LINE Console**: Go to the Provider page and click "Create" ![Enter LINE Console](images/2.jpg).
2. **Enter Provider Name**: Input the new provider name ![Enter Provider Name](images/3.jpg).
3. **Select Channel Type**: Choose "Create a Messaging API Channel" ![Select Channel Type](images/4.jpg).
4. **Input Company or Owner's Country or Region**: Enter the relevant information ![Input Company or Owner's Country or Region](images/5.jpg).
5. **Set Channel Icon, Name, and Description**: Provide the necessary details ![Set Channel Icon, Name, and Description](images/6.jpg).
6. **Set Category and Subcategory**: Select the appropriate options ![Set Category and Subcategory](images/7.jpg).
7. **Remember Channel Secret**: Note down the Channel Secret ![Remember Channel Secret](images/8.jpg).
8. **Go to LINE Bot Page**: Navigate to the Basic settings page ![Go to LINE Bot Page](images/9.jpg).
9. **Webhook Not Set Yet**: Webhook will be configured later ![Webhook Not Set Yet](images/10.jpg).
10. **Channel Access Token**: Scroll down and note the Channel Access Token ![Channel Access Token](images/11.jpg).

### Step 2: Set Up Google Cloud Platform
1. **Open GCP Console**: Click on "Create Project" ![Open GCP Console](images/12.jpg).
2. **Set Project Name**: Enter the project name ![Set Project Name](images/13.jpg).
3. **Go Back to Project Home Page**: Navigate back ![Go Back to Project Home Page](images/14.jpg).
4. **Select Newly Created Project**: Choose the project ![Select Newly Created Project](images/15.jpg).
5. **Search App Engine**: Find App Engine ![Search App Engine](images/16.jpg).
6. **App Engine Home Page**: Click "Create" ![App Engine Home Page](images/17.jpg).
7. **Select Region and Set API Access**: Choose the region and set permissions ![Select Region and Set API Access](images/18.jpg).
8. **Go Back to App Engine**: Click "Get Started" ![Go Back to App Engine](images/19.jpg).
9. **App Engine Page**: Select the Go runtime ![App Engine Page](images/20.jpg).

### Step 3: Prepare Local Directory
1. **Create Folder**: The folder should contain the following files ![Create Folder](images/dir.jpg):
   - `.gcloudignore`: Specify files and directories to ignore when deploying.
   - `app.yaml`: Configuration file for App Engine.
   - `cloudbuild.yaml`: Build configuration for Cloud Build.
   - `Dockerfile`: Instructions to build the Docker image.
   - `go.mod`: Module definition file.
   - `go.sum`: Checksums for module dependencies.
   - `main.go`: Main application code.

### Step 4: Deploy Application
1. **Deploy Application**: Use PowerShell to deploy the application ![Deploy Application](images/gcp-deploy-1.jpg).
2. **Upload to Google Cloud Storage**: Run `gcloud app deploy` to upload ![Upload to Google Cloud Storage](images/gcp-deploy-2.jpg).
3. **Get Webhook URL**: Obtain the URL using `gcloud app browse --no-launch-browser` ![Get Webhook URL](images/get-webhookurl.jpg).

### Step 5: Set Up Go Modules
1. **Initialize Go Module**: Run `go mod init line-gpt-bot` ![Initialize Go Module](images/go-1.jpg).
2. **Install Dependencies**:
   - `go get github.com/line/line-bot-sdk-go/v7/linebot` ![Install Dependency](images/go-2.jpg).
   - `go get github.com/sashabaranov/go-openai` ![Install Dependency](images/go-3.jpg).
   - `go get cloud.google.com/go/secretmanager/apiv1` ![Install Dependency](images/go-4.jpg).
   - `go get google.golang.org/api/option` ![Install Dependency](images/go-5.jpg).
   - `go get cloud.google.com/go/iam/apiv1/iampb@v1.1.8` ![Install Dependency](images/go-6.jpg).

### Step 6: Final Deployment
1. **Deploy with gcloud**: Run `gcloud app deploy` ![Deploy with gcloud](images/go-7-gcloud-deploy-1.jpg).
2. **Confirm Deployment**: Ensure deployment is successful ![Confirm Deployment](images/go-7-gcloud-deploy-2.jpg).
3. **Deployment Complete**: Deployment process ![Deployment Complete](images/go-7-gcloud-deploy-3.jpg).

### Step 7: Configure LINE Bot
1. **Go Back to LINE Console**: Navigate back to the main page ![Go Back to LINE Console](images/Line01.jpg).
2. **Edit Group Chat Settings**: Disable group chat invitations ![Edit Group Chat Settings](images/Line02.jpg, images/Line03.jpg).
3. **Edit Auto-Reply Messages**: Disable welcome message and enable webhook ![Edit Auto-Reply Messages](images/Line04.jpg, images/Line05.jpg).
4. **Manual Chat Response**: Set response mode to manual ![Manual Chat Response](images/Line07.jpg).
5. **Enter Webhook URL**: Input the webhook URL and update ![Enter Webhook URL](images/Line-input-webhook-url.jpg).
6. **Verify Webhook**: Verify the URL ![Verify Webhook](images/Line-Verify01.jpg).

### Step 8: Set Up Secret Manager
1. **Open Secret Manager**: Search for Secret Manager in GCP Console ![Open Secret Manager](images/SecretManager-1.jpg).
2. **Enable Secret Manager**: Activate the service ![Enable Secret Manager](images/SecretManager-2.jpg).
3. **Create Secrets**: Add new secrets ![Create Secrets](images/SecretManager-3.jpg, images/SecretManager-4.jpg).
4. **Secrets Created**: Ensure the following secrets are created ![Secrets Created](images/SecretManager-5.jpg):
   - LINE_CHANNEL_ACCESS_TOKEN
   - LINE_CHANNEL_SECRET
   - OPENAI_API_KEY

## File Descriptions
- **.gcloudignore**: Lists files and directories to ignore during deployment.
- **app.yaml**: Configuration file for deploying the app to App Engine.
- **cloudbuild.yaml**: Defines the build steps for Cloud Build.
- **Dockerfile**: Contains the instructions to build a Docker image for the app.
- **go.mod**: Defines the module's path and its dependencies.
- **go.sum**: Records the exact versions of dependencies used.
- **main.go**: Main application code.

## PowerShell Commands
- `gcloud app deploy`: Deploy the application to Google App Engine.
- `gcloud app browse --no-launch-browser`: Retrieve the Webhook URL without launching the browser.

## GCP Services Used
- **App Engine**: Host and run the application.
- **Cloud Storage**: Store application credentials.
- **Secret Manager**: Securely manage and access secrets.
- **Cloud Build**: Build and deploy the application.
- **IAM**: Manage permissions and access control.

## Testing and Verification
- **LINE Chat Response**: Verify that the LINE Bot responds correctly to user messages ![LINE Chat Response](images/Line-chat.jpg).

## Conclusion
This repository serves as a complete guide to setting up a LINE Bot that interacts with OpenAI's API using GCP services.
