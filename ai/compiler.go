package ai

import (
	"log"
)
type Task struct {
	Language string
	Source string
	// success msg is exe path
	Callback func(success bool, msg string)
	Canceled bool
}

type Worker struct {
	tasks chan *Task
}

func (worker *Worker) Start() {
	go worker.work()
}

var defaultWorker Worker

func init() {
	log.Println("Compiler is starting...")
	defaultWorker.tasks = make(chan *Task, 128)
	defaultWorker.Start()
}
