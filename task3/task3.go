package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("请输入一段英文（中间用空格空开）：")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	input = strings.ToLower(input)
	//words := strings.Fields(input)
	words := skip(input)
	fmt.Println(words)
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}
	fmt.Println("一共", len(words), "个单词，其中不同的单词共", len(counts), "个。")
	var max string
	for word, n := range counts {
		fmt.Printf("%s : %d次\n", word, n)
		if n > counts[max] {
			max = word
		}
	}
	fmt.Printf("其中出现次数最多的单词是: %s (出现了 %d 次)\n", max, counts[max])
}

func skip(a string) []string {
	var word string
	var s []string
	for _, str := range a {
		if str == ' ' || str == '\n' || str == '\t' {
			if word == "" {
				continue
			}
			s = append(s, word)
			word = ""
			continue
		} else {
			word = word + string(str)
		}
	}
	if word != "" {
		s = append(s, word)
	}
	return s
}
