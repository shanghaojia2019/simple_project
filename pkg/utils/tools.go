package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//IsEmpty 是否为空
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

//HostName ...
func HostName() string {
	name, err := os.Hostname()
	if err != nil {
		return ""
	}
	return name
}

// GetPWD 当前所在的目录位置，编译过后只有一个bin文件
func GetPWD() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	// fmt.Printf("cwd:%s\n", path)
	// _, filename, _, _ := runtime.Caller(0)
	// dir, _ := filepath.Abs(filename)
	return filepath.Dir(path)
}
