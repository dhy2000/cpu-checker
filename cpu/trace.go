package cpu

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Trace []Behavior

func ReadTraceFromFile(filename string) (Trace, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	buffer := bufio.NewReader(file)
	trace := make(Trace, 0)
	for {
		line, err := buffer.ReadString('\n')
		line = strings.TrimSpace(line)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		behavior, err := ParseBehavior(line)
		if err != nil {
			return nil, err
		}
		if behavior != nil {
			trace = append(trace, behavior)
		}
	}
	return trace, nil
}
