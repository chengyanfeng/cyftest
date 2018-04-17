package utils

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/graphics-go/graphics"
)
//1.读取文件输出字符串
func ReadFile(path string) string {
	return string(ReadFileBytes(path))
}
//2.读取文件输出字节
func ReadFileBytes(path string) []byte {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		Error("ReadFile", err)
	}
	return c
}
//3.像文件中输入字符串
func WriteFile(path string, body []byte) error {
	err := ioutil.WriteFile(path, body, 0644)
	if err != nil {
		Error(err)
		return err
	}
	return nil
}
//4.像文件中添加文本
func AppendFile(file string, text string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		Error(err)
	}
	if _, err = f.WriteString(text); err != nil {
		Error(err)
	}
}
//5.删除文件
func DeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		Error(err)
	}
}
//5.读取文件的一段字符
func ReadLine(fileName string, limit int, offset int) (r string, e error) {
	f, err := os.Open(fileName)
	if err != nil {
		e = err
		return
	}
	buf := bufio.NewReader(f)
	for i := 0; i < offset+limit; i++ {
		line, err := buf.ReadString('\n')
		if i >= offset {
			r = r + line
		}
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}
	}
	return
}
//6.替换文本中的某一行的文本
func ReplaceLine(fileName string, line int, with string) (string, error) {
	if line < 1 {
		return "", errors.New(JoinStr("无效的行号", line))
	}
	return Exec(fmt.Sprintf(`sed -i '' '%vs/.*/%v/' %v`, line, with, fileName))
}
//7.把url解析成map
func Pathinfo(url string) P {
	p := P{}
	url = strings.Replace(url, "\\", "/", -1)
	if strings.Index(url, "/") < 0 {
		url = JoinStr("./", url)
	}
	re := regexp.MustCompile("(.*)/([^/]*)\\.([^.]*)")
	match := re.FindAllStringSubmatch(url, -1)
	if len(match) > 0 {
		m0 := match[0]
		fmt.Println(m0)
		if len(m0) == 4 {
			p["basename"] = m0[0]
			p["dirname"] = m0[1]
			p["filename"] = m0[2]
			p["extension"] = strings.ToLower(m0[3])
		}
	}
	return p
}
//8.判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
//9.创建一个文件
func Mkdir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
//10.取出文件中的包含的某段文本
func ExtractFile(path string, target string, ext string) {
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(f.Name(), ext) {
			Copy(path, target+"/"+f.Name())
		}
		return nil
	})
}
//11.遍历文件夹中的层级
func DirTree(path string, ext string, limit int) (files []P) {
	files = []P{}
	i := 0
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if i >= limit {
			return errors.New("reach limit")
		}
		i++
		if f != nil && !f.IsDir() {
			if strings.HasSuffix(f.Name(), ext) {
				files = append(files, P{"file": path})
			}
		}
		return nil
	})
	return
}
//12.删除文本中的某行
func FileRemoveLine(file string, start int, lines int) {
	cmd := fmt.Sprintf("sed -i '%v,%vd' %v", start, start+lines-1, file)
	Exec(cmd)
}
//13.删除文本中的某些字符
func RemoveSpaceLine(file string, filter interface{}) {
	cmd := fmt.Sprintf("sed -i '/%v/d' %v", filter, file)
	Exec(cmd)
}
//14.像文件中插入某些字符
func FileInsertLine(file string, start int, txt string) {
	cmd := fmt.Sprintf("sed -i '%vi %v' %v", start, txt, file)
	Exec(cmd)
}
//14.图片的大小
func ResizeImage(file string, width int) error {
	src, err := LoadImage(file)
	if err != nil {
		return err
	}
	bound := src.Bounds()
	dx := bound.Dx()
	dy := bound.Dy()
	// 缩略图的大小
	dst := image.NewRGBA(image.Rect(0, 0, width, width*dy/dx))
	// 产生缩略图,等比例缩放
	err = graphics.Scale(dst, src)
	if err != nil {
		return err
	}
	//保存文件
	err = SaveImage(file, dst)
	if err != nil {
		return err
	}
	return nil
}
//15.下载图片
func LoadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, err
}
//15.保存图片
func SaveImage(path string, img image.Image) (err error) {
	imgfile, err := os.Create(path)
	defer imgfile.Close()
	err = png.Encode(imgfile, img)
	if err != nil {
		log.Fatal(err)
	}
	return
}
