// go:build windows
package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

const HGDI_ERROR HGDIOBJ = HGDIOBJ(^uintptr(0))

// 三元栅格操作
const (
	SRCCOPY        uint32 = 0x00CC0020 // dest = source
	SRCPAINT       uint32 = 0x00EE0086 // dest = source OR dest
	SRCAND         uint32 = 0x008800C6 // dest = source AND dest
	SRCINVERT      uint32 = 0x00660046 // dest = source XOR dest
	SRCERASE       uint32 = 0x00440328 // dest = source AND (NOT dest )
	NOTSRCCOPY     uint32 = 0x00330008 // dest = (NOT source)
	NOTSRCERASE    uint32 = 0x001100A6 // dest = (NOT src) AND (NOT dest)
	MERGECOPY      uint32 = 0x00C000CA // dest = (source AND pattern)
	MERGEPAINT     uint32 = 0x00BB0226 // dest = (NOT source) OR dest
	PATCOPY        uint32 = 0x00F00021 // dest = pattern
	PATPAINT       uint32 = 0x00FB0A09 // dest = DPSnoo
	PATINVERT      uint32 = 0x005A0049 // dest = pattern XOR dest
	DSTINVERT      uint32 = 0x00550009 // dest = (NOT dest)
	BLACKNESS      uint32 = 0x00000042 // dest = BLACK
	WHITENESS      uint32 = 0x00FF0062 // dest = WHITE
	NOMIRRORBITMAP uint32 = 0x80000000 // Do not Mirror the bitmap in this call
	CAPTUREBLT     uint32 = 0x40000000 // Include layered windows
)

/*
	BitBlt

原型：
BOOL BitBlt(

	[in] HDC hdc,
	[in] int x,
	[in] int y,
	[in] int cx,
	[in] int cy,
	[in] HDC hdcsrc,
	[in] int x1,
	[in] int y1,
	[in] DWORD rop

);
*/
func BitBlt(
	hdcDest HDC,
	nXDest int32,
	nYDest int32,
	nWidth int32,
	nHeight int32,
	hdcSrc HDC,
	nXSrc int32,
	nYSrc int32,
	Rop uint32,
) error {
	r1, _, e1 := syscall.Syscall9(
		procBitBlt.Addr(),
		9,
		uintptr(hdcDest),
		uintptr(nXDest),
		uintptr(nYDest),
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(hdcSrc),
		uintptr(nXSrc),
		uintptr(nYSrc),
		uintptr(Rop),
	)
	if r1 == 0 {
		if e1 != 0 {
			return error(e1)
		} else {
			return errors.New("BitBlt failed.")
		}
	} else {
		return nil
	}
}

/*
	DeleteObject

原型：
BOOL DeleteObject(

	[in] HGDIOBJ ho

);
*/
func DeleteObject(obj HGDIOBJ) error {
	r1, _, _ := syscall.Syscall(procDeleteObject.Addr(), 1, uintptr(obj), 0, 0)
	if r1 != 0 {
		return nil
	} else {
		return errors.New("The specified handle is not valid or is currently selected into a DC.")
	}
}

/*
BITMAP结构体

	typedef struct tagBITMAP {
		LONG bmType;
		LONG bmWidth;
		LONG bmHeight;
		LONG bmWidthBytes;
		WORD bmPlanes;
		WORD bmBitsPixel;
		LPVOID bmBits;
	} BITMAP, *PBITMAP, *NPBITMAP, *LPBITMAP;
*/
type BITMAP struct {
	BmType       int32
	BmWidth      int32
	BmHeight     int32
	BmWidthBytes int32
	BmPlanes     uint16
	BmBitsPixel  uint16
	BmBits       uintptr
}

/*
	GetObject

原型：
int GetObject(

	[in] HANDLE h,
	[in] int c,
	[out] LPVOID pv

);
*/
func GetObject(gdiObj HGDIOBJ, cbBuffer int32, pv *byte) int32 {
	r1, _, _ := syscall.Syscall(procGetObject.Addr(), 3, uintptr(gdiObj), uintptr(cbBuffer), uintptr(unsafe.Pointer(pv)))
	return int32(r1)
}

/*
	PAINTSTRUCT结构体

来源：winuser.h

	typedef struct tagPAINTSTRUCT {
		HDC hdc;
		BOOL fErase;
		RECT rcPaint;
		BOOL fRestore;
		BOOL fIncUpdate;
		BYTE rgbReserved[32];
	} PAINTSTRUCT, *PPAINTSTRUCT, *NPAINTSTRUCT, *LPPAINTSTRUCT;
*/
type PAINTSTRUCT struct {
	Hdc         HDC
	FErase      int32
	RcPaint     RECT
	FRestore    int32
	FIncUpdate  int32
	RGBReserved [32]byte
}

/*
	BeginPaint

来源：winuser.h
原型：
HDC BeginPaint(

	[in] HWND hWnd,
	[out] LPPAINTSTRUCT lpPaint

);
*/
func BeginPaint(hWnd HWND, ps *PAINTSTRUCT) (hdc HDC, err error) {
	r1, _, _ := syscall.Syscall(procBeginPaint.Addr(), 2, uintptr(hWnd), uintptr(unsafe.Pointer(ps)), 0)
	if r1 == 0 {
		err = errors.New("BeginPaint failed.")
	} else {
		hdc = HDC(r1)
	}
	return
}

/*
	EndPaint

来源：winuser.h
原型：
BOOL EndPaint(

	[in] HWND hWnd,
	[in] const PAINTSTRUCT *lpPaint

);
*/
func EndPaint(hWnd HWND, ps *PAINTSTRUCT) {
	syscall.Syscall(procEndPaint.Addr(), 2, uintptr(hWnd), uintptr(unsafe.Pointer(ps)), 0)
}

/*
	CreateCompatibleDC

原型：
HDC CreateCompatibleDC(

	[in] HDC hdc

);
*/
func CreateCompatibleDC(dc HDC) (HDC, error) {
	r1, _, _ := syscall.Syscall(procCreateCompatibleDC.Addr(), 1, uintptr(dc), 0, 0)
	if r1 == 0 {
		return 0, errors.New("CreateCompatibleDC failed.")
	} else {
		return HDC(r1), nil
	}
}

/*
	SelectObject

原型：
HGDIOBJ SelectObject(

	[in] HDC hdc,
	[in] HGDIOBJ h

);
*/
func SelectObject(hdc HDC, hgdiobj HGDIOBJ) (robj HGDIOBJ, err error) {
	r1, _, _ := syscall.Syscall(procSelectObject.Addr(), 2, uintptr(hdc), uintptr(hgdiobj), 0)
	if r1 == 0 {
		err = errors.New("An error occurs and the selected object is not a region.")
	} else if HGDIOBJ(r1) == HGDI_ERROR {
		err = errors.New("SelectObject failed.")
	} else {
		robj = HGDIOBJ(r1)
	}
	return
}

/*
	DeleteDC

原型：
BOOL DeleteDC(

	[in] HDC hdc

);
*/
func DeleteDC(dc HDC) error {
	r1, _, _ := syscall.Syscall(procDeleteDC.Addr(), 1, uintptr(dc), 0, 0)
	if r1 == 0 {
		return errors.New("DeleteDC failed.")
	} else {
		return nil
	}
}
