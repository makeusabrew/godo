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

type GistFile struct {
    Filename string
    Url string `json:"raw_url"`
    Body []byte
}

type Gist struct {
    Files map[string] GistFile
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

    token, err := getToken(authlist)

    if err != nil {
        return
    }

    if token == "" {
        return errors.New("Could not retrieve token")
    }

    currentUser.authenticate(token)

    return

}

func Authed() bool {
    return currentUser.authed
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

    done := make(chan bool)
    // _ = key, which we don't yet care about
    for _, file := range gist.Files {
        // re-assigning file within the for {} scope
        // just means we can avoid a closure
        file := file
        go readGist(&file, done)
    }

    for _, file := range gist.Files {
        <-done
        // @FIXME: the files we passed by reference were *still*
        // just copies. We can't get the address of a map. Hmm
        fmt.Println(file.Body)
    }
}

func readGist(file *GistFile, done chan<- bool) {
    fmt.Println(file.Url)
    request, _ := http.NewRequest("GET", file.Url, nil)
    body, err := doRequest(request)

    if err == nil {
        file.Body = body
        done <- true
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

func getToken(list []GithubAuthorization) (token string, err error) {
    clientId := os.Getenv("CLIENT_ID")
    if clientId == "" {
        return token, errors.New("No CLIENT_ID environment variable set")
    }

    for _, auth := range list {
        if auth.App.ClientId == clientId {
            token = auth.Token
            return
        }
    }
    return
}
