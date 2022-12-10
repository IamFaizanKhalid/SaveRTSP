package main

import (
	"fmt"
	"github.com/IamFaizanKhalid/SaveRTSP/ffmpeg"
	"github.com/IamFaizanKhalid/pointer"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	OutputPath string   `yaml:"output_path"`
	Streams    []Stream `yaml:"streams"`
}

type Stream struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"stream_url"`
	Split int    `yaml:"split"`
}

func Run(cfg Config) error {
	if len(cfg.Streams) == 0 {
		return fmt.Errorf("no stream provided in the config")
	}

	names := make(map[string]bool)
	for _, stream := range cfg.Streams {
		if _, ok := names[stream.Name]; ok {
			return fmt.Errorf("repeating stream name: %w", stream.Name)
		}
		names[stream.Name] = true
	}

	pwd, _ := os.Getwd()

	if !filepath.IsAbs(cfg.OutputPath) {
		cfg.OutputPath = filepath.Join(pwd, cfg.OutputPath)
	}

	err := os.MkdirAll(cfg.OutputPath, 0755)
	if err != nil {
		return fmt.Errorf("unable to create the output directory: %w", err)
	}

	var wg sync.WaitGroup

	for _, stream := range cfg.Streams {
		wg.Add(1)
		go func(stream Stream) {
			err := saveStream(cfg.OutputPath, stream)
			if err != nil {
				log.Printf("%s: %s", stream.Name, err)
			}
			wg.Done()
		}(stream)
	}

	wg.Wait()

	return nil
}

func saveStream(outDir string, stream Stream) error {
	outDir += "/" + stream.Name
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		return fmt.Errorf("unable to create the output directory: %w", err)
	}

	fmt.Printf(`
%s
--------
Input Stream: %s
Output Directory: %s
Video Length: %d minutes

`,
		stream.Name, stream.Url, outDir, stream.Split,
	)

	return ffmpeg.New().
		Input(stream.Url).
		Output(fmt.Sprintf("%s/%s-%%Y%%m%%d-%%H%%M.mp4", outDir, stream.Name)).
		Options(&ffmpeg.Options{
			RTSPTransport:        pointer.String("tcp"),
			NativeFramerateInput: pointer.Bool(true),
			VideoCodec:           pointer.String("h264"),
			AudioCodec:           pointer.String("aac"),
			MapStreamId:          pointer.Int(0),
			Format:               pointer.String("segment"),
			SegmentAtClockTime:   pointer.Int(1),
			SegmentTime:          pointer.Int(stream.Split * 60),
			SegmentFormat:        pointer.String("mp4"),
			SegmentNameByTime:    pointer.Int(1),
		}).
		Start()
}
