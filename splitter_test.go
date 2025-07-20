package splitter

import (
	"fmt"
	"testing"
)

func TestSplit(t *testing.T) {
	// fmt.Println(numDigits(18446744073709551615))
	v := numDigits(int64(1234567890))
	fmt.Println(v)

	d := int64(15) / int64(7)
	fmt.Println(d + 1)
}

func TestBytes(t *testing.T) {
	// dirPath := "/home/alex/Development/learning-go/splitter/run/file-data_zoom_amd64.deb"

}

func TestCleanFile(t *testing.T) {
	dirPath, fileName := "a/b/c", "testFileName"
	expectedFilePath := dirPath + "/" + fileName

	cleanedFilePath := cleanFilePath(dirPath, fileName)
	if cleanedFilePath != dirPath+"/"+fileName {
		t.Errorf("incorrectly cleaned filepath - expected '%s' - actual '%s'",
			expectedFilePath, cleanedFilePath)
	}

	dirPath = "a/b/c/"
	expectedFilePath = dirPath + fileName
	cleanedFilePath = cleanFilePath(dirPath, fileName)
	if cleanedFilePath != dirPath+fileName {
		t.Errorf("incorrectly cleaned filepath - expected '%s' - actual '%s'",
			expectedFilePath, cleanedFilePath)
	}
}

func TestNumDigits(t *testing.T) {
	number := int64(12345)
	result := numDigits(number)
	if result != 5 {
		t.Errorf("incorrectly counted number of digits - expected %d - actual %d", 5, result)
	}

	number = int64(999999999999999999)
	result = numDigits(number)
	if result != 18 {
		t.Errorf("incorrectly counted number of digits - expected %d - actual %d", 5, result)
	}
}

func TestGetFileFormatSize(t *testing.T) {
	kilo := 1
	result, err := getFileFormatSize("B")
	if err != nil || result != kilo {
		t.Errorf("expected: %d - actual: %d with error: %s", kilo, result, err)
	}

	kilo *= 1024
	result, err = getFileFormatSize("KB")
	if err != nil || result != kilo {
		t.Errorf("expected: %d - actual: %d with error: %s", kilo, result, err)
	}

	kilo *= 1024
	result, err = getFileFormatSize("MB")
	if err != nil || result != kilo {
		t.Errorf("expected: %d - actual: %d with error: %s", kilo, result, err)
	}

	kilo *= 1024
	result, err = getFileFormatSize("GB")
	if err != nil || result != kilo {
		t.Errorf("expected: %d - actual: %d with error: %s", kilo, result, err)
	}

	result, err = getFileFormatSize("TB")
	if err == nil || result != 0 {
		t.Errorf("expected: %d - actual: %d with error: %s", kilo, result, err)
	}

	result, err = getFileFormatSize("test")
	if err == nil || result != 0 {
		t.Errorf("expected: %d - actual: %d with error: %s", kilo, result, err)
	}
}
