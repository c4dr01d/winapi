// go:build windows
package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	MF_INSERT = 0x00000000
	MF_CHANGE = 0x00000080
	MF_APPEND = 0x00000100
	MF_DELETE = 0x00000200
	MF_REMOVE = 0x00001000

	MF_BYCOMMAND  = 0x00000000
	MF_BYPOSITION = 0x00000400

	MF_SEPARATOR = 0x00000800

	MF_ENABLED  = 0x00000000
	MF_GRAYED   = 0x00000001
	MF_DISABLED = 0x00000002

	MF_UNCHECKED       = 0x00000000
	MF_CHECKED         = 0x00000008
	MF_USECHECKBITMAPS = 0x00000200

	MF_STRING    = 0x00000000
	MF_BITMAP    = 0x00000004
	MF_OWNERDRAW = 0x00000100

	MF_POPUP        = 0x00000010
	MF_MENUBARBREAK = 0x00000020
	MF_MENUBREAK    = 0x00000040

	MF_UNHILITE = 0x00000000
	MF_HILITE   = 0x00000080

	MF_DEFAULT = 0x00001000

	MF_SYSMENU      = 0x00002000
	MF_HELP         = 0x00004000
	MF_RIGHTJUSTIFY = 0x00004000

	MF_MOUSESELECT = 0x00008000
	MF_END         = 0x00000080 // 只适用于旧版资源文件
)

/*
	AppendMenu

来源：winuser.h
原型：
BOOL AppendMenu(

	[in] HMENU hMenu,
	[in] UINT uFlags,
	[in] UINT_PTR uIDNewItem,
	[in, optional] LPCWSTR lpNewItem

);
*/
func AppendMenu(hMenu HMENU, Flags uint32, IdNewItem uintptr, NewItem string) error {
	var err error
	var pNewItem *uint16
	if NewItem != "" {
		pNewItem, err = syscall.UTF16PtrFromString(NewItem)
		if err != nil {
			return err
		}
	}
	return _AppendMenu(hMenu, Flags, IdNewItem, pNewItem)
}

func _AppendMenu(hMenu HMENU, Flags uint32, IdNewItem uintptr, NewItem *uint16) (err error) {
	r1, _, e1 := syscall.SyscallN(procAppendMenu.Addr(), 4, uintptr(hMenu), uintptr(Flags), IdNewItem, uintptr(unsafe.Pointer(NewItem)), 0, 0)
	if r1 == 0 {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("winapi: AppendMenu failed")
		}
	}
	return
}

/*
	CreateMenu

来源：winuser.h
原型：
HMENU CreateMenu();
*/
func CreateMenu() (hMenu HMENU, err error) {
	r1, _, e1 := syscall.SyscallN(procCreateMenu.Addr(), 0, 0, 0, 0)
	if r1 == 0 {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("winapi: CreateMenu failed")
		}
	} else {
		hMenu = HMENU(r1)
	}
	return
}

/*
	CreatePopupMenu

来源：winuser.h
原型：
HMENU CreatePopupMenu();
*/
func CreatePopupMenu() (hMenu HMENU, err error) {
	r1, _, e1 := syscall.SyscallN(procCreatePopupMenu.Addr(), 0, 0, 0, 0)
	if r1 == 0 {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("winapi: CreatePopupMenu failed")
		}
	} else {
		hMenu = HMENU(r1)
	}
	return
}

/*
	DestroyMenu

来源：winuser.h
原型：
BOOL DestroyMenu(

	[in] HMENU hMenu

);
*/
func DestroyMenu(hMenu HMENU) (err error) {
	r1, _, e1 := syscall.SyscallN(procDestroyMenu.Addr(), 1, uintptr(hMenu), 0, 0)
	if r1 == 0 {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("winapi: DestroyMenu failed")
		}
	}
	return
}
