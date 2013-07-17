# GODO

Command line TODO management, written in Go

## Install

```go build```

## Add a todo

```./godo "buy bananas"```

## List todos

```./godo```

```
1) [✗] - learn go
2) [✗] - make godo better
3) [✗] - buy bananas
```

## Mark a task as done

```./godo -d n``` - where ```n``` is a number representing the order of the
task as it appears in your todo list
