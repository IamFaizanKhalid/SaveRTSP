package ffmpeg

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func New() *ffmpeg {
	return &ffmpeg{}
}

type ffmpeg struct {
	input   *string
	output  *string
	options *Options
}

func (f *ffmpeg) Input(inputPath string) *ffmpeg {
	f.input = &inputPath
	return f
}

func (f *ffmpeg) Output(outputPath string) *ffmpeg {
	f.output = &outputPath
	return f
}

func (f *ffmpeg) Options(options *Options) *ffmpeg {
	f.options = options
	return f
}

func (f *ffmpeg) Start() error {
	if f.input ==
		nil {
		if f.options.InputStream == nil {
			return fmt.Errorf("no input provided")
		}
	} else {
		f.options.InputStream = f.input
	}

	output := "output.mp4"
	if f.output != nil {
		output = *f.output
	}

	if f.options == nil {
		f.options = DefaultOptions()
	}

	args := f.options.GetStrArguments()
	args = append(args, output)

	cmd := exec.Command("ffmpeg", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errLines := strings.Split(stderr.String(), "\n")

		var errStr string
		if len(errLines) > 11 {
			errStr = strings.Join(errLines[11:], "\n")
		}

		return fmt.Errorf("%s:\n%s", err, errStr)
	}

	return nil
}
