package main

import (
    "flag"
    "fmt"
    "github.com/makeusabrew/godo/src/godo"
    "code.google.com/p/gopass"
)

func main() {
    newTask := flag.String("a", "", "add a new todo")
    list := flag.Bool("l", false, "list outstanding tasks")
    sync := flag.Bool("s", false, "sync tasks")
    doneTask := flag.Int("d", -1, "mark a task as done")

    // parse all flags into their respective variables
    flag.Parse()

    godo.LoadTasks()

    if *newTask != "" {
        addTask(*newTask)
        return
    }

    if *list {
        listTasks()
        return
    }

    if *sync {
        syncTasks()
        return
    }

    if *doneTask != -1 {
        godo.MarkTaskDone(*doneTask)
        return
    }

    // if we got here, start looking at some shortcut commands

    switch args := flag.Args(); len(args) {
    case 0:
        listTasks()
        return
    case 1:
        addTask(args[0])
        return
    }
}

func addTask(newTask string) {
    fmt.Println("Creating new task '"+newTask+"'")
    godo.AddTask(newTask)
}

func listTasks() {

    tasks := godo.GetTasks()

    fmt.Println()
    for _, task := range(tasks) {
        fmt.Printf("%d) [%s] - %s\n", task.Order, task.Status(), task.Text)
    }
    fmt.Println()
}

func syncTasks() {
    if !godo.Authed() {
        username := getInput("Please enter your username: ")
        password, err := gopass.GetPass("Please enter your password: ")

        // @TODO how *do* we properly handle errors?
        if err != nil {
            return
        }

        err = godo.Authenticate(username, password); if err != nil {
            fmt.Println("Could not authenticate!")
            return
        }
    }

    godo.FetchRemoteTasks()
    godo.PushRemoteTasks()
}

func colour(c int) string {
    return fmt.Sprintf("\x1b[%dm", c)
}

func getInput(prompt string) (input string) {
    fmt.Print(prompt)
    fmt.Scanln(&input)
    return
}
