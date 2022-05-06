package global

import (
	"github.com/spf13/pflag"
	"log"
)

var (
	StdFile string
	AnsFile string
)

var Option struct {
	Separate       bool // separate register and memory
	RegWriteOrigin bool // write-back origin value to register
	MemWriteOrigin bool // write-back origin value to memory
	RegWAW         bool // register write-after-write
	MemWAW         bool // memory write-after-write
	ReOrder        bool // allow re-order (separately check per register/memory address)
	CheckByCpu     bool // check by cpu state (by behavior compare default)
	Show           bool // show output after check result
}

func init() {
	pflag.StringVarP(&StdFile, "std", "s", "", "path of std file")
	pflag.StringVarP(&AnsFile, "ans", "a", "", "path to ans file")
	pflag.BoolVar(&Option.Separate, "sep", false, "separate register and memory")
	pflag.BoolVar(&Option.RegWriteOrigin, "reg-origin", false, "allow write-back origin value to register")
	pflag.BoolVar(&Option.MemWriteOrigin, "mem-origin", false, "allow write-back origin value to memory")
	//pflag.BoolVar(&Option.RegWAW, "reg-waw", false, "register write after write")
	//pflag.BoolVar(&Option.MemWAW, "mem-waw", false, "memory write after write")
	//pflag.BoolVar(&Option.ReOrder, "reorder", false, "allow instruction re-order")
	pflag.BoolVar(&Option.CheckByCpu, "cpu", false, "check by cpu state (default by behavior compare)")
	pflag.BoolVar(&Option.Show, "show", false, "show listed output after check result")
}

func ParseArgv() {
	pflag.Parse()
	if len(StdFile) == 0 {
		log.Fatalln("please specify standard file.")
	}
}
