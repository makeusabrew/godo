package client

import (
    "encoding/json"
    "net/http"
    "fmt"
    "io/ioutil"
)

type GithubAuthorization struct {
    Id int
    Token string
    App struct {
        ClientId string `json:"client_id"`
        Name string
    }
}

func GetAuthorizations(username string, password string) {
    client := &http.Client{}
    request, _ := http.NewRequest("GET", "https://api.github.com/authorizations", nil)

    request.SetBasicAuth(username, password)

    response, err := client.Do(request)

    defer response.Body.Close()

    if err != nil {
        fmt.Println("err", err)
        return
    }

    body, err := ioutil.ReadAll(response.Body)

    if err != nil {
        fmt.Println("Could not read response body")
        return
    }

    var authlist []GithubAuthorization

    err = json.Unmarshal(body, &authlist)

    if err != nil {
        fmt.Println("JSON unmarshal error", err)
        return
    }

    if response.StatusCode != 200 {
        fmt.Println("Bad response code", response.StatusCode)
        return
    }

    for i, auth := range authlist {
        fmt.Println(i, auth.App.Name)
        // pick out auth.App.Id, check it against our config client ID
        // if we have one, great - use the token and move on
        // else POST /authorizations for our app and use the returned token
    }
}

func Authed() bool {
    return false
}

func FetchRemoteTasks() {}
func PushRemoteTasks() {}
