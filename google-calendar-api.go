package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/calendar/v3"
    "google.golang.org/api/option"
)

// Credentials from Google Cloud Console
type Credentials struct {
    ClientID     string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
    RedirectURL  string `json:"redirect_url"`
}

func main() {
    // Read credentials from file
    creds, err := os.ReadFile("credentials.json")
    if err != nil {
        log.Fatalf("Unable to read credentials file: %v", err)
    }

    // Configure OAuth2
    config, err := google.ConfigFromJSON(creds, calendar.CalendarReadonlyScope)
    if err != nil {
        log.Fatalf("Unable to parse credentials: %v", err)
    }

    // Get token from web
    token := getTokenFromWeb(config)

    // Create client
    client := config.Client(context.Background(), token)

    // Create Calendar service
    srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
    if err != nil {
        log.Fatalf("Unable to retrieve Calendar client: %v", err)
    }

    // Example: List next 10 events
    events, err := srv.Events.List("primary").MaxResults(10).Do()
    if err != nil {
        log.Fatalf("Unable to retrieve events: %v", err)
    }

    for _, event := range events.Items {
        fmt.Printf("Event: %v (%v)\n", event.Summary, event.Start.DateTime)
    }
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)

    fmt.Print("Enter the authorization code: ")
    var authCode string
    if _, err := fmt.Scan(&authCode); err != nil {
        log.Fatalf("Unable to read authorization code: %v", err)
    }

    token, err := config.Exchange(context.Background(), authCode)
    if err != nil {
        log.Fatalf("Unable to retrieve token from web: %v", err)
    }
    return token
}
