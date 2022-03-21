package main

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"log"

	"github.com/andybalholm/brotli"
	"github.com/gabriel-vasile/mimetype"

	"go.sancus.dev/web/embed"
)

const (
	ReadChunkSize = 1024
)

// CloserFunc converts a close handle into an io.Closer
type CloserFunc func() error

func (f CloserFunc) Close() error {
	return f()
}

type BlobContent struct {
	io.Writer
	io.Closer

	Encoding string
	buffer   *bytes.Buffer
	Bytes    []byte
}

func (bc *BlobContent) finish() {
	if bc.buffer != nil {
		bc.Bytes = bc.buffer.Bytes()
		bc.buffer = nil
	}
}

type BlobDigest struct {
	hash.Hash

	Name string
}

func (bd BlobDigest) String() string {
	hash := base64.RawStdEncoding.EncodeToString(bd.Hash.Sum([]byte{}))
	return fmt.Sprintf("%s=%s", bd.Name, hash)
}

type Blob struct {
	ContentType string

	Digest  []*BlobDigest
	Content []*BlobContent
	writers []io.Writer
}

func (blob *Blob) Export(fi fs.FileInfo) (*embed.File, error) {

	out := &embed.File{
		Name:        fi.Name(),
		Size:        fi.Size(),
		ModTime:     fi.ModTime().UTC(),
		ContentType: blob.ContentType,
	}

	for _, bd := range blob.Digest {
		out.Digest = append(out.Digest, bd.String())
	}

	for _, bc := range blob.Content {
		if len(blob.Content) == 1 || bc.Encoding != "" {
			content := embed.Content{
				Encoding: bc.Encoding,
				Bytes:    bc.Bytes,
			}
			out.Content = append(out.Content, content)
		}
	}

	return out, nil
}

func (blob *Blob) ReadFrom(r io.Reader) (int64, error) {
	var total int

	// read loop
	b := make([]byte, ReadChunkSize)
	w := io.MultiWriter(blob.writers...)

	for {
		n, err := r.Read(b)
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}

		w.Write(b[:n])
		total += n
	}

	// close encoders
	for i := range blob.Content {
		blob.Content[i].Close()
	}

	return int64(total), nil
}

func (blob *Blob) addDigest(name string, hasher hash.Hash) {
	c := &BlobDigest{
		Name: name,
		Hash: hasher,
	}

	blob.Digest = append(blob.Digest, c)
	blob.writers = append(blob.writers, c)
}

func (blob *Blob) addEncoding(name string, wrap func(io.Writer) io.WriteCloser) {

	b := &bytes.Buffer{}
	c := &BlobContent{
		Encoding: name,
		buffer:   b,
	}

	if wrap != nil {
		w := wrap(b)
		c.Writer = w
		c.Closer = CloserFunc(func() error {
			err := w.Close()
			c.finish()
			return err
		})
	} else {
		c.Writer = b
		c.Closer = CloserFunc(func() error {
			c.finish()
			blob.ContentType = mimetype.Detect(c.Bytes).String()
			return nil
		})
	}

	blob.Content = append(blob.Content, c)
	blob.writers = append(blob.writers, c)
}

func NewBlob() *Blob {
	blob := &Blob{}

	// digest
	blob.addDigest("sha-256", sha256.New())
	blob.addDigest("sha", sha1.New())
	blob.addDigest("md5", md5.New())

	// content-encoding
	blob.addEncoding("", nil)
	blob.addEncoding("br", func(dst io.Writer) io.WriteCloser {
		return brotli.NewWriterLevel(dst, brotli.BestCompression)
	})
	blob.addEncoding("gzip", func(dst io.Writer) io.WriteCloser {
		w, err := gzip.NewWriterLevel(dst, gzip.BestCompression)
		if err != nil {
			log.Fatal(err)
		}
		return w
	})

	return blob
}
