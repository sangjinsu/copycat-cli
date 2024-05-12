package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func copyDir(src, dst, old, new string) error {
	// 소스 디렉토리 정보 가져오기
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 디렉토리 생성
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	// 디렉토리 내용 읽기
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 디렉토리인 경우 재귀적으로 복사
			wg.Add(1)
			go func(srcPath, dstPath, old, new string) {
				defer wg.Done()
				copyDir(srcPath, dstPath, old, new)
			}(srcPath, dstPath, old, new)
		} else {
			// 파일인 경우 복사
			wg.Add(1)
			go func(srcPath, dstPath, old, new string) {
				defer wg.Done()
				copyFile(srcPath, dstPath, old, new)
			}(srcPath, dstPath, old, new)
		}
	}

	wg.Wait()

	return nil
}

func copyFile(src, dst, old, new string) error {
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

	// 파일 내용 읽기
	scanner := bufio.NewScanner(srcFile)
	for scanner.Scan() {
		// 특정 문자열을 다른 문자열로 대체하고 새 파일에 쓰기
		newText := strings.ReplaceAll(scanner.Text(), old, new)
		_, err := fmt.Fprintln(dstFile, newText)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	var src, dst, old, new string
	inputs := bufio.NewScanner(os.Stdin)

	fmt.Print("소스 디렉토리: ")
	inputs.Scan()
	src = inputs.Text()

	fmt.Print("대상 디렉토리: ")
	inputs.Scan()
	dst = inputs.Text()

	fmt.Print("찾을 문자열: ")
	inputs.Scan()
	old = inputs.Text()

	fmt.Print("바꿀 문자열: ")
	inputs.Scan()
	new = inputs.Text()

	err := copyDir(src, dst, old, new)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("폴더 복사 완료")
}