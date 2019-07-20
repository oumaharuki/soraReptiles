package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func SaveImg(url, imgPath, path string, str []string) {
	///upload/vod/20190712-1/37028a8a314e23ed79ef7e4c31dd14b4.jpg

	fmt.Println(str)
	fmt.Println(len(str))

	url = "http://pilipali.cc" + url
	bol := Exists(path)

	if !bol {
		err1 := os.Mkdir(path, os.ModePerm) //创建文件夹
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		resp, err := httpClient.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		respbyte, _ := ioutil.ReadAll(resp.Body)

		err = ioutil.WriteFile(imgPath, respbyte, 06444)
		if err != nil {
			fmt.Println("下载失败")
		} else {
			fmt.Println("下载成功")
		}

	} else {
		resp, err := httpClient.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		respbyte, _ := ioutil.ReadAll(resp.Body)

		err = ioutil.WriteFile(imgPath, respbyte, 06444)
		if err != nil {
			fmt.Println("下载失败")
		} else {
			fmt.Println("下载成功")
		}
	}
}
