#TimerGo

##概述
timergo是用golang开发的一个定时器命令行工具。可以使用其完成定时任务。
本项目主要用于练习golang，实际可以使用systemd的timer服务。

##主要功能
- 设定一个倒计时，当时间计时结束，执行自设定的命令
- 设定一个频率计，在固定的时间段，循环执行命令
- 设定一个按照日期，星期，每天来执行的日历计时器
- 从配置文件中，读取参数执行

##使用方法
1. 设定倒计时

设定倒计时20s，显示helloworld
```
timergo -c 20 -cmd "echo helloworld"
```

设定倒计时1min20s,显示helloworld
```
timergo -c 1:20 -cmd "echo helloworld"
```
设定倒计时1小时20分钟30秒
```
timergo -c 1:20:30 -cmd "echo hellworld"
```
设定到2020年3月4日 1：00执行
```
timergo -c "2020-3-4 1:00" -cmd "echo hellowold"
```
设定到2019年2月12日 12时22分22秒执行
```
timergo -c "2019-2-12 12:22:22" -cmd "echo helloworld"
```
2. 设定频率计

设定时间间隔为20s，循环显示helloworld
```
timergo -s 20 -cmd "echo helloworld"
```
其他时间格式与定时器相同

3. 设定日历计时器
设定每天定时执行：
```
timergo -cal "daily 10:20" -cmd "echo helloworld"
```
设定每周执行定时任务：
```
timergo -cal "weekly Sunday 11:00" -cmd "echo hellowrld"
```
设定每个月执行定时任务：
```
timergo -cal "monthly 10 15:00" -cmd "echo helloworld"
```

4. 从配置中读取
默认的配置是ini格式，示例如下：

```
[Timer]
CountdownTime = 20
#FreqTime = 1:05
#Calendar = daily 12:10
Exec = echo helloworld
```

##已经完成的功能

- [x] 实现倒计时执行任务
- [x] 实现频率计执行任务
- [x] 实现日期的倒计时
- [x] 实现每日定时，每周定时，每月定时功能



