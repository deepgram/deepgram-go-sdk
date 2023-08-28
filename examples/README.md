# Examples for Testing Features Locally

The example projects are meant to be used to test features locally by contributors working on this SDK.

## Steps to Test Your Code

If you are contributing changes to this SDK, you can test those changes by using the `prerecorded` or `streaming` projects in the `examples` folder. Here are the steps to follow:

### Add Your Code

Make your changes to the SDK (be sure you are on a branch you have created to do this work).

### Install dependencies

If the project requires third-party dependencies and not just standard library dependencies, you will have to install them. Make sure you are in the folder of the specific project (for example: `streaming`) and then use this command: 

```
go mod tidy
```

### Edit the API key, the file, and the mimetype (as needed)

Replace the API key where it says "YOUR_DEEPGRAM_API_KEY"

```go
DEEPGRAM_API_KEY = "YOUR_DEEPGRAM_API_KEY"
```

### Run the project

Make sure you're in the directory with the `main.go` file and run the project with the following command.

```
go run main.go
```