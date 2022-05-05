package check

import (
	"cpu-checker/cpu"
	"cpu-checker/global"
)

func CpuCheck(std cpu.Trace, ans cpu.Trace) (cpu.Trace, error) {
	// 允许写回原值
	if global.Option.MemWriteOrigin {
		std = MergeMemWriteOrigin(std)
		ans = MergeMemWriteOrigin(ans)
	}
	if global.Option.RegWriteOrigin {
		std = MergeRegWriteOrigin(std)
		ans = MergeRegWriteOrigin(ans)
	}
	// 按寄存器和内存分开检查
	if global.Option.Separate {
		stdReg, stdMem := SeparateRegMem(std)
		ansReg, ansMem := SeparateRegMem(ans)
		correctReg, errReg := SimpleCheck(stdReg, ansReg, global.Option.CheckByCpu)
		correctMem, errMem := SimpleCheck(stdMem, ansMem, global.Option.CheckByCpu)
		if errReg != nil {
			return correctReg, errReg
		}
		if errMem != nil {
			return correctMem, errMem
		}
		return append(correctReg, correctMem...), nil
	} else {
		return SimpleCheck(std, ans, global.Option.CheckByCpu)
	}
}
