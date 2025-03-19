package main

import (
	"fmt"
	"time"

	"github.com/boulevard-h/rsa_accumulator/accumulator"
)

func main() {
	// 1. 初始化：生成 RSA accumulator 的参数
	setup := accumulator.TrustedSetup()
	fmt.Println("Param N:\n", setup.N.String())

	// 2. 定义初始元素集合（例如使用姓名作为输入）
	elements := []string{"Alice", "Bob", "Charlie"}
	// 选择编码方式，这里使用 HashToPrimeFromSha256（其他方式已删除）
	var encodeType accumulator.EncodeType = accumulator.HashToPrimeFromSha256

	// 3. 生成代表数和证明，同时计算累加器：
	acc, proofs := accumulator.AccAndProve(elements, encodeType, setup)
	fmt.Println("Initial accumulator value:\n", acc.String())

	// 4. 验证每个元素的证明：
	representatives := accumulator.GenRepresentatives(elements, encodeType)
	for i, elem := range elements {
		computed := accumulator.AccumulateNew(proofs[i], representatives[i], setup.N)
		valid := computed.Cmp(acc) == 0
		fmt.Printf("Element %s proof verification result: %v\n", elem, valid)
	}

	// 5. 添加一个新元素 "David"
	newElement := "David"
	newRep := accumulator.HashToPrime([]byte(newElement))
	updatedAcc := accumulator.AccumulateNew(acc, newRep, setup.N)
	fmt.Println("Updated accumulator value:\n", updatedAcc.String())

	// 6. 重新生成累加器和证明以验证更新正确性
	newElements := append(elements, newElement)
	updatedAcc2, updatedProofs := accumulator.AccAndProve(newElements, encodeType, setup)
	fmt.Println("Updated accumulator value:\n", updatedAcc2.String())
	if updatedAcc.Cmp(updatedAcc2) == 0 {
		fmt.Println("Accumulator update consistency, new element addition and proof verification successful!")
	} else {
		fmt.Println("Accumulator update inconsistency, new element addition or proof verification failed!")
	}

	// 7. 对新加入的元素单独验证证明
	newReps := accumulator.GenRepresentatives(newElements, encodeType)
	lastIdx := len(newElements) - 1
	computedNew := accumulator.AccumulateNew(updatedProofs[lastIdx], newReps[lastIdx], setup.N)
	if computedNew.Cmp(updatedAcc2) == 0 {
		fmt.Printf("New element %s proof verification successful!\n", newElement)
	} else {
		fmt.Printf("New element %s proof verification failed!\n", newElement)
	}

	// 8. 测试放入1000个元素（不生成证明）
	elements = make([]string, 1000)
	for i := 0; i < 1000; i++ {
		elements[i] = fmt.Sprintf("Element %d", i)
	}
	fmt.Println("Testing 1000 elements without proof...")
	startingTime := time.Now().UTC()
	acc, _ = accumulator.AccWithoutProve(elements, encodeType, setup)
	endingTime := time.Now().UTC()
	duration := endingTime.Sub(startingTime)
	fmt.Printf("Running AccWithoutProve Takes [%.3f] Seconds \n", duration.Seconds())
	fmt.Println("Initial accumulator value:\n", acc.String())

	// 9. 测试放入100个元素（生成证明）
	elements = make([]string, 100)
	for i := 0; i < 100; i++ {
		elements[i] = fmt.Sprintf("Element %d", i)
	}
	fmt.Println("Testing 100 elements with proof...")
	startingTime = time.Now().UTC()
	acc, _ = accumulator.AccAndProve(elements, encodeType, setup)
	endingTime = time.Now().UTC()
	duration = endingTime.Sub(startingTime)
	fmt.Printf("Running AccAndProve Takes [%.3f] Seconds \n", duration.Seconds())
	fmt.Println("Initial accumulator value:\n", acc.String())

}
