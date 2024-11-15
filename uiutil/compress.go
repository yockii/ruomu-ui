package uiutil

import (
	"bytes"
	"encoding/base64"
	"github.com/klauspost/compress/zstd"
	"io"
)

// CompressString 字符串压缩，并返回字符串
func CompressString(s string) (string, error) {
	return CompressBytes([]byte(s))
}

func CompressBytes(s []byte) (string, error) {
	in := bytes.NewReader(s)
	buf := new(bytes.Buffer)
	err := compress(in, buf)
	if err != nil {
		return "", err
	}
	// base64
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecompressString(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	decompressed, err := decompress(data)
	return decompressed, err
}

func DecompressBytes(s []byte) ([]byte, error) {
	return DecompressString(string(s))
}

func compress(in io.Reader, out io.Writer) error {
	enc, err := zstd.NewWriter(out)
	if err != nil {
		return err
	}
	_, err = io.Copy(enc, in)
	if err != nil {
		enc.Close()
		return err
	}
	return enc.Close()
}

var decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))

func decompress(src []byte) ([]byte, error) {
	return decoder.DecodeAll(src, nil)
}
