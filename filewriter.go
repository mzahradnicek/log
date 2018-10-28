package log

import (
	"os"
	"time"
)

type FileWriter struct {
	name string
	hour int
	file *os.File
}

func (f *FileWriter) Write(p []byte) (n int, err error) {
	if f.hour != time.Now().Hour() {
		if err := f.Open(f.name); err != nil {
			return 0, err
		}
	}

	if _, err := f.file.Write(p); err != nil {
		return 0, err
	}

	// add newline
	return f.file.Write([]byte{10})
}

func (f *FileWriter) Open(name string) error {
	if f.file != nil {
		f.file.Close()
	}

	t := time.Now()

	var err error
	f.file, err = os.OpenFile(name+"-"+t.Format("2006-01-02-15")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	f.hour = t.Hour()
	f.name = name

	return nil
}

func (f *FileWriter) Close() error {
	err := f.file.Close()
	f.file = nil
	return err
}

func NewFileWriter(name string) (*FileWriter, error) {
	fw := &FileWriter{}
	if err := fw.Open(name); err != nil {
		return nil, err
	}

	return fw, nil
}
