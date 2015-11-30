package cacheddownloader

import (
	"os"
	"runtime"
)

type CachedFile struct {
	*os.File

	ExpandedDirectoryPath string
	onClose               func(string)
}

func NewFileCloser(file *os.File, directory string, onClose func(string)) *CachedFile {
	fc := &CachedFile{
		File: file,
		ExpandedDirectoryPath: directory,
		onClose:               onClose,
	}

	runtime.SetFinalizer(fc, func(f *CachedFile) {
		f.Close()
	})

	return fc
}

func (fw *CachedFile) Close() error {
	err := fw.File.Close()
	if err != nil {
		return err
	}

	fw.onClose(fw.File.Name())
	runtime.SetFinalizer(fw, nil)

	return nil
}
