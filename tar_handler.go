// https://github.com/kubernetes/kubectl/blob/master/pkg/cmd/cp/cp.go

package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnTarAll(reader io.Reader, destDir string, prefix string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		if strings.Contains(header.Name, "..") {
			return fmt.Errorf("file may be performing arbitary write")
		}
		mode := header.FileInfo().Mode()
		destFileName := filepath.Join(destDir, header.Name[len(prefix):])
		baseName := filepath.Dir(destFileName)
		if err := os.MkdirAll(baseName, 0750); err != nil {
			return err
		}
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(destFileName, 0750); err != nil {
				return err
			}
			continue
		}

		evalPath, err := filepath.EvalSymlinks(baseName)
		if err != nil {
			return err
		}

		if mode&os.ModeSymlink != 0 {
			linkname := header.Linkname
			if !filepath.IsAbs(linkname) {
				_ = filepath.Join(evalPath, linkname)
			}
			_, err := filepath.EvalSymlinks(filepath.Join(linkname, destFileName))
			if err != nil {
				return fmt.Errorf("symlinks not in folder")
			}
			if err := os.Symlink(linkname, destFileName); err != nil {
				return err
			}
		} else {
			outFile, err := os.Create(destFileName)
			if err != nil {
				return err
			}
			for {
				if _, err := io.CopyN(outFile, tarReader, 1024*1024*4); err != nil {
					if err == io.EOF {
						break
					} else {
						return err
					}
				}
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
