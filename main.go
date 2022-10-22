package main

import (
	"cpu-checker/check"
	"cpu-checker/cpu"
	"cpu-checker/global"
	"fmt"
	"log"
)

func main() {
	global.ParseArgv()
	traceStd, err := cpu.ReadTraceFromFile(global.StdFile)
	if err != nil {
		log.Fatalln(err)
	}
	traceAns, err := cpu.ReadTraceFromFile(global.AnsFile)
	if err != nil {
		log.Fatalln(err)
	}
	correct, err := check.CpuCheck(traceStd, traceAns)
	if err == nil {
		fmt.Println("Your answer is correct.")
	} else {
		fmt.Println(err)
		if len(correct) > 0 {
			lastCorrect := correct[len(correct)-1]
			fmt.Printf("Wrong answer after correct execution of [%d] %v\n", lastCorrect.LineNumber(), lastCorrect)
		}
	}
	// formatted output
	if err = traceStd.FormattedOutput(global.StdFormat); err != nil {
		log.Println(err)
	}
	if err = traceAns.FormattedOutput(global.AnsFormat); err != nil {
		log.Println(err)
	}
}
