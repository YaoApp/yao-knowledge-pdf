package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	pdf := &PDF{}
	text, err := pdf.Text(filepath.Join(testPath(t), "tests", "software.pdf"))
	if err != nil {
		fmt.Printf("获取文本内容发生错误: %v\n", err)
		t.Fatal(err)
		return
	}
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "application framework")
}

func TestContent(t *testing.T) {
	pdf := &PDF{}
	rows, err := pdf.Content(filepath.Join(testPath(t), "tests", "software.pdf"))
	if err != nil {
		fmt.Printf("获取文本内容发生错误: %v\n", err)
		t.Fatal(err)
		return
	}
	assert.Equal(t, 20, len(rows))
}

func testPath(t *testing.T) string {
	// 获取调用者的文件路径和行号
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("获取调用者的文件路径和行号发生错误")
	}

	// 获取源码路径的绝对路径
	srcFile, err := filepath.Abs(file)
	if err != nil {
		t.Fatal(err)
	}

	return filepath.Dir(srcFile)
}
