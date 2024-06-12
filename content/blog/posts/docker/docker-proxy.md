---
title: "DockerHub 加速镜像部署 - 使用cloudflare 代理"
date: 2024-06-12
draft: false
authors: "jobcher"
featuredImage: "/images/docker.png"
featuredImagePreview: "/images/docker.png"
images: ['/images/docker.png']
tags: ["docker"]
categories: ["docker"]
series: ["docker进阶系列"]
---
## 背景
6 月 6 日，上海交大的 Docker Hub 镜像加速器宣布因收到通知要求被下架。声明称：“即时起我们将中止对 dockerhub 仓库的镜像。docker 相关工具默认会自动处理失效镜像的回退，如果对官方源有访问困难问题，建议尝试使用其他仍在服务的镜像源。我们对给您带来的不便表示歉意，感谢您的理解与支持。”Docker Hub 是目前最大的容器镜像社区，去年 5 月起国内用户报告 Docker Hub 官网无法访问，其网址解析返回了错误 IP 地址。  
![sjtug](/images/dockerhub-sjtug.jpeg)  
因为不能直接访问国外的镜像仓库，下载国外的docker镜像速度一直很慢, 国内从 Docker Hub 拉取镜像有时会遇到困难，此时可以配置镜像加速器。  
## 使用
我这边已经部署好了加速镜像节点，同学们如果不想自己部署，可以使用我的加速节点，但是，不能保证节点长期有效。  
```sh
https://dockerhub.jobcher.com
```
### 第一步：代理拉取镜像
假如我们下载node镜像，那么我们可以这样写：
```sh
docker pull dockerhub.jobcher.com/library/node:latest
```
### 第二步：重命名镜像
```sh
docker tag dockerhub.jobcher.com/library/node:latest node:latest
```
### 第三步：删除代理镜像
```sh
docker rmi dockerhub.jobcher.com/library/node:latest
```
### 或者直接配置到镜像仓库
```sh
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
    "registry-mirrors": [
        "https://dockerhub.jobcher.com",
    ]
}
EOF
```
重新加载docker
```sh
sudo systemctl daemon-reload
sudo systemctl restart docker
```
## 配置 cloudflare 代理
我们通过 cloudflare 的全球 CDN 节点，以 Workers 中转代理的方式来访问国外镜像仓库，从而加速镜像的下载。  
在 Cloudflare Workers 部署一个 Worker 时，它会在30秒之内部署到 Cloudflare 的整个边缘网络，全世界95个国家/200个城市节点。域中的每个请求都会由离用户更近地点的 Worker 来处理，基于此来实现代码的 “随处运行”。  
### 1.创建 Worker
进入`Cloudflare Workers`，点击 `Create Worker`  
![create-worker](/images/docker-proxy-1.png)  
### 2.部署 Worker
填写你的环境参数
>const hub_host = 'registry-1.docker.io'
>const auth_url = 'https://auth.docker.io'
填写你的worker page url地址
>const workers_url = 'https://dockerhub.jobcher.com'  
  
粘贴一下到cloudeflare的worker.js中
```js
'use strict'

const hub_host = 'registry-1.docker.io'
const auth_url = 'https://auth.docker.io'
const workers_url = 'https://dockerhub.jobcher.com' //填写你自己的域名地址，不要填写默认的
/**
 * static files (404.html, sw.js, conf.js)
 */

/** @type {RequestInit} */
const PREFLIGHT_INIT = {
    status: 204,
    headers: new Headers({
        'access-control-allow-origin': '*',
        'access-control-allow-methods': 'GET,POST,PUT,PATCH,TRACE,DELETE,HEAD,OPTIONS',
        'access-control-max-age': '1728000',
    }),
}

/**
 * @param {any} body
 * @param {number} status
 * @param {Object<string, string>} headers
 */
function makeRes(body, status = 200, headers = {}) {
    headers['access-control-allow-origin'] = '*'
    return new Response(body, {status, headers})
}


/**
 * @param {string} urlStr
 */
function newUrl(urlStr) {
    try {
        return new URL(urlStr)
    } catch (err) {
        return null
    }
}


addEventListener('fetch', e => {
    const ret = fetchHandler(e)
        .catch(err => makeRes('cfworker error:\n' + err.stack, 502))
    e.respondWith(ret)
})


/**
 * @param {FetchEvent} e
 */
async function fetchHandler(e) {
  const getReqHeader = (key) => e.request.headers.get(key);

  let url = new URL(e.request.url);

  if (url.pathname === '/token') {
      let token_parameter = {
        headers: {
        'Host': 'auth.docker.io',
        'User-Agent': getReqHeader("User-Agent"),
        'Accept': getReqHeader("Accept"),
        'Accept-Language': getReqHeader("Accept-Language"),
        'Accept-Encoding': getReqHeader("Accept-Encoding"),
        'Connection': 'keep-alive',
        'Cache-Control': 'max-age=0'
        }
      };
      let token_url = auth_url + url.pathname + url.search
      return fetch(new Request(token_url, e.request), token_parameter)
  }

  url.hostname = hub_host;
  
  let parameter = {
    headers: {
      'Host': hub_host,
      'User-Agent': getReqHeader("User-Agent"),
      'Accept': getReqHeader("Accept"),
      'Accept-Language': getReqHeader("Accept-Language"),
      'Accept-Encoding': getReqHeader("Accept-Encoding"),
      'Connection': 'keep-alive',
      'Cache-Control': 'max-age=0'
    },
    cacheTtl: 3600
  };

  if (e.request.headers.has("Authorization")) {
    parameter.headers.Authorization = getReqHeader("Authorization");
  }

  let original_response = await fetch(new Request(url, e.request), parameter)
  let original_response_clone = original_response.clone();
  let original_text = original_response_clone.body;
  let response_headers = original_response.headers;
  let new_response_headers = new Headers(response_headers);
  let status = original_response.status;

  if (new_response_headers.get("Www-Authenticate")) {
    let auth = new_response_headers.get("Www-Authenticate");
    let re = new RegExp(auth_url, 'g');
    new_response_headers.set("Www-Authenticate", response_headers.get("Www-Authenticate").replace(re, workers_url));
  }

  if (new_response_headers.get("Location")) {
    return httpHandler(e.request, new_response_headers.get("Location"))
  }

  let response = new Response(original_text, {
            status,
            headers: new_response_headers
        })
  return response;
  
}


/**
 * @param {Request} req
 * @param {string} pathname
 */
function httpHandler(req, pathname) {
    const reqHdrRaw = req.headers

    // preflight
    if (req.method === 'OPTIONS' &&
        reqHdrRaw.has('access-control-request-headers')
    ) {
        return new Response(null, PREFLIGHT_INIT)
    }

    let rawLen = ''

    const reqHdrNew = new Headers(reqHdrRaw)

    const refer = reqHdrNew.get('referer')

    let urlStr = pathname
    
    const urlObj = newUrl(urlStr)

    /** @type {RequestInit} */
    const reqInit = {
        method: req.method,
        headers: reqHdrNew,
        redirect: 'follow',
        body: req.body
    }
    return proxy(urlObj, reqInit, rawLen, 0)
}


/**
 *
 * @param {URL} urlObj
 * @param {RequestInit} reqInit
 */
async function proxy(urlObj, reqInit, rawLen) {
    const res = await fetch(urlObj.href, reqInit)
    const resHdrOld = res.headers
    const resHdrNew = new Headers(resHdrOld)

    // verify
    if (rawLen) {
        const newLen = resHdrOld.get('content-length') || ''
        const badLen = (rawLen !== newLen)

        if (badLen) {
            return makeRes(res.body, 400, {
                '--error': `bad len: ${newLen}, except: ${rawLen}`,
                'access-control-expose-headers': '--error',
            })
        }
    }
    const status = res.status
    resHdrNew.set('access-control-expose-headers', '*')
    resHdrNew.set('access-control-allow-origin', '*')
    resHdrNew.set('Cache-Control', 'max-age=1500')
    
    resHdrNew.delete('content-security-policy')
    resHdrNew.delete('content-security-policy-report-only')
    resHdrNew.delete('clear-site-data')

    return new Response(res.body, {
        status,
        headers: resHdrNew
    })
}

```
### 3. 确认修改
确认修改后，点击 `Deploy` 按钮，等待部署完成。之后就可以正常使用你的Docker Hub 加速了。
