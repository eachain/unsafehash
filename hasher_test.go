package unsafehash

import "testing"

func TestMakeInt(t *testing.T) {
	crc32 := Make[int](CRC32)
	t.Logf("crc32 123: %v", crc32(123))
}

func TestMakeString(t *testing.T) {
	mapping := Map[string]()
	t.Logf("mapping 123: %v", mapping("123"))
}

type foo struct {
	a int
	b string
}

func TestStruct(t *testing.T) {
	fnv := Make[foo](FNV64)
	f := foo{a: 123, b: "abc"}
	t.Logf("fnv %+v: %v", f, fnv(f))
}

func TestAny(t *testing.T) {
	adler := Make[any](Adler32)
	t.Logf("adler hello: %v", adler("hello"))
	t.Logf("adler 123: %v", adler(123))
	f := foo{a: 123, b: "abc"}
	t.Logf("adler %+v: %v", f, adler(f))
}
