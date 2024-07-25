package memory

import (
	"errors"
	"sync"
	"unsafe"
)

//  first-fit allocation strategy

/* example
* [1000] -> nil
* 300 (request space)
* [300 (in use (1000 - 300))] -> [700 left (not in use)] -> nil
* 500 (request space)
* [300 (in use)] -> [500 (in use)] -> [200 left (not in use)] -> nil
 */
const MemorySize = 1024 * 1024

type Block struct {
	Start unsafe.Pointer // staring address of memory block , helps to find where each block is located in memory
	Size  int            // size of memory block
	InUse bool           // is in use
	Next  *Block
}

type MemoryManager struct {
	memory []byte
	head   *Block
	mu     sync.Mutex
}

func NewMemoryManager() *MemoryManager {
	mem := make([]byte, MemorySize)
	return &MemoryManager{
		memory: mem,
		head: &Block{
			Start: unsafe.Pointer(&mem[0]),
			Size:  MemorySize,
			InUse: false,
		},
	}
}

func (mm *MemoryManager) Allocate(size int) (unsafe.Pointer, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	block := mm.head

	for block != nil {
		if !block.InUse && block.Size >= size {
			if block.Size > size {
				newBlock := &Block{
					Start: unsafe.Pointer(uintptr(block.Size) + uintptr(size)),
					Size:  block.Size - size, // remaining space
					InUse: false,             // because it is not in use , so it is a free memory (allocating new memory block)
					Next:  block.Next,
				}
				block.Size = size
				block.Next = newBlock
			}
			block.InUse = true
			return block.Start, nil
		}
		block = block.Next
	}
	return nil, errors.New("out of memory")
}

func (mm *MemoryManager) Free(ptr unsafe.Pointer) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	block := mm.head
	for block != nil {
		if block.Start == ptr {
			block.InUse = false
			mm.coalesce()
		}
		block = block.Next
	}
	return nil
}

func (mm *MemoryManager) coalesce() {
	block := mm.head
	for block != nil && block.Next != nil {
		if !block.InUse && !block.Next.InUse {
			block.Size += block.Next.Size
			block.Next = block.Next.Next
		} else {
			block = block.Next
		}
	}
}
