package ffmpeg

import (
	"fmt"
	"os/exec"
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
	if f.input == nil {
		return fmt.Errorf("no input provided")
	}
	inputArg := fmt.Sprintf("-i %s", *f.input)

	outputArg := "output.mp4"
	if f.output != nil {
		outputArg = *f.output
	}

	if f.options == nil {
		f.options = DefaultOptions()
	}

	args := []string{inputArg}
	args = append(args, f.options.GetStrArguments()...)
	args = append(args, outputArg)

	cmd := exec.Command("ffmpeg", args...)

	return cmd.Run()
}
