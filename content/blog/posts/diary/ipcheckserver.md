---
title: "搭建ip地址检索服务"
date: 2024-10-11
draft: false
featuredImage: "/images/ipcheck-1.png"
featuredImagePreview: "/images/ipcheck-1.png"
images: ["/images/ipcheck-1.png"]
authors: "jobcher"
tags: ["daliy"]
categories: ["日常"]
series: ["日常系列"]
---
## 背景
很多时候，我们需要查询一个IP地址，都得通过百度，谷歌，或者其他搜索引擎，非常麻烦。教大家一个使用cloudflare worker搭建一个只属于我们自己的ip地址检索服务。  
### 条件
- 需要一个cloudflare账号（自行注册，必选）
- 域名（自行购买，可选）

## 步骤
### 1. 注册cloudflare账号，并登录。  
![注册](/images/ipcheck-1.png)  
  
### 2. 在cloudflare的dashboard中，点击`workers`，点击`create a worker`。  
![创建worker](/images/ipcheck-2.png)  
![创建worker](/images/ipcheck-3.png)  
![创建worker](/images/ipcheck-4.png)  

### 3. 创建一个worker
![创建worker](/images/ipcheck-5.png)  
![创建worker](/images/ipcheck-6.png)  

### 4. 复制代码并粘贴到`worker.js`中。
```js
addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request));
});

/**
 * Handle the incoming request and return formatted IP information
 * @param {Request} request
 */
async function handleRequest(request) {
  const url = new URL(request.url);
  const queryIp = url.searchParams.get('ip');
  const clientIp = queryIp || request.headers.get('cf-connecting-ip');
  const path = url.pathname;

  // 获取 IP 信息
  const ipInfo = await getIpInfo(clientIp);

  // 根据路径选择返回格式
  if (path === '/table') {
    const tableFormat = formatAsTable(ipInfo);
    return new Response(tableFormat, {
      headers: { 'Content-Type': 'text/plain; charset=utf-8' },
      status: 200,
    });
  } else {
    return new Response(JSON.stringify(ipInfo), {
      headers: { 'Content-Type': 'application/json' },
      status: 200,
    });
  }
}

/**
 * Get IP information and format it according to the required structure
 * @param {string} ip
 */
async function getIpInfo(ip) {
  try {
    const response = await fetch(`http://ip-api.com/json/${ip}`);
    const data = await response.json();

    // Format the returned data
    return {
      ip: data.query || ip,
      city: data.city || "None",
      province: data.regionName || "None",
      country: data.country || "None",
      continent: data.continent || "None",
      isp: data.isp || "None",
      time_zone: data.timezone || "None",
      latitude: data.lat || 0,
      longitude: data.lon || 0,
      postal_code: data.zip || "None",
      iso_code: data.countryCode || "None",
      notice: "Let it be",
      provider: "Powered by Jobcher",
      blog: "https://www.jobcher.com",
      data_updatetime: new Date().toISOString().slice(0, 10).replace(/-/g, '')
    };
  } catch (error) {
    return {
      ip: ip,
      city: "None",
      province: "None",
      country: "None",
      continent: "None",
      isp: "None",
      time_zone: "None",
      latitude: 0,
      longitude: 0,
      postal_code: "None",
      iso_code: "None",
      notice: "Let it be",
      provider: "Powered by Jobcher",
      blog: "https://www.jobcher.com",
      data_updatetime: new Date().toISOString().slice(0, 10).replace(/-/g, '')
    };
  }
}

/**
 * Format the IP information as a table-like structure
 * @param {Object} info
 */
function formatAsTable(info) {
  const pad = (str, length) => str.toString().padEnd(length, ' ');
  const keyWidth = 20;  // Adjusted width for keys
  const valueWidth = 60; // Adjusted width for values

  const lineBorder = (keyLength, valueLength) =>
    `┏${'━'.repeat(keyLength)}┳${'━'.repeat(valueLength)}┓`;
  const lineMiddle = (keyLength, valueLength) =>
    `┡${'━'.repeat(keyLength)}╇${'━'.repeat(valueLength)}┩`;
  const lineEnd = (keyLength, valueLength) =>
    `└${'─'.repeat(keyLength)}┴${'─'.repeat(valueLength)}┘`;

  const borderTop = lineBorder(keyWidth+2, valueWidth+2);
  const borderMiddle = lineMiddle(keyWidth+2, valueWidth+2);
  const borderBottom = lineEnd(keyWidth+2, valueWidth+2);

  const lines = [
    '                                  ip.jobcher.com                               ',
    borderTop,
    `┃ ${pad('key', keyWidth)} ┃ ${pad('value', valueWidth)} ┃`,
    borderMiddle,
    `│ ${pad('ip', keyWidth)} │ ${pad(info.ip, valueWidth)} │`,
    `│ ${pad('city', keyWidth)} │ ${pad(info.city, valueWidth)} │`,
    `│ ${pad('province', keyWidth)} │ ${pad(info.province, valueWidth)} │`,
    `│ ${pad('country', keyWidth)} │ ${pad(info.country, valueWidth)} │`,
    `│ ${pad('continent', keyWidth)} │ ${pad(info.continent, valueWidth)} │`,
    `│ ${pad('isp', keyWidth)} │ ${pad(info.isp, valueWidth)} │`,
    `│ ${pad('time_zone', keyWidth)} │ ${pad(info.time_zone, valueWidth)} │`,
    `│ ${pad('latitude', keyWidth)} │ ${pad(info.latitude, valueWidth)} │`,
    `│ ${pad('longitude', keyWidth)} │ ${pad(info.longitude, valueWidth)} │`,
    `│ ${pad('postal_code', keyWidth)} │ ${pad(info.postal_code, valueWidth)} │`,
    `│ ${pad('iso_code', keyWidth)} │ ${pad(info.iso_code, valueWidth)} │`,
    `│ ${pad('notice', keyWidth)} │ ${pad(info.notice, valueWidth)} │`,
    `│ ${pad('', keyWidth)} │ ${pad('©2020-01-01->now', valueWidth)} │`,
    `│ ${pad('provider', keyWidth)} │ ${pad(info.provider, valueWidth)} │`,
    `│ ${pad('blog', keyWidth)} │ ${pad(info.blog, valueWidth)} │`,
    `│ ${pad('data_updatetime', keyWidth)} │ ${pad(info.data_updatetime, valueWidth)} │`,
    borderBottom,
    ` `,
  ];

  return lines.join('\n');
}
```

### 5. 点击`save and deploy`，然后点击`部署`。
![创建worker](/images/ipcheck-7.png)  

## 使用
### json返回
```sh
# 查询本地ip
curl ip.jobcher.com
```

```sh
# 查询指定ip
curl ip.jobcher.com?ip=1.1.1.1
```

### table返回
```sh
curl ip.jobcher.com/table
```

```sh
curl ip.jobcher.com/table?ip=1.1.1.1
```

