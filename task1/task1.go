package main

import "fmt"

func main() {
	fmt.Println("下面是一些问卷调查")
	const thisYear = 2026
	var name string
	fmt.Print("请输入你的名字:")
	fmt.Scanln(&name)
	var birthYear int
	fmt.Print("请输入你的出生年份(比如 2007):")
	fmt.Scanln(&birthYear)
	age := thisYear - birthYear
	fmt.Printf("你好,%s!\n", name)
	fmt.Printf("%d 年出生的你,今年大约 %d 岁。\n", birthYear, age)
	var city string
	fmt.Print("请输入你所在的城市:")
	fmt.Scanln(&city)
	fmt.Printf("你来自%s\n", city)
}
