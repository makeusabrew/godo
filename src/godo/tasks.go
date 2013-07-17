package godo

import (
    "os"
    "log"
    "os/user"
    "bufio"
)

type Task struct {
    Text string
    Order int
}

func getTaskPath() (path string) {
    user, err := user.Current()
    if err != nil {
        log.Fatal(err)
    }

    path = user.HomeDir + "/.godo/tasks"

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

        AddTask(line[:len(line)-1])

    }

    return tasks
}

func GetTasks() ([]Task) {
    return tasks
}


func AddTask(text string) {
    order := len(tasks) + 1
    tasks = append(tasks, Task{text, order})
}

func WriteTasks() {
    file, err := os.Create(getTaskPath())

    writer := bufio.NewWriter(file)

    defer file.Close()

    if err != nil {
        log.Fatal(err)
    }

    for _, task := range(tasks) {
        print(task.Text)
        _, err := writer.WriteString(task.Text + "\n")

        if err != nil {
            log.Fatal(err)
        }

    }

    writer.Flush()
}
