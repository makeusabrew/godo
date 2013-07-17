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

    doneTask := flag.Int("d", -1, "mark a task as done")

    // parse all flags into their respective variables
    flag.Parse()

    godo.LoadTasks()

    if newTask != "" {
        fmt.Println("Creating new task '"+newTask+"'")
        godo.AddTask(newTask)
        return
    }

    if list {
        listTasks()
        return
    }

    if *doneTask != -1 {
        godo.MarkTaskDone(*doneTask)
        return
    }
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
