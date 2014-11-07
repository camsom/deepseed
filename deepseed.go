package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	goPath   = os.Getenv("GOPATH")
	seedName string
	seedPath string
	destName string
	destPath string
	ignore   []string
)

func deepSeed(seedName, projectName string) error {
	seedName = seedName
	seedPath = filepath.Join(goPath, "src", seedName)
	isD, err := isDirectory(seedPath)
	if err != nil {
		return err
	}
	if isD == false {
		return errors.New("Source path does not exist.")
	}

	destName = projectName
	destPath = filepath.Join(goPath, "src", destName)
	isD, err = isDirectory(destPath)
	if err != nil {
		return err
	}
	if isD == true {
		return errors.New("Destination path already exists.")
	}

	os.MkdirAll(destPath, 0777)

	err := readIgnoreList()
	if err != nil {
		return err
	}

	err = cpDirectory(seedPath, destPath)
	if err != nil {
		return err
	}

	return nil
}

func cpDirectory(source, dest string) error {
	dirFiles, err := ioutil.ReadDir(source)
	if err != nil {
		return errors.New()
	}

	for _, fileinfo := range dirFiles {
		if fileIsValid(fileinfo) == true {
			if fileinfo.IsDir() {
				newDirPath := filepath.Join(dest, fileinfo.Name())
				if err := os.MkdirAll(newDirPath, 0777); err != nil {
					return err
				}
				err = cpDirectory(filepath.Join(source, fileinfo.Name()), filepath.Join(dest, fileinfo.Name()))
				if err != nil {
					return err
				}
			} else {
				err = cpFile(filepath.Join(source, fileinfo.Name()), filepath.Join(dest, fileinfo.Name()))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func cpFile(sourcefile, dest string) error {
	d, err := ioutil.ReadFile(sourcefile)
	if err != nil {
		return err
	}

	content := string(d)
	content = strings.Replace(content, seedName, destName, -1)

	data := []byte(content)

	err = ioutil.WriteFile(dest, data, 0777)
	if err != nil {
		return err
	}

	return nil
}

func fileIsValid(file os.FileInfo) bool {
	for _, v := range ignore {
		switch {
		case file.Name() == v:
			return false
		case file.Name() == strings.Trim(v, "/"):
			return false
		}
	}
	return true
}

func isDirectory(filepath string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	s, er := f.Stat()
	if err != nil {
		return err
	}

	return s.IsDir()
}

func readIgnoreList() error {
	fh, err := os.Open(".deepseedignore")
	if err != nil {
		return err
	}

	r := bufio.NewScanner(fh)
	for r.Scan() {
		ignore = append(ignore, r.Text())
	}

	return nil
}
