package winapi

import (
	"errors"
	"syscall"
	"time"
	"unsafe"
)

// CreateNamedPipe dwOpenMode值
const (
	PIPE_ACCESS_INBOUND  = 0x00000001
	PIPE_ACCESS_OUTBOUND = 0x00000002
	PIPE_ACCESS_DUPLEX   = 0x00000003
)

// GetNamedPipeInfo的NamedPipeEnd flags值
const (
	PIPE_CLIENT_END = 0x00000000
	PIPE_SERVER_END = 0x00000001
)

// CreateNamedPipe dwPipeMode值
const (
	PIPE_WAIT                  = 0x00000000
	PIPE_NOWAIT                = 0x00000001
	PIPE_READMODE_BYTE         = 0x00000000
	PIPE_READMODE_MESSAGE      = 0x00000002
	PIPE_TYPE_BYTE             = 0x00000000
	PIPE_TYPE_MESSAGE          = 0x00000004
	PIPE_ACCEPT_REMOTE_CLIENTS = 0x00000000
	PIPE_REJECT_REMOTE_CLIENTS = 0x00000008
)

// CreateNamedPipe nMaxInstances理论最大命名管道数
const PIPE_UNLIMITED_INSTANCES = 255

/*
	ConnectNamedPipe

原型：
BOOL ConnectNamedPipe(

	[in] HANDLE hNamedPipe,
	[in, out, optional] LPOVERLAPPED lpOverlapped

);
*/
func ConnectNamedPipe(hNamedPipe HANDLE, po *OVERLAPPED) (err error) {
	r1, _, e1 := syscall.SyscallN(procConnectNamedPipe.Addr(), 2, uintptr(hNamedPipe), uintptr(unsafe.Pointer(po)), 0)
	if r1 == 0 {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("ConnectNamedPipe failed")
		}
	}
	return
}

/*
	CreateNamedPipe

原型：
HANDLE CreateNamedPipe(

	[in] LPCWSTR lpName,
	[in] DWORD dwOpenMode,
	[in] DWORD dwPipeMode,
	[in] DWORD nMaxInstances,
	[in] DWORD nOutBufferSize,
	[in] DWORD nInBufferSize,
	[in] DWORD nDefaultTimeOut,
	[in, optional] LPSECURITY_ATTRIBUTES lpSecurityAttributes

);
*/
func CreateNamedPipe(
	name string,
	openMode uint32,
	pipeMode uint32,
	maxInstances uint32,
	outBufferSize uint32,
	inBufferSize uint32,
	defaultTimeOut time.Duration,
	sa *SECURITY_ATTRIBUTES,
) (h HANDLE, err error) {
	pName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return
	}
	dto := uint32(uint64(defaultTimeOut) / 1e6)
	h, err = _CreateNamedPipe(pName, openMode, pipeMode, maxInstances, outBufferSize, inBufferSize, dto, sa)
	return
}

func _CreateNamedPipe(
	pName *uint16,
	dwOpenMode uint32,
	dwPipeMode uint32,
	nMaxInstances uint32,
	nOutBufferSize uint32,
	nInBufferSize uint32,
	nDefaultTimeOut uint32,
	pSecurityAttributes *SECURITY_ATTRIBUTES,
) (h HANDLE, err error) {
	r1, _, e1 := syscall.SyscallN(
		procCreateNamedPipe.Addr(),
		8,
		uintptr(unsafe.Pointer(pName)),
		uintptr(dwOpenMode),
		uintptr(dwPipeMode),
		uintptr(nMaxInstances),
		uintptr(nOutBufferSize),
		uintptr(nInBufferSize),
		uintptr(nDefaultTimeOut),
		uintptr(unsafe.Pointer(pSecurityAttributes)),
		0,
	)
	if h == INVALID_HANDLE_VALUE {
		wec := WindowsErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("CreateNamedPipe failed")
		}
	} else {
		h = HANDLE(r1)
	}
	return
}
