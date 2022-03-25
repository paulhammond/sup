package remote

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"

	"github.com/paulhammond/sup/internal/object"
	"go.etcd.io/bbolt"
)

type fake struct {
	db *bbolt.DB
}

func (r *fake) Close() error {
	return r.db.Close()
}

func (r *fake) Set() (object.Set, error) {
	set := object.Set{}

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("values"))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			set[string(k)] = fakeObject{db: r.db, path: string(k)}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return set, nil
}

func (r *fake) Upload(set object.Set) error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		for path, o := range set {
			err := r.uploadFile(tx, path, o)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *fake) uploadFile(tx *bbolt.Tx, path string, o object.Object) error {
	valueBucket := tx.Bucket([]byte("values"))
	metadataBucket := tx.Bucket([]byte("metadata"))

	err := o.Open(func(r io.Reader) error {
		// upload file
		val, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		err = valueBucket.Put([]byte(path), val)
		if err != nil {
			return err
		}

		// upload metadata
		m, err := o.Metadata()
		if err != nil {
			return err
		}
		if m.ContentType != nil {
			err = metadataBucket.Put([]byte(path+":contenttype"), []byte(*m.ContentType))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

type fakeObject struct {
	db   *bbolt.DB
	path string
}

func (o fakeObject) get() ([]byte, error) {
	var value []byte
	err := o.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("values"))
		v := b.Get([]byte(o.path))
		if v == nil {
			return errors.New("tktk")
		}
		value = make([]byte, len(v))
		copy(value, v)
		return nil
	})
	return value, err
}

func (o fakeObject) getMetadataValue(key string) *string {
	var value *string
	err := o.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("metadata"))
		v := b.Get([]byte(o.path + ":" + key))
		if v == nil {
			value = nil
			return nil
		}
		tmp := make([]byte, len(v))
		copy(tmp, v)
		tmp2 := string(tmp)
		value = &tmp2
		return nil
	})
	if err != nil {
		panic("impossible error")
	}
	return value
}

func (o fakeObject) Hash() (string, error) {
	v, err := o.get()
	if err != nil {
		return "", err
	}
	h := md5.Sum(v)
	return fmt.Sprintf("%d%x", len(v), h[:]), nil
}

func (o fakeObject) Metadata() (*object.Metadata, error) {
	metadata := object.Metadata{
		ContentType: o.getMetadataValue("contenttype"),
	}
	return &metadata, nil
}

func (o fakeObject) Open(fnc func(io.Reader) error) error {
	val, err := o.get()
	if err != nil {
		return err
	}
	r := bytes.NewReader(val)
	return fnc(r)
}

func openFake(path string) (Remote, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("values"))
		if b == nil {
			return errors.New("fake remote not initialized")
		}
		b = tx.Bucket([]byte("metadata"))
		if b == nil {
			return errors.New("fake remote not initialized")
		}
		return nil
	})

	return &fake{db}, err
}

func CreateFake(path string) error {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		valueBucket := tx.Bucket([]byte("values"))
		var err error
		if valueBucket != nil {
			return errors.New("fake remote already initialized")
		}
		metadataBucket := tx.Bucket([]byte("metadata"))
		if metadataBucket != nil {
			return errors.New("fake remote already initialized")
		}
		valueBucket, err = tx.CreateBucket([]byte("values"))
		if err != nil {
			return err
		}
		metadataBucket, err = tx.CreateBucket([]byte("metadata"))
		if err != nil {
			return err
		}
		err = valueBucket.Put([]byte("a.txt"), []byte("42"))
		if err != nil {
			return err
		}
		err = metadataBucket.Put([]byte("a.txt:contenttype"), []byte("text/plain"))
		if err != nil {
			return err
		}
		err = valueBucket.Put([]byte("b.txt"), []byte("b\n"))
		if err != nil {
			return err
		}
		err = metadataBucket.Put([]byte("b.txt:contenttype"), []byte("text/plain"))
		if err != nil {
			return err
		}
		err = valueBucket.Put([]byte("d.txt"), []byte("d\n"))
		if err != nil {
			return err
		}
		err = metadataBucket.Put([]byte("d.txt:contenttype"), []byte("text/plain"))
		if err != nil {
			return err
		}

		return nil

	})

	db.Close()

	return err
}
