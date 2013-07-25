package godo

import (
    "os/user"
    "io/ioutil"
    "encoding/json"
    "fmt"
)

func baseDir() (path string) {
    user, err := user.Current()
    if err != nil {
        // just an empty path
        return
    }

    path = user.HomeDir + "/.godo"

    return
}

type config struct {
    GistId string `json:"gistId"`
    Token string `json:"token"`
}

var c config

func LoadConfig() error {
    body, err := ioutil.ReadFile(baseDir() + "/.meta")

    if err != nil {
        return err
    }

    err = json.Unmarshal(body, &c)

    if err != nil {
        fmt.Println(err)
        return err
    }

    // I don't like LoadConfig applying all this stuff; needs
    // splitting out
    if c.Token != "" {
        currentUser.authenticate(c.Token)
    }

    if c.GistId != "" {
        currentUser.GistId = c.GistId
    }

    return nil
}
