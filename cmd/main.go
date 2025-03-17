package main

import (
	"fmt"

	"github.com/boulevard-h/rsa_accumulator/accumulator"
)

func main() {
	// 1. 初始化：生成 RSA accumulator 的参数
	setup := accumulator.TrustedSetup()
	fmt.Println("累加器参数 N =", setup.N.String())

	// 2. 定义初始元素集合（例如使用姓名作为输入）
	elements := []string{"Alice", "Bob", "Charlie"}
	// 选择编码方式，这里使用 HashToPrimeFromSha256（其他方式已删除）
	var encodeType accumulator.EncodeType = accumulator.HashToPrimeFromSha256

	// 3. 生成代表数和证明，同时计算累加器：
	acc, proofs := accumulator.AccAndProve(elements, encodeType, setup)
	fmt.Println("初始累加器值：", acc.String())

	// 4. 验证每个元素的证明：
	representatives := accumulator.GenRepresentatives(elements, encodeType)
	for i, elem := range elements {
		computed := accumulator.AccumulateNew(proofs[i], representatives[i], setup.N)
		valid := computed.Cmp(acc) == 0
		fmt.Printf("元素 %s 的证明验证结果: %v\n", elem, valid)
	}

	// 5. 添加一个新元素 "David"
	newElement := "David"
	newRep := accumulator.HashToPrime([]byte(newElement))
	updatedAcc := accumulator.AccumulateNew(acc, newRep, setup.N)
	fmt.Println("添加新元素后的累加器值：", updatedAcc.String())

	// 6. 重新生成累加器和证明以验证更新正确性
	newElements := append(elements, newElement)
	updatedAcc2, updatedProofs := accumulator.AccAndProve(newElements, encodeType, setup)
	fmt.Println("通过重新生成得到的累加器值：", updatedAcc2.String())
	if updatedAcc.Cmp(updatedAcc2) == 0 {
		fmt.Println("累加器更新一致，新元素的添加及证明验证成功！")
	} else {
		fmt.Println("累加器更新不一致，新元素的添加或证明验证失败！")
	}

	// 7. 对新加入的元素单独验证证明
	newReps := accumulator.GenRepresentatives(newElements, encodeType)
	lastIdx := len(newElements) - 1
	computedNew := accumulator.AccumulateNew(updatedProofs[lastIdx], newReps[lastIdx], setup.N)
	if computedNew.Cmp(updatedAcc2) == 0 {
		fmt.Printf("新元素 %s 的证明验证成功！\n", newElement)
	} else {
		fmt.Printf("新元素 %s 的证明验证失败！\n", newElement)
	}
}
