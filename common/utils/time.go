package utils

import (
	"fmt"
	"time"
)

func GetMonthTbName(tableNamePrefix string, fieldVal int) string {
	if fieldVal == 0 {
		return fmt.Sprintf("%s_%d", tableNamePrefix, GetYYYYMM())
	}
	yyyy := fieldVal / 10000
	mm := fieldVal % 10000 / 100
	return fmt.Sprintf("%s_%04d%02d", tableNamePrefix, yyyy, mm)
}

func GetYYYYMMDDHH() string {
	now := time.Now()
	return fmt.Sprintf("%04d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour())
}
func GetYYYYMMDDHHMMSS() string {
	now := time.Now()
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}
func GetYYYYMMDD() int {
	now := time.Now()
	return StrToInt(fmt.Sprintf("%04d%02d%02d", now.Year(), now.Month(), now.Day()))
}
func GetYYYYMM() int {
	now := time.Now()
	return StrToInt(fmt.Sprintf("%04d%02d", now.Year(), now.Month()))
}
