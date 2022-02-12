# 算法Tips

**求中位数**

求中位数为了防止操作最大整型：
```golang
mid := (max-min)>>2 + min
```

**求奇偶个数**

在一堆偶数个数字里面找唯一的一个出现奇数次数的数:
```golang
// 全部异或到一起
func FindOddTime(all []int) int {
	var sumxor int
	for _, n := range all {
		sumxor = sumxor ^ n
	}
	return sumxor
}
```

取一个不等于0的数的最右测的 1 的数：
```golang
// 取反+1再做并集
func FindRightMostOne(number int)int{
	return number & (^number + 1)
}
```

不用额外变量的元素交换：
```golang
// 通过异或自己等于0的特性
func Swap(a, b int) {
	a = a ^ b
	b = a ^ b
	a = a ^ b
}
```