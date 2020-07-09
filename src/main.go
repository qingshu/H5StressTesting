//设计思路：所有机器人每隔一段时间登陆，，登陆后的机器人每隔一段时间随机执行一些操作
//设计思想：主要是模板设计模式思想
//时间：2019/11/15
//作者：lyp
package main

import (
	"H5Test/log"
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"
)

var testPeople = 5                    //测试玩家数
var timeSpaceOfPlayerLogin = 1000     //每个玩家登陆的时间间隔（毫秒）
var minTimeSpaceOfPlayerAction = 800  //每个玩家执行下一个操作等待最小时间(毫秒)
var maxTimeSpaceOfPlayerAction = 2000 //每个玩家执行下一个操作等待最大时间（毫秒）
var isSaveLog2LocalText = false       //是否日志保存到本地文本
var serverUlr = "http://120.77.174.186:90/game/s1/mywork/index.php?r=game/index"
var PlayerNameFlag = "Test"

func main() {
	//设置cpu核数
	cpuNum := runtime.NumCPU() * 2 / 3
	if cpuNum < 1 {
		cpuNum = 1
	}
	runtime.GOMAXPROCS(cpuNum)

	//接收参数
	if len(os.Args) >= 8 {
		testPeople, _ = strconv.Atoi(os.Args[1])
		timeSpaceOfPlayerLogin, _ = strconv.Atoi(os.Args[2])
		minTimeSpaceOfPlayerAction, _ = strconv.Atoi(os.Args[3])
		maxTimeSpaceOfPlayerAction, _ = strconv.Atoi(os.Args[4])
		isSaveLog2LocalText, _ = strconv.ParseBool(os.Args[5])
		serverUlr = os.Args[6]
		PlayerNameFlag = os.Args[7]
	}
	log.Init(isSaveLog2LocalText)
	log.Info(fmt.Sprintf("测试人数:%d,每个玩家登陆间隔%d毫秒,每个玩家执行下一个行为等待间隔%d-%d毫秒", testPeople,
		timeSpaceOfPlayerLogin, minTimeSpaceOfPlayerAction, maxTimeSpaceOfPlayerAction))

	//开始压测
	for i := 1; i <= testPeople; i++ {
		go func(j int) {
			defer func(k int) {
				if r := recover(); nil != r {
					log.Error("线程异常，玩家序号：" + strconv.Itoa(k))
					log.Error(r)
				}
			}(j)
			player := new(Player)
			player.Init(j)
		}(i)
		time.Sleep(time.Millisecond * time.Duration(timeSpaceOfPlayerLogin))
	}

	//主线程等待键盘输入
	cmdReader := bufio.NewReader(os.Stdin)
	cmdReader.ReadString('\n')
}

//url.Values{"key": {"Value"}, "id": {"123"}}
func HttpPostForm(postURL string, urlValues url.Values) (error, string) {
	resp, err := http.PostForm(postURL, urlValues)

	if err != nil {
		return err, ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, ""
	}

	return nil, string(body)
}
