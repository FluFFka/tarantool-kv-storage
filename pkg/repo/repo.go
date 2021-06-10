package repo

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/tarantool/go-tarantool"
)

type Repository struct {
	Conn *tarantool.Connection
}

var (
	ErrWrongValue = errors.New("wrong value returned from database")
	ErrNoContent  = errors.New("no content database returned")
	ErrKeyFound   = errors.New("key found in database, can't insert")
)

func (r *Repository) GetByKey(key string) (string, error) {
	resp, err := r.Conn.Select("storage", "primary", 0, 100, tarantool.IterEq, []interface{}{key})
	if err != nil {
		return "", err
	}
	for _, item := range resp.Data {
		if reflect.TypeOf(item).Kind() == reflect.Slice {
			itemSlice := reflect.ValueOf(item)
			if itemSlice.Len() != 2 {
				return "", ErrWrongValue
			}
			return fmt.Sprintf("%v", itemSlice.Index(1)), nil
		}
		return "", ErrWrongValue
	}
	return "", ErrNoContent
}

func (r *Repository) InsertValue(key string, value string) error {
	resp, err := r.Conn.Insert("storage", []interface{}{key, value})
	if err != nil {
		if resp != nil {
			if resp.Code == tarantool.ErrTupleFound {
				return ErrKeyFound
			}
		}
		return err
	}
	return nil
}
