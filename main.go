package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"regexp"
)

var (
	ver            bool
	countdown_time string
	freq_time      string
	cmd            string
	config_file    string
	calendar       string
)

const (
	Version = "0.5"
)
var Week = map[string]int{
		"Sunday":    0,
		"sun":       0,
		"Monday":    1,
		"mon":       1,
		"Tuesday":   2,
		"tue":       2,
		"Wednesday": 3,
		"wed":       3,
		"Thursday":  4,
		"thu":       4,
		"Friday":    5,
		"fri":       5,
		"Saturday":  6,
		"sat":       6,
	}


func init() {
	flag.BoolVar(&ver, "v", false, "show version")
	flag.BoolVar(&ver, "version", false, "show version")
	flag.StringVar(&countdown_time, "c", "", "countdown timer's time")
	flag.StringVar(&countdown_time, "count", "", "countdown timer's time")
	flag.StringVar(&freq_time, "set", "", "timer's time")
	flag.StringVar(&freq_time, "s", "", "timer's time")
	flag.StringVar(&calendar, "calendar", "", "calendar event")
	flag.StringVar(&calendar, "cal", "", "calendar event")
	flag.StringVar(&config_file, "f", "", "config file")
	flag.StringVar(&cmd, "cmd", "", "Command that you want to exec")
	flag.StringVar(&cmd, "command", "", "Command that you want to exec")
}

func main() {
	logfile, err := os.OpenFile("TimerGo.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer logfile.Close()
	if err != nil {
		log.Fatalf("open file err :%v\n", err)

	}
	LogInit(logfile)

	parserArg()
	if config_file != "" {
		countdown_time, freq_time,calendar, cmd = GetConfig(config_file)
		fmt.Printf("count :%v fre :%v cal :%v cmd :%v\n", countdown_time, freq_time, calendar,cmd)
	}

	if calendar != "" {
		onCalendar()
		os.Exit(0)
	}

	if countdown_time != "" {
		doCountdown()
		os.Exit(0)
	}
	// 死循环，ctrl + c 退出
	if freq_time != "" {
		freq := strings.Split(freq_time, " ")
		var ftime time.Duration
		var err error
		ftime, err = timeFormat(freq[len(freq)-1])
		if err != nil {
			LogErr("timeFormat error :%v\n", err)
			os.Exit(-1)
		}
		s := frequencyTimer(ftime)
		for {
			ExecCmd(s, true)
		}
	} else if cmd != "" {
		// execute command immediately
		sign := make(chan byte)
		go func() {
			sign <- 1
			<-sign
		}()
		ExecCmd(sign, true)
	}
}

func doCountdown() {
	var ctime time.Duration
	var err error
	ctime, err = timeFormat(countdown_time)
	if err != nil {
		LogErr("timeFormat error :%v \n", err)
		os.Exit(-1)
	}
	fmt.Printf("coutdown time set to %v\n", ctime)
	ExecCmd(countdownTimer(ctime), true)
}

// 处理参数
func parserArg() {
	flag.Usage = usage
	flag.Parse()

	if 0 == flag.NFlag() {
		usage()
		os.Exit(0)
	}
	if ver {
		fmt.Println("Version :", Version)
		os.Exit(0)
	}
}

//处理时间格式
func timeFormat(t string) (ret time.Duration, err error) {
	//long_format = "2016-6-23 12:13:14"
	//short_format = "12:13:14"
	timeSlice := make([]string, 6)
	//var t1 timeVar.Time
	var dateVar, timeVar []string

	if len(t) > 8 {
		//确保输入时空格不会影响到slice
		timeSlice = regexp.MustCompile(" +").Split(strings.TrimSpace(t),-1)
		dateVar = strings.Split(timeSlice[0], "-")
		timeVar = strings.Split(timeSlice[1], ":")
		targetTime, err := timeTransform(dateVar, timeVar)
		if err != nil {
			LogErr("time transformat err :%v\n", err)
			return -1, err
		}
		now := time.Now()
		if !now.Before(targetTime) {
			return -1, errors.New("The time you input should after now")
		}
		return targetTime.Sub(now), nil
	} else {
		timeVar = strings.Split(t, ":")

		//使用不同的时间格式，1:20:30 ; 1:20 ; 10
		len := len(timeVar)

		switch len {
		case 3:
			return computeTime(timeVar, 3)
		case 2:
			return computeTime(timeVar, 2)
		case 1:
			a, e := strconv.Atoi(t)
			return time.Duration(a) * time.Second, e
		default:
			return -1, errors.New("The error in timeFormat.len(timeSlice) < 1")
		}
	}
}

func timeTransform(d []string, t []string) (ret time.Time, err error) {
	var year, month, day, hour, min, sec int = 0, 0, 0, 0, 0, 0

	if d != nil {
		year, err = strconv.Atoi(d[0])
		month, err = strconv.Atoi(d[1])
		day, err = strconv.Atoi(d[2])
	}
	if t != nil {
		switch len(t) {
		case 2:
			hour, err = strconv.Atoi(t[0])
			min, err = strconv.Atoi(t[1])
		default:
			hour, err = strconv.Atoi(t[0])
			min, err = strconv.Atoi(t[1])
			sec, err = strconv.Atoi(t[2])
		}
	}
	if err != nil {
		LogErr("time strconv atoi err :%v\n", err)
		return time.Time{}, err
	}
	ret = time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)
	err = nil
	return
}

//把时间hh:mm:ss 转化为秒数

func computeTime(s []string, n int) (time.Duration, error) {

	var sum time.Duration
	for i := 0; i < n; i++ {
		tmp, err := strconv.Atoi(s[i])
		if err != nil {
			LogErr("convert string to int err :%v\n", err)
			return -1, err
		}
		sum = time.Duration(tmp)*time.Second + sum*60
	}
	return sum, nil

}

func onCalendar() {
	var now time.Time
	var time_diff time.Duration
	cSlice := regexp.MustCompile(" +").Split(strings.TrimSpace(calendar),-1)
	switch cSlice[0] {
	case "daily":
		t, err := timeTransform(nil, strings.Split(cSlice[1], ":"))
		if err != nil {
			LogErr("daily time transform error :%v\n", err)
			os.Exit(-1)
		}
		now = time.Now()
		time_diff = time.Duration((t.Hour()-now.Hour())*3600+(t.Minute()-now.Minute())*60+(t.Second()-now.Second())) * time.Second
		fmt.Printf("now :%v \n t : %v\n, time_diff :%v\n", now, t, time_diff)

		if time_diff > 0 {
			ExecCmd(frequencyTimer(time_diff), false)
		} else {
			ExecCmd(frequencyTimer(time_diff+time.Hour*24), false)
		}
		for {
			ExecCmd(frequencyTimer(time.Hour*24), false)
		}
	case "weekly":
		t, err := timeTransform(nil, strings.Split(cSlice[2], ":"))
		if err != nil {
			LogErr("weekly time transform error :%v\n", err)
			os.Exit(-1)
		}
		now = time.Now()
		time_diff = time.Duration((t.Hour()-now.Hour())*3600+(t.Minute()-now.Minute())*60+(t.Second()-now.Second())) * time.Second
		fmt.Printf("now :%v \n t : %v\n, time_diff :%v\n", now, t, time_diff)
		if int(now.Weekday()) == Week[cSlice[1]] {
			if  time_diff > 0 {
				go ExecCmd(countdownTimer(time_diff),true)
			}
		}
		for {
			sign := countdownTimer(time.Hour * 24)
			<-sign
			now =  time.Now()
			//提前一天检查，防止时间过掉
			if int (now.Weekday())== Week[cSlice[1]]- 1  {
				if time_diff > 0{
					go ExecCmd(countdownTimer(time_diff + time.Hour * 24),true)
				}else{
					go ExecCmd(countdownTimer(time_diff + time.Hour * 24 * 2),true)
				}
			}
		}

	case "monthly":
		t, err := timeTransform(nil, strings.Split(cSlice[2], ":"))
		if err != nil {
			LogErr("weekly time transform error :%v\n", err)
			os.Exit(-1)
		}
		now = time.Now()
		time_diff = time.Duration((t.Hour()-now.Hour())*3600+(t.Minute()-now.Minute())*60+(t.Second()-now.Second())) * time.Second
		fmt.Printf("now :%v \n t : %v\n, time_diff :%v\n", now, t, time_diff)
		day , _ := strconv.Atoi(cSlice[1])
		if now.Day() == day {
			if  time_diff > 0 {
				go ExecCmd(countdownTimer(time_diff),true)
			}
		}
		for{
			sign := countdownTimer(time.Hour * 24)
			<-sign
			now =  time.Now()
			//提前一天，防止过时无法执行
			if now.Day()== day - 1 {
				if time_diff > 0{
					go ExecCmd(countdownTimer(time_diff + time.Hour * 24),true)
				}else{
					go ExecCmd(countdownTimer(time_diff + time.Hour * 24 * 2),true)
				}
			}
		}
	default:
		fmt.Println("calendar's format is error")
		os.Exit(-1)

	}

}

func countdownTimer(t time.Duration) chan byte {
	fmt.Printf("countdown time :%v\n",t)
	sign := make(chan byte)
	go func() {
		timer := timerSec(t)
		defer timer.Stop()
		<-timer.C
		sign <- 1
		<-sign //接收command完成信息
	}()
	return sign
}

func frequencyTimer(t time.Duration) chan byte {

	fmt.Printf("the frequency time : %v \n", t)
	sign := make(chan byte)
	//假设是数字，不是时间格式
	go func() {
		timer := timerSec(t)
		for {
			select {
			case <-timer.C:
				sign <- 1
				timer.Reset(time.Duration(t))
			}
			<-sign //非缓冲的channel 会阻塞等待,在命令执行完成后，再开始计时
		}
	}()
	return sign
}

func timerSec(n time.Duration) *time.Timer {
	return time.NewTimer(n)
}
// sign :接受信号，开始执行命令
// wait：定时器是否需要等待命令结束
func ExecCmd(sign chan byte, wait bool) {
	if sig := <-sign; sig == 1 {
		if !wait {
			sign <- 0
		}
		if cmd == "" {
			return
		}
		fmt.Printf("recevie a sign : 1\n")
//		cmdAndArgs:= regexp.MustCompile(" +").Split(strings.TrimSpace(cmd),-1)
//		fmt.Printf("cmd :%v args :%v\n", cmdAndArgs[0], cmdAndArgs[1:])
		c := exec.Command("bash","-c",cmd)
		c.Stdout = os.Stdout
		if err := c.Run(); err != nil {
			LogErr("command runing error :%v\n", err)
		}
		fmt.Printf("*******command end*******")
		if wait {
			sign <- 0
		}
	}
}

const usageString string = `TimerGo is a timer tools
Usage :
          timergo [flags] [timer] [flags][commands]
flags
     -v,--version : show the Version Nubmer
     -h,--help    : show this help
     -c,--count   : set the countdown time to execute command
                     example:
                    -c 10   	10 sec
                    -c 1:20 	1min 20 sec
                    -c 1:23:33	1hour 23 min 33 sec
                    -c "2016-12-30 12:13" execute the command until that time
     -s,--set     : set the time interval of the command
     	            example:
     	            -s 10       execute the command for every 10 seconds
     	            -s 1:20     execute the command for every 1 min 20 sec
     	            -s 1:23:33  execute the command for every 1 hoour 23 min 33 sec
     -cal, --calendar    : set calendar clock to execute command
                    example :
                    -calendar "daily 10:20"
                    -calendar "weekly Sunday 22:30"
                    -calendar "monthly 20 12:20"
                    notes: time is the 24-hour clock
     -cmd,--command :  the command and it's arguments that you want to execute
     -f           : use config.ini
`
func usage() {
	fmt.Println(usageString)
}
