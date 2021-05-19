//
// filerepository.go
//
// May 2021, Prashant Desai
//

package repository

import (
	"bufio"
	"fmt"
	"os"

	"path/filepath"

	"WardrobeManagerMS/pkg/api"
)

type fileImageRepo struct {
	Dir string
}

func NewFileImageRepository(path string) (api.ImageRepository, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	imageRepo := &fileImageRepo{
		Dir: path,
	}

	return imageRepo, nil
}

func (m *fileImageRepo) AddFile(name string, file []byte) error {

	path := filepath.Join(m.Dir, name)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error create file %s : %w", path, err)
	}
	defer f.Close()

	_, err = f.Write(file)

	if err != nil {
		return fmt.Errorf("Error writing to file %s : %w", path, err)
	}

	return nil
}

func (m *fileImageRepo) GetFile(name string) ([]byte, error) {

	filename := filepath.Join(m.Dir, name)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []byte{}, api.NoSuchFileOrDirectory{
			File: filename,
		}
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening file %s : %w", filename, err)
	}
	defer file.Close()

	stats, err1 := file.Stat()
	if err1 != nil {
		return nil, fmt.Errorf("Error getting stats forfile %s : %w", filename, err1)
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("Error while reading bytes from file %s : %w", filename, err)
	}

	return bytes, err

}

func (m *fileImageRepo) UpdateFile(name string, file []byte) error {
	return nil
}

func (m *fileImageRepo) DeleteFile(name string) error {

	return nil
}
