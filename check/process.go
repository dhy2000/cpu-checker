package check

import (
	"cpu-checker/cpu"
	"reflect"
)

func SeparateRegMem(trace cpu.Trace) (reg cpu.Trace, mem cpu.Trace) {
	reg = make(cpu.Trace, 0)
	mem = make(cpu.Trace, 0)
	for _, t := range trace {
		tReg, isReg := t.(*cpu.RegBehavior)
		if isReg {
			reg = append(reg, tReg)
		}
		tMem, isMem := t.(*cpu.MemBehavior)
		if isMem {
			mem = append(mem, tMem)
		}
	}
	return
}

// 合并指定种类(写寄存器/写内存)的写回原值
// tp: 为 cpu.Behavior 的一种特定实现的实例
func mergeWriteOrigin(trace cpu.Trace, tp cpu.Behavior) cpu.Trace {
	merged := make(cpu.Trace, 0)
	cpuModel := cpu.NewCpu()
	cpuModel.Reset()
	for _, t := range trace {
		modified := t.Simulate(cpuModel)
		if modified {
			merged = append(merged, t)
		} else { // 写回原值了
			// 通过反射判断类型，如果类型和需要合并的 tp 类型一致则合并了，不加入 merged
			if reflect.TypeOf(t) != reflect.TypeOf(tp) {
				merged = append(merged, t)
			}
		}
	}
	return merged
}

func MergeRegWriteOrigin(trace cpu.Trace) cpu.Trace {
	return mergeWriteOrigin(trace, &cpu.RegBehavior{})
}

func MergeMemWriteOrigin(trace cpu.Trace) cpu.Trace {
	return mergeWriteOrigin(trace, &cpu.MemBehavior{})
}
