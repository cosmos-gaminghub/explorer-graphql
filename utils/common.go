package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"crypto/md5"
	"encoding/hex"
)

const (
	StatusFailed = "failed"
	StatusFail   = "fail"
)

func RemoveDuplicationStrArr(list []string) []string {
	unique_set := make(map[string]bool, len(list))
	for _, x := range list {
		unique_set[x] = true
	}
	result := make([]string, 0, len(unique_set))
	for x := range unique_set {
		result = append(result, x)
	}
	return result
}

func ParseInt(text string) (i int64, b bool) {
	i, err := strconv.ParseInt(text, 10, 0)
	if err != nil {
		return i, false
	}
	return i, true
}

func ParseIntBase(text string) (i int, b bool) {
	i, err := strconv.Atoi(text)
	if err != nil {
		return i, false
	}
	return i, true
}

func ParseUint(text string) (uint64, bool) {
	i, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return i, false
	}
	return i, true
}

func RoundFloat(num float64, bit int) (i float64, b bool) {
	format := "%" + fmt.Sprintf("0.%d", bit) + "f"
	s := fmt.Sprintf(format, num)
	i, err := strconv.ParseFloat(s, 0)
	if err != nil {
		return i, false
	}
	return i, true
}

func Round(x float64) int64 {
	return int64(math.Floor(x + 0.5))
}

func ParseStringToFloat(x string) (float64, error) {
	f, err := strconv.ParseFloat(x, 0)
	return f, err
}

func Copy(src interface{}, dest interface{}) error {
	bz, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bz, dest)
	if err != nil {
		return err
	}
	return nil
}

func NewRatFromFloat64(f float64) *big.Rat {
	return new(big.Rat).SetFloat64(f)
}

func ParseStringFromFloat64(data float64) string {
	return strconv.FormatFloat(data, 'f', -1, 64)
}

func Md5Encryption(data []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}
func FailtoFailed(status string) string {

	if status == StatusFail {
		return StatusFailed
	}
	return status
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
