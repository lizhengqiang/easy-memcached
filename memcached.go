package easy_memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
)
type EasyMemcached struct {
	Prefix string
	Client *memcache.Client
}

func (m *EasyMemcached) Set(key, value string, ttl int32) error {
	err := m.Client.Set(&memcache.Item{
		Key:        m.Prefix + key,
		Value:      []byte(value),
		Expiration: ttl,
	})
	return err
}

func (m *EasyMemcached) Get(key string) (value string, err error) {
	item, err := m.Client.Get(m.Prefix + key)

	if err == nil {
		value = string(item.Value)
	}
	return

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
