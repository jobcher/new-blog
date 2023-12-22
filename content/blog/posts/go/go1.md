---
title: "go 基础知识"
date: 2022-04-25
draft: false
author: "jobcher"
tags: ["golang"]
categories: ["go"]
series: ["go系列"]
---

# go 基础知识

## 目录结构

```sh
├─ code  -- 代码根目录
│  ├─ bin
│  ├─ pkg
│  ├─ src
│     ├── hello
│         ├── hello.go
```

- bin 存放编译后可执行的文件。
- pkg 存放编译后的应用包。
- src 存放应用源代码。

Hello World 代码

```go
//在 hello 目录下创建 hello.go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World!")
}
```

## 基础命令

```sh
go build hello
#在src目录或hello目录下执行 go build hello，只在对应当前目录下生成文件。
go install hello
#在src目录或hello目录下执行 go install hello，会把编译好的结果移动到 $GOPATH/bin。
go run hello
#在src目录或hello目录下执行 go run hello，不生成任何文件只运行程序。
go fmt hello
#在src目录或hello目录下执行 go run hello，格式化代码，将代码修改成标准格式。
```

## 数据类型

| 类型   | 表示                                                       | 备注                                                                     |
| :----- | :--------------------------------------------------------- | :----------------------------------------------------------------------- |
| 字符串 | string                                                     | 只能用一对双引号（""）或反引号（``）括起来定义，不能用单引号（''）定义！ |
| 布尔   | bool                                                       | 只有 true 和 false，默认为 false。                                       |
| 整型   | int8 uint8 int16 uint16 int32 uint32 int64 uint64 int uint | 具体长度取决于 CPU 位数。                                                |
| 浮点型 | float32 float64                                            |                                                                          |

## 常量声明

`常量`，在程序编译阶段就确定下来的值，而程序在运行时无法改变该值。

### 1. 单个常量声明

第一种：const 变量名称 数据类型 = 变量值  
如果不赋值，使用的是该数据类型的默认值。  
第二种：const 变量名称 = 变量值  
根据变量值，自行判断数据类型。

### 2. 多个常量声明

第一种：const 变量名称,变量名称 ... ,数据类型 = 变量值,变量值 ...  
第二种：const 变量名称,变量名称 ... = 变量值,变量值 ...

### 3. 代码

```go
//demo_1.go
package main

import (
	"fmt"
)

func main() {
	const name string = "Tom"
	fmt.Println(name)

	const age = 30
	fmt.Println(age)

	const name_1, name_2 string = "Tom", "Jay"
	fmt.Println(name_1, name_2)

	const name_3, age_1 = "Tom", 30
	fmt.Println(name_3, age_1)
}
```

## 变量声明

### 单个变量声明

第一种：var 变量名称 数据类型 = 变量值  
如果不赋值，使用的是该数据类型的默认值。  
第二种：var 变量名称 = 变量值  
根据变量值，自行判断数据类型。  
第三种：变量名称 := 变量值  
省略了 var 和数据类型，变量名称一定要是未声明过的。

### 多个变量声明

第一种：var 变量名称,变量名称 ... ,数据类型 = 变量值,变量值 ...  
第二种：var 变量名称,变量名称 ... = 变量值,变量值 ...  
第三种：变量名称,变量名称 ... := 变量值,变量值 ...

### 代码

```go
//demo_2.go
package main

import (
	"fmt"
)

func main() {
	var age_1 uint8 = 31
	var age_2 = 32
	age_3 := 33
	fmt.Println(age_1, age_2, age_3)

	var age_4, age_5, age_6 int = 31, 32, 33
	fmt.Println(age_4, age_5, age_6)

	var name_1, age_7 = "Tom", 30
	fmt.Println(name_1, age_7)

	name_2, is_boy, height := "Jay", true, 180.66
	fmt.Println(name_2, is_boy, height)
}
```

## 输出方法

fmt.Print：输出到控制台（仅只是输出）

fmt.Println：输出到控制台并换行

fmt.Printf：仅输出格式化的字符串和字符串变量（整型和整型变量不可以）

fmt.Sprintf：格式化并返回一个字符串，不输出。

### 代码

```go
//demo_3.go
package main

import (
	"fmt"
)

func main() {
	fmt.Print("输出到控制台不换行")
	fmt.Println("---")
	fmt.Println("输出到控制台并换行")
	fmt.Printf("name=%s,age=%d\n", "Tom", 30)
	fmt.Printf("name=%s,age=%d,height=%v\n", "Tom", 30, fmt.Sprintf("%.2f", 180.567))
}
```

## 数组

数组是一个由固定长度的特定类型元素组成的序列，一个数组可以由零个或多个元素组成，一旦声明了，数组的长度就固定了，不能动态变化。

`len()` 和 `cap()` 返回结果始终一样。

### 声明数组

```go
package main

import (
	"fmt"
)

func main() {
	//一维数组
	var arr_1 [5] int
	fmt.Println(arr_1)

	var arr_2 =  [5] int {1, 2, 3, 4, 5}
	fmt.Println(arr_2)

	arr_3 := [5] int {1, 2, 3, 4, 5}
	fmt.Println(arr_3)

	arr_4 := [...] int {1, 2, 3, 4, 5, 6}
	fmt.Println(arr_4)

	arr_5 := [5] int {0:3, 1:5, 4:6}
	fmt.Println(arr_5)

	//二维数组
	var arr_6 = [3][5] int {{1, 2, 3, 4, 5}, {9, 8, 7, 6, 5}, {3, 4, 5, 6, 7}}
	fmt.Println(arr_6)

	arr_7 :=  [3][5] int {{1, 2, 3, 4, 5}, {9, 8, 7, 6, 5}, {3, 4, 5, 6, 7}}
	fmt.Println(arr_7)

	arr_8 :=  [...][5] int {{1, 2, 3, 4, 5}, {9, 8, 7, 6, 5}, {0:3, 1:5, 4:6}}
	fmt.Println(arr_8)
}

```

### 注意事项

一、数组不可动态变化问题，一旦声明了，其长度就是固定的。

```go
var arr_1 = [5] int {1, 2, 3, 4, 5}
arr_1[5] = 6
fmt.Println(arr_1)
```

运行会报错：`invalid array index 5 (out of bounds for 5-element array)`

二、数组是值类型问题，在函数中传递的时候是传递的值，如果传递数组很大，这对内存是很大开销。

```go
//demo_5.go
package main

import (
	"fmt"
)

func main() {
	var arr =  [5] int {1, 2, 3, 4, 5}
	modifyArr(arr)
	fmt.Println(arr)
}

func modifyArr(a [5] int) {
	a[1] = 20
}
```

```go
//demo_6.go
package main

import (
	"fmt"
)

func main() {
	var arr =  [5] int {1, 2, 3, 4, 5}
	modifyArr(&arr)
	fmt.Println(arr)
}

func modifyArr(a *[5] int) {
	a[1] = 20
}
```

三、`数组赋值`问题，同样类型的数组（长度一样且每个元素类型也一样）才可以相互赋值，反之不可以。

```go
var arr =  [5] int {1, 2, 3, 4, 5}
var arr_1 [5] int = arr
var arr_2 [6] int = arr
```

运行会报错：`cannot use arr (type [5]int) as type [6]int in assignment`
