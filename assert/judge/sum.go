package main

import (
	"fmt"
	"bytes"
	"os"
	"encoding/binary"
)

type record struct {
	ai1point, ai2point, ai1expect, ai2expect int
}

func main() {
	records := []record{}
	buff := bytes.NewBuffer(make([]byte, 0, 128 * 1024))
	rec := record{}
	for round := 1; ;round++ {
		buff.Reset()
		fmt.Fprintf(buff, "%d\n", round)
		for i := 0; i < round - 1; i++ {
			fmt.Fprintf(buff, "%d %d %d %d\n", records[i].ai1point, records[i].ai1expect, records[i].ai2point, records[i].ai2expect)
		}
		binary.Write(os.Stdout, binary.LittleEndian, int16(buff.Len()))
		buff.WriteTo(os.Stdout)
		buff.Reset()
		fmt.Fprintln(buff, round)
		for i := 0; i < round - 1; i++ {
			fmt.Fprintf(buff, "%d %d %d %d\n",  records[i].ai2point, records[i].ai2expect, records[i].ai1point, records[i].ai1expect)
		}
		binary.Write(os.Stdout, binary.LittleEndian, int16(buff.Len()))
		buff.WriteTo(os.Stdout)

		var outLen int16
		binary.Read(os.Stdin, binary.LittleEndian, &outLen)
		buffer := make([]byte, outLen)
		os.Stdin.Read(buffer)
		fmt.Fscanf(bytes.NewReader(buffer), "%d %d\n", &rec.ai1point, &rec.ai1expect)

		binary.Read(os.Stdin, binary.LittleEndian, &outLen)
		buffer = make([]byte, outLen)
		os.Stdin.Read(buffer)
		fmt.Fscanf(bytes.NewReader(buffer), "%d %d\n", &rec.ai2point, &rec.ai2expect)

		sum := rec.ai1point + rec.ai2point
		if sum == rec.ai1expect && sum == rec.ai2expect {
			fmt.Fprintf(os.Stderr, "%d\n", 0)
		} else if sum == rec.ai1expect {
			fmt.Fprintf(os.Stderr, "%d\n", 1)
		} else if sum == rec.ai2expect {
			fmt.Fprintf(os.Stderr, "%d\n", 2)
		} else {
			records = append(records, rec)
			continue
		}
		return
	}
}