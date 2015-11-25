package audru

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Date(2015, 11, 24, 0, 42, 3, 3, time.UTC)
	//startTime = time.Now()
}

type WriterManager struct {
	tmpDir    string
	startTime *time.Time
	pid       string
	counter   int
	lock      *sync.RWMutex
}

func (piper *WriterManager) Close() error {
	return os.Stdout.Close()
}

//func NewWriterManagerLong(n int, args string, dir string) (*WriterManager, error) {

func NewWriterManager(concurrency int, dir string) (*WriterManager, error) {
	err := writeHeader(concurrency)
	if err != nil {
		return nil, err
	}

	newWriterManager := new(WriterManager)
	newWriterManager.counter = 0
	newWriterManager.lock = new(sync.RWMutex)
	newWriterManager.pid = strconv.Itoa(os.Getpid())
	if dir == "" {
		newWriterManager.tmpDir = os.TempDir()
	} else {
		newWriterManager.tmpDir = dir
	}
	newWriterManager.startTime = &startTime
	return newWriterManager, nil
}

func (wm *WriterManager) newWriter(prefix string, suffix string) (*Writer, error) {
	wm.lock.Lock()
	newWriter := Writer{
		file:          nil,
		writer:        nil,
		writerManager: wm,
		prefix:        prefix,
		suffix:        suffix,
		id:            wm.counter,
	}
	wm.counter = wm.counter + 1
	wm.lock.Unlock()
	return &newWriter, nil
}

func (wm *WriterManager) NewWriterPre(prefix string, suffix string) (*Writer, error) {
	return wm.newWriter(prefix, suffix)
}

const DEFAULT_PREFIX = ".fifo"
const DEFAULT_SUFFIX = ".named"

func (wm *WriterManager) NewWriter() (*Writer, error) {
	return wm.newWriter(DEFAULT_PREFIX, DEFAULT_SUFFIX)
}

type Writer struct {
	file          *os.File
	writer        *bufio.Writer
	writerManager *WriterManager
	prefix        string
	suffix        string
	id            int
}

func (w *Writer) init() error {
	fifo, err := w.writerManager.newNamedPipe(w.id, w.prefix, w.suffix)
	if err != nil {
		return err
	}

	w.file = fifo
	w.writer = bufio.NewWriter(fifo)

	log.Println(w.file.Name())
	//_, _ = os.Stdout.Write([]byte(newWriter.file.Name() + "\n"))

	return nil
}

func (w *Writer) Write(p []byte) (int, error) {
	if w.file == nil {
		err := w.init()
		if err != nil {
			return 0, nil
		}
	}
	return w.Write(p)
}

func (w *Writer) Flush() error {
	return w.Flush()
}

func (w *Writer) Close() error {
	err := w.Flush()
	if err != nil {
		return err
	}
	err = w.file.Close()
	if err != nil {
		return err
	}
	return nil

}

func (wm *WriterManager) newNamedPipe(id int, prefix string, suffix string) (*os.File, error) {
	//filename := wm.tmpDir + "/" + prefix + "_" + wm.pid + "_" + strconv.Itoa(wm.startTime.Nanosecond()) + suffix
	filename := wm.tmpDir + "/" + prefix + "_" + wm.pid + "_" + strconv.Itoa(id) + "_" + strconv.FormatInt(time.Since(*wm.startTime).Nanoseconds(), 10) + suffix

	err := syscall.Mkfifo(filename, 0600)
	if err != nil {
		return nil, err
	}
	writeComment("")
	writeNamedPipeName(filename)

	file, err := os.OpenFile(filename, os.O_WRONLY, os.ModeNamedPipe)
	_, _ = os.Stdout.Write([]byte("NewFifo \n"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	_, _ = os.Stdout.Write([]byte("NewFifo \n"))
	return file, nil
}

func writeHeader(concurrency int) error {
	return write(strconv.Itoa(concurrency))
}

func writeComment(s string) error {
	return write("# " + s)
}

func writeNamedPipeName(s string) error {
	return write(s)
}

func write(s string) error {
	_, err := os.Stdout.Write([]byte(s + "\n"))
	return err
}
