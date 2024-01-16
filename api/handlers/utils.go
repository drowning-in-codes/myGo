package handlers

import (
	"bufio"
	"errors"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

const pixivision_site = "https://www.pixivision.net/zh/c/illustration"
const safebooru_site = "https://safebooru.org/index.php?page=post&s=list"
const booru_site = "https://danbooru2.booru.org/index.php?page=post&s=list"
const mygo_site = "https://anime-pictures.net/posts?page=0&search_tag=bang+dream!+it%27s+mygo!!!!!&lang=en"

var conf = loadConf("./configure.yaml")
var download_root_folder, ok = conf["root_folder"].(string)
var proxy = conf["proxy"].(map[string]interface{})

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"

func checkAllowSite() []string {
	var allow_img_site = make([]string, 0)

	value := reflect.ValueOf(conf["allow_site"])
	if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
		// implement
		for i := 0; i < value.Len(); i++ {
			// fmt.Println(value.Index[i])
			allow_img_site = append(allow_img_site, value.Index(i).Interface().(string))
		}
	} else {
		log.Panic("请填写允许爬取网站域名字符串")
	}
	return allow_img_site
}
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
