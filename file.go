package utils

import (
	"io/ioutil"
	"os"
)

func DownloadFile(link, filePath string) (err error) {
	http := HttpClient{}
	data, err := http.Get(link)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filePath, data, 0666)
	return
}

func SafeCreateDir(path string) {
	if exists := PathExists(path); exists {
		return
	}
	os.Mkdir(path, os.ModePerm)
}

// 判断所给路径文件/文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 读取文件
func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}
