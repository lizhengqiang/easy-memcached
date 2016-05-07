// Redis风格的Memcached客户端
package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

// 支持前缀的Memcached客户端
type EasyMemcached struct {
	Prefix string
	Client *memcache.Client
}

func (m *EasyMemcached) Set(key string, value []byte, expire int32) error {
	err := m.Client.Set(&memcache.Item{
		Key:        m.Prefix + key,
		Value:      value,
		Expiration: expire,
	})
	return err
}

func (m *EasyMemcached) Add(key string, value []byte, expire int32) error {
	err := m.Client.Add(&memcache.Item{
		Key:        m.Prefix + key,
		Value:     value,
		Expiration: expire,
	})
	return err
}

func (m *EasyMemcached) Get(key string) (value []byte, err error) {
	item, err := m.Client.Get(m.Prefix + key)

	if err != nil {
		return
	}
	value = item.Value
	return

}

func (m *EasyMemcached) Del(key string) (err error) {
	return m.Client.Delete(m.Prefix + key)
}

func (m *EasyMemcached) Rename(source string, target string) error {
	item, errGet := m.Client.Get(m.Prefix + source)

	if errGet != nil {
		return errGet
	}

	item.Key = target

	errSet := m.Client.Set(item)
	if errSet != nil {
		return errSet
	}
	errDel := m.Client.Delete(source)

	if errDel != nil {
		return errDel
	}
	return nil
}

func (m *EasyMemcached) Ping() error {
	errSet := m.Client.Set(&memcache.Item{
		Key:        "foo",
		Value:      []byte("bar"),
		Expiration: 3600,
	})
	if errSet != nil {
		return errSet
	}
	if _, errGet := m.Client.Get("foo"); errGet != nil {
		return errGet
	}
	return nil
}

// 分布式锁
// 广博写的, 不要喷我!
type DistributeLock struct {
	Memcached *EasyMemcached
	Item      *memcache.Item
	Ttl       int32
}

func (lock *DistributeLock) CheckLock(ch chan bool) {
	RETRY:
	if err := lock.Memcached.Add(lock.Item.Key, lock.Item.Value, lock.Item.Expiration); err != nil {
		time.Sleep(time.Duration(lock.Ttl) * time.Second)
		goto RETRY
	}
	ch <- true
}

func (m *EasyMemcached) GetLock(key string, value []byte, expiration, ttl int32) (lock *DistributeLock) {
	lock = &DistributeLock{
		Memcached: m,
		Item: &memcache.Item{
			Key:        key,
			Value:      value,
			Expiration: expiration,
		},
		Ttl: ttl,
	}
	return lock
}

func (lock *DistributeLock) Lock() {
	ch := make(chan bool)
	go lock.CheckLock(ch)
	_ = <-ch
}

func (lock *DistributeLock) Unlock() (err error) {
	err = lock.Memcached.Del(lock.Item.Key)
	return
}
