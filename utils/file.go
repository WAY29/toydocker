package utils

import (
	"archive/tar"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FileHash(filepath string) (string, error) {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	sum := hash.Sum(nil)

	return fmt.Sprintf("%Xn", sum), nil
}

// 系统命令实现的tar
func Tar(dirPath, outputPath string) error {
	_, err := exec.Command("tar", "-cf", outputPath, "-C", dirPath, ".").CombinedOutput()
	return err
}

// 系统命令实现的untar
func Untar(tarFilePath, decompressionPath string) error {
	_, err := exec.Command("tar", "-xf", tarFilePath, "-C", decompressionPath).CombinedOutput()
	return err
}

// go实现的untar，可能存在bug
func BadUntar(tarFilePath, decompressionPath string) error {
	file, err := os.Open(tarFilePath)

	if err != nil {
		return err
	}

	tarBallReader := tar.NewReader(file)
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		filepath := path.Join(decompressionPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			err = os.MkdirAll(filepath, os.FileMode(header.Mode)) // or use 0755 if you prefer

			if err != nil {
				return err
			}
		case tar.TypeLink:
			// handle link file
			writer, err := os.Create(filepath)

			if err != nil {
				return err
			}

			reader, err := os.Open(path.Join(decompressionPath, header.Linkname))
			if err != nil {
				return err
			}

			io.Copy(writer, reader)

			err = os.Chmod(filepath, os.FileMode(header.Mode))

			if err != nil {
				return err
			}

			writer.Close()
		case tar.TypeReg:
			// handle normal file
			writer, err := os.Create(filepath)

			if err != nil {
				return err
			}

			io.Copy(writer, tarBallReader)

			err = os.Chmod(filepath, os.FileMode(header.Mode))

			if err != nil {
				return err
			}

			writer.Close()
		default:
			return fmt.Errorf("Unknown file[%s] type: %c", filepath, header.Typeflag)
		}
	}

	return nil
}
