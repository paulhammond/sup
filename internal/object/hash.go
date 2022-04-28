package object

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

type Hash struct {
	Size     int64
	PartSize int64
	Hash     string
}

func (h Hash) String() string {
	return fmt.Sprintf("%d.%d.%s", h.Size, h.PartSize, h.Hash)
}

func GenerateHash(o Object) (*Hash, error) {
	var h Hash
	err := o.Open(func(r io.Reader) error {
		m := md5.New()
		size, err := io.Copy(m, r)
		if err != nil {
			return err
		}
		h = Hash{size, 0, hex.EncodeToString(m.Sum(nil))}
		return nil
	})
	return &h, err
}

func matchHash(o1, o2 Object) (bool, error) {

	h1, err := o1.Hash()
	if err != nil {
		return false, err
	}
	h2, err := o2.Hash()
	if err != nil {
		return false, err
	}
	if h2.PartSize != 0 {
		panic("unimplemented")
	}

	if h1.Size != h2.Size {
		return false, nil
	}

	if h1.PartSize == 0 {
		return h1.Hash == h2.Hash, nil
	}

	var hash2 string
	err = o2.Open(func(r io.Reader) error {
		var err error
		// this is an implementation of the AWS S3 multipart upload hashing
		// algorithm, documented at
		// https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html#large-object-checksums

		partCount := int(h1.Size / h1.PartSize)
		i := 0
		h2 := md5.New()
		for err == nil {
			i++
			h := md5.New()
			_, err = io.CopyN(h, r, h1.PartSize)
			if err == nil || err == io.EOF {
				sum := h.Sum(nil)

				h2.Write(sum)
			}
			// don't loop forever
			if err == nil && i > partCount+1 {
				err = io.EOF
			}
		}
		if err != io.EOF {
			return err
		}
		hash2 = fmt.Sprintf("%x-%d", h2.Sum(nil), i)

		return nil
	})

	if err != nil {
		return false, err
	}

	return h1.Hash == hash2, nil

}
