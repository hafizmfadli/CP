package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	MaxInt32 = 1 << 31 - 1
)

func main(){
	cin := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(cin, &t)
	for t > 0 {
		var n int
		fmt.Fscan(cin, &n)
		a := make([]int, n)
		min := MaxInt32 
		for i := 0; i < n; i++ {
			fmt.Fscan(cin, &a[i])
			if a[i] < min {
				min = a[i]
			}
		}
		ans := 0
		for _, v := range a {
			ans += v - min
		}
		fmt.Println(ans)
		 t -= 1
	}
}