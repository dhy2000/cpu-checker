package check

import (
	"cpu-checker/cpu"
)

func IsWriteRegZero(behavior cpu.Behavior) bool {
	bReg, ok := behavior.(*cpu.RegBehavior)
	if !ok {
		return false
	}
	return bReg.RegNumber == 0
}

// SimpleCheck 检查 CPU 行为序列的正确性(基础版)
// std: 标准CPU输出(作为参考)
// ans: 待测CPU输出
// byCpu: true 按 CPU 状态检查，false 按行为输出一致性检查
func SimpleCheck(std cpu.Trace, ans cpu.Trace, byCpu bool) (cpu.Trace, error) {
	correct := make(cpu.Trace, 0) // record correct executions
	i, j := 0, 0
	stdCpu, ansCpu := cpu.NewCpu(), cpu.NewCpu()
	stdCpu.Reset()
	ansCpu.Reset()
	var err error = nil
	for i < len(std) && j < len(ans) {
		// skip ignored behavior
		for i < len(std) && IsWriteRegZero(std[i]) {
			i++
		}
		for j < len(ans) && IsWriteRegZero(ans[j]) {
			j++
		}
		if !(i < len(std) && j < len(ans)) {
			break
		}
		// check behavior
		bStd, bAns := std[i], ans[j]
		bStd.Simulate(stdCpu)
		bAns.Simulate(ansCpu)
		if byCpu {
			if err = cpu.CheckCpu(stdCpu, ansCpu); err != nil {
				break
			}
		} else {
			if !cpu.IsEqualBehavior(bStd, bAns) {
				err = cpu.WrongBehavior{Got: bAns, Expected: bStd}
				break
			}
		}
		correct = append(correct, bAns)
		i++
		j++
	}
	// skip ignored behavior
	for i < len(std) && IsWriteRegZero(std[i]) {
		i++
	}
	for j < len(ans) && IsWriteRegZero(ans[j]) {
		j++
	}
	if err == nil && i < len(std) && j >= len(ans) {
		err = cpu.WrongBehavior{Got: nil, Expected: std[i]}
	}
	if err == nil && i >= len(std) && j < len(ans) {
		err = cpu.WrongBehavior{Got: ans[j], Expected: nil}
	}
	return correct, err
}
