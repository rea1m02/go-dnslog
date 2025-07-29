package dns

import (
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/spf13/viper"

	"github.com/rea1m/go-dnslog/database"
	"github.com/rea1m/go-dnslog/models"
)

var (
	dnsDomain string
	serverIP  string
	ns1Domain string
	ns2Domain string
	logQueue  = make(chan *models.DNSLog, 1000)
	wg        sync.WaitGroup
)

// Init 初始化DNS服务器配置
func Init() {
	dnsDomain = viper.GetString("dns.domain")
	serverIP = viper.GetString("dns.server_ip")
	ns1Domain = viper.GetString("dns.ns1")
	ns2Domain = viper.GetString("dns.ns2")

	// 启动日志处理协程
	go processLogs()
}

// Start 启动DNS服务器
func Start() error {
	port := viper.GetInt("dns.port")
	handler := dns.HandlerFunc(handleDNSRequest)

	// 启动UDP服务器
	udpServer := &dns.Server{
		Addr:    net.JoinHostPort("0.0.0.0", strconv.Itoa(port)),
		Net:     "udp",
		Handler: handler,
	}

	// 启动TCP服务器
	tcpServer := &dns.Server{
		Addr:    net.JoinHostPort("0.0.0.0", strconv.Itoa(port)),
		Net:     "tcp",
		Handler: handler,
	}

	log.Printf("Starting DNS server on port %d (UDP)", port)

	// 使用goroutine启动服务器，避免阻塞
	go func() {
		if err := udpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start DNS(UDP) server: %v", err)
		}
	}()

	go func() {
		if err := tcpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start DNS(TCP) server: %v", err)
		}
	}()

	return nil
}

// handleDNSRequest 处理DNS查询请求
func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	// 权威配置
	msg.Authoritative = true
	msg.Compress = false

	clientIP, _, _ := net.SplitHostPort(w.RemoteAddr().String())

	for _, q := range r.Question {
		// 处理不同类型的DNS查询
		switch q.Qtype {
		case dns.TypeA:
			// 处理A记录查询
			if err := handleAQuery(msg, q, clientIP); err != nil {
				log.Printf("Failed to handle A query: %v", err)
			}
		case dns.TypeAAAA:
			// 忽略AAAA记录查询
			continue
		case dns.TypeNS:
			// 处理NS记录查询
			if err := handleNSQuery(msg, q); err != nil {
				log.Printf("Failed to handle NS query: %v", err)
			}
		default:
			// 其他类型查询返回空响应
			continue
		}
	}

	// 发送DNS响应
	_ = w.WriteMsg(msg)
}

// handleAQuery 处理A记录查询
func handleAQuery(msg *dns.Msg, q dns.Question, clientIP string) error {
	qName := strings.ToLower(q.Name)
	baseDomain := dnsDomain
	if !strings.HasSuffix(baseDomain, ".") {
		baseDomain += "."
	}

	// 检查是否为负责的域名
	if !strings.HasSuffix(qName, baseDomain) {
		msg.SetRcode(msg, dns.RcodeNameError)
		return nil
	}

	// 提取用户子域名
	userDomain, subName := extractUserDomain(qName, baseDomain)

	// 处理DNS Rebind功能
	if strings.Contains(qName, ".e.") {
		// 这里应该查询对应域名绑定的两个ip，并随机返回两个ip
		ip := rebindIP(qName)
		if ip == "" {
			msg.SetRcode(msg, dns.RcodeNameError)
			return nil
		}
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 1},
			A:   net.ParseIP(ip),
		})

		// 记录DNS日志
		// 重绑定功能的话，貌似不记录日志也可以
		// logDNSQuery(userDomain, clientIP, qName, "A", subName)
		return nil
	}

	// 正常返回服务器IP
	msg.Answer = append(msg.Answer, &dns.A{
		Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 300},
		A:   net.ParseIP(serverIP),
	})

	// 记录DNS日志
	logDNSQuery(userDomain, clientIP, qName, "A", subName)

	return nil
}

// handleNSQuery 处理NS记录查询
func handleNSQuery(msg *dns.Msg, q dns.Question) error {
	qName := strings.ToLower(q.Name)
	baseDomain := dnsDomain + "."

	if !strings.HasSuffix(qName, baseDomain) {
		msg.SetRcode(msg, dns.RcodeNameError)
		return nil
	}

	// 添加NS记录响应
	msg.Answer = append(msg.Answer, &dns.NS{
		Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 3600},
		Ns:  ns1Domain + ".",
	})
	msg.Answer = append(msg.Answer, &dns.NS{
		Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 3600},
		Ns:  ns2Domain + ".",
	})

	return nil
}

// extractUserDomain 从查询域名中提取用户域名和子域名
func extractUserDomain(qName, baseDomain string) (userDomain, subName string) {
	// 移除末尾的点
	qName = strings.TrimSuffix(qName, ".")
	baseDomain = strings.TrimSuffix(baseDomain, ".")

	// 提取用户域名部分
	prefix := strings.TrimSuffix(qName, baseDomain)
	prefix = strings.TrimSuffix(prefix, ".")

	parts := strings.Split(prefix, ".")
	if len(parts) > 0 {
		userDomain = parts[len(parts)-1]
		if len(parts) > 1 {
			subName = strings.Join(parts[:len(parts)-1], ".")
		}
	}

	return userDomain, subName
}

// generateRebindIP 生成DNS Rebind攻击用的随机IP
func rebindIP(qName string) string {
	// 去除末尾的点
	qName = strings.TrimSuffix(qName, ".")
	var rebind models.Rebind
	database.DB.Model(&models.Rebind{}).Where("domain = ?", qName).First(&rebind)
	if rebind.ID == 0 {
		return ""
	}
	if time.Now().UnixNano()%2 == 0 {
		return rebind.FirstIP
	} else {
		return rebind.SecondIP
	}
	
}

// logDNSQuery 将DNS查询记录添加到日志队列
func logDNSQuery(userDomain, clientIP, host, queryType, subName string) {
	// 查询用户
	var user models.User
	
	if err := database.DB.Where("user_domain = ?", userDomain).First(&user).Error; err != nil {
		log.Println("User not found for domain:", userDomain)
		return
	}

	host = strings.TrimSuffix(host, ".")

	// 创建DNS日志记录
	dnsLog := &models.DNSLog{
		UserID:  user.ID,
		Host:    host,
		SubName: subName,
		Type:    queryType,
		IP:      clientIP,
		City:    "", // 预留IP地理位置字段
	}

	// 添加到日志队列
	select {
	case logQueue <- dnsLog:
	default:
		log.Println("Log queue is full, dropping log entry")
	}
}

// processLogs 处理日志队列，将日志写入数据库
func processLogs() {
	for entry := range logQueue {
		wg.Add(1)
		go func(dnsLog *models.DNSLog) {
			defer wg.Done()
			// 使用事务保存日志
			if err := database.DB.Create(dnsLog).Error; err != nil {
				// 使用全局的log包输出错误信息
				log.Println("Failed to save DNS log:", err)
			}
		}(entry)
	}
}

// Shutdown 优雅关闭DNS服务器
func Shutdown() {
	close(logQueue)
	wg.Wait()
	log.Println("DNS server shutdown successfully")
}
