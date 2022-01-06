package utils_v1

import (
	"bufio"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

// 1、结构体 -------------------------------------------------------------------------
type uFile struct {
}

// 2、全局变量 -------------------------------------------------------------------------

// 3、初始化函数 -------------------------------------------------------------------------

// 4、开放函数 -------------------------------------------------------------------------

// 对外函数1
func File() *uFile {
	return &uFile{}
}

/**-------------------------
// 名称：检查文件是否存在
***-----------------------*/
func (Me uFile) CheckFileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 获取文件大小
func (Me uFile) FileSize(fileAddr string) (int64, error) {
	// 保存记录
	fi, err := os.Stat(fileAddr)
	var fileSize int64 = 0
	if err != nil {
		return 0, err
	} else {
		fileSize = fi.Size()
	}
	return fileSize, nil
}

/**-------------------------
// 名称：读取文件全部内容
***-----------------------*/
func (Me uFile) ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	return ioutil.ReadAll(f)
}

// 辅助函数5：读取指定长度数据
func (Me *uFile) FileReadLength(fp *os.File, Length int) ([]byte, error) {
	// 读取一批
	var retBuff = make([]byte, 0)
	for {
		ContentTran := make([]byte, Length) // 建立一个slice
		// 读取数据
		if n, err := fp.Read(ContentTran); err == nil {
			if n != Length {
				retBuff = append(retBuff, ContentTran[:n]...)
				Length -= n
			} else {
				retBuff = append(retBuff, ContentTran[:n]...)
				break
			}
		} else {
			return nil, err
		}
	}
	return retBuff, nil
}

/**-------------------------
// 名称：分片读取内容
***-----------------------*/
func (Me uFile) ReadBlock(filePth string, bufSize int, hookFn func(buff []byte) bool) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	buf := make([]byte, bufSize) // 一次读取多少个字节
	bfRd := bufio.NewReader(f)
	for {
		// 读取数据
		n, err := bfRd.Read(buf)

		// 回调
		if n > 0 {
			continueRead := hookFn(buf[:n]) // 放在错误处理前面，即使发生错误，也会处理已经读取到的数据。
			if !continueRead {              // n 是成功读取字节数
				return nil
			}
		}
		// 遇到任何错误立即返回， 如果是EOF，返回成功
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

	}
}

/**-------------------------
// 名称：逐行读取内容
err := zfile.ReadLine("tbl_resource_main.txt", func(data []byte) bool {
	if len(data) > 0 {
		rowNum++
	}
	if rowNum >= 2000 {
		return false
	}
	return true
})
***-----------------------*/
func (Me uFile) ReadLine(filePth string, backFunc func(buff []byte) bool) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		if len(line) > 0 {
			continueRead := backFunc(line) // 放在错误处理前面，即使发生错误，也会处理已经读取到的数据。
			if !continueRead {
				return nil
			}
		}
		// 遇到任何错误立即返回， 如果是EOF，返回成功
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

/**-------------------------
// 名称：全部写入
***-----------------------*/
func (Me uFile) WriteAll(filename string, B []byte) error {
	var (
		f    *os.File
		err1 error
	)
	f, err1 = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644) // 打开文件
	if err1 != nil {
		return err1
	}
	defer func() { _ = f.Close() }()

	_, err1 = f.Write(B)
	// _, err1 = io.WriteString(f, stringWrite) // 写入文件(字符串)
	if err1 != nil {
		return err1
	}
	return nil
}

/**-------------------------
// 名称：追加写入
***-----------------------*/
func (Me uFile) WriteAppend(filename string, buffWrite []byte) error {
	var (
		f    *os.File
		err1 error
	)
	f, err1 = os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644) // 打开文件
	if err1 != nil {
		return err1
	}
	defer func() { _ = f.Close() }()

	// 写入文件(字符串)
	if _, err := f.Write(buffWrite); err != nil {
		return err
	}
	return nil
}

/**-------------------------
// 名称：删除文件
***-----------------------*/
func (Me uFile) DelFile(filename string) (bool, error) {
	var (
		err error
	)
	err = os.Remove(filename)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// 下载文件到指定目录，保持相对路径
// urlFile:     http://static.anoah.com/uploads//bookcover/46/ab/5cb936539507a.jpg
// floder:      /app/nlsper/static/
// 文件保存位置: /app/nlsper/static/uploads/bookcover/46/ab/5cb936539507a.jpg
// replace:   文件存在替换掉他
// 返回
//  1、errorCode:  1:目录不存在 2:解析url出错 3:创建文件目录失败 100:【文件已存在】 5:下载失败
//  2、filePath: 文件保存到的位置
//  3、error
func (Me uFile) DownLoadUrl2FolderFresh(urlFile string, inFolder string, replace bool) (int, string, error) {
	if !Me.CheckFileExist(inFolder) {
		return 1, "", errors.New(inFolder + " 目录不存在")
	}
	// 解析url
	p, err := url.Parse(urlFile)
	if err != nil {
		return 2, "", err
	}
	folder := path.Dir(p.Path)

	// 创建目录
	err = os.MkdirAll(Me.FolderJoin(inFolder, folder), os.ModePerm)
	if err != nil {
		return 3, "", err
	}

	// 文件存在判断
	fileAddr := Me.FolderJoin(inFolder, folder) + path.Base(p.Path)
	if Me.CheckFileExist(fileAddr) && !replace {
		return 100, fileAddr, errors.New("文件已存在:" + fileAddr)
	}

	// 执行下载
	err = Me.DownLoadUrlFile(urlFile, fileAddr)
	if err != nil {
		return 5, fileAddr, err
	}
	return 0, fileAddr, nil
}

// 名称：下载网络文件到指定(完整路径)位置
//   1、保存到的目录不存在：报错
//   2、目录存在，文件不存在：覆盖
func (Me uFile) DownLoadUrlFile(urlAddr string, fileAddr string) error {
	// 下载文件
	res, err := http.Get(urlAddr)
	if err != nil {
		return err
	}

	// 创建文件
	f, err := os.Create(fileAddr)
	if err != nil {
		return err
	}
	defer f.Close()

	// 保存文件
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}
	return nil
}

// 下载文件到指定目录，保持相对路径
// urlFile:     http://static.anoah.com/uploads//bookcover/46/ab/5cb936539507a.jpg
// floder:      /app/nlsper/static/
// 文件保存位置: /app/nlsper/static/uploads/bookcover/46/ab/5cb936539507a.jpg
// replace:   文件存在替换掉他
// 返回
//  1、errorCode:  1:目录不存在 2:解析url出错 3:创建文件目录失败 100:【文件已存在】 5:下载失败
//  2、filePath: 文件保存到的位置
//  3、error
func (Me uFile) DownLoadUrl2FolderFreshWGet(urlFile string, inFolder string, replace bool) (int, string, error) {
	if !Me.CheckFileExist(inFolder) {
		return 1, "", errors.New(inFolder + " 目录不存在")
	}
	// 解析url
	p, err := url.Parse(urlFile)
	if err != nil {
		return 2, "", err
	}
	folder := path.Dir(p.Path)

	// 创建目录
	err = os.MkdirAll(Me.FolderJoin(inFolder, folder), os.ModePerm)
	if err != nil {
		return 3, "", err
	}

	// 文件存在判断
	fileAddr := Me.FolderJoin(inFolder, folder) + path.Base(p.Path)
	if Me.CheckFileExist(fileAddr) && !replace {
		return 100, fileAddr, errors.New("文件已存在:" + fileAddr)
	}

	// 执行下载
	err = Me.DownLoadUrlFileWGet(urlFile, fileAddr)
	if err != nil {
		return 5, fileAddr, err
	}
	return 0, fileAddr, nil
}

// 下载文件
// 说明
//   1、保存到的目录不存在：报错
//   2、目录存在，文件不存在：覆盖
func (Me uFile) DownLoadUrlFileWGet(urlAddr string, fileAddr string) error {
	// 单引号处理
	urlAddrCmd := strings.ReplaceAll(urlAddr, "'", "%27")
	urlAddrCmd = strings.ReplaceAll(urlAddrCmd, "\"", "%22")
	fileAddrCmd := strings.ReplaceAll(fileAddr, "'", "'\"'\"'")
	// wget命令下载
	cmd := exec.Command("bash", "-c", `wget -4 --no-check-certificate '`+urlAddrCmd+`' -T1800 -t1  -O '`+fileAddrCmd+`' `)
	buf, err := cmd.Output()
	// 下载出错
	if err != nil {
		log.Error("下载出错", err, string(buf))
		return err
	}
	// 下载的文件不存在
	if !Me.CheckFileExist(fileAddr) {
		return errors.New("下载的文件不存在:" + fileAddr)
	}
	return nil
}

// 创建目录
// MkdirAll("/app/services","/abc/def")
func (Me uFile) MkdirAll(inFolder, folder string) error {
	err := os.MkdirAll(Me.FolderJoin(inFolder, folder), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// 文件路径拼接
func (Me uFile) FolderJoin(inFolder, folder string) string {
	return strings.TrimRight(inFolder, "/") + "/" + strings.Trim(folder, "/") + "/"
}

// 查找文本内容里的url地址
func (Me uFile) UrlsInContent(content string, fReplace func(Url string) string) (contentReplaced string) {

	// 匹配正则
	reg := regexp.MustCompile(`(?i)(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)

	// 执行正则匹配
	return reg.ReplaceAllStringFunc(content, func(s string) string {

		// 回调并以回调的内容返回替换
		return fReplace(s)

	})

}
