package unsafehash

import (
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"hash/maphash"
	"reflect"
	"unsafe"
)

type BytesHasher interface {
	HashBytes([]byte) uint64
}

type HashBytesFunc func([]byte) uint64

func (f HashBytesFunc) HashBytes(b []byte) uint64 {
	return f(b)
}

type HashFunc[T comparable] func(T) uint64

var (
	CRC32     HashBytesFunc
	CRC64ISO  HashBytesFunc
	CRC64ECMA HashBytesFunc
	Adler32   HashBytesFunc
	FNV64     HashBytesFunc
	FNV64a    HashBytesFunc
)

func Make[T comparable](hasher BytesHasher) HashFunc[T] {
	return hashOfType[T](reflect.TypeOf((*T)(nil)).Elem(), hasher)
}

func Map[T comparable]() HashFunc[T] {
	seed := maphash.MakeSeed()
	return Make[T](HashBytesFunc(func(b []byte) uint64 {
		return maphash.Bytes(seed, b)
	}))
}

func init() {
	CRC32 = HashBytesFunc(func(b []byte) uint64 {
		return uint64(crc32.ChecksumIEEE(b))
	})

	iso := crc64.MakeTable(crc64.ISO)
	CRC64ISO = HashBytesFunc(func(b []byte) uint64 {
		return crc64.Checksum(b, iso)
	})

	ecma := crc64.MakeTable(crc64.ECMA)
	CRC64ECMA = HashBytesFunc(func(b []byte) uint64 {
		return crc64.Checksum(b, ecma)
	})

	Adler32 = HashBytesFunc(func(b []byte) uint64 {
		return uint64(adler32.Checksum(b))
	})

	FNV64 = HashBytesFunc(func(b []byte) uint64 {
		h := fnv.New64()
		h.Write(b)
		return h.Sum64()
	})

	FNV64a = HashBytesFunc(func(b []byte) uint64 {
		h := fnv.New64a()
		h.Write(b)
		return h.Sum64()
	})
}

func hashOfType[T comparable](typ reflect.Type, hasher BytesHasher) HashFunc[T] {
	switch typ.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.Pointer, reflect.Uintptr, reflect.UnsafePointer,
		reflect.Array,
		reflect.Struct:
		n := int(typ.Size())
		return func(t T) uint64 {
			return hasher.HashBytes(unsafe.Slice((*byte)(unsafe.Pointer(&t)), n))
		}

	case reflect.String:
		return func(t T) uint64 {
			s := *(*string)(unsafe.Pointer(&t))
			return hasher.HashBytes(unsafe.Slice(unsafe.StringData(s), len(s)))
		}

	case reflect.Interface:
		return hashOfIface[T](hasher)

	default:
		panic(fmt.Errorf("unsafehash: type %v is not comparable", typ))
	}
}

type iface struct {
	_    unsafe.Pointer
	data unsafe.Pointer
}

func hashOfIface[T comparable](hasher BytesHasher) HashFunc[T] {
	return func(t T) uint64 {
		typ := reflect.TypeOf(t)
		if typ == nil {
			return hasher.HashBytes(nil)
		}
		switch typ.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.Complex64, reflect.Complex128,
			reflect.Pointer, reflect.Uintptr, reflect.UnsafePointer,
			reflect.Array,
			reflect.Struct:
			n := int(typ.Size())
			return hasher.HashBytes(unsafe.Slice((*byte)((*iface)(unsafe.Pointer(&t)).data), n))

		case reflect.String:
			s := *(*string)((*iface)(unsafe.Pointer(&t)).data)
			return hasher.HashBytes(unsafe.Slice(unsafe.StringData(s), len(s)))

		default:
			panic(fmt.Errorf("unsafehash: type %v is not comparable", typ))
		}
	}
}
