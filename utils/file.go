package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	defer f.Close()
	return ioutil.ReadAll(f)
}

type FileList struct {
	FilePath string
	FileDir  string
}

func prepareDirs(dirs []string) ([]string, [][]os.FileInfo) {
	resultDir := make([]string, 0)
	resultDirFileInfo := make([][]os.FileInfo, 0)
	for _, dir := range dirs {
		if fi, err := os.Stat(dir); err != nil {
			if !os.IsNotExist(err) {
				continue
			}
			if err = os.MkdirAll(dir, 0700); err != nil {
				continue
			}
		} else if !fi.IsDir() {
			continue
		}
		if fis, err := ioutil.ReadDir(dir); err != nil {
		} else {
			resultDir = append(resultDir, dir)
			resultDirFileInfo = append(resultDirFileInfo, fis)
		}
	}
	return resultDir, resultDirFileInfo
}

func GetFileList(fileDir, suffix string) []*FileList {
	arrFileLists := make([]*FileList, 0)
	suffixUp := strings.ToUpper(suffix)
	arrDirs, arrInfos := prepareDirs([]string{fileDir})
	for idx, dbDir := range arrDirs {
		for _, fi := range arrInfos[idx] {
			fileName := fi.Name()
			// try match suffix and `ordinal_pubKey_bitLength.suffix`
			if !strings.HasSuffix(strings.ToUpper(fileName), suffixUp) {
				continue
			}
			filePath := filepath.Join(dbDir, fileName)
			arrFileLists = append(arrFileLists, &FileList{
				FilePath: filePath,
				FileDir:  dbDir,
			})
		}
	}
	return arrFileLists
}
