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
	Cameras    []Camera `yaml:"cameras"`
}

type Camera struct {
	Name      string `yaml:"name"`
	StreamUrl string `yaml:"stream_url"`
	Split     int    `yaml:"split"`
}

func Run(cfg Config) error {
	if len(cfg.Cameras) == 0 {
		return fmt.Errorf("no camera provided in the config")
	}

	names := make(map[string]bool)
	for _, camera := range cfg.Cameras {
		if _, ok := names[camera.Name]; ok {
			return fmt.Errorf("repeating camera name: %w", camera.Name)
		}
		names[camera.Name] = true
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

	for _, camera := range cfg.Cameras {
		wg.Add(1)
		go func(camera Camera) {
			err := saveStream(cfg.OutputPath, camera)
			if err != nil {
				log.Printf("%s: %s", camera.Name, err)
			}
			wg.Done()
		}(camera)
	}

	wg.Wait()

	return nil
}

func saveStream(outDir string, camera Camera) error {
	outDir += "/" + camera.Name
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
		camera.Name, camera.StreamUrl, outDir, camera.Split,
	)

	return ffmpeg.New().
		Input(camera.StreamUrl).
		Output(fmt.Sprintf("%s/%s-%%Y%%m%%d-%%H%%M.mp4", outDir, camera.Name)).
		Options(&ffmpeg.Options{
			RTSPTransport:        pointer.String("tcp"),
			NativeFramerateInput: pointer.Bool(true),
			VideoCodec:           pointer.String("h264"),
			AudioCodec:           pointer.String("aac"),
			MapStreamId:          pointer.Int(0),
			Format:               pointer.String("segment"),
			SegmentAtClockTime:   pointer.Int(1),
			SegmentTime:          pointer.Int(camera.Split * 60),
			SegmentFormat:        pointer.String("mp4"),
			SegmentNameByTime:    pointer.Int(1),
		}).
		Start()
}
