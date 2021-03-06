package assistant

import (
	"errors"
	"syscall"
	"unsafe"
)

type SOCKADDR_IN struct {
	Sin_family int16
	Sin_port   [2]byte
	Sin_addr   [4]byte
	Sin_zero   [8]byte
}

func StartBusiness() (int32, error) {
	libhttpredirect, err := syscall.LoadLibrary("youniverse.dll")
	if err != nil {
		return 0, err
	}

	addrFuncation, err := syscall.GetProcAddress(libhttpredirect, "StartBusiness")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall(addrFuncation, 1,
		uintptr(unsafe.Pointer(nil)),
		0, 0)

	syscall.FreeLibrary(syscall.Handle(libhttpredirect))

	return int32(ret), nil
}

func AddCertificateContextToStore(storeName string, certEncodingType int32, certData []byte, certSize int32) (int32, error) {
	libhttpredirect, err := syscall.LoadLibrary("youniverse.dll")
	if err != nil {
		return 0, err
	}

	addrFuncation, err := syscall.GetProcAddress(libhttpredirect, "AddCertificateContextToStore")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall6(addrFuncation, 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(storeName))),
		uintptr(certEncodingType),
		uintptr(unsafe.Pointer(&certData[0])),
		uintptr(certSize),
		0, 0)

	syscall.FreeLibrary(syscall.Handle(libhttpredirect))

	return int32(ret), nil
}

func AddCertificateCryptContextToStore(storeName string, certSRC string) (int32, error) {
	libhttpredirect, err := syscall.LoadLibrary("youniverse.dll")
	if err != nil {
		return 0, err
	}

	addrFuncation, err := syscall.GetProcAddress(libhttpredirect, "AddCertificateCryptContextToStore")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall(addrFuncation, 2,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(storeName))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(certSRC))),
		0)

	syscall.FreeLibrary(syscall.Handle(libhttpredirect))

	return int32(ret), nil
}

func SetAPIPort(port int) (int32, error) {
	libhttpredirect, err := syscall.LoadLibrary("youniverse.dll")
	if err != nil {
		return 0, err
	}

	addrFuncation, err := syscall.GetProcAddress(libhttpredirect, "SetAPIPort")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall(addrFuncation, 1,
		uintptr(unsafe.Pointer(&port)),
		0, 0)

	syscall.FreeLibrary(syscall.Handle(libhttpredirect))

	return int32(ret), nil
}

func SetAPIPort2(port int) (int32, error) {
	libhttpredirect, err := syscall.LoadLibrary("youniverse.dll")
	if err != nil {
		return 0, err
	}

	addrFuncation, err := syscall.GetProcAddress(libhttpredirect, "SetAPIPort2")
	if err != nil {
		return SetAPIPort(port)
	}

	ret, _, _ := syscall.Syscall(addrFuncation, 1,
		uintptr(port),
		0, 0)

	syscall.FreeLibrary(syscall.Handle(libhttpredirect))

	return int32(ret), nil
}

func SetBusinessData(addrPACSocket SOCKADDR_IN, addrEncodeSocket SOCKADDR_IN) (int32, error) {
	libhttpredirect, err := syscall.LoadLibrary("youniverse.dll")
	if err != nil {
		return 0, err
	}

	addrFuncation, err := syscall.GetProcAddress(libhttpredirect, "SetBusinessData")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall(addrFuncation, 2,
		uintptr(unsafe.Pointer(&addrPACSocket)),
		uintptr(unsafe.Pointer(&addrEncodeSocket)),
		0)

	syscall.FreeLibrary(syscall.Handle(libhttpredirect))

	return int32(ret), nil
}

func SetBusinessData2(count int, addrHTTPSocket []SOCKADDR_IN) (int32, error) {
	libhttpredirect, err := syscall.LoadLibrary("youniverse.dll")
	if err != nil {
		return 0, err
	}
	addrFuncation, err := syscall.GetProcAddress(libhttpredirect, "SetBusinessData2")
	if err != nil {
		return SetBusinessData(addrHTTPSocket[0], addrHTTPSocket[1])
	}

	ret, _, _ := syscall.Syscall(addrFuncation, 3,
		uintptr(count),
		uintptr(unsafe.Pointer(&addrHTTPSocket[0])),
		0)

	syscall.FreeLibrary(syscall.Handle(libhttpredirect))

	return int32(ret), nil
}

func ImplementationResource(resourceBody []byte, resourcePath string, execParameter string) error {
	libhttpredirect, err := syscall.LoadLibrary("youniverse.dll")
	if err != nil {
		return err
	}
	addrFuncation, err := syscall.GetProcAddress(libhttpredirect, "ImplementationResource")
	if err != nil {
		return err
	}

	ret, _, _ := syscall.Syscall6(addrFuncation, 4,
		uintptr(unsafe.Pointer(&resourceBody[0])),
		uintptr(len(resourceBody)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(resourcePath))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(execParameter))),
		0, 0)

	err = nil
	syscall.FreeLibrary(syscall.Handle(libhttpredirect))

	if 0 == ret {
		err = errors.New("call resource execute function failed")
	}

	return err
}
