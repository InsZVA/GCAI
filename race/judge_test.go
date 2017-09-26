package race

import (
	"testing"
	"log"
	"time"
)

func TestJudge(t *testing.T) {
	AddTask(&Task{
		GameId: 1,
		AI1Path: "C:/Go/path/src/github.com/inszva/GCAI/assert/ai/test/1.exe",
		AI2Path: "C:/Go/path/src/github.com/inszva/GCAI/assert/ai/test/1.exe",
		Callback: func(result int) {
			log.Println("result:", result)
		},
	})
	time.Sleep(11 * time.Second)
}