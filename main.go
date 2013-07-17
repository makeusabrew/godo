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

    godo.LoadTasks()

    if newTask != "" {
        fmt.Println("Creating new task...")

        createTask(newTask)
        return
    }

    if list {
        fmt.Println("listing todos...\n")

        listTasks()
        return
    }
}

func createTask(text string) bool {

    godo.AddTask(text)

    godo.WriteTasks()

    return true
}

func listTasks() {

    tasks := godo.GetTasks()

    for _, task := range(tasks) {
        fmt.Println(task.Order, ") " + task.Text)
    }
}
