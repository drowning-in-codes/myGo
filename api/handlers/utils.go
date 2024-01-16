package handlers

import (
	"bufio"
	"errors"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const pixivision_site = "https://www.pixivision.net/zh/c/illustration"
const safebooru_site = "https://safebooru.org/index.php?page=post&s=list"
const booru_site = "https://danbooru2.booru.org/index.php?page=post&s=list"

var conf = loadConf("./configure.yaml")
var download_root_folder, ok = conf["root_folder"].(string)

func PicExist(file string) (err error) {
	// 查看该目录下是否有图片
	entries, err := os.ReadDir(file)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return errors.New("没有图片")
	}
	return nil
}

func RandomFile(dirPath string) string {
	// 传入目录，根据静态资源的配置 返回该目录下随机一张图片的绝对路径
	entries, _ := os.ReadDir(dirPath)
	// 随机获取一张图片
	random_file := entries[rand.Intn(len(entries))]
	random_file_abs := filepath.Join(dirPath, random_file.Name())
	return random_file_abs
}

func RandomFiletoURL(dirPath string) (string, error) {
	// 传入目录，根据静态资源的配置 返回该目录下随机一张图片的绝对路径 然后再转为url
	mid_dir := filepath.Base(dirPath)
	entries, err := os.ReadDir(filepath.Join("./imgs", mid_dir))
	// 读取文件夹失败
	if err != nil {
		log.Panic(err.Error())
		return "", err
	}
	// 该目录下没有图片
	if len(entries) == 0 {
		return "", errors.New("没有图片")
	}
	// 随机获取一张图片
	random_file := entries[rand.Intn(len(entries))]
	random_file_abs := strings.Join([]string{"/img", mid_dir, random_file.Name()}, "/")
	return random_file_abs, nil
}

type URL_path interface {
	isPic() bool
	picName() string
}
type url_path string

func (url url_path) isPic() bool {
	ext := url[strings.LastIndex(string(url), ".")+1:]
	if ext == "jpg" || ext == "png" || ext == "gif" || ext == "jpeg" || ext == "webp" {
		return true
	}
	return false
}

func (url url_path) picName() string {
	img_name := strings.Split(string(url), "/")[len(strings.Split(string(url), "/"))-1]
	return img_name
}

func loadConf(configure_path string) map[interface{}]interface{} {
	file, err := os.Open(configure_path)
	if err != nil {
		log.Panic(err.Error())
	}
	defer file.Close()
	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		log.Panic(err.Error())
	}

	bytes := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bytes)
	if err != nil {
		log.Panic(err.Error())
	}
	results := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(bytes), &results)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return results
}
