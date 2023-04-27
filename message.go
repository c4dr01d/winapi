// go:build windows
package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	WM_NULL    uint32 = 0x0000
	WM_CREATE  uint32 = 0x0001
	WM_DESTROY uint32 = 0x0002
	WM_MOVE    uint32 = 0x0003
	WM_SIZE    uint32 = 0x0005

	WM_ACTIVATE uint32 = 0x0006

	WM_PAINT uint32 = 0x000F

	// WM_ACTIVATE状态值
	WA_INACTIVE    = 0
	WA_ACTIVE      = 1
	WA_CLICKACTIVE = 2

	WM_CLOSE uint32 = 0x0010
	WM_QUIT  uint32 = 0x0012

	WM_GETMINMAXINFO uint32 = 0x0024

	WM_COMMAND uint32 = 0x0111
)

/*
MSG结构体

	typedef struct tagMSG {
		HWND hwnd;
		UINT message;
		WPARAM wParam;
		LPARAM lParam;
		DWORD time;
		POINT pt;
		DWORD lPrivate;
	} MSG, *PMSG, *NPMSG, *LPMSG;
*/
type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

/*
	GetMessage

来源：winuser.h
原型：
BOOL GetMessage(

	[out] LPMSG lpMsg,
	[in, optional] HWND hWnd,
	[in] UINT wMsgFilterMin,
	[in] UINT wMsgFilterMax

);
*/
func GetMessage(pMsg *MSG, hWnd HWND, wMsgFilterMin uint32, wMsgFilterMax uint32) int32 {
	r1, _, _ := syscall.SyscallN(procGetMessage.Addr(), 4, uintptr(unsafe.Pointer(pMsg)), uintptr(hWnd), uintptr(wMsgFilterMin), uintptr(wMsgFilterMax), 0, 0)
	return int32(r1)
}

/*
	TranslateMessage

来源：winuser.h
原型：
BOOL TranslateMessage(

	[in] const MSG *lpMsg

);
*/
func TranslateMessage(pMsg *MSG) error {
	r1, _, _ := syscall.SyscallN(procTranslateMessage.Addr(), 1, uintptr(unsafe.Pointer(pMsg)), 0, 0)
	if r1 == 0 {
		return errors.New("winapi: TranslateMessage failed")
	} else {
		return nil
	}
}

/*
	DispatchMessage

来源：winuser.h
原型：
LRESULT DispatchMessage(

	[in] const MSG *lpMsg

);
*/
func DispatchMessage(pMsg *MSG) uintptr {
	r1, _, _ := syscall.SyscallN(procDispatchMessage.Addr(), 1, uintptr(unsafe.Pointer(pMsg)), 0, 0)
	return r1
}

/*
	PostQuitMessage

来源：winuser.h
原型：
void PostQuitMessage(

	[in] int nExitCode

);
*/
func PostQuitMessage(ExitCode int32) {
	syscall.SyscallN(procPostQuitMessage.Addr(), 1, uintptr(ExitCode), 0, 0)
}

/*
	RegisterWindowMessage

来源：winuser.h
原型：
UINT RegisterWindowMessage(

	[in] LPCSTR lpString

);
*/
func RegisterWindowMessage(str string) (message uint32, err error) {
	p, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		return
	}
	r1, _, e1 := syscall.SyscallN(procRegisterWindowMessage.Addr(), 1, uintptr(unsafe.Pointer(p)), 0, 0)
	if r1 == 0 {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("winapi: RegisterWindowMessage failed")
		}
	} else {
		message = uint32(r1)
	}
	return
}
