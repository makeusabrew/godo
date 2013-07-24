package godo

import (
    "os"
    "log"
    "os/user"
    "bufio"
    "encoding/json"
)

type Task struct {
    Text string
    Order int
    Done bool
}

func getTaskPath() (path string) {
    user, err := user.Current()
    if err != nil {
        log.Fatal(err)
    }

    path = user.HomeDir + "/.godo/lists/tasks"

    return
}

var tasks []Task

func LoadTasks() ([]Task) {

    file, err := os.Open(getTaskPath())

    if err != nil {
        log.Fatal(err)
    }

    reader := bufio.NewReader(file)

    defer file.Close()

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }

        var task Task

        data := line[:len(line)-1]

        json.Unmarshal([]byte(data), &task)

        tasks = append(tasks, task)

    }

    return tasks
}

func GetTasks() ([]Task) {
    return tasks
}

func AddTask(text string) {
    order := len(tasks) + 1
    tasks = append(tasks, Task{text, order, false})

    writeTasks()
}

func MarkTaskDone(order int) {
    task := findTask(order)

    task.Done = true

    writeTasks()
}

func findTask(order int) (*Task) {
    for i := range tasks {
        // not sure if this is very idiomatic; but we need a reference here
        // rather than a copy which `_, task := range tasks` would give us
        task := &tasks[i]
        if task.Order == order {
            return task
        }
    }

    return nil
}

func writeTasks() {
    file, err := os.Create(getTaskPath())

    writer := bufio.NewWriter(file)

    defer file.Close()

    if err != nil {
        log.Fatal(err)
    }

    for _, task := range tasks {
        line, _ := json.Marshal(task)
        _, err := writer.WriteString(string(line) + "\n")

        if err != nil {
            log.Fatal(err)
        }

    }

    writer.Flush()
}

func (task Task) Status() (string) {
    if task.Done {
        return "✓"
    }

    return "✗"
}
