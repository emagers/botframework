# BotFramework

A simple package for working with [Azure's Bot Framework](https://dev.botframework.com/).

## Current Functionality

1. Gets and stores access token using the flow documented [here](https://docs.microsoft.com/en-us/azure/bot-service/rest-api/bot-framework-rest-connector-authentication?view=azure-bot-service-4.0#bot-to-connector)

## Usage

```go
configuration := Configuration {
	<Application ID>,
	<App Password>,
}

authClient := AuthenticationClient.Init(configuration)

token, err := authClient.GetAccessToken()
```

The first call to `GetAccessToken()` will make a call to retrieve an access token and will cache that token until it expires. At that point, the next call to `GetAccessToken()` will get a new token.