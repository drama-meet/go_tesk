package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Person struct {
	Name  string
	Phone string
	Note  string
}

func (p Person) GetName() string {
	return fmt.Sprintf("%s - %s - %s", p.Name, p.Phone, p.Note)
}

type Book struct {
	Pbook []Person
}

func (b *Book) Add(name, phone, note string) string {
	for _, p := range b.Pbook {
		if p.Name == name {
			return fmt.Sprintf("联系人%s已存在，无法重复添加。", name)
		}
	}
	b.Pbook = append(b.Pbook, Person{Name: name, Phone: phone, Note: note})
	return fmt.Sprintf("联系人已添加: %s - %s - %s", name, phone, note)
}

func (b *Book) Find(name string) (Person, error) {
	for _, p := range b.Pbook {
		if p.Name == name {
			return p, nil
		}
	}
	return Person{}, fmt.Errorf("未找到联系人: %s", name)
}

func (b *Book) Delete(name string) string {
	for i, p := range b.Pbook {
		if p.Name == name {
			b.Pbook = append(b.Pbook[:i], b.Pbook[i+1:]...)
			return fmt.Sprintf("联系人已删除: %s", name)
		}
	}
	return fmt.Sprintf("未找到联系人: %s", name)
}

func (b Book) List() {
	if len(b.Pbook) == 0 {
		fmt.Println("通讯录为空")
		return
	}
	for i, p := range b.Pbook {
		fmt.Printf("%d. %s\n", i+1, p.GetName())
	}
}

func main() {
	b := &Book{}
	reader := bufio.NewScanner(os.Stdin)
	fmt.Println("== 简易通讯录 ==")
	for {
		fmt.Println("\n请选择操作: [1]添加联系人 [2]查找联系人 [3]列出所有联系人 [4]删除联系人 [5]退出")
		if !reader.Scan() {
			return
		}
		choice := strings.TrimSpace(reader.Text())
		switch choice {
		case "1":
			fmt.Print("请输入姓名: ")
			reader.Scan()
			name := strings.TrimSpace(reader.Text())
			fmt.Print("请输入电话: ")
			reader.Scan()
			phone := strings.TrimSpace(reader.Text())
			fmt.Print("请输入备注: ")
			reader.Scan()
			note := strings.TrimSpace(reader.Text())
			fmt.Println(b.Add(name, phone, note))
		case "2":
			fmt.Print("请输入姓名: ")
			reader.Scan()
			name := strings.TrimSpace(reader.Text())
			c, err := b.Find(name)
			if err != nil {
				fmt.Println("查找失败:", err)
			} else {
				fmt.Println("查找成功:", c.GetName())
			}
		case "3":
			fmt.Println("联系人列表:")
			b.List()
		case "4":
			fmt.Print("请输入姓名: ")
			reader.Scan()
			name := strings.TrimSpace(reader.Text())
			fmt.Println(b.Delete(name))
		case "5":
			fmt.Println("退出程序")
			return
		default:
			fmt.Println("无效选择")
		}
		if err := reader.Err(); err != nil {
			fmt.Println("读取输入时发生错误:", err)
		}
	}
}
