package cpu

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

type Trace []Behavior

func ReadTraceFromFile(filename string) (Trace, error) {
	var file *os.File
	if len(filename) == 0 {
		file = os.Stdin
	} else {
		var err error
		file, err = os.OpenFile(filename, os.O_RDONLY, 0444)
		defer func() { _ = file.Close() }()
		if err != nil {
			return nil, err
		}
	}
	buffer := bufio.NewReader(file)
	trace := make(Trace, 0)
	var num uint = 0
	for {
		num++
		line, err := buffer.ReadString('\n')
		line = strings.TrimSpace(line)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		behavior, err := ParseBehavior(line, num)
		if err != nil {
			return nil, err
		}
		if behavior != nil {
			trace = append(trace, behavior)
		}
	}
	return trace, nil
}

func (trace Trace) FormattedOutput(filename string) error {
	var file *os.File
	if len(filename) == 0 {
		return nil
	}
	var err error
	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer func() { _ = file.Close() }()
	if err != nil {
		return err
	}
	for i := 0; i < len(trace); i++ {
		bh := trace[i]
		// skip reg[0] writing
		if bhReg, isReg := bh.(*RegBehavior); isReg && bhReg.RegNumber == 0 {
			continue
		}
		if _, err = file.WriteString(bh.String() + "\n"); err != nil {
			log.Println(err)
		}
	}
	return nil
}
