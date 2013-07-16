package main

import (
    "flag"
    "fmt"
    "github.com/makeusabrew/godo/src/godo"
)

func main() {
    var newTask string
    var list bool
    flag.StringVar(&newTask, "a", "", "add a new todo")
    flag.BoolVar(&list, "l", false, "list outstanding tasks")

    // parse all flags into their respective variables
    flag.Parse()

    if newTask != "" {
        fmt.Println("Creating new task...")

        createTask(newTask)
        return
    }

    if list {
        fmt.Println("listing todos...")

        listTasks()
        return
    }
}

func createTask(text string) bool {
    fmt.Println(text)

    task := godo.Task{"foo", 1}

    fmt.Println(task)

    return true
}

func listTasks() (tasks []godo.Task) {

    tasks = godo.LoadTasks()

    return
}
