package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type PkgLockDependency struct {
	//package-lock.json每个元素对应结构
	Version   string            `json:"version"`
	Resolved  string            `json:"resolved"`
	Integrity string            `json:"integrity"`
	Dev       bool              `json:"dev"`
	Requires  map[string]string `json:"requires"`
}

type PkgLock struct {
	//package-lock.json文件对应结构
	Dependencies map[string]PkgLockDependency
}

type Pkg struct {
	Dependencies       map[string]string `json:"dependencies"`
	BundleDependencies []string          `json:"bundleDependencies"`
	DevDependencies    map[string]string `json:"devDependencies"`
}

//解析json文件
func praseJSON(path string) {
	//打开json文件
	var err error
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	//获取内容
	var content []byte
	content, err = io.ReadAll(io.LimitReader(f, 16*1024*1024))
	if err != nil {
		fmt.Println("open file error: " + err.Error())
		return
	}
	//如果是package-lock.json，用PkgLock解析
	if strings.HasSuffix(path, "package-lock.json") {
		//println("----")
		var pkglock PkgLock
		err = json.Unmarshal([]byte(content), &pkglock)
		if err != nil {
			fmt.Println("package-lock prase ERROR: ", err.Error())
			return
		}
		fmt.Println(pkglock)

	} else if strings.HasSuffix(path, "yarn.lock") {
		//对yarn.lock解析
		lockfile, err := ParseLockFileData(content)
		if err != nil {
			fmt.Println("yarn.lock prase ERROR: ", err.Error())
			return
		}
		fmt.Println(len(lockfile))
	} else {
		//如果是其他.json文件，用Pkg解析
		var pkg Pkg
		err = json.Unmarshal([]byte(content), &pkg)
		if err != nil {
			fmt.Println("other json file prase ERROR: ", err.Error())
			return
		}
		fmt.Println(pkg)
	}
	println()
}

// SearchJSON 搜索所有的json文件,放到PathSet中
func SearchJSON(Path string, PathSet *[]string) {
	fileinf, err := ioutil.ReadDir(Path)
	if err != nil {
		fmt.Println("dir read error!")
		return
	}
	for _, file := range fileinf {
		//如果是目录，递归搜索
		//filepath.Join()
		if file.IsDir() == true {
			SearchJSON(filepath.Join(Path, file.Name()), PathSet)
		} else if strings.HasSuffix(file.Name(), ".json") {
			*PathSet = append(*PathSet, filepath.Join(Path, file.Name()))
		} else if strings.HasSuffix(file.Name(), "yarn.lock") {
			*PathSet = append(*PathSet, filepath.Join(Path, file.Name()))
		}

	}
}

func main() {
	var PathSet []string
	SearchJSON("C:\\Users\\10507\\Downloads\\Compressed\\h5-Dooring-master", &PathSet)
	for _, path := range PathSet {
		fmt.Println(path)
		praseJSON(path)
	}

}
