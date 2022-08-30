## 兼容 Redis-Cluster 和 Redis-Standalone 的实现方案

使用到的包

> https://github.com/go-redis/redis

## 测试方法

先启动redis 测试环境

使用 docker composer 来启动一个redis-cluster集群和 redis-stanalone 单机版

```shell
https://github.com/go-redis/redis
```

其中，`redis` 版本号可以更换至 `7.0`，也可以调整到 `6.2`

```yaml
image: docker.io/bitnami/redis-cluster:7.0
```

