package remote

import (
	"crypto/md5"
	"errors"
	"fmt"

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

func (o fakeObject) Hash() (string, error) {
	v, err := o.get()
	if err != nil {
		return "", err
	}
	h := md5.Sum(v)
	return fmt.Sprintf("%d%x", len(v), h[:]), nil
}

func openFake(path string) (Remote, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("values"))
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
		b := tx.Bucket([]byte("values"))
		var err error
		if b != nil {
			return errors.New("fake remote already initialized")
		}
		b, err = tx.CreateBucket([]byte("values"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("a.txt"), []byte("42"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("b.txt"), []byte("b\n"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("d.txt"), []byte("d\n"))
		if err != nil {
			return err
		}

		return nil

	})

	db.Close()

	return err
}
