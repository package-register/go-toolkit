package timetable

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// DataParser 结构体用于处理数据解析
type DataParser struct{}

func NewDataParser() *DataParser {
	return &DataParser{}
}

// ParseData 解析单个课程信息
func (dp *DataParser) ParseData(info string) map[string]string {
	spData := dp.SplitString(info, ",")
	return map[string]string{
		"courseName": spData[0],
		"classroom":  spData[1],
		"class":      spData[2],
		"teacher":    spData[3],
	}
}

// SplitString 分割字符串
func (dp *DataParser) SplitString(s string, sep string) []string {
	return strings.Split(s, sep)
}

// GetOneData 获取一周的数据
func (dp *DataParser) GetOneData(indexData []string) map[string]any {
	oneData := make(map[string]any)
	for index, d := range indexData {
		if d == "" {
			oneData[fmt.Sprintf("第%d节", index+1)] = "noClass"
		} else {
			oneData[fmt.Sprintf("第%d节", index+1)] = dp.ParseData(d)
		}
	}
	return oneData
}

// ReadFile 读取文件内容
func (dp *DataParser) ReadFile(filePath string) ([]any, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return dp.UnmarshalData(data)
}

// ReadBytes 读取字节流内容
func (dp *DataParser) ReadBytes(data []byte) ([]any, error) {
	return dp.UnmarshalData(data)
}

// UnmarshalData 将字节数据解析为 []any
func (dp *DataParser) UnmarshalData(data []byte) ([]any, error) {
	var jsonData []any
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// ProcessData 处理文件中的课程数据
func (dp *DataParser) ProcessData(filePath string) (map[int]map[string]any, error) {
	data, err := dp.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return dp.processGenericData(data)
}

// ProcessBytes 处理字节流中的课程数据
func (dp *DataParser) ProcessBytes(data []byte) (map[int]map[string]any, error) {
	unmarshalledData, err := dp.ReadBytes(data)
	if err != nil {
		return nil, err
	}
	return dp.processGenericData(unmarshalledData)
}

// processGenericData 处理通用的数据 (无论是从文件还是字节流读取的)
func (dp *DataParser) processGenericData(data []any) (map[int]map[string]any, error) {
	length := len(data)
	fullData := make(map[int]map[string]any)

	for i := range length {
		indexData, ok := data[i].([]any)
		if !ok {
			fmt.Printf("Invalid data format at index %d\n", i)
			continue
		}

		stringData := make([]string, len(indexData))
		for j, v := range indexData {
			str, ok := v.(string)
			if !ok {
				fmt.Printf("Invalid data type in inner array at index %d, element %d\n", i, j)
				stringData[j] = "" // Or handle the error differently
				continue
			}
			stringData[j] = str
		}

		d := dp.GetOneData(stringData)
		fullData[i] = d
	}

	return fullData, nil
}
