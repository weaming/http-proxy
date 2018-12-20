# 如何写一个 HTTP 代理

![](https://www.ibm.com/support/knowledgecenter/ssw_ibm_i_73/rzaie/rzal8502.gif)

![](https://www.ibm.com/support/knowledgecenter/ssw_ibm_i_73/rzaie/rzal8504.gif)

* 只有服务端的透明代理
  * 问题
    * 如何做到重写 http 域名
* CS架构的代理
  * 客户端和服务器协作，完成流量的传输。整个“服务器、客户端配对”是对外透明的。
  * 问题
    * 如何复用客户端到服务器之间的连接
  * 核心
    * 如何把进来的请求，以什么格式转化为出去的请求
    * 如何把进来的响应，以什么格式转化为出去的响应
  * 请求内容
    * Schema
    * Host
    * Protocol Version
    * Headers
    * Body
  * 响应内容
    * Status Code
    * Headers
    * Body
  * 转发可能会丢失的头部
    * Host
    * IP
    * Proto (http or https)
  * [代理服务器应该添加的头部](https://tools.ietf.org/html/rfc7239#section-4)
    * Forwarded
      * 标准头部
      * 格式：`Forwarded: by=<identifier>; for=<identifier>; host=<host>; proto=<http|https>`
    * X-Forwarded-By
    * X-Forwarded-For
    * X-Forwarded-Host
    * X-Forwarded-Proto
    * Via
* 代理类型
  * forward proxies
    * implicit
      * 使用自签证书
    * explict
      * must be specifically configured within the client application to shuttle the request
      * by establishing an `HTTP CONNECT` tunnel
        * authorize the client saying `HTTP/1.1 200 OK`
        * client will start streaming TCP packets which are routed through to the remote/upstream server specified in the `CONNECT` verb
        * `CONNECT` is a hop-by-hop method.
        * > If you're writing a proxy server, all you need to do for allowing your clients to connect to HTTPS servers is read in the CONNECT request, make a connection from the proxy to the end server (given in the CONNECT request), send the client with a 200 OK reply and then forward everything that you read from the client to the server, and vice versa.
  * reverse proxies
* 其他代理软件
  * Squid
  * mitmproxy
  * tinyproxy
  * Apache Traffic Server
  * privoxy
  * nginx
    * [How to use Nginx as a HTTP/HTTPS proxy server? - Server Fault](https://serverfault.com/questions/298392/how-to-use-nginx-as-a-http-https-proxy-server)
    * If you want to use an HTTP/HTTPS proxy, you should use Squid. It was written to do exactly that. Nginx was written to act as a reverse proxy and load balancer, but not a forward proxy.
* 其他特性
  * 支持更多协议
    * socks5
    * socks4
    * http
    * https
  * pac
  * proxy authentication
  * bypass list
  * include origin IP
* References
  * [GitHub - smartystreets/cproxy: A simple, explicit forward proxy written in Go to facilitate HTTP CONNECT Tunneling.](https://github.com/smartystreets/cproxy)
  * [Proxy servers and tunneling | MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Proxy_servers_and_tunneling)
  * [apache - CONNECT request to a forward HTTP proxy over an SSL connection? - Stack Overflow](https://stackoverflow.com/questions/6594604/connect-request-to-a-forward-http-proxy-over-an-ssl-connection)
  * [Debugging problems with the network proxy - The Chromium Projects](https://www.chromium.org/developers/design-documents/network-stack/debugging-net-proxy)
