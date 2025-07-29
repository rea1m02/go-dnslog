# DNSlog 平台
## 前言

本项目参考[eyes.sh](https://github.com/lijiejie/eyes.sh)项目，由Golang开发。

## 项目功能
- ✅ DNS查询日志记录
- ✅ DNS Rebind
- ✅ 自定义域名与IP绑定
- ✅ 支持日志导出


## 部署需求
- 公网ip
- 公网域名（无需备案，有条件可以准备2个）


## 环境要求
- 后端：golang
- 前端：node+vue
- 数据库：mysql
- 端口：53（DNS）

## 部署指南

目前项目仅支持手动部署

1. 后端部署
+ 有go环境
```bash
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 本地编译
go build -o dnslog
# 跨平台编译
## linux/amd64
GOOS=linux GOARCH=amd64 go build -o dnslog-linux-amd64
## linux/arm64
GOOS=linux GOARCH=arm64 go build -o dnslog-linux-arm64
```

+ 无go环境
直接下载对应release的二进制可执行文件

2. 前端部署
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 编译
npm run build
## 编译完成后，将dist目录上传到服务器
```

## 配置说明
核心配置文件    `backend/config.yaml`
```yaml
dns:
    domain: dnslog.example.com
    ns1: ns1.dnslog.example.com
    ns2: ns2.dnslog.example.com
    server_ip: 127.0.0.1
    port: 53
```
### 后端服务配置
```service
[Unit]
Description=DNSLog Platform
After=network.target mysql.service

[Service]
User=root
WorkingDirectory=/your/path/backend/
ExecStart=/your/path/go-dnslog
Restart=always

[Install]
WantedBy=multi-user.target
```
设置后台运行
```bash
systemctl daemon-reload

# 假设服务文件名为 dnslog.service
# 启动服务
systemctl start dnslog.service

# 开机自启
systemctl enable dnslog.service
```

### nginx配置参考
```conf
	server {
		listen 80;
		server_name your-ip/your-domain;
		root /your/path/frontend/dist;
		index index.html;

		location / {
			try_files $uri $uri/ /index.html;
		}
		
		location /api {
			proxy_pass http://127.0.0.1:8081;  # 后端服务地址
			proxy_set_header Host $host;
		}
	}
```


## 测试
正常情况下，在配置好之后，由于dns传播特性，需要24-48小时才能生效，所以可以通过在本地自行测试是否部署成功
```bash
# 测试域名解析
dig @localhost dnslog.example.com
```
可以测试是否正常解析，是否成功存储到数据库

## 常见问题
- **Q: 为什么DNS服务启动失败?**
  A: 确保端口53未被系统DNS服务占用，可使用`lsof -i:53`检查，若存在53端口占用，在关闭对应的服务后，该系统可能存在无法正常解析域名的情况，需要在`/etc/resolv.conf`文件中添加`nameserver 8.8.8.8`

