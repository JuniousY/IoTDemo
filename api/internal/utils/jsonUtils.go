package utils

import (
	"bytes"
	"encoding/json"
	"gorm.io/datatypes"
	"io"
)

// Marshal 序列化对象为JSON字节
func Marshal[T any](v T) []byte {
	b, _ := json.Marshal(v)
	return b
}

// MarshalIndent 带缩进的序列化
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal 反序列化JSON字节
func Unmarshal[T any](data []byte) T {
	var obj T
	json.Unmarshal(data, &obj)
	return obj
}

func UnmarshalWithErr[T any](data []byte) (T, error) {
	var obj T
	err := json.Unmarshal(data, &obj)
	return obj, err
}

// MarshalString 序列化为字符串
func MarshalString[T any](v T) string {
	return string(Marshal(v))
}

// UnmarshalString 从字符串反序列化
func UnmarshalString[T any](data string, v T) {
	json.Unmarshal([]byte(data), v)
}

// 高级操作

// Valid 验证JSON有效性
func Valid(data []byte) bool {
	return json.Valid(data)
}

// PrettyPrint JSON美化输出
func PrettyPrint(data []byte) ([]byte, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "  "); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// Stream处理

// DecodeStream 流式解码
func DecodeStream(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// EncodeStream 流式编码
func EncodeStream(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// gorm
func StructToDatatypesJSON[T any](v T) datatypes.JSON {
	jsonData, _ := json.Marshal(v) // 序列化为 []byte
	return datatypes.JSON(jsonData)
}
