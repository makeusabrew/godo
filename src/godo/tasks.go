package godo

import (
    "os"
    "log"
    "os/user"
    "fmt"
)

type Task struct {
    Text string
    Order int
}

func LoadTasks() (tasks []Task) {
    user, err := user.Current()
    if err != nil {
        log.Fatal(err)
    }

    file, err := os.Open(user.HomeDir + "/.godo/tasks")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(file)

    return
}
