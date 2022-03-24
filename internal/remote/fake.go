package remote

import (
	"errors"

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
		return b.Put([]byte("a.txt"), []byte("42"))
	})

	return err
}
