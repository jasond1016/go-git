package main

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const dotDir = ".ggit"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ggit <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		err := InitRepo()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Initialized empty Git repository in .ggit/")
	case "status":
		err := Status()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ggit add <filename>")
			os.Exit(1)
		}
		err := Add(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "log":
		fmt.Println("Not implemented yet")
	case "commit":
		fmt.Println("Not implemented yet")
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func Sha1Hash(t string, input string) string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s %d\x00%s", t, len(input), input)))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func InitRepo() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error in get working directory")
	}

	// Check if .ggit folder already exists
	dotDir := joinPath(currentDir, ".ggit")
	if fileExists(dotDir) {
		fmt.Printf("Reinitialized existing Git repository in %s\n", dotDir)
		return nil
	}

	// Create .ggit folder
	if err := os.Mkdir(dotDir, 0755); err != nil {
		return err
	}

	// Create config file
	if _, err := os.Create(joinPath(dotDir, "config")); err != nil {
		return err
	}

	// Create HEAD file
	headFile, err := os.Create(joinPath(dotDir, "HEAD"))
	if err != nil {
		return err
	}
	_, err = headFile.WriteString("ref: refs/heads/master\n")
	if err != nil {
		return err
	}

	// Create objects folder
	if err := os.Mkdir(joinPath(dotDir, "objects"), 0755); err != nil {
		return err
	}

	// Create refs folder
	if err := os.Mkdir(joinPath(dotDir, "refs"), 0755); err != nil {
		return err
	}

	if err := os.Mkdir(joinPath(dotDir, "refs", "heads"), 0755); err != nil {
		return err
	}

	if err := os.Mkdir(joinPath(dotDir, "refs", "tags"), 0755); err != nil {
		return err
	}
	return nil
}

func Status() error {
	branch, err := getCurrentBranchName()
	if err != nil {
		return err
	}
	fmt.Printf("On branch %s\n", branch)
	fmt.Println()
	fmt.Printf("No commits yet")
	fmt.Println()
	return nil
}

func getCurrentBranchName() (string, error) {
	headFile, err := os.Open(joinPath(dotDir, "HEAD"))
	if err != nil {
		return "", err
	}
	defer headFile.Close()

	scanner := bufio.NewScanner(headFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ref: ") {
			return strings.TrimPrefix(line, "ref: refs/heads/"), nil
		}
	}

	return "", nil
}

type indexRecord struct {
	permission string
	sha1       string
	fileType   string
	filename   string
}

func updateIndex(sha1Str string, filename string) error {
	dotDir, err := getDotDir()
	if err != nil {
		return err
	}
	index := joinPath(dotDir, "index")

	if !fileExists(index) {
		indexFile, err := os.Create(index)
		if err != nil {
			return err
		}
		err = updateIndexFile(indexFile, filename, sha1Str)
		if err != nil {
			return err
		}
		defer indexFile.Close()
	} else {
		indexFile, err := os.OpenFile(index, os.O_RDWR, 0644)
		if err != nil {
			return err
		}

		err = updateIndexFile(indexFile, filename, sha1Str)
		if err != nil {
			return err
		}
		defer indexFile.Close()
	}

	return nil
}

func updateIndexFile(file *os.File, filename string, sha1Str string) error {
	// 逐行读取文件内容，并更新 SHA-1 编码值
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		idx := indexRecord{permission: fields[0], sha1: fields[1], fileType: fields[2], filename: fields[3]}
		if idx.filename == filename && idx.sha1 != sha1Str {
			fmt.Fprintf(file, "%s %s %s %s\n", idx.permission, idx.sha1, idx.fileType, idx.filename)
			return nil
		}
	}
	return nil
}

func Add(files []string) error {
	for _, src := range files {
		if !fileExists(src) {
			fmt.Printf("file %s not exists.\n", src)
			continue
		}

		sha1Str := Sha1Hash("blob", src)
		dir, err := getObjectsDir()
		if err != nil {
			fmt.Println("error in get objects folder.")
			return err
		}

		dst := joinPath(dir, sha1Str[0:2], sha1Str[2:])
		fmt.Println(dst)
		if fileExists(dst) {
			// duplicate add
			continue
		}

		err = copyFile(src, dst)
		if err != nil {
			fmt.Println("error in copy file.")
			return err
		}
	}
	return nil
}

func getDotDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for dir != "" {
		dotDir := joinPath(dir, dotDir)
		if _, err := os.Stat(dotDir); !os.IsNotExist(err) {
			return dotDir, nil
		}
		dir = filepath.Dir(dir)
		// fmt.Println("dir: " + dir)
		// fmt.Println("filepath.Dir(dir)：" + filepath.Dir(dir))
		if dir == filepath.Dir(dir) {
			// fmt.Println(".ggit folder not found")
			return "", errors.New(".ggit folder not found")
		}
	}
	// fmt.Println(".ggit folder not found")
	return "", errors.New(".ggit folder not found")
}

func getObjectsDir() (string, error) {
	dir, err := getDotDir()
	if err != nil {
		return "", err
	}

	return joinPath(dir, "objects"), nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 确保目标目录存在
	err = os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return err
	}

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}

func joinPath(parts ...string) string {
	return filepath.Join(parts...)
}
