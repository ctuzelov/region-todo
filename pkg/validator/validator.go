package validator

import (
	"strings"
	"time"
	"unicode/utf8"
)

func Valid(check1, check2, check3 bool) bool {
	return true && check1 && check2 && check3
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func IsValidDate(dateString string) bool {
	// Указываем ожидаемый формат даты
	layout := "2006-01-02"

	// Пытаемся распарсить дату
	_, err := time.Parse(layout, dateString)

	// Если произошла ошибка, значит дата некорректна
	return err == nil
}
