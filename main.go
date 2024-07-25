package main

import (
	"fmt"
	"main/file"
	"main/memory"
)

func main() {
	fs := file.NewFileSystem()

	fs.CreateFile("/hello.txt", []byte("Hello, World!"))
	fs.CreateFile("/animals.txt", []byte("lion\ntiger\nmouse\ncow\npeacock\n"))

	fmt.Printf("%+v \n", fs.Root)
	content, _ := fs.ReadFile("/hello.txt")
	birdsContent, _ := fs.ReadFile("/animals.txt")

	fmt.Println(string(content))
	fmt.Println(string(birdsContent))

	mem := memory.NewMemoryManager()
	addr1, _ := mem.Allocate(1024)
	addr2 , _ := mem.Allocate(1024)

	fmt.Println("allocated memory address",addr1)
	fmt.Println("allocated memory address",addr2)
	mem.Free(addr1)
	mem.Free(addr2)
}
