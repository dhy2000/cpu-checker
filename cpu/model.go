package cpu

import "sort"

// CPU Model

type WordType uint

const NumRegisters = 32

type RegisterFile [NumRegisters]WordType

func (r *RegisterFile) Read(index int) WordType {
	if index >= NumRegisters || index < 0 {
		return 0
	}
	if index == 0 {
		return 0
	}
	return r[index]
}

func (r *RegisterFile) Write(index int, data WordType) {
	if index >= NumRegisters || index <= 0 {
		return
	}
	r[index] = data
}

type DataMemory map[WordType]WordType

func (m DataMemory) Read(address WordType) WordType {
	data, exists := m[address]
	if !exists {
		return 0
	}
	return data
}

func (m DataMemory) Write(address WordType, data WordType) {
	if data == 0 {
		delete(m, address)
		return
	}
	m[address] = data
}

type Cpu struct {
	reg RegisterFile
	mem DataMemory
}

func (c *Cpu) Reset() {
	c.reg = RegisterFile{}
	c.mem = make(DataMemory)
}

func NewCpu() *Cpu {
	cpu := &Cpu{}
	cpu.Reset()
	return cpu
}

func CheckCpu(std, ans *Cpu) error {
	// register
	for i := 0; i < NumRegisters; i++ {
		if std.reg[i] != ans.reg[i] {
			return WrongRegState{Number: WordType(i), Got: ans.reg[i], Expected: std.reg[i]}
		}
	}
	// memory
	addresses := make(map[WordType]bool) // union addresses set
	for k := range std.mem {
		addresses[k] = true
	}
	for k := range ans.mem {
		addresses[k] = true
	}
	addressList := make([]WordType, 0) // make addresses in order
	for k := range addresses {
		addressList = append(addressList, k)
	}
	sort.Slice(addressList, func(i, j int) bool {
		return addressList[i] < addressList[j]
	})
	for _, a := range addressList {
		if mStd, mAns := std.mem.Read(a), ans.mem.Read(a); mStd != mAns {
			return WrongMemState{Address: a, Got: mAns, Expected: mStd}
		}
	}
	return nil
}

func IsEqualCpu(this, that *Cpu) bool {
	return CheckCpu(this, that) == nil
}
