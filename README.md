## Keepalived-exporter
`keepalived-exporter` 用来监控 keepalived 进程及当前 VIP 状态。

## 背景
prometheus 支持多样的 exporter，并且开发者也可以定制自己的 exporter，在项目中发现有需要对 keepalived 进程及状态的监控，由于在网上没有找到对应或者
不错的监控，因此做了个对应的监控。

## 设计思路
参考了社区部分 exporter 的实现逻辑，使用自建 collector 的方式。

## 前置需求
需求 kubernetes 集群安装 prometheus

## 安装方法
支持 kubernetes 安装
```
cd <your-proj-path>/keepalived-exporter && docker build -t keepalived-exporter:<Tag-you-prefered> .
```

## 使用方法
Prometheus 页面查询语句支持查找 `keepalived_status_<running,sleeping,waiting,zombie,other>`, `keepalived_vip_ready`