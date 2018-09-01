package main

import (
    "os"
    "fmt"
    "log"

    "golang.org/x/net/context"
    "google.golang.org/api/admin/directory/v1"
)

var cliOptionsMsg = `Possible options:
- gmapper token

`
var confFile string = "config.json"
var tokFile string  = "token.json"
var credFile string = "credentials.json"


func main() {
    // Check command line arguments
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "token":
            genOauthToken()
        default:
            fmt.Print(cliOptionsMsg)
        }
    } else {
        startMapper()
    }
}


func genOauthToken() {
    fmt.Println("Will generate file " + tokFile + " for Google Directory Admin API")
    // Load oauth.Config (e.g. Google oauth endpoint, client_id, client_secret)
    oauthConf := getOauthConfig(credFile)
    // Start oauth process on the web to get oauth token 
    err := getTokenFromWeb(oauthConf, tokFile)
    if err != nil {
        log.Fatalf("Unable to create oauth token: %v", err)
    } else {
        fmt.Println(tokFile + " created!")
    }
}


func startMapper() {
    // Load config
    config := getConfig(confFile)
    fmt.Println("CF API Endpoint: " + config.CFApi)
    // Load oauth.Config (e.g. Google oauth endpoint)
    oauthConf := getOauthConfig(credFile)
    // Load existing oauth token (access_key and resfresh_key)
    oauthTok, err := tokenFromFile(tokFile)
    // Create 'Service' so Google Directory (Admin) can be requested
    httpClient := oauthConf.Client(context.Background(), oauthTok)
    googleService, err := admin.New(httpClient)
    if err != nil {
        log.Fatalf("Unable to retrieve directory Client: %v", err)
    }
    // Search for all Google Groups matching the search pattern
    groupsRes, err := googleService.Groups.List().Customer("my_customer").Query("email:snpaas__*").MaxResults(10).Do()
    if err != nil {
        log.Fatalf("Unable to retrieve Google Groups: %v", err)
    }
    if len(groupsRes.Groups) == 0 {
        fmt.Println("No groups found.\n")
    } else {
        for _, gr := range groupsRes.Groups {
            fmt.Printf("GROUP EMAIL: %s\n", gr.Email)

            membersRes, err := googleService.Members.List(gr.Email).MaxResults(10).Do()
            if err != nil {
                log.Fatalf("Unable to retrieve members in group: %v", err)
            }
            if len(membersRes.Members) == 0 {
                fmt.Println("No members found.\n")
            } else {
                fmt.Println("MEMBERS:")
                for _, m := range membersRes.Members {
                    fmt.Printf("%s\n", m.Email)
                }
            }
        } // End for
    } // End else
} // End startMapper