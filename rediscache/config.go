package rediscache

type Config struct {
  Host      string `conf:"host"`
  Port      int    `conf:"port"`
  DBNo      int    `conf:"db_no"`
  TimeoutMs int    `conf:"timeout,unit:ms"`
}
