package godo

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
    GistId string

    token string
    authed bool
}

func (u *User) authenticate(token string) {
    u.authed = true
    u.token = token
}

var currentUser = &User{}

func Authenticate(username string, password string) (err error) {
    request, _ := http.NewRequest("GET", "https://api.github.com/authorizations", nil)
    request.SetBasicAuth(username, password)

    body, err := doRequest(request)

    if err != nil {
        return
    }

    var authlist []GithubAuthorization

    err = json.Unmarshal(body, &authlist)

    if err != nil {
        return
    }

    token := getToken(authlist)

    if token == "" {
        return errors.New("Could not authenticate")
    }

    currentUser.authenticate(token)

    return

}

func Authed() bool {
    return currentUser.authed
}

type Gist struct {
    Files map[string] struct {
        Filename string
        Url string `json:"raw_url"`
    }
}

func FetchRemoteTasks() {
    fmt.Println("Fetching remote tasks")
    body, err := authedGithubRequest("gists/" + currentUser.GistId)

    if err != nil {
        return
    }

    var gist Gist

    err = json.Unmarshal(body, &gist)

    if err != nil {
        return
    }

    // _ = key
    gists := make(chan string, len(gist.Files))
    for _, file := range gist.Files {
        go readGist(file.Url, gists)
    }

    for i := 0; i < len(gist.Files); i++ {
        fmt.Println(i, <-gists)
    }
}

func readGist(url string, gists chan<- string) {
    request, _ := http.NewRequest("GET", url, nil)
    body, err := doRequest(request)

    if err == nil {
        gists <- string(body)
    }
}

func PushRemoteTasks() {}

func authedGithubRequest(url string) (body []byte, err error) {
    request, _ := http.NewRequest("GET", "https://api.github.com/" + url, nil)
    request.Header.Add("Authorization", "bearer "+currentUser.token)

    return doRequest(request)
}

func doRequest(request *http.Request) (body []byte, err error) {
    client := &http.Client{}

    response, err := client.Do(request)

    defer response.Body.Close()

    if err != nil {
        return
    }

    if response.StatusCode != 200 {
        return body, errors.New("Bad response code " + response.Status)
    }

    body, err = ioutil.ReadAll(response.Body)

    if err != nil {
        return
    }

    return
}

func getToken(list []GithubAuthorization) (token string) {
    clientId := os.Getenv("CLIENT_ID")
    for _, auth := range list {
        if auth.App.ClientId == clientId {
            token = auth.Token
            return
        }
    }
    return
}
