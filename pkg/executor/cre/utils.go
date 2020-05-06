package cre

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/utils"
	"github.com/sirupsen/logrus"
)

// logger is the logging helper for the executor
type logger struct {
	// pipelineName contains the name of the pipeline we are working on.
	pipelineName string

	// pipelineID contains the run ID of the pipeline in context.
	pipelineID string

	// task contains the name of the task
	task string

	step string
}

// GetLogWriter return io.WriteCloser for the command logs
func (l *logger) getLogWriter() io.WriteCloser {
	taskDir := fmt.Sprintf("%s/%s/%s/%s", defaults.AgentLogsDir, l.pipelineName, l.pipelineID, l.task)

	filename := fmt.Sprintf("%s/%s", taskDir, l.step)
	if !utils.DirExists(taskDir) {
		if err := os.MkdirAll(taskDir, os.ModePerm); err != nil {
			return nil
		}
	}

	var (
		f   *os.File
		err error
	)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		f, err = os.Create(filename)
		if err != nil {
			logrus.Errorf("error while creating log file: %s", err)
			return nil
		}
	} else {
		f, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			logrus.Errorf("error while opening log file: %s", err)
			return nil
		}
	}

	return &loggerWriter{f, bufio.NewWriter(f)}
}

func newLogger(pName, pID, tName, sName string) *logger {
	return &logger{
		pName,
		pID,
		tName,
		sName,
	}
}

type loggerWriter struct {
	f *os.File

	*bufio.Writer
}

func (w *loggerWriter) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	return w.f.Close()
}

func parseImageCanonicalURL(image string) string {
	if strings.Contains(image, "/") {
		return fmt.Sprintf("docker.io/%s", image)
	}

	return fmt.Sprintf("docker.io/library/%s", image)
}
