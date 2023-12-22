---
title: "Navicat 查看导出连接的密码 | navicat查看密码方案"
date: 2022-02-08
draft: false
authors: "jobcher"
tags: ["mysql"]
categories: ["数据库"]
series: ["mysql系列"]
---

# Navicat 查看密码方案

## 解决问题：

我们经常使用 navicat 连接数据库，有时候时间久了之后，会忘记之前的密码，那么现在我们有办法获得只要正常连接的数据库的密码

## 步骤：

1. 导出连接 connections.ncx，拿到保存到本地的 connections.ncx 文件中的 Password，粘贴到下面的代码中
2. 登陆[https://tool.lu/coderunner/](https://tool.lu/coderunner/)，使用 PHP 在线运行工具，粘贴下面添加密码后的代码  
   备用工具网址（https://zixuephp.net/tool-runcode.html）

```php
<?php
class NavicatPassword
{
    protected $version = 0;
    protected $aesKey = 'libcckeylibcckey';
    protected $aesIv = 'libcciv libcciv ';
    protected $blowString = '3DC5CA39';
    protected $blowKey = null;
    protected $blowIv = null;

    public function __construct($version = 12)
    {
        $this->version = $version;
        $this->blowKey = sha1('3DC5CA39', true);
        $this->blowIv = hex2bin('d9c7c3c8870d64bd');
    }

    public function encrypt($string)
    {
        $result = FALSE;
        switch ($this->version) {
            case 11:
                $result = $this->encryptEleven($string);
                break;
            case 12:
                $result = $this->encryptTwelve($string);
                break;
            default:
                break;
        }

        return $result;
    }

    protected function encryptEleven($string)
    {
        $round = intval(floor(strlen($string) / 8));
        $leftLength = strlen($string) % 8;
        $result = '';
        $currentVector = $this->blowIv;

        for ($i = 0; $i < $round; $i++) {
            $temp = $this->encryptBlock($this->xorBytes(substr($string, 8 * $i, 8), $currentVector));
            $currentVector = $this->xorBytes($currentVector, $temp);
            $result .= $temp;
        }

        if ($leftLength) {
            $currentVector = $this->encryptBlock($currentVector);
            $result .= $this->xorBytes(substr($string, 8 * $i, $leftLength), $currentVector);
        }

        return strtoupper(bin2hex($result));
    }

    protected function encryptBlock($block)
    {
        return openssl_encrypt($block, 'BF-ECB', $this->blowKey, OPENSSL_RAW_DATA|OPENSSL_NO_PADDING);
    }

    protected function decryptBlock($block)
    {
        return openssl_decrypt($block, 'BF-ECB', $this->blowKey, OPENSSL_RAW_DATA|OPENSSL_NO_PADDING);
    }

    protected function xorBytes($str1, $str2)
    {
        $result = '';
        for ($i = 0; $i < strlen($str1); $i++) {
            $result .= chr(ord($str1[$i]) ^ ord($str2[$i]));
        }

        return $result;
    }

    protected function encryptTwelve($string)
    {
        $result = openssl_encrypt($string, 'AES-128-CBC', $this->aesKey, OPENSSL_RAW_DATA, $this->aesIv);
        return strtoupper(bin2hex($result));
    }

    public function decrypt($string)
    {
        $result = FALSE;
        switch ($this->version) {
            case 11:
                $result = $this->decryptEleven($string);
                break;
            case 12:
                $result = $this->decryptTwelve($string);
                break;
            default:
                break;
        }

        return $result;
    }

    protected function decryptEleven($upperString)
    {
        $string = hex2bin(strtolower($upperString));

        $round = intval(floor(strlen($string) / 8));
        $leftLength = strlen($string) % 8;
        $result = '';
        $currentVector = $this->blowIv;

        for ($i = 0; $i < $round; $i++) {
            $encryptedBlock = substr($string, 8 * $i, 8);
            $temp = $this->xorBytes($this->decryptBlock($encryptedBlock), $currentVector);
            $currentVector = $this->xorBytes($currentVector, $encryptedBlock);
            $result .= $temp;
        }

        if ($leftLength) {
            $currentVector = $this->encryptBlock($currentVector);
            $result .= $this->xorBytes(substr($string, 8 * $i, $leftLength), $currentVector);
        }

        return $result;
    }

    protected function decryptTwelve($upperString)
    {
        $string = hex2bin(strtolower($upperString));
        return openssl_decrypt($string, 'AES-128-CBC', $this->aesKey, OPENSSL_RAW_DATA, $this->aesIv);
    }
};


//需要指定navacat版本两种，11或12
$navicatPassword = new NavicatPassword(12);

//解密，括号里面写入navicat加密后的密码
$decode = $navicatPassword->decrypt('E75BF077AB8BAA3AC2D5');
echo $decode."\n";
?>
```

3. 点击执行之后，就会得到真实密码
