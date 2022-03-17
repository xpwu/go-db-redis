package rediscache

import (
  "fmt"
  "github.com/go-redis/redis"
  "net"
  "sync"
  "time"
)

var (
  clients  = make(map[Config]*redis.Client)
  clientMu sync.RWMutex
)

/**

一、底层redis Client的逻辑：
  1、一个 Client 对应一个 Options
  2、一个 Client 有一个 ConnPool, 一个 ConnPool 可能有多个连接，所有的连接都是指向同一个服务器
  3、没一个连接建立后，都会自动执行(如果有必要) Auth、SelectDB (不能选为0)、OnlyRead 等操作
  4、一个 Client 的命令执行时都是从ConnPool 获取一个Conn，并独占这个连接直到命令结束，然后返回Conn给ConnPool

二、这里的思考
  1、不应该每次都使用底层的NewClient，不然对同一个服务，没有共用到底层的 connpool
  2、Client生成后，没有去修改Options的接口，所以一旦Client生成了，就没法修改参数了。

三、这里的设计
  基于以上分析，这里基于 Config 缓存 Client，

*/

func Get(conf Config) *redis.Client {

  clientMu.RLock()
  client, ok := clients[conf]
  clientMu.RUnlock()
  if ok {
    return client
  }

  clientMu.Lock()
  defer clientMu.Unlock()

  timeout := time.Duration(conf.TimeoutMs) * time.Millisecond
  client = redis.NewClient(&redis.Options{
    Addr:         net.JoinHostPort(conf.Host, fmt.Sprintf("%d", conf.Port)),
    DB:           conf.DBNo,
    DialTimeout:  timeout,
    ReadTimeout:  timeout,
    WriteTimeout: timeout,
    MaxRetries:   3,
  })

  clients[conf] = client

  return client
}
