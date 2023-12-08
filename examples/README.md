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

### Set your API Key as an Environment Variable named "DEEPGRAM_API_KEY"

If using bash, this could be done in your `~/.bash_profile` like so:

```bash
export DEEPGRAM_API_KEY = "YOUR_DEEPGRAM_API_KEY"
```

or this could also be done by a simple export before executing your Go application:

```bash
DEEPGRAM_API_KEY="YOUR_DEEPGRAM_API_KEY" go run main.go
```

### Run the project

If you chose to set an environment variable in your shell profile (ie `.bash_profile`) you can change directory into each example folder and run the example like so:

```bash
go run main.go
```