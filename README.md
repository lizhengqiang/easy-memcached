# easy-memcached

```golang
type EasyMemcached struct {
    Prefix string           // 前缀, 如果和ThinkPHP这种国产框架共享缓存的话,可能会需要这个东西
    Client *memcache.Client // github.com/bradfitz/gomemcache/memcache 这个客户端就可以
}
```

## 示例代码

```golang
func main(){
    client := memcache.New("127.0.0.1:11211")
    
    easyClient := &easy.EasyMemcached{
        Client: client,
        Prefix: "/v1/",
    }
    
    if err := easyClient.Ping(); err != nil {
        log.Fatal("Memcached Conn Failed", err.Error())
    }
    
    easyClient.Set("foo", string(mpJson), 3600)
    
    
}
```