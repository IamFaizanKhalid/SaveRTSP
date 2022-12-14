package ffmpeg

import (
	"fmt"
	"reflect"

	"github.com/IamFaizanKhalid/pointer"
)

// Options defines allowed FFmpeg arguments
type Options struct {
	RTSPTransport        *string `flag:"-rtsp_transport"`      // protocol to use to capture input streams
	NativeFramerateInput *bool   `flag:"-re"`                  // read input at native frame rate
	InputStream          *string `flag:"-i"`                   // input file path
	VideoCodec           *string `flag:"-vcodec"`              // force video codec ('copy' to copy stream)
	AudioCodec           *string `flag:"-acodec"`              // force audio codec ('copy' to copy stream)
	MapStreamId          *int    `flag:"-map"`                 // set stream mapping from input streams to output streams. Just enumerate the input streams in the order you want them in the output
	Format               *string `flag:"-f"`                   // force format
	SegmentAtClockTime   *int    `flag:"-segment_atclocktime"` // if set to "1" split at regular clock time intervals starting from 00:00 o’clock
	SegmentTime          *int    `flag:"-segment_time"`        // set segment duration (in seconds)
	SegmentFormat        *string `flag:"-segment_format"`      // force format for the segments
	SegmentNameByTime    *int    `flag:"-strftime"`            // if set to "1" segments will be named by time of file creation
}

func DefaultOptions() *Options {
	return &Options{
		RTSPTransport:        pointer.String("tcp"),
		NativeFramerateInput: pointer.Bool(true),
		VideoCodec:           pointer.String("h264"),
		AudioCodec:           pointer.String("aac"),
	}
}

// GetStrArguments ...
func (opts Options) GetStrArguments() []string {
	f := reflect.TypeOf(opts)
	v := reflect.ValueOf(opts)

	values := []string{}

	for i := 0; i < f.NumField(); i++ {
		flag := f.Field(i).Tag.Get("flag")
		value := v.Field(i).Interface()

		if !v.Field(i).IsNil() {

			if _, ok := value.(*bool); ok {
				values = append(values, flag)
			}

			if vs, ok := value.(*string); ok {
				values = append(values, flag, *vs)
			}

			if va, ok := value.([]string); ok {

				for i := 0; i < len(va); i++ {
					item := va[i]
					values = append(values, flag, item)
				}
			}

			if vm, ok := value.(map[string]string); ok {
				for k, v := range vm {
					values = append(values, flag, fmt.Sprintf("%v:%v", k, v))
				}
			}

			if vi, ok := value.(*int); ok {
				values = append(values, flag, fmt.Sprintf("%d", *vi))
			}

		}
	}

	return values
}
