// go:build windows
package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

/*
	GetLastError

来源：errhandlingapi.h
原型：
_Post_equals_last_error_ DWORD GetLastError();
*/
func GetLastError() WindowsErrorCode {
	ret, _, _ := syscall.SyscallN(procGetLastError.Addr(), 0, 0, 0, 0)
	return WindowsErrorCode(uint32(ret))
}

/*
	ExitProcess

来源：processthreadsapi.h
原型：
void ExitProcess(

	[in] UINT uExitCode

);
*/
func ExitProcess(ExitCode uint32) {
	syscall.SyscallN(procExitProcess.Addr(), 1, uintptr(ExitCode), 0, 0)
}

/*
	GetModuleHandle

来源：libloaderapi.h
原型：
HMODULE GetModuleHandle(

	[in, optional] LPCWSTR lpModuleName

);
*/
func GetModuleHandle(ModuleName string) (h HINSTANCE, err error) {
	var a uintptr
	var pStr *uint16
	if ModuleName == "" {
		a = 0
	} else {
		pStr, err = syscall.UTF16PtrFromString(ModuleName)
		if err != nil {
			return
		} else {
			a = uintptr(unsafe.Pointer(pStr))
		}
	}
	r1, _, e1 := syscall.SyscallN(procGetModuleHandle.Addr(), 1, a, 0, 0)
	if r1 == 0 {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("GetModuleHandle failed")
		}
	} else {
		h = HINSTANCE(r1)
	}
	return
}

/*
	CloseHandle

来源：handleapi.h
原型：
BOOL CloseHandle(

	[in] HANDLE hObject

);
*/
func CloseHandle(h HANDLE) (err error) {
	r1, _, e1 := syscall.SyscallN(procCloseHandle.Addr(), 1, uintptr(h), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("CloseHandle failed")
		}
	}
	return
}
