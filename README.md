# dcache

dcache 是一个高性能的存储服务器,致力于分布式,每个 key 只会存在于一台dcache中,

2019.8.10 添加多路复用单链接性能提升三倍
2019.8.15 添加限制最大连接数

分布式群集完成，待测试

> 压力测试 

在开启debug日志的时候性能会大幅度降低,生产环境禁止开启


 2C 2G的服务器
### set
```$xslt

./client -h=10.10.175.145:7777  -c=256 -n=200000
2.463920 seconds total
rps is 81171.467496
throughput is 81.171467 MB/s

```

### get
```$xslt

./client -h=10.10.175.145:7777  -c=256 -n=200000 -t=get
2.568034 seconds total
rps is 77880.592175
throughput is 77.880592 MB/s
```