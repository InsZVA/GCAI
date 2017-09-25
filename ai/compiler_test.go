package ai

import (
	"testing"
	"time"
)

func TestCompiler(t *testing.T) {
	defaultWorker.tasks <- &Task{
		Language: "go",
		Source: "xxxx",
		Callback: func(success bool, msg string) {
			if success != false {
				t.Error("Except false")
			}
		},
	}
	defaultWorker.tasks <- &Task{
		Language: "go",
		Source: `
		package main

		import "fmt"

		func main() {
			fmt.Println("Hello world!")
		}
		`,
		Callback: func(success bool, msg string) {
			if success != true {
				t.Error("Except true")
			}
		},
	}
	defaultWorker.tasks <- &Task{
		Language: "cpp",
		Source: "xxxx",
		Callback: func(success bool, msg string) {
			if success != false {
				t.Error("Except false")
			}
		},
	}
	defaultWorker.tasks <- &Task{
		Language: "cpp",
		Source: `
		#include <iostream>

		using namespace std;

		int main() {
			cout << "Hello world!" << endl;
		}
		`,
		Callback: func(success bool, msg string) {
			if success != true {
				t.Error("Except true")
			}
		},
	}
	time.Sleep(8 * time.Second)
}
