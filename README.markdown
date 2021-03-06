Alfred-Youdao-Saber
===

Youdao Dict for Alfred (coded by saber).    

![预览](./assets/demo.gif)

## Features

+ 使用方法 `yd {query}`
+ 支持屏幕取词功能, 需要先在Alfred Workflow中设置热键
+ 使用`Shift+Enter`在有道词典网页中打开
+ 使用`Cmd+Enter`发音
+ 按`Enter`键复制翻译结果
+ 支持设置源语言和目标语言，支持中文、日语、英语等相互翻译。语法是`yd zh=>ja 我爱你`
+ `yd #`查询历史记录，最多支持20条历史记录

支持相互翻译的语言列表如下:

+ `zh` - 中文
+ `ja`     - 日文
+ `en`     - 英文
+ `ko`     - 韩文
+ `fr`     - 法文
+ `ru`     - 俄文
+ `pt`     - 葡萄牙文
+ `es`     - 西班牙文
+ `auto`   - 自动

## Dependencies

+ Go 1.6+

## TODO

+ [x] 一个好的自动更新机制
+ [x] 添加语音
+ [x] 修改查询中文翻译发音为中文原文的Bug，并修改其他相关语音逻辑
+ [x] 添加历史记录功能
+ [ ] 自定义有道API Key和Secret
+ [ ] 拼写提示功能

## CHANGELOG (Edited by saber)

### 1.5.1-saber

+ 修改查询中文翻译发音为中文原文的Bug，并修改其他相关语音逻辑
+ 添加历史记录功能
+ 发现MacOS从10.11开始的版本对于QuickLookURL的支持不太兼容的问题

## CHANGELOG (Depreciated)

### 1.5.1

+ 发音添加了对日语、法语等语言的支持

### 1.5.0

+ 添加发音功能，调用Mac自带的发音软件
+ 优化了打开词典的网页
+ 其他小修改

### 1.4.0

+ 添加手动设置翻译语言
+ 修改结果展示，让展示的结果更全面

### 1.3.2

+ 修复缓存失效问题
+ 移除没有什么用处的自动更新

### 1.3.0

+ 添加在后台自动更新的机制

### 1.2.4

+ 修改Alfred执行队列行为，让查询更平滑

### 1.2.3

+ 添加按`Shift+Enter`键在有道网站打开功能

### 1.2.2

+ 添加按Enter键复制功能
+ 可以在Workflows里面设置快捷键实现屏幕取词功能

### 1.2.1

+ 延长缓存时间
+ 延长查询超时时间

### 1.2.0

+ 添加快捷键可以选词查询

#### 1.1.0

+ 添加了缓存功能

#### 1.0.0

+ 完成查询和翻译功能
