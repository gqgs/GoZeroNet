package site

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/bigfile"
	"github.com/gqgs/go-zeronet/pkg/lib/crypto"
)

func hash(reader io.Reader, dir, relativePath string) (map[string]File, error) {
	var size int
	var hashes []string
	files := make(map[string]File)
	buf := make([]byte, config.PieceSize)
	for {
		n, err := reader.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		buf = buf[:n]
		hashes = append(hashes, crypto.Sha512_256(buf))
		size += n
	}

	if len(hashes) == 0 {
		return nil, errors.New("empty file")
	}

	if len(hashes) == 1 {
		files[relativePath] = File{
			Sha512: hashes[0],
			Size:   size,
		}
		return files, nil
	}

	piecemap, err := bigfile.MarshalPieceMap(filepath.Base(relativePath), hashes)
	if err != nil {
		return nil, err
	}
	piecemapRelativePath := relativePath + ".piecemap.msgpack"
	if err := os.WriteFile(path.Join(dir, piecemapRelativePath), piecemap, os.ModePerm); err != nil {
		return nil, err
	}

	files[relativePath] = File{
		Sha512:    bigfile.MerkleRoot(hashes),
		Size:      size,
		PieceSize: config.PieceSize,
		Piecemap:  piecemapRelativePath,
	}

	files[piecemapRelativePath] = File{
		Sha512: crypto.Sha512_256(piecemap),
		Size:   len(piecemap),
	}

	return files, nil
}
