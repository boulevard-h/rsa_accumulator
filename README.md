## RSA Accumulator

修改自 [GarryFCR/RSA_ACCUMULATOR](https://github.com/GarryFCR/RSA_ACCUMULATOR)

### 使用

参考 cmd/main.go 中的代码

``` shell
go run cmd/main.go
```

### 原理

基于RSA-2048实现，参数$N,G$采用固定常数。

`GenRepresentatives(set, encodeType)` 为集合 `set` 中的每个元素生成一个代表数。对每个元素，采用 `SHA-256` 哈希，然后利用 `big.Int` 的 `ProbablyPrime` 方法找到不小于该哈希值的下一个质数。

`ProveMembership(g, N, reps)` 为每个代表数计算 membership proof。对于每个代表数$r$，其 proof = $g^{（全体代表数乘积）/ r} \mod N$。

随后将 `proofs[0], reps[0]` 作为参数传入 `AccumulateNew(proofs[0], reps[0], N)` 计算累加器。由于 `proofs[0]` = $g^{（全体代表数乘积）/ r_0} \mod N$，所以 `AccumulateNew(proofs[0], reps[0], N)` = $g^{（全体代表数乘积）/ r_0} * g^{r_0} \mod N$ = $g^{（全体代表数乘积）} \mod N$。

accumulator 包中目前没有实现 Verify 接口，证明可以和 `cmd/main.go` 中一样，对需要验证的元素生成代表数，然后计算累加器，比较新累加器和旧累加器是否相等。