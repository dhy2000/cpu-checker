package cpu

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 检查的 CPU 行为: 写寄存器和写内存

type Behavior interface {
	fmt.Stringer
	Simulate(c *Cpu) bool // return: CPU
	LineNumber() uint
}

type RegBehavior struct {
	Line           uint     // 行号
	Timestamp      uint     // 时间戳 (0 为不存在)
	ProgramCounter WordType // 指令计数器
	RegNumber      int      // 寄存器编号
	Data           WordType // 写入的值
}

type MemBehavior struct {
	Line           uint     // 行号
	Timestamp      uint     // 时间戳 (0 为不存在)
	ProgramCounter WordType // 指令计数器
	Address        WordType // 内存地址
	Data           WordType // 写入的值
}

func (r RegBehavior) String() string {
	return fmt.Sprintf("@%08x: $%2d <= %08x", r.ProgramCounter, r.RegNumber, r.Data)
}

func (r RegBehavior) Simulate(c *Cpu) bool {
	if r.RegNumber == 0 {
		return false
	}
	if c.reg.Read(r.RegNumber) == r.Data {
		return false
	}
	c.reg.Write(r.RegNumber, r.Data)
	return true
}

func (r RegBehavior) LineNumber() uint {
	return r.Line
}

func (m MemBehavior) String() string {
	return fmt.Sprintf("@%08x: *%08x <= %08x", m.ProgramCounter, m.Address, m.Data)
}

func (m MemBehavior) Simulate(c *Cpu) bool {
	if c.mem.Read(m.Address) == m.Data {
		return false
	}
	c.mem.Write(m.Address, m.Data)
	return true
}

func (m MemBehavior) LineNumber() uint {
	return m.Line
}

func IsEqualBehavior(this, that Behavior) bool {
	reg1, isReg1 := this.(*RegBehavior)
	reg2, isReg2 := that.(*RegBehavior)
	mem1, isMem1 := this.(*MemBehavior)
	mem2, isMem2 := that.(*MemBehavior)
	if isReg1 && isReg2 {
		return reg1.ProgramCounter == reg2.ProgramCounter && reg1.RegNumber == reg2.RegNumber && reg1.Data == reg2.Data
	}
	if isMem1 && isMem2 {
		return mem1.ProgramCounter == mem2.ProgramCounter && mem1.Address == mem2.Address && mem1.Data == mem2.Data
	}
	return false
}

var (
	// RegPattern 100@00003000: $ 9 <= 12345678
	RegPattern = regexp.MustCompile(`(?P<time>\d*)\s*@\s*(?P<pc>[\dA-Fa-fXxZz]+)\s*:\s*\$\s*(?P<reg>[\dXxZz]+)\s*<=\s*(?P<data>[\dA-Fa-fXxZz]+)`)
	MemPattern = regexp.MustCompile(`(?P<time>\d*)\s*@\s*(?P<pc>[\dA-Fa-fXxZz]+)\s*:\s*\*\s*(?P<address>[\dA-Fa-fXxZz]+)\s*<=\s*(?P<data>[\dA-Fa-fXxZz]+)`)
)

func detectTriState(value string, line string) error {
	if strings.ContainsAny(value, "XzZz") {
		return TriStateError{Value: value, Output: line}
	}
	return nil
}

func parseHex(s string) (uint, error) {
	i64, err := strconv.ParseUint(s, 16, 32)
	return uint(i64), err
}

func parseDec(s string) (uint, error) {
	i64, err := strconv.ParseUint(s, 10, 32)
	return uint(i64), err
}

func ParseRegBehavior(s string, line uint) (*RegBehavior, error) {
	match := RegPattern.FindStringSubmatch(s)
	if len(match) < 5 {
		return nil, nil
	}
	sTimestamp, sPc, sReg, sData := match[1], match[2], match[3], match[4]
	// timestamp
	timestamp := uint(0)
	if len(sTimestamp) > 0 {
		timestamp, _ = parseDec(sTimestamp)
	}
	for _, ss := range []string{sPc, sReg, sData} {
		if err := detectTriState(ss, fmt.Sprintf("@%s: $%s <= %s", sPc, sReg, sData)); err != nil {
			return nil, err
		}
	}
	var pc, reg, data uint
	var err error
	if pc, err = parseHex(sPc); err != nil {
		return nil, err
	}
	if reg, err = parseDec(sReg); err != nil {
		return nil, err
	}
	if data, err = parseHex(sData); err != nil {
		return nil, err
	}
	behavior := &RegBehavior{
		Line:           line,
		Timestamp:      timestamp,
		ProgramCounter: WordType(pc),
		RegNumber:      int(reg),
		Data:           WordType(data),
	}
	if pc%4 != 0 {
		return behavior, UnalignedProgramCounter{ProgramCounter: WordType(pc)}
	}
	if reg < 0 || reg >= NumRegisters {
		return behavior, RegNumberError{Num: WordType(reg)}
	}
	return behavior, nil
}

func ParseMemBehavior(s string, line uint) (*MemBehavior, error) {
	match := MemPattern.FindStringSubmatch(s)
	if len(match) < 5 {
		return nil, nil
	}
	sTimestamp, sPc, sAddress, sData := match[1], match[2], match[3], match[4]
	// timestamp
	timestamp := uint(0)
	if len(sTimestamp) > 0 {
		timestamp, _ = parseDec(sTimestamp)
	}
	for _, ss := range []string{sPc, sAddress, sData} {
		if err := detectTriState(ss, fmt.Sprintf("@%s: *%s <= %s", sPc, sAddress, sData)); err != nil {
			return nil, err
		}
	}
	var pc, address, data uint
	var err error
	if pc, err = parseHex(sPc); err != nil {
		return nil, err
	}
	if address, err = parseHex(sAddress); err != nil {
		return nil, err
	}
	if data, err = parseHex(sData); err != nil {
		return nil, err
	}
	behavior := &MemBehavior{
		Line:           line,
		Timestamp:      timestamp,
		ProgramCounter: WordType(pc),
		Address:        WordType(address),
		Data:           WordType(data),
	}
	if pc%4 != 0 {
		return behavior, UnalignedProgramCounter{ProgramCounter: WordType(pc)}
	}
	if address%4 != 0 {
		return behavior, UnalignedMemoryAddress{Address: WordType(address)}
	}
	return behavior, nil
}

func ParseBehavior(s string, line uint) (Behavior, error) {
	if behavior, err := ParseRegBehavior(s, line); behavior != nil || err != nil {
		return behavior, err
	}
	if behavior, err := ParseMemBehavior(s, line); behavior != nil || err != nil {
		return behavior, err
	}
	return nil, nil
}
