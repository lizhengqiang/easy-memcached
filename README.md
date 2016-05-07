# easy-memcached

## Struct

```golang 
type EasyMemcached struct {
    Prefix string           // 前缀, 如果和ThinkPHP这种国产框架共享缓存的话,可能会需要这个东西
    Client *memcache.Client // github.com/bradfitz/gomemcache/memcache 这个客户端就可以
}
```

## Example

```golang
func main(){
    client := memcache.New("127.0.0.1:11211")
    
    easyClient := &easy.EasyMemcached{
        Client: client,
        Prefix: "/v1/",
    }
    
    var err error
    
    err = easyClient.Ping()
    if err != nil {
        log.Fatal("Memcached Conn Failed", err)
    }
    
    err = easyClient.Set("foo", []byte("bar"), 1)
    if err != nil {
        log.Println(err)
    }
    
    value, err := easyClient.Get("foo")
    if err != nil {
        log.Println(err)
    }
    
    time.Sleep(2 * time.Second)
    
    value, err = easyClient.Get("foo")
    if err == memcache.ErrCacheMiss {
        // cache miss
    } else if err != nil {
        log.Error(err)
    }
    
}
```