# ggt-dist-cache
[English](README.md) | [中文](README_zh.md)


ggt-dist-cache is a project that aims to handwrite a distributed cache system as an imitation of [groupcache](https://github.com/golang/groupcache). The goal is to gain a deeper understanding of the underlying principles and design patterns of groupcache and eventually implement a simplified version of a distributed cache system called "gdc".

The "ggt" in the project name stands for "Gallifrey's GoTutorial," and "gdc" is the abbreviation for "ggt-dist-cache."

The primary reference for this project is the blog post by GeekTutu: [Building a Distributed Cache in 7 Days with Go](https://geektutu.com/post/geecache.html). For more details and considerations regarding the program design, please refer to the original blog post.

## Development Plan

- [x] Implement LRU cache eviction algorithm.
- [ ] Implement concurrent control of LRU cache using sync.Mutex.
- [ ] Implement the core data structure, Group, which calls the callback function to obtain the source data when the cache is not present.
- [ ] Set up an HTTP Server.
- [ ] Start the HTTP Server and test the API.
- [ ] Implement consistent hashing.
- [ ] Register nodes and select nodes using the consistent hashing algorithm.
- [ ] Prevent cache breakdown using singleflight.
- [ ] Use protobuf for inter-node communication and message encoding.