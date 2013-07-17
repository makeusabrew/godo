package main

import (
    "flag"
    "fmt"
    "github.com/makeusabrew/godo/src/godo"
)

func main() {
    newTask := flag.String("a", "", "add a new todo")
    list := flag.Bool("l", false, "list outstanding tasks")
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

func colour(c int) string {
    return fmt.Sprintf("\x1b[%dm", c)
}
