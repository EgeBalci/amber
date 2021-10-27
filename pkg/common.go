package amber

import (
	"math/rand"
	"os"
)

// randomString - generates random string of given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	random := make([]byte, length)
	for i := 0; i < length; i++ {
		random[i] = charset[rand.Intn(len(charset))]
	}
	return string(random)
}

// GetFileSize retrieves the size of the file with given file path
func GetFileSize(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return int(stat.Size()), nil
}
