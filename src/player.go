package main

import (
	"H5Test/log"
	"encoding/json"
	"math/rand"
	"net/url"
	"strconv"

	"strings"
	"time"
)

//玩家行为
const (
	Action_CreateRole          = 0                             //创角
	Action_Guide               = 1                             //引导
	Action_SelectStartRole     = 2                             //选择初始神将
	Action_NiuDan              = 3                             //扭蛋
	Action_DrawCard            = 4                             //抽卡
	Action_2Battle             = 5                             //上阵
	Action_HangupDup           = 6                             //挂机副本
	Action_OneKeyDress         = 7                             //一键穿戴
	Action_NewPlayerActionMax  = 8                             //特殊！！！！！！！！！！！！！！！！相当于分割线，上面的全是新玩家必须请求的引导类消息
	Action_FightBoss           = Action_NewPlayerActionMax + 0 //打boss
	Action_AllEquipMelting     = Action_NewPlayerActionMax + 1 //装备全部熔炼
	Action_HnagupOnLine        = Action_NewPlayerActionMax + 2 //在线挂机
	Action_RequestNiuDanInfo   = Action_NewPlayerActionMax + 3 //请求扭蛋信息
	Action_GetDiscountShopInfo = Action_NewPlayerActionMax + 4 //获取折扣商店信息
	Action_Max                 = Action_NewPlayerActionMax + 5
)

//引导步骤
const (
	GuideStep_First           = 1 //第一个引导
	GuideStep_SelectStartRole = 2 //引导选择初始神将
	GuideStep_Fourth          = 4 //第四个引导
	GuideStep_NiuDan          = 5 //引导请求扭蛋
	GuideStep_DrawCard        = 6 //引导抽卡
	GuideStep_2Battle         = 7 //引导上阵
	GuideStep_HeightNiuDan    = 8 //引导高级扭蛋
	GuideStep_2Battle2        = 9 //引导第二次上阵
)

var ActionCmd map[int]string  //玩家行为对应的指令
var ActionName map[int]string //玩家行为名字

func init() {
	//初始化玩家指令
	ActionCmd = make(map[int]string)
	ActionCmd[Action_CreateRole] = "20001"
	ActionCmd[Action_Guide] = "20002"
	ActionCmd[Action_SelectStartRole] = "10010"
	ActionCmd[Action_NiuDan] = "10184"
	ActionCmd[Action_DrawCard] = "10186"
	ActionCmd[Action_2Battle] = "10060"
	ActionCmd[Action_HangupDup] = "10088"
	ActionCmd[Action_OneKeyDress] = "40076"
	ActionCmd[Action_FightBoss] = "10088"
	ActionCmd[Action_AllEquipMelting] = "10116"
	ActionCmd[Action_HnagupOnLine] = "10008"
	ActionCmd[Action_RequestNiuDanInfo] = "10184"
	ActionCmd[Action_GetDiscountShopInfo] = "10126"

	//初始化指令名字
	ActionName = make(map[int]string)
	ActionName[Action_CreateRole] = "创角"
	ActionName[Action_Guide] = "引导"
	ActionName[Action_SelectStartRole] = "选择初始神将"
	ActionName[Action_NiuDan] = "扭蛋"
	ActionName[Action_DrawCard] = "抽卡"
	ActionName[Action_2Battle] = "上阵"
	ActionName[Action_HangupDup] = "挂机副本"
	ActionName[Action_OneKeyDress] = "一键穿戴"
	ActionName[Action_FightBoss] = "挑战boss"
	ActionName[Action_AllEquipMelting] = "装备全部熔炼"
	ActionName[Action_HnagupOnLine] = "在线挂机"
	ActionName[Action_RequestNiuDanInfo] = "请求扭蛋数据"
	ActionName[Action_GetDiscountShopInfo] = "获取折扣商店信息"
}

//玩家对象
type Player struct {
	session              string                //会话标识（时间戳）
	guid                 string                //服务器返回的唯一id
	serverId             string                //服务器id
	loginName            string                //玩家名字
	curGuideStep         int                   //当前引导步骤
	curActionType        int                   //当前执行的行为
	RequestDataDealFuncs []RequestDataDealFunc //请求数据处理函数
	ResultDealFuncs      []ResultDealFunc      //结果处理函数
}
type RequestDataDealFunc func(this *Player, urlValues *url.Values)
type ResultDealFunc func(this *Player, result string)

//初始化一个玩家并登陆
func (this *Player) Init(index int) {
	this.loginName = PlayerNameFlag + strconv.Itoa(index)
	this.serverId = "1"

	//对请求的数据进行处理
	this.RequestDataDealFuncs = make([]RequestDataDealFunc, Action_Max)
	this.RequestDataDealFuncs[Action_CreateRole] = CreateRole_DealRequstData
	this.RequestDataDealFuncs[Action_Guide] = Guide_DealRequstData
	this.RequestDataDealFuncs[Action_SelectStartRole] = SelectStartRole_DealRequstData
	this.RequestDataDealFuncs[Action_NiuDan] = DealNothing2RequestData
	this.RequestDataDealFuncs[Action_DrawCard] = DrawCard_DealRequstData
	this.RequestDataDealFuncs[Action_2Battle] = _2Battle_DealRequstData
	this.RequestDataDealFuncs[Action_HangupDup] = DealNothing2RequestData
	this.RequestDataDealFuncs[Action_OneKeyDress] = OneKeyDress_DealRequstData
	this.RequestDataDealFuncs[Action_FightBoss] = DealNothing2RequestData
	this.RequestDataDealFuncs[Action_AllEquipMelting] = AlLEquipMelting_DealRequstData
	this.RequestDataDealFuncs[Action_HnagupOnLine] = DealNothing2RequestData
	this.RequestDataDealFuncs[Action_RequestNiuDanInfo] = DealNothing2RequestData
	this.RequestDataDealFuncs[Action_GetDiscountShopInfo] = DealNothing2RequestData

	//对结果进行处理
	this.ResultDealFuncs = make([]ResultDealFunc, Action_Max)
	this.ResultDealFuncs[Action_CreateRole] = DealNothing2Result
	this.ResultDealFuncs[Action_Guide] = Guide_DealResult
	this.ResultDealFuncs[Action_SelectStartRole] = DealNothing2Result
	this.ResultDealFuncs[Action_NiuDan] = DealNothing2Result
	this.ResultDealFuncs[Action_DrawCard] = DealNothing2Result
	this.ResultDealFuncs[Action_2Battle] = DealNothing2Result
	this.ResultDealFuncs[Action_HangupDup] = DealNothing2Result
	this.ResultDealFuncs[Action_OneKeyDress] = DealNothing2Result
	this.ResultDealFuncs[Action_FightBoss] = FightBoss_DealResult
	this.ResultDealFuncs[Action_AllEquipMelting] = DealNothing2Result
	this.ResultDealFuncs[Action_HnagupOnLine] = DealNothing2Result
	this.ResultDealFuncs[Action_RequestNiuDanInfo] = DealNothing2Result
	this.ResultDealFuncs[Action_GetDiscountShopInfo] = DealNothing2Result

	//开始登陆，构造登陆数据
	urlValues := url.Values{}
	urlValues.Add("fm[cmd]", "10000")
	urlValues.Add("fm[loginname]", this.loginName)
	urlValues.Add("fm[serverid]", this.serverId)
	urlValues.Add("fm[channel]", "test")
	urlValues.Add("fm[token]", "")
	urlValues.Add("fm[gameid]", "")
	urlValues.Add("fm[login_type]", "")

	//向服务器请求数据
	log.Waring(this.loginName, "登陆")
	err, rawResult := HttpPostForm(serverUlr, urlValues)
	if nil != err {
		log.Error(this.loginName, "登陆异常：", err.Error())
		return
	}

	//服务器两个结果拼在一起，json串格式不对，特殊处理
	results := strings.Split(rawResult, "},")
	result := results[0] + "}"
	result = result[1:]

	//解析登陆结果
	var loginResultData LoginResult
	errLoginResultJson := json.Unmarshal([]byte(result), &loginResultData)
	if nil != errLoginResultJson {
		log.Error(this.loginName, "登陆结果解析异常：", errLoginResultJson.Error(), ",result:", result)
		return
	}
	this.curGuideStep, _ = strconv.Atoi(loginResultData.M.Guidestep)
	this.session = strconv.FormatInt(loginResultData.M.Session, 10)
	this.guid = loginResultData.M.Guid
	if loginResultData.M.Createrole == 1 {
		//创角
		this.ActionByType(Action_CreateRole)
	} else if this.curGuideStep < GuideStep_2Battle2 {
		//引导未走完
		this.ActionByType(Action_Guide)
	} else {
		this.DoSomethingByRand()
	}
}

//随机执行一些行为
func (this *Player) DoSomethingByRand() {
	testFuncIndex := rand.Intn(Action_Max - Action_NewPlayerActionMax)
	this.ActionByType(testFuncIndex + Action_NewPlayerActionMax)
}

//按顺序执行新玩家行为
func (this *Player) DoNewPlayerAction() {
	this.curActionType++
	if this.curActionType < Action_NewPlayerActionMax {
		this.ActionByType(this.curActionType)
	} else {
		this.DoSomethingByRand()
	}
}

//根据类型执行行为
func (this *Player) ActionByType(actionType int) {
	if actionType >= Action_Max {
		log.Error("行为错误：", actionType)
		return
	}

	//休眠一段时间再执行操作
	timeSpace := rand.Intn(maxTimeSpaceOfPlayerAction-minTimeSpaceOfPlayerAction) + minTimeSpaceOfPlayerAction
	time.Sleep(time.Millisecond * time.Duration(timeSpace))

	//构造通用请求数据
	this.curActionType = actionType
	urlValues := url.Values{}
	urlValues.Add("fm[guid]", this.guid)
	urlValues.Add("fm[cmd]", ActionCmd[actionType])
	urlValues.Add("fm[session]", this.session)

	//处理每个行为自己的请求数据
	this.RequestDataDealFuncs[actionType](this, &urlValues)

	//向服务器请求数据
	if actionType == Action_Guide {
		log.Waring(this.loginName, ActionName[actionType], ",guideStep:", this.curGuideStep)
	} else {
		log.Waring(this.loginName, ActionName[actionType])
	}
	err, rawResult := HttpPostForm(serverUlr, urlValues)
	if nil != err {
		log.Error(this.loginName, ActionName[actionType]+"异常：", err.Error())
		return
	}

	//结果处理函数
	this.ResultDealFuncs[actionType](this, rawResult)
}

//*******************************************请求数据处理开始************************************
//创角请求数据特殊处理
func CreateRole_DealRequstData(this *Player, urlValues *url.Values) {
	//名字就为角色名处理
	urlValues.Add("fm[rolename]", this.loginName)
	urlValues.Add("fm[type]", "0")
	urlValues.Add("fm[serverid]", this.serverId)
}

//引导请求数据特殊处理
func Guide_DealRequstData(this *Player, urlValues *url.Values) {
	this.curGuideStep += 1
	urlValues.Add("fm[step]", strconv.Itoa(this.curGuideStep))
}

//选择初始神将请求数据特殊处理
func SelectStartRole_DealRequstData(this *Player, urlValues *url.Values) {
	urlValues.Add("fm[id]", "2")
}

//抽卡请求数据特殊处理
func DrawCard_DealRequstData(this *Player, urlValues *url.Values) {
	if this.curGuideStep == GuideStep_DrawCard {
		urlValues.Add("fm[type]", "1")
	} else {
		urlValues.Add("fm[type]", "3")
	}
}

//上阵请求数据特殊处理
func _2Battle_DealRequstData(this *Player, urlValues *url.Values) {
	if this.curGuideStep == GuideStep_2Battle {
		urlValues.Add("fm[pos]", "1")
		urlValues.Add("fm[mercenaryid]", "101032")
		urlValues.Add("fm[type]", "-1")
	} else {
		urlValues.Add("fm[pos]", "2")
		urlValues.Add("fm[mercenaryid]", "101052")
		urlValues.Add("fm[type]", "-1")
	}
}

//一键穿戴请求数据特殊处理
func OneKeyDress_DealRequstData(this *Player, urlValues *url.Values) {
	urlValues.Add("fm[reelpos]", "0")
}

//装备全部熔炼请求数据特殊处理
func AlLEquipMelting_DealRequstData(this *Player, urlValues *url.Values) {
	urlValues.Add("fm[star]", "1;2;3;4;")
}

//请求数据共用处理（不用处理）
func DealNothing2RequestData(this *Player, urlValue *url.Values) {
}

//*******************************************请求数据处理结束************************************

//*******************************************结果处理开始************************************
//引导结果特殊处理
func Guide_DealResult(this *Player, result string) {
	if this.curGuideStep == GuideStep_First || this.curGuideStep == GuideStep_Fourth {
		//继续下一个引导
		this.ActionByType(Action_Guide)
	} else if this.curGuideStep == GuideStep_NiuDan {
		this.ActionByType(Action_NiuDan)
	} else if this.curGuideStep == GuideStep_DrawCard {
		this.ActionByType(Action_DrawCard)
	} else if this.curGuideStep == GuideStep_2Battle || this.curGuideStep == GuideStep_2Battle2 {
		this.ActionByType(Action_2Battle)
	} else if this.curGuideStep == GuideStep_HeightNiuDan {
		this.ActionByType(Action_DrawCard)
	} else {
		this.DoNewPlayerAction()
	}
}

//挑战副本boss结果特殊处理
func FightBoss_DealResult(this *Player, result string) {
	//解析结果
	var resultData []interface{}
	errResultJson := json.Unmarshal([]byte(result), &resultData)
	if nil != errResultJson {
		log.Error(this.loginName, "挑战boss结果解析异常：", errResultJson.Error(), result)
		this.DoSomethingByRand()
		return
	}
	resultMap, ok := resultData[0].(map[string]interface{})
	if !ok {
		this.DoSomethingByRand()
		return
	}
	if resultMap["cmd"].(float64) == 10050 {
		//熔炼
		this.ActionByType(Action_AllEquipMelting)
	} else {
		this.DoSomethingByRand()
	}
}

//共用结果处理函数（不用处理）
func DealNothing2Result(this *Player, result string) {
	if this.curActionType < Action_NewPlayerActionMax {
		//新玩家指令处理
		if this.curActionType == Action_SelectStartRole ||
			this.curActionType == Action_NiuDan ||
			this.curActionType == Action_DrawCard {
			//执行引导
			this.ActionByType(Action_Guide)
		} else if this.curActionType == Action_2Battle {
			if this.curGuideStep == GuideStep_2Battle {
				this.ActionByType(Action_NiuDan)
			} else {
				this.DoNewPlayerAction()
			}
		} else {
			this.DoNewPlayerAction()
		}
	} else {
		//老玩家随机执行一些指令
		this.DoSomethingByRand()
	}

}

//*******************************************结果处理结束************************************
