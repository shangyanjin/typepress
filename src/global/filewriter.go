package global

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func NewFileWriter(filename string) (io.WriteCloser, error) {
	path := filepath.Base(filename)
	fw := &filewriter{path: path, filename: filename}
	err := fw.rotate()
	if err != nil {
		return nil, err
	}
	OnShutDown(fw.close)
	return fw, nil
}

type filewriter struct {
	sync.Mutex
	path     string
	filename string
	fd       *os.File
}

func (w *filewriter) rotate() (err error) {
	w.fd, err = os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	return
}

func (w *filewriter) Write(b []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	if w.fd == nil {
		return 0, errors.New("filewriter.fd is nil")
	}
	return w.fd.Write(b)
}

func (w *filewriter) close() {
	w.Close()
}

func (w *filewriter) Close() (err error) {
	if w.fd == nil {
		return nil
	}
	w.Lock()
	defer w.Unlock()
	err = w.fd.Close()
	w.fd = nil
	return
}
