package godo

import (
    "os"
    "log"
    "io/ioutil"
    "bufio"
    "encoding/json"
    "regexp"
    "strings"
)

type Task struct {
    Text string
    Order int
    Done bool
}

func getTaskPath() (path string) {
    return baseDir() + "/lists/" + c.List
}

var tasks []Task

func LoadTasks() ([]Task) {

    body, err := ioutil.ReadFile(getTaskPath())

    if err != nil {
        log.Fatal(err)
    }

    lines := strings.Split(string(body), "\n")

    regexp := regexp.MustCompile(`-\s*\[(x|\s)\]\s?(.+)`)

    for i, line := range lines {

        matches := regexp.FindStringSubmatch(line)

        if len(matches) == 0 {
            continue
        }

        task := Task{}
        task.Text = matches[2]
        task.Order = i
        if matches[1] == "x" {
            task.Done = true
        }

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
