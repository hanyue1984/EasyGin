package tools

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"math/rand"
	"reflect"
	"strings"
)

type common struct {
}

var Common common

func (c common) GenerateUniqueUid() string {
	// 生成一个随机的唯一标识符作为Uid
	uid := uuid.New().String()
	return uid
}

// GenerateUniqueHash 生成一个指定长度的hash
func (c common) GenerateUniqueHash(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteByte(charset[rand.Intn(len(charset))])
	}
	return result.String()
}

// GenerateMD5Hash 生成一个md5
func (c common) GenerateMD5Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// MergeObjects 函数用于将对象 B 的值覆盖对象 A 的值
func MergeObjects(a interface{}, b interface{}) {
	valA := reflect.ValueOf(a).Elem()
	valB := reflect.ValueOf(b)

	for i := 0; i < valA.NumField(); i++ {
		fieldA := valA.Field(i)
		fieldB := valB.Field(i)

		if fieldB.Kind() == reflect.Zero(fieldB.Type()).Kind() || fieldB.Kind() == reflect.String && fieldB.String() == "" {
			continue // 如果对象 B 的字段是零值或者是空，则跳过
		}

		fieldA.Set(fieldB)
	}
}
