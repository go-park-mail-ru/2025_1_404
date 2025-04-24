package s3

import "io"

type Upload struct {
	Bucket      string
	File        io.Reader
	Filename    string
	Size        int64
	ContentType string
}