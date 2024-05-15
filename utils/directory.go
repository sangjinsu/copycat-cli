package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var mu sync.Mutex
var copyErr error

func CopyDir(src, dst string, replacer *strings.Replacer) error {
	// 소스 디렉토리 정보 가져오기
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 디렉토리 생성
	_, err = os.Stat(dst)
	if !os.IsNotExist(err) {
		return fmt.Errorf("destination directory already exists: %s", dst)
	}

	log.Println("create directory", dst)

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	// 디렉토리 내용 읽기
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		wg.Add(1)
		if entry.IsDir() {
			// 디렉토리인 경우 재귀적으로 복사
			go copyFileSystem(srcPath, dstPath, replacer, CopyDir)
		} else {
			// 파일인 경우 복사
			go copyFileSystem(srcPath, dstPath, replacer, copyFile)
		}
	}

	wg.Wait()

	return copyErr
}

func copyFileSystem(srcPath, dstPath string, replacer *strings.Replacer, copyFunc func(srcPath, dstPath string, replacer *strings.Replacer) error) {
	defer wg.Done()
	if err := copyFunc(srcPath, dstPath, replacer); err != nil {
		mu.Lock()
		copyErr = err
		mu.Unlock()
	}
}

func copyFile(src, dst string, replacer *strings.Replacer) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = srcFile.Seek(0, 0) // 파일 포인터를 처음으로 이동
	if err != nil {
		return err
	}

	// 압축 파일인 경우 복사
	if isArchived, isArchiveErr := isArchive(srcFile); isArchived {
		if copyArchiveErr := copyArchive(src, dst); copyArchiveErr != nil {
			return copyArchiveErr
		}
		return nil
	} else if isArchiveErr != nil {
		return isArchiveErr
	}

	// 파일 내용 읽기 및 교체
	log.Println("copy file", src, "to", dst)
	_, err = srcFile.Seek(0, 0) // 파일 포인터를 처음으로 이동
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(srcFile)
	writer := bufio.NewWriter(dstFile)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()

		// 교체
		line = replacer.Replace(line)

		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

func copyArchive(src, dst string) error {
	// 원본 파일 열기
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 대상 파일 열기
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 파일 복사
	if _, err = io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}

	return nil
}

func isArchive(file *os.File) (bool, error) {
	// 첫 512 바이트를 읽어 MIME 타입을 판별합니다.
	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return false, err
	}

	if _, err := file.Seek(0, 0); err != nil { // 파일 포인터를 처음으로 되돌립니다.
		return false, err
	}

	mimeType := http.DetectContentType(buf)

	// 일반적인 압축 파일의 MIME 타입
	archiveTypes := []string{
		"application/zip",
		"application/x-tar",
		"application/gzip",
		"application/x-gzip",
		"application/x-bzip2",
		"application/x-7z-compressed",
		"application/x-rar-compressed",
		"application/x-xz",
	}

	for _, archiveType := range archiveTypes {
		if mimeType == archiveType {
			return true, nil
		}
	}

	return false, nil
}
