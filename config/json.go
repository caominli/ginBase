package config //定义在commons共享包下
//json主要用于处理json相关，并且解析带有注释的JSON字符串
import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

// JsonGet 解析带有注释的 JSON 字符串
func JsonGet(data []byte, v interface{}) error {
	processedData, err := removeComments(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(processedData, v) // 使用标准库解析 JSON
}

// 移除 JSON 字符串中的注释
func removeComments(data []byte) ([]byte, error) {
	lines := bytes.Split(data, []byte("\n"))
	var outputLines [][]byte

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(string(line))
		if strings.HasPrefix(trimmedLine, "//") {
			continue // 跳过以 `//` 开头的注释行
		}
		outputLines = append(outputLines, line)
	}

	if len(outputLines) == 0 {
		return nil, errors.New("empty json after removing comments")
	}

	return bytes.Join(outputLines, []byte("\n")), nil
}
