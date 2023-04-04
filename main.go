package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
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
		Status()
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ggit add <filename>")
			os.Exit(1)
		}
		Add(os.Args[2:])
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
	dotDir, err := os.Getwd()
	// Check if .ggit folder already exists
	if err != nil {
		return fmt.Errorf(".ggit folder already exists")
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

func Status() {
	branch := getCurrentBranchName()
	fmt.Printf("On branch %s\n", branch)
	fmt.Println()
	fmt.Printf("No commits yet")
	fmt.Println()
}

func getCurrentBranchName() string {
	headFile, err := os.Open(joinPath(dotDir, "HEAD"))
	if err != nil {
		return ""
	}
	defer headFile.Close()

	scanner := bufio.NewScanner(headFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ref: ") {
			return strings.TrimPrefix(line, "ref: refs/heads/")
		}
	}

	return ""
}

func updateIndex(sha1Str string, filename string) {
	index := joinPath(getDotDir(), "index")
	if !fileExists(index) {
		indexFile, err := os.Create(index)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer indexFile.Close()
	} else {
		indexFile, err := os.Open(index)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer indexFile.Close()
	}

	_, err := ioutil.ReadFile(index)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// read from index file, check if file has been added before
	// if so, update sha-1
	// if not, add one

}

func Add(files []string) {
	for _, src := range files {
		if !fileExists(src) {
			fmt.Printf("file %s not exists.\n", src)
			continue
		}

		sha1Str := Sha1Hash("blob", src)
		dst := joinPath(getObjectsDir(), sha1Str[0:2], sha1Str[2:])
		fmt.Println(dst)
		if fileExists(dst) {
			// duplicate add
			continue
		}

		err := copyFile(src, dst)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func getDotDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for dir != "" {
		dotDir := joinPath(dir, dotDir)
		if _, err := os.Stat(dotDir); !os.IsNotExist(err) {
			return dotDir
		}
		dir = filepath.Dir(dir)
		fmt.Println("dir: " + dir)
		fmt.Println("filepath.Dir(dir)：" + filepath.Dir(dir))
		if dir == filepath.Dir(dir) {
			fmt.Println(".ggit folder not found")
			os.Exit(1)
		}
	}
	fmt.Println(".ggit folder not found")
	os.Exit(1)
	return ""
}

func getObjectsDir() string {
	dir := getDotDir()
	return joinPath(dir, "objects")
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
