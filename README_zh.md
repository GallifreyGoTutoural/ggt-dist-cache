# ggt-dist-cache
[English](README.md) | [中文](README_zh.md)

ggt-dist-cache是模仿[groupcache](https://github.com/golang/groupcache)手写一个分布式缓存系统的项目，旨在更深入地了解groupcache的底层原理和设计模式，最终实现一个简易版本的分布式缓存系统——"gdc"。

项目的名称中的 "ggt" 是 "Gallifrey's GoTutoural" 的简写，"gdc"是"ggt-dist-cache"的简写。

项目的主要参考来源是是[极客兔兔](https://geektutu.com/)大佬的博客：[7天用Go从零实现分布式缓存GeeCache](https://geektutu.com/post/geecache.html)，如果想了解更多程序设计细节和考量，请查阅原博客。

## 开发计划

- [x] 实现 LRU 缓存淘汰算法
- [x] 利用 sync.Mutex 互斥锁，实现 LRU 缓存的并发控制
- [x] 实现核心数据结构 Group，缓存不存在时，调用回调函数获取源数据
- [x] 搭建 HTTP Server
- [x] 启动 HTTP Server 测试API
- [ ] 实现一致性哈希代码
- [ ] 注册节点，借助一致性哈希算法选择节点
- [ ] 使用 singleflight 防止缓存击穿
- [ ] 使用 protobuf 进行节点间通信，编码报文

