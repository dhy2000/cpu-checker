package global

import (
	"github.com/spf13/pflag"
	"log"
)

var (
	StdFile   string
	AnsFile   string
	StdFormat string
	AnsFormat string
)

var Option struct {
	Separate       bool // separate register and memory
	RegWriteOrigin bool // write-back origin value to register
	MemWriteOrigin bool // write-back origin value to memory
	CheckByCpu     bool // check by cpu state (by behavior compare default)
	RegWAW         bool // register write-after-write = todo
	MemWAW         bool // memory write-after-write   = todo
	ReOrder        bool // allow re-order (separately check per register/memory address) = todo
}

func init() {
	pflag.StringVarP(&StdFile, "std", "s", "", "path to std file")
	pflag.StringVarP(&AnsFile, "ans", "a", "", "path to ans file")
	pflag.StringVar(&StdFormat, "std-format", "", "output path to formatted std file")
	pflag.StringVar(&AnsFormat, "ans-format", "", "output path to formatted ans file")
	pflag.BoolVar(&Option.Separate, "sep", false, "separate register and memory")
	pflag.BoolVar(&Option.RegWriteOrigin, "reg-origin", false, "allow write-back origin value to register")
	pflag.BoolVar(&Option.MemWriteOrigin, "mem-origin", false, "allow write-back origin value to memory")
	pflag.BoolVar(&Option.CheckByCpu, "cpu", false, "check by cpu state (default by behavior compare)")
}

func ParseArgv() {
	pflag.Parse()
	if len(StdFile) == 0 {
		log.Fatalln("please specify standard file.")
	}
}
