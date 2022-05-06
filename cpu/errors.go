package cpu

import "fmt"

type RegNumberError struct {
	Num WordType // 错误的寄存器号
}

func (r RegNumberError) Error() string {
	return fmt.Sprintf("invalid register number: %v", r.Num)
}

type TriStateError struct {
	Value  string // 出现三态的字面值
	Output string // 出错的原始输出
}

func (t TriStateError) Error() string {
	return fmt.Sprintf("unexpected tri-state literal: %v in %v", t.Value, t.Output)
}

type UnalignedProgramCounter struct {
	ProgramCounter WordType
}

func (u UnalignedProgramCounter) Error() string {
	return fmt.Sprintf("unaligned program counter: %08x", u.ProgramCounter)
}

type UnalignedMemoryAddress struct {
	Address WordType // 未字对齐的地址
}

func (u UnalignedMemoryAddress) Error() string {
	return fmt.Sprintf("unaligned memory address: %08x", u.Address)
}

type WrongBehavior struct {
	Got      Behavior
	Expected Behavior
}

func (w WrongBehavior) Error() string {
	switch {
	case w.Got != nil && w.Expected != nil:
		return fmt.Sprintf("Wrong behavior, Got: %v, Expected %v", w.Got, w.Expected)
	case w.Got == nil && w.Expected != nil:
		return fmt.Sprintf("Output less than expected, Expected %v", w.Expected)
	case w.Got != nil && w.Expected == nil:
		return fmt.Sprintf("Output more than expected, Got:%v", w.Got)
	default:
		return ""
	}
}

type WrongRegState struct {
	Number   WordType
	Got      WordType
	Expected WordType
}

func (w WrongRegState) Error() string {
	return fmt.Sprintf("Wrong register value grf[%d], Got: %08x, Expected %08x", w.Number, w.Got, w.Expected)
}

type WrongMemState struct {
	Address  WordType
	Got      WordType
	Expected WordType
}

func (w WrongMemState) Error() string {
	return fmt.Sprintf("Wrong memory value at %08x, Got: %08x, Expected %08x", w.Address, w.Got, w.Expected)
}
