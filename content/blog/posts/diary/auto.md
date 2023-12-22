---
title: "获取用户浏览器默认语言设置，自动判断跳转不同网站"
date: 2021-12-16
draft: false
author: "jobcher"
tags: ["javascript"]
categories: ["日常"]
series: ["日常系列"]
---

# 自动判断跳转不同网站

1. 根据用户目前的浏览器配置语言进行显示
2. 供语言切换按钮，用户自定义选择不同的语言显示

根据识别用户的浏览器语言，自动判断并跳转到相应的语言网页，让你的网站更加灵动。  
以下需要将代码放在 HTML 的内即可，然后自行制作多语言页面。  
代码如下：

```html
<script type="text/javascript">
  //获取用户语言的顺序是
  //1.获取本地缓存里的内容
  //2.用户浏览器的语言设置
  //如果上面2个都没有获取到，就直接使用'en'作为用户选择的语言
  var language =
    localStorage.getItem("locale") ||
    window.navigator.language.toLowerCase() ||
    "en";
  //把用户的语言写入缓存，供下次获取使用
  localStorage.setItem("locale", language);
  //判断用户的语言，跳转到不同的地方
  if (language.indexOf("zh-") !== -1) {
    window.location = "/zh-cn/index.html";
  } else if (language.indexOf("en") !== -1) {
    window.location = "/en/index.html";
  } else {
    //其它的都使用英文
    window.location = "/en/index.html";
  }
</script>
```

核心代码  
其实核心代码就是利用 navigator 的 language 属性

```code
navigator.language
```

## 第二种解决方案

可以通过获取用户的 IP，然后把 IP 放到 IP 库里查询所在地，从而加载对应的资源，这样的方案回更加准确！有的第三方会直接返回所在国家的编码，比如 cn / en 等就更好了

但是这样的方案也有一个弊端：如果用户通过科学上网，全局模式下，会被认为属于美国 / 日本等等（看梯子的 IP 而定了），那么会导致访问非常慢；但是这种偏差，很多翻墙的人都是了解的，没人会故意用美国的 IP 访问国内的淘宝 / 百度等网站的，除非是忘记切换回来了；

IP 判断  
市场上有很多 IP 判断的，拿 IP 倒是非常好做的一件事；比如我现在可以拿到用户访问本网站时候的 IP；

欢迎关注我的博客[www.jobcher.com](https://www.jobcher.com/)
