package client

import (
    "encoding/json"
    "net/http"
    "io/ioutil"
    "fmt"
    "os"
    "errors"
)

type GithubAuthorization struct {
    Id int
    Token string
    App struct {
        ClientId string `json:"client_id"`
        Name string
    }
}

type User struct {
    token string
    authed bool
}

func (u *User) authenticate(token string) {
    u.authed = true
    u.token = token
}

var user = &User{}

func Authenticate(username string, password string) error {
    client := &http.Client{}
    request, _ := http.NewRequest("GET", "https://api.github.com/authorizations", nil)

    request.SetBasicAuth(username, password)

    response, err := client.Do(request)

    defer response.Body.Close()

    if err != nil {
        return err
    }

    body, err := ioutil.ReadAll(response.Body)

    if err != nil {
        return err
    }

    var authlist []GithubAuthorization

    err = json.Unmarshal(body, &authlist)

    if err != nil {
        return err
    }

    if response.StatusCode != 200 {
        return errors.New("Bad response code " + response.Status)
    }

    token := getToken(authlist)

    if token == "" {
        return errors.New("Could not authenticate")
    }

    user.authenticate(token)

    return nil

}

func Authed() bool {
    return user.authed
}

func FetchRemoteTasks() {
    fmt.Println("Fetching remote tasks")
}

func PushRemoteTasks() {}

func getToken(list []GithubAuthorization) (token string) {
    clientId := os.Getenv("CLIENT_ID")
    for _, auth := range list {
        if auth.App.ClientId == clientId {
            token = auth.Token
            return;
        }
    }
    return
}
