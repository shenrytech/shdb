package shdb

import (
	"bytes"
	"net/url"

	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	BUCKET_OBJ = []byte("obj")
	db         *bbolt.DB
	log        *zap.Logger
)

func Init(logger *zap.Logger, dbFile string) {

	db, err := bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bbolt.Tx) error {
		tx.CreateBucketIfNotExists(BUCKET_OBJ)
		return nil
	})
	log = logger
}

func Close() {
	if err := db.Close(); err != nil {
		log.Error("error closing database", zap.Error(err))
	} else {
		log.Debug("closed database", zap.String("dbFile", db.Path()))
	}
}

func Put[T IObject](val ...T) error {
	if len(val) == 0 {
		return nil
	}
	for _, v := range val {
		v.GetMetadata().UpdatedAt = timestamppb.Now()
	}
	kv, err := Marshal(val...)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BUCKET_OBJ)
		for _, v := range kv {
			err = b.Put(v.Key(), v.Value)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func Get[T IObject](typeKey TypeKey, id []byte) (T, error) {
	var t T
	kv := KeyVal{DbKey: *NewDbKey(typeKey, id)}
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BUCKET_OBJ)
		kv.Value = b.Get(kv.Key())
		if kv.Value == nil {
			return ErrNotFound
		}
		var err error
		t, err = Unmarshal[T](kv)
		return err
	})
	return t, err
}

func Delete[T IObject](typeKey TypeKey, id []byte) (T, error) {
	val, err := Get[T](typeKey, id)
	if err != nil {
		return val, err
	}
	dbk := NewDbKey(typeKey, id)
	return val, db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BUCKET_OBJ)
		return b.Delete(dbk.Key())
	})

}
func GetAllKV(typeKey TypeKey) ([]KeyVal, error) {
	allKvs := []KeyVal{}
	return allKvs, db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket(BUCKET_OBJ).Cursor()
		for k, v := c.Seek(typeKey[:]); k != nil && bytes.HasPrefix(k, typeKey[:]); k, v = c.Next() {
			kv := KeyVal{DbKey: *MarshalDbKey(k), Value: v}
			allKvs = append(allKvs, kv)
		}
		return nil
	})
}

func GetAll[T IObject](typeKey TypeKey) ([]T, error) {
	allKvs, err := GetAllKV(typeKey)
	if err != nil {
		return nil, err
	}
	return UnmarshalMany[T](allKvs)
}

func List[T IObject](typeKey TypeKey, pageSize int32, pageToken string) (res []T, nextPageToken string, err error) {
	last := &DbKey{}
	firstIdx := 0
	lastIdx := 0
	if pageToken != "" {
		ptVal, err := url.ParseQuery(pageToken)
		if err != nil {
			return nil, "", err
		}
		last, err = DbKeyFromString(ptVal.Get("last"))
		if err != nil {
			return nil, "", err
		}
	}

	allKvs, err := GetAllKV(typeKey)
	if err != nil {
		return nil, "", err
	}
	if pageToken != "" {
	findLastLoop:
		for k, v := range allKvs {
			if v.Equal(last) {
				firstIdx = k + 1
				lastIdx = k + 1 + int(pageSize)
				break findLastLoop
			}
		}
		return nil, "", nil // All items returned
	}
	if lastIdx > len(allKvs) {
		lastIdx = len(allKvs)
	}
	ret := allKvs[firstIdx:lastIdx]
	ptVal := url.Values{}
	ptVal.Set("last", ret[len(ret)-1].String())
	nextPageToken = ptVal.Encode()
	ts, err := UnmarshalMany[T](ret)
	return ts, nextPageToken, err
}
