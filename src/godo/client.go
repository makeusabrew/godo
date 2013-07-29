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

    for key, file := range gist.Files {
        // re-assigning key here just means we don't have to pass it to
        // the closure
        key := key

        go func() {
            request, _ := http.NewRequest("GET", gist.Files[key].Url, nil)
            body, err := doRequest(request)

            if err == nil {
                // can't do a read-modify-write, so have to get the current
                // struct in its entirity and update that then re-assign
                file.Body = body
                gist.Files[key] = file
            }

            done <- true
        }()

    }

    // is there a neater way of waiting for all the functions to return rather
    // than blocking len(files) times? feels a bit clunky...
    for _ = range gist.Files {
        <-done
    }

    for key, file := range gist.Files {
        writeTaskList(key, file.Body)
    }

    fmt.Println(string(len(gist.Files)) + " files synced")
}

func writeTaskList(filename string, data []byte) {
    ioutil.WriteFile(baseDir() + "/lists/" + filename, data, 0600)
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
