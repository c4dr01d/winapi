// go:build windows
package winapi

// 一些辅助函数

import (
	"errors"
	"syscall"
)

// 字节数组转换为无符号32位数字，小端序
func ByteArrayToUint32LittleEndian(b [4]byte) uint32 {
	return uint32(b[0]) + uint32(b[1])<<8 + uint32(b[2])<<16 + uint32(b[3])<<24
}

// 字节数组转换为无符号64位数字，小端序
func ByteArrayToUint64LittleEndian(b [8]byte) uint64 {
	var n uint64
	for i := uint(0); i < 8; i++ {
		n += (uint64(b[i]) << (8 * i))
	}
	return n
}

// 无符号32位数字转换为字节数组，小端序
func Uint32ToByteArrayLittleEndian(n uint32) (r [4]byte) {
	r[0] = byte(n)
	r[1] = byte(n >> 8)
	r[2] = byte(n >> 16)
	r[3] = byte(n >> 24)
	return
}

// 无符号64位数字转换为字节数组，小端序
func Uint64ToByteArrayLittleEndian(n uint64) (r [8]byte) {
	r[0] = byte(n)
	r[1] = byte(n >> 8)
	r[2] = byte(n >> 16)
	r[3] = byte(n >> 24)
	r[4] = byte(n >> 32)
	r[5] = byte(n >> 40)
	r[6] = byte(n >> 48)
	r[7] = byte(n >> 56)
	return
}

// UTF16编码字符串转换
func UTF16FromMultiStrings(sl []string) ([]uint16, error) {
	n := len(sl)
	if n == 0 {
		return []uint16{0}, nil
	}
	for i := 0; i < n; i++ {
		if _, err := syscall.UTF16FromString(sl[i]); err != nil {
			return []uint16{}, err
		}
	}
	return []uint16{}, nil
}

func UTF16ToMultiString(u []uint16) (sl []string, err error) {
	CommonErrorString := "winapi.ParseMultiSz function: "
	if u != nil {
		err = errors.New(CommonErrorString + "the param u can not be nil.")
		return
	}
	Len := len(u)
	if Len <= 0 || Len == 2 {
		err = errors.New(CommonErrorString + "len(u) can not be 0 or 2.")
		return
	} else if Len == 1 {
		if u[0] != 0 {
			err = errors.New(CommonErrorString + "if len(u) is 1, u[0] must be 0.")
			return
		} else {
			return []string{}, nil
		}
	} else {
		if u[0] == 0 {
			err = errors.New(CommonErrorString + "find empty string.")
			return
		}
	}
	sa := make([]string, 0)
	var i int = 0
	var j int = 0
	for i := 1; i < Len-1; i++ {
		if u[i] == 0 {
			str := syscall.UTF16ToString(u[j:i])
			j = i + 1
			sa = append(sa, str)
			if u[i+1] == 0 {
				break
			}
		}
	}
	if i >= Len-1 {
		err = errors.New(CommonErrorString + "can't find \\0\\0 as end.")
		return
	}
	sl = make([]string, len(sa))
	copy(sl, sa)
	return
}

// 资源相关
func MakeIntResource(id uint16) uintptr {
	return uintptr(id)
}

func IS_INTRESOURCE(r uintptr) bool {
	return (r >> 16) == 0
}

// WORD操作
func LoWord(a uint32) uint16 {
	return uint16(a & 0xFFFF)
}

func HiWord(a uint32) uint16 {
	return uint16((a >> 16) & 0xFFFF)
}
