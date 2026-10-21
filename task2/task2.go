package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("无聊至极的你进入了彩票店")
	fmt.Println("这家彩票店引进了一种新的刮刮乐，大奖是国庆13天连放")
	fmt.Println("彩票号码从1到99")
	fmt.Println("但你的钱只能够买七张彩票")
	fmt.Println("选一张你想要的号码：")
	var number int
	win := rand.Intn(99) + 1
	for i := 0; i <= 7; i++ {
		if i == 7 {
			fmt.Println("你已经用完了七张彩票，很遗憾你错过了13天假期连放")
			fmt.Println("但是这里有一个神秘咒语，输入y即可重新开始")
			fmt.Println("你想重新开始吗？(y)")
			var word string
			fmt.Scanln(&word)
			if word == "y" {
				i = -1
				fmt.Println("你重新获取了7次机会")
				fmt.Println("现在请选一张你想要的号码：")
				continue
			} else {
				fmt.Println("客人慢走，欢迎下次光临")
				break
			}
		}
		fmt.Scanln(&number)
		if number < 1 || number > 99 {
			fmt.Println("emmmm....咱是不是拿错彩票了？")
			fmt.Println("咱去重新拿一张吧")
			i--
		}
		switch {
		case number == win:
			fmt.Println("运气这么好！")
			fmt.Println("快去迎接你的13天假期吧。")
		case number < win:
			fmt.Println("也许大奖的号码要大一点？")
		default:
			fmt.Println("也许大奖的号码要小一点？")
		}
	}
}
