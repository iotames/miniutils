package miniutils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// GetPidByPort. 查找端口所属进程的PID。
// 返回端口号对应的进程PID，若没有找到相关进程，返回-1
func GetPidByPort(portNumber int) int {
	supportOSs := map[string]bool{
		"windows": true,
		"linux":   true,
		"darwin":  false,
	}
	support, ok := supportOSs[runtime.GOOS]
	if !ok || !support {
		panic("GetPidByPort Not Support " + runtime.GOOS)
	}

	res := -1
	var outBytes bytes.Buffer
	var cmdStr string
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// fmt.Printf("---resStr---%s---resStr---", resStr)
		// 以换行符结尾`TCP    127.0.0.1:9222         0.0.0.0:0              LISTENING       9700
		//`
		cmdStr = fmt.Sprintf("netstat -ano -p tcp | findstr %d", portNumber)
		cmd = exec.Command("cmd", "/c", cmdStr)
	}

	if runtime.GOOS == "linux" {
		// processInfo := exec.Command("/bin/sh", "-c",`lsof -i:8299 | awk '{print $2}' | awk  'NR==2{print}'`)
		// 直接返回端口号, 但是以换行符结尾。直接转换字符串为数字会出BUG。使用 strings.TrimSpace 函数转换之。
		cmdStr = fmt.Sprintf(`lsof -i:%d | awk '{print $2}' | awk  'NR==2{print}'`, portNumber)
		cmd = exec.Command("/bin/sh", "-c", cmdStr)
	}

	cmd.Stdout = &outBytes
	cmd.Run()
	resStr := outBytes.String()
	log.Printf("----Executed---For---GetPidByPort---%s----", cmdStr)
	log.Printf("----Result----outString:%s-------", resStr)
	if len(outBytes.Bytes()) == 0 {
		return res
	}

	pidStr := ""
	if runtime.GOOS == "linux" {
		pidStr = strings.TrimSpace(resStr)
	}
	if runtime.GOOS == "windows" {
		r := regexp.MustCompile(`\s\d+\s`).FindAllString(resStr, -1)
		if len(r) == 0 {
			return res
		}
		pidStr = strings.TrimSpace(r[0])
	}
	pid, err := strconv.Atoi(pidStr)
	if err == nil {
		res = pid
	}
	return res
}

// StartBrowserByUrl 打开系统默认浏览器
func StartBrowserByUrl(url string) error {
	var cmdMap = map[string]*exec.Cmd{
		"windows": exec.Command("cmd", "/c", "start", url),
		"darwin":  exec.Command("/bin/bash", "-c", "open", url),
		"linux":   exec.Command("/bin/bash", "-c", "xdg-open", url),
	}
	eCmd, ok := cmdMap[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}
	return eCmd.Start()
}

// KillPid 杀死运行中的进程
func KillPid(pid string) error {
	fmt.Println("Kill---PID:", pid)
	unixCmd := exec.Command("/bin/bash", "-c", "kill "+pid)
	cmdMap := map[string]*exec.Cmd{
		"windows": exec.Command("cmd", "/c", fmt.Sprintf("taskkill -pid %s -F", pid)), // MUST ADD arg -F
		"linux":   unixCmd,
		"darwin":  unixCmd,
	}
	cmd, ok := cmdMap[runtime.GOOS]
	if !ok {
		return fmt.Errorf("not support platform: %s", runtime.GOOS)
	}
	return cmd.Run()
}

// RunCmd 直接执行操作系统中的命令 RunCmd("go", "version")
func RunCmd(name string, arg ...string) ([]byte, error) {
	var bf bytes.Buffer
	cmd := exec.Command(name, arg...)
	// cmd.Dir = "" // 设置在哪个目录执行命令
	// cmd.Stdout = os.Stdout
	// cmd.Env = os.Environ()
	cmd.Stdout = &bf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	result := bf.Bytes()
	// os.Stdout.Write(result)
	return result, err
}
