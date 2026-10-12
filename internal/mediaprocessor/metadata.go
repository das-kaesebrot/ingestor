package mediaprocessor

import "time"

type MediaMetadata struct {
	// all immutable
	Width            int
	Height           int
	Codec            string
	OriginalFilename string
	ContentType      string
	CaptureTime      time.Time

	// mutable
	CaptureCorrectionOffset time.Duration
}

type ImageMetadata struct {
	MediaMetadata
}

type VideoMetadata struct {
	MediaMetadata
	Duration   time.Duration
	Container  string
	AudioCodec string
}
