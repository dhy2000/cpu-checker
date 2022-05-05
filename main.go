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
		fmt.Printf("Wrong answer after correct execution of %v\n", correct[len(correct)-1])
		fmt.Println(err)
	}
}
