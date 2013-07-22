package client

import (
    "encoding/json"
    "net/http"
    "fmt"
    "io/ioutil"
)

func ReadRemoteTasks() {
    fmt.Println("syncing")
    res, err := http.Get("http://a/url/sample.json")

    if err != nil {
        fmt.Println("could not sync")
        return
    }

    defer res.Body.Close()

    var data interface{}

    body, err := ioutil.ReadAll(res.Body)

    if err != nil {
        fmt.Println("err", err)
        return
    }

    json.Unmarshal(body, &data)

    fmt.Println("data", data)
}
