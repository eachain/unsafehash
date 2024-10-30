# unsafehash

unsafehash基于unsafe.Pointer计算comparable变量的uint64 hash值。

## 示例

```go
package main

import (
	"fmt"
	"hash/crc32"

	"github.com/eachain/unsafehash"
)

type foo struct {
	a int
	b string
}

func main() {
	var hello [5]byte
	copy(hello[:], "hello")
	fmt.Println(crc32.ChecksumIEEE(hello[:]))
	fmt.Println(unsafehash.Make[string](unsafehash.CRC32)(string(hello[:])))
	fmt.Println(unsafehash.Make[[5]byte](unsafehash.CRC32)(hello))
	// Output:
	// 907060870
	// 907060870
	// 907060870

	index := unsafehash.Map[any]()
	m := make(map[uint64]any)
	m[index(123)] = 123
	m[index("hello")] = "hello"
	f := foo{a: 123, b: "abc"}
	m[index(f)] = f

	fmt.Printf("123: %v\n", m[index(123)])
	fmt.Printf("hello: %v\n", m[index("hello")])
	fmt.Printf("foo: %+v\n", m[index(f)])
	// Output:
	// 123: 123
	// hello: hello
	// foo: {a:123 b:abc}
}
```
