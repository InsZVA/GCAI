package ai

import "os/exec"
import (
	"github.com/satori/go.uuid"
	"io/ioutil"
	"time"
	"bytes"
	"log"
	"os"
)

var assertPath = "C:/Go/path/src/github.com/inszva/GCAI/assert/ai"

func (worker *Worker) work() {
	for task := range worker.tasks {
		fname := uuid.NewV4().String()
		switch task.Language {
		case "go":
			if task.Canceled {
				break
			}

			ioutil.WriteFile(assertPath + "/" + fname + ".go", []byte(task.Source), 0666)
			stdout := bytes.NewBuffer(make([]byte, 0, 1024))
			stderr := bytes.NewBuffer(make([]byte, 0, 1024))
			cmd := exec.Command("C:/Go/bin/go.exe", "build", "-o", assertPath + "/" + fname + ".exe",
				assertPath + "/" + fname + ".go")
			cmd.Stdout = stdout
			cmd.Stderr = stderr
			if err := cmd.Start(); err != nil {
				log.Println("Compiler:", "ERROR:", err)
			}
			log.Println("Compiler:", "PID:", cmd.Process.Pid)

			stop := make(chan struct{})
			go func() {
				cmd.Wait()
				stop <- struct{}{}
			} ()
			select {
			case <- time.After(1 * time.Second):
				cmd.Process.Kill()
				cmd.Process.Release()
			case <- stop:
			}

			if task.Canceled {
				break
			}

			log.Println("Compiler:", "STDOUT:", stdout.String())
			log.Println("Compiler:", "STDERR:", stderr.String())
			if task.Callback == nil {
				break
			}
			if stderr.String() != "" {
				task.Callback(false, stderr.String())
			} else {
				task.Callback(true, assertPath + "/" + fname + ".exe")
			}
		case "cpp":
			if task.Canceled {
				break
			}

			ioutil.WriteFile(assertPath + "/" + fname + ".cpp", []byte(task.Source), 0666)
			stdout := bytes.NewBuffer(make([]byte, 0, 1024))
			stderr := bytes.NewBuffer(make([]byte, 0, 1024))
			cmd := exec.Command("C:/cygwin64/bin/g++.exe", "-o", assertPath + "/" + fname + ".exe",
				assertPath + "/" + fname + ".cpp")
			cmd.Stdout = stdout
			cmd.Stderr = stderr
			if err := cmd.Start(); err != nil {
				log.Println("Compiler:", "ERROR:", err)
			}
			log.Println("Compiler:", "PID:", cmd.Process.Pid)

			stop := make(chan struct{})
			go func() {
				cmd.Wait()
				stop <- struct{}{}
			} ()
			select {
			case <- time.After(1 * time.Second):
				cmd.Process.Kill()
				cmd.Process.Release()
			case <- stop:
			}

			if task.Canceled {
				break
			}

			log.Println("Compiler:", "STDOUT:", stdout.String())
			log.Println("Compiler:", "STDERR:", stderr.String())
			if task.Callback == nil {
				break
			}
			if stderr.String() != "" {
				task.Callback(false, stderr.String())
			} else {
				task.Callback(true, assertPath + "/" + fname + ".exe")
			}
		}

		// Clean up
		os.Remove(assertPath + "/" + fname + "." + task.Language)
	}
}
