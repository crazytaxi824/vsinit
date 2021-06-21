package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

// 项目根目录下生成 test 文件夹
const JestFolder = "test"

// json 序列化，写入文件
func WrieFile(packageFile *os.File, root *jsonvalue.V) error {
	// json 写入字段顺序
	orders := jsonvalue.Opt{
		MarshalKeySequence: []string{"name", "version", "description",
			"main", "directories", "scripts", "jest", "keywords", "author", "license",
			"dependencies", "devDependencies"},
	}

	// json 序列化
	marshbytes, err := root.Marshal(orders)
	if err != nil {
		return err
	}

	// 格式化
	var buf bytes.Buffer
	err = json.Indent(&buf, marshbytes, "", "  ")
	if err != nil {
		return err
	}

	// unescape json string
	result, err := UnescapeStringInJSON(buf.String())
	if err != nil {
		return err
	}

	// 重置 offset 清空文件然后重新写入内容
	return truncAndWrieFile(packageFile, result)
}

// 清空文件，重新写入内容
func truncAndWrieFile(packageFile *os.File, content string) error {
	// 将读写符重置到文件的起始位置
	_, err := packageFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// 清空文件
	err = packageFile.Truncate(0)
	if err != nil {
		return err
	}

	// 写入新内容
	_, err = packageFile.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// 将文件内容 json 反序列化到 map 中。
func ReadFileToJsonvalue(packageFile *os.File) (*jsonvalue.V, error) {
	// 读取 package.json 内容
	packageContent, err := io.ReadAll(packageFile)
	if err != nil {
		return nil, err
	}

	// json 反序列化
	return jsonvalue.Unmarshal(packageContent)
}

// 将字段添加到 packageMap
func CheckPackageFile(packageRootV, packageConfig *jsonvalue.V, key string) error {
	// 判断 key 是否存在
	value, err := packageRootV.Get(key)
	if err != nil && !errors.Is(err, jsonvalue.ErrNotFound) {
		return err
	} else if !errors.Is(err, jsonvalue.ErrNotFound) && !value.IsObject() {
		// 如果 key 存在，但是不是 object 的情况
		er := packageRootV.Delete(key)
		if er != nil {
			return er
		}
	}

	cfgV, err := packageConfig.Get(key)
	if err != nil {
		return err
	}

	// 插入数据
	cfgV.RangeObjects(func(k string, v *jsonvalue.V) bool {
		_, er := packageRootV.Set(v).At(key, k)
		if er != nil {
			err = er
			return false
		}
		return true
	})
	if err != nil {
		return err
	}

	return nil
}
