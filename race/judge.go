package race

import (
	"os/exec"
	"time"
	"bytes"
	"encoding/binary"
	"log"
	"errors"
	"fmt"
)

type Task struct {
	GameId int
	AI1Path string
	AI2Path string
	// 0-draw 1-ai1win 2-ai2win 3-error
	Callback func(result int)
}

type Judge struct {
	tasks chan *Task
}

func (judge *Judge) Work() {
	for task := range judge.tasks {
		gameInfo, err := GetGameInfo(task.GameId)
		if err != nil {
			task.Callback(3)
			continue
		}

		stop := make(chan struct{}, 16)
		cmdJudge := exec.Command(gameInfo.JudgePath)
		stdinJudge, err := cmdJudge.StdinPipe()
		if err != nil {
			task.Callback(3)
			continue
		}
		stdoutJudge, err := cmdJudge.StdoutPipe()
		if err != nil {
			task.Callback(3)
			continue
		}
		stderrJudge := bytes.NewBuffer(make([]byte, 0, 16*1024))
		cmdJudge.Stderr = stderrJudge
		cmdJudge.Start()

		// TODO: 解释性语言的处理
		cmdAI1 := exec.Command(task.AI1Path)
		cmdAI2 := exec.Command(task.AI2Path)

		var result int
		var finish bool
		go func () {
			cmdJudge.Wait()
			stop <- struct{}{}
		} ()

		go func () {
			for !finish {
				var outLen int16
				if binary.Read(stdoutJudge, binary.LittleEndian, &outLen) != nil {
					result = 3
					stop <- struct{}{}
					return
				}
				ai1in := make([]byte, outLen)
				n, e := stdoutJudge.Read(ai1in)
				if e != nil || int16(n) != outLen {
					result = 3
					stop <- struct{}{}
					return
				}
				log.Println("AI1 IN:", string(ai1in))

				if binary.Read(stdoutJudge, binary.LittleEndian, &outLen) != nil {
					result = 3
					stop <- struct{}{}
					return
				}
				ai2in := make([]byte, outLen)
				n, e = stdoutJudge.Read(ai2in)
				if e != nil || int16(n) != outLen {
					result = 3
					stop <- struct{}{}
					return
				}
				log.Println("AI2 IN:", string(ai2in))

				cmdAI1.Stdin = bytes.NewBuffer(ai1in)
				stdoutAI1 := bytes.NewBuffer(make([]byte, 0, 16*1024))
				cmdAI1.Stdout = stdoutAI1
				cmdAI1.Start()

				cmdAI2.Stdin = bytes.NewBuffer(ai2in)
				stdoutAI2 := bytes.NewBuffer(make([]byte, 0, 16*1024))
				cmdAI2.Stdout = stdoutAI2
				cmdAI2.Start()

				var ai1done, ai2done bool
				go func() {
					cmdAI1.Wait()
					ai1done = true
				}()
				go func() {
					cmdAI2.Wait()
					ai2done = true
				}()

				nodeStartTime := time.Now()
				for nodeStartTime.Add(time.Duration(gameInfo.TimeLimit) * time.Millisecond).After(time.Now()) {
					// TODO: Check RAM Usage
					time.Sleep(10 * time.Millisecond)
					if ai1done && ai2done {
						break
					}
				}
				cmdAI1.Process.Kill()
				cmdAI2.Process.Kill()

				ai1out := stdoutAI1.Bytes()
				n, e = stdoutAI1.Read(ai1out)
				binary.Write(stdinJudge, binary.LittleEndian, int16(len(ai1out)))
				stdinJudge.Write(ai1out)
				log.Println("AI1 OUT:", string(ai1out))

				ai2out := stdoutAI2.Bytes()
				n, e = stdoutAI2.Read(ai2out)
				binary.Write(stdinJudge, binary.LittleEndian, int16(len(ai2out)))
				stdinJudge.Write(ai2out)
				log.Println("AI2 OUT:", string(ai2out))

				cmdAI1 = exec.Command(task.AI1Path)
				cmdAI2 = exec.Command(task.AI2Path)
			}
		} ()

		select {
		case <-time.After(100 * time.Second):
			cmdJudge.Process.Kill()
		case <-stop:
		}
		finish = true
		if stderrJudge.Len() != 0 {
			fmt.Fscanf(stderrJudge, "%d\n", &result)
		}
		if task.Callback != nil {
			task.Callback(result)
		}
	}
}

func (judge *Judge) Start() {
	go judge.Work()
}

var defaultJudge Judge

func init() {
	log.Println("Judge is starting...")
	defaultJudge.tasks = make(chan *Task, 1024)
	defaultJudge.Start()
}

func AddTask(task *Task) error {
	select {
	case defaultJudge.tasks <- task:
		return nil
	default:
		return errors.New("task queue is full")
	}
}