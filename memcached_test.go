package memcached

import (
	"testing"
)
import (
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

var client *memcache.Client = memcache.New("127.0.0.1:11211")

var easyClient *EasyMemcached = &EasyMemcached{
	Client: client,
	Prefix: "/v1/",
}

func Test_Conn(t *testing.T) {
	var err error
	err = easyClient.Ping()
	if err != nil {
		t.Fatal("Memcached Conn Failed", err)
		t.Fail()
	}
}

func Test_Client(t *testing.T) {

	var err error

	err = easyClient.Set("foo", []byte("bar"), 1)

	if err != nil {
		t.Error(err)
		t.Fail()
	}
	value, err := easyClient.Get("foo")

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log("Get Passed", string(value))

	time.Sleep(2 * time.Second)

	value, err = easyClient.Get("foo")

	if err == memcache.ErrCacheMiss {
		t.Log("Expired Passed", string(value))
	} else if err != nil {
		t.Error(err)
		t.Fail()
	}

}

func Benchmark_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		easyClient.Client.Set(&memcache.Item{
			Key:        "foo",
			Value:      []byte("bar"),
			Expiration: 3600,
		})
	}
}

func Benchmark_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		easyClient.Client.Get("foo")
	}
}

func Benchmark_EasySet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		easyClient.Set("foo", []byte("bar"), 3600)
	}
}

func Benchmark_EasyGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		easyClient.Get("foo")
	}
}
