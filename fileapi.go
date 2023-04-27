// go:build windows
package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

// 可读取最大字节数
const _MAX_READ = 1 << 30

// CreateFile flags
const (
	CREATE_NEW        = 1
	CREATE_ALWAYS     = 2
	OPEN_EXISTING     = 3
	OPEN_ALWAYS       = 4
	TRUNCATE_EXISTING = 5
)

/*
	CreateFile

原型：
HANDLE CreateFileA(

	[in] LPCSTR lpFileName,
	[in] DWORD dwDesiredAccess,
	[in] DWORD dwSharedMode,
	[in, optional] LPSECURITY_ATTRIBUTES lpSecurityAttributes,
	[in] DWORD dwCreationDisposition,
	[in] DWORD dwFlagsAndAttributes,
	[in, optional] HANDLE hTemplateFile

);
*/
func CreateFile(
	FileName string,
	DesiredAccess uint32,
	ShareMode uint32,
	SecurityAttributes *SECURITY_ATTRIBUTES,
	CreationDisposition uint32,
	FlagsAndAttributes uint32,
	TemplateFile HANDLE,
) (HANDLE, error) {
	pFileName, err := syscall.UTF16PtrFromString(FileName)
	if err != nil {
		return 0, err
	}
	r1, _, e1 := syscall.SyscallN(
		procCreateFile.Addr(),
		7,
		uintptr(unsafe.Pointer(pFileName)),
		uintptr(DesiredAccess),
		uintptr(ShareMode),
		uintptr(unsafe.Pointer(SecurityAttributes)),
		uintptr(CreationDisposition),
		uintptr(FlagsAndAttributes),
		uintptr(TemplateFile),
		0, 0,
	)
	if r1 == 0 {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			return 0, wec
		} else {
			return 0, errors.New("CreateFile failed")
		}
	} else {
		return HANDLE(r1), nil
	}
}

/*
	ReadFile

原型：
BOOL ReadFile(

	[in] HANDLE hFile,
	[out] LPVOID lpBuffer,
	[in] DWORD nNumberOfBytesToRead,
	[out, optional] LPDWORD lpNumberOfBytesRead,
	[in, out, optional] LPOVERLAPPED lpOverlapped

);
*/
func ReadFile(hFile HANDLE, buf []byte, pOverlapped *OVERLAPPED) (uint32, error) {
	if buf == nil {
		return 0, errors.New("ReadFile: 必须提供有效的缓冲区")
	}
	var Len int = len(buf)
	if Len <= 0 || Len > _MAX_READ {
		return 0, errors.New("ReadFile: 缓冲区长度必须大于0且不超过_MAX_READ")
	}
	var NumberOfBytesRead uint32 = 0
	err := _ReadFile(hFile, &buf[0], uint32(Len), &NumberOfBytesRead, pOverlapped)
	if err != nil {
		return 0, err
	} else {
		return NumberOfBytesRead, nil
	}
}

func _ReadFile(
	hFile HANDLE,
	pBuffer *byte,
	NumberOfBytesToRead uint32,
	pNumberOfBytesRead *uint32,
	pOverlapped *OVERLAPPED,
) error {
	r1, _, e1 := syscall.SyscallN(
		procReadFile.Addr(),
		5,
		uintptr(hFile),
		uintptr(unsafe.Pointer(pBuffer)),
		uintptr(NumberOfBytesToRead),
		uintptr(unsafe.Pointer(pNumberOfBytesRead)),
		uintptr(unsafe.Pointer(&pOverlapped)),
		0,
	)
	if r1 == 0 {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			return wec
		} else {
			return errors.New("winapi: ReadFile failed")
		}
	} else {
		return nil
	}
}

/*
	WriteFile

原型：
BOOL WriteFile(

	[in] HANDLE hFile,
	[in] LPCVOID lpBuffer,
	[in] DWORD nNumberOfBytesToWrite,
	[out, optional] LPDWORD lpNumberOfBytesWritten,
	[in, out, optional] LPOVERLAPPED lpOverlapped

);
*/
func WriteFile(hFile HANDLE, buf []byte, pOverlapped *OVERLAPPED) (uint32, error) {
	if buf == nil {
		return 0, errors.New("WriteFile: 必须提供有效的缓冲区")
	}
	var Len int = len(buf)
	if Len <= 0 || Len > _MAX_READ {
		return 0, errors.New("WriteFile: 缓冲区大小必须大于0且不超过_MAX_READ")
	}
	var NumberOfBytesWritten uint32
	err := _WriteFile(hFile, &buf[0], uint32(Len), &NumberOfBytesWritten, pOverlapped)
	if err != nil {
		return 0, err
	} else {
		return NumberOfBytesWritten, nil
	}
}

func _WriteFile(
	hFile HANDLE,
	pBuffer *byte,
	NumberOfBytesToWrite uint32,
	pNumberOfBytesWritten *uint32,
	pOverlapped *OVERLAPPED,
) error {
	r1, _, e1 := syscall.SyscallN(
		procWriteFile.Addr(),
		5,
		uintptr(hFile),
		uintptr(unsafe.Pointer(pBuffer)),
		uintptr(NumberOfBytesToWrite),
		uintptr(unsafe.Pointer(pNumberOfBytesWritten)),
		uintptr(unsafe.Pointer(pOverlapped)),
		0,
	)
	if r1 == 0 {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			return wec
		} else {
			return errors.New("winapi: WriteFile failed")
		}
	} else {
		return nil
	}
}

// SetFilePointer flags
const (
	FILE_BEGIN   = 0
	FILE_CURRENT = 1
	FILE_END     = 2
)

/*
	SetFilePointer

原型：
DWORD SetFilePointer(

	[in] HANDLE hFile,
	[in] LONG lDistanceToMove,
	[in, out, optional] PLONG lpDistanceToMoveHigh,
	[in] DWORD dwMoveMethod

);
*/
func SetFilePointer(hFile HANDLE, DistanceToMove int64, MoveMethod uint32) (NewPointer int64, err error) {
	var np int64
	r1, _, e1 := syscall.SyscallN(
		procSetFilePointer.Addr(),
		4,
		uintptr(hFile),
		uintptr(DistanceToMove),
		uintptr(unsafe.Pointer(&np)),
		uintptr(MoveMethod),
		0, 0,
	)
	if r1 == 0 {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("winapi: SetFilePointer failed")
		}
	} else {
		NewPointer = np
	}
	return
}
