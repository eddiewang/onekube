package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func CopyKubeConfig(newKubeConfigPath string) error {
	kubeconfigPath, err := kubeconfigPath()
	if err != nil {
		log.Fatal(err)
	}

	_, err = copyFile(newKubeConfigPath, kubeconfigPath)

	return err
}

func BackupNonOneKubeConfig() error {
	kubeconfigPath, err := kubeconfigPath()
	if err != nil {
		log.Fatal(err)
	}

	backupFilePath := fmt.Sprintf("%s-onekube-backup", kubeconfigPath)

	file, err := os.Open(kubeconfigPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		if line == "# Managed by onekube" {
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(backupFilePath); err == nil {
		log.Fatalf("tried to backup current unmanaged kubeconfig, but file at %s already exists", backupFilePath)
	}

	_, err = copyFile(kubeconfigPath, backupFilePath)
	return err
}

func kubeconfigPath() (string, error) {
	// KUBECONFIG env var
	if v := os.Getenv("KUBECONFIG"); v != "" {
		list := filepath.SplitList(v)
		if len(list) > 1 {
			// TODO KUBECONFIG=file1:file2 currently not supported
			return "", errors.New("multiple files in KUBECONFIG are currently not supported")
		}
		return v, nil
	}

	// default path
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("HOME environment variable not set")
	}
	return filepath.Join(home, ".kube", "config"), nil
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
