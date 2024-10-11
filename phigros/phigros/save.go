package phigros

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
)

func ReadZip(path string) (m map[string][]byte, err error) {
	m = map[string][]byte{}
	// 打开 zip 文件
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// 遍历 zip 文件中的文件
	for _, file := range reader.File {
		// 打开文件
		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()
		// 读取文件内容
		buf := make([]byte, file.FileInfo().Size())
		_, _ = f.Read(buf)
		m[file.Name] = buf
	}
	return m, nil
}

func Decrypt(in []byte) (out []byte, err error) {
	// CBCDecrypt AES-CBC 解密
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(in) < aes.BlockSize {
		return nil, fmt.Errorf("cipherText too short")
	}
	out = make([]byte, len(in))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(out, in[1:])

	return append(in[0:1], unpad(out)...), nil
	//return out, nil
}

// 去填充函数的示例实现
func unpad(data []byte) []byte {
	padding := data[len(data)-1]
	return data[:len(data)-int(padding)]
}

func DecoderWithStruct[T PhigrosStruct](in []byte) *T {
	var ps T
	v := reflect.ValueOf(&ps).Elem()
	reader := NewBytesReader(in)
	for i := 0; i < v.NumField(); i++ {
		set(v.Field(i), reader)
	}
	return &ps
}

func set(rv reflect.Value, reader *Bytes) {
	switch rv.Kind() {
	case reflect.Bool:
		rv.SetBool(reader.ReadBool())
	case reflect.String:
		rv.SetString(reader.ReadString())
	case reflect.Float32:
		rv.SetFloat(float64(reader.ReadFloat32()))
	case reflect.Int16:
		rv.SetInt(int64(reader.ReadShort()))
	case reflect.Uint8:
		rv.SetUint(uint64(reader.ReadByte1()))
	case reflect.Array:
		for i, j := 0, rv.Len(); i < j; i++ {
			if rv.Index(i).Kind() == reflect.Struct {
				for ii, k := 0, rv.Index(i).NumField(); ii < k; ii++ {
					set(rv.Index(i).Field(ii), reader)
				}
			} else {
				//非结构体数组
				set(rv.Index(i), reader)
			}
		}
	default:
	}
}

func DecoderGameRecord(in []byte) []ScoreAcc {
	records := []ScoreAcc{}
	reader := NewBytesReader(in)

	for i, s := byte(0), reader.ReadVarShort(); i < s; i++ {
		t := reader.ReadString()
		songId := t[:len(t)-2]
		record := reader.ReadRecord(songId)
		records = append(records, record...)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Rks > records[j].Rks
	})
	return records

}

// 前19成绩,取最高成绩放第一位
func B19(records []ScoreAcc) []ScoreAcc {
	return BN(records, 19)
}

// 取前n成绩,取最高成绩放第一位
func BN(records []ScoreAcc, n int) []ScoreAcc {
	var maxRecord ScoreAcc
	for _, r := range records {
		if r.Score == 1000000 {
			if r.Difficulty > maxRecord.Difficulty {
				maxRecord = r
			}
		}
	}
	bn := []ScoreAcc{maxRecord}
	if n <= 0 {
		return append(bn, records...)
	}
	// 将records中的前19个记录加入b19
	if len(records) >= n {
		bn = append(bn, records[:n]...)
	} else {
		bn = append(bn, records...)
	}
	return bn
}

// 通过zip文件读取所有云端内容
func ParseSave(path string) (map[string]any, error) {
	m, err := ReadZip(path)
	if err != nil {
		return nil, err
	}
	for k, v := range m {
		out, err := Decrypt(v)
		if err != nil {
			return nil, fmt.Errorf("Decrypt file %s Error %s", k, err.Error())
		}
		m[k] = out
	}
	if m["gameRecord"][0] != byte(0x01) {
		return nil, errors.New("版本号不正确，可能协议已更新。")
	}
	//json
	jsons := make(map[string]any)
	jsons["gameRecord"] = B19(DecoderGameRecord(m["gameRecord"][1:]))
	jsons["settings"] = *DecoderWithStruct[Settings](m["settings"][1:])
	jsons["user"] = *DecoderWithStruct[User](m["user"][1:])
	return jsons, nil
}

// 通过url获取战绩,其余内容丢弃
func ParseStatsByUrl(url string) ([]ScoreAcc, error) {
	d, err := GetGameRecordData(url)
	f, _ := os.Create("./gamesave")
	f.Write(d)
	f.Close()
	if err != nil {
		return nil, err
	}
	d, err = Decrypt(d)
	if err != nil {
		return nil, fmt.Errorf("Decrypt file gameRecord Error %s", err.Error())
	}
	if d[0] != byte(0x01) {
		return nil, errors.New("版本号不正确，可能协议已更新。")
	}
	return DecoderGameRecord(d[1:]), nil
}

func ProcessSummary(sum string) (s *Summary) {
	if sum == "" {
		return nil
	}
	b, err := base64.RawStdEncoding.DecodeString(sum)
	if err != nil {
		return nil
	}
	s = DecoderWithStruct[Summary](b)
	// for range 4 {
	// 	s.ScoreAcc = append(s.ScoreAcc, SummaryScoreAcc{
	// 		Cleared:   reader.ReadShort(),
	// 		FullCombo: reader.ReadShort(),
	// 		Phi:       reader.ReadShort(),
	// 	})
	// }
	return
}
