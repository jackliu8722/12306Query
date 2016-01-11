package main 

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"time"
)

type cheCiInfo struct{
	Train_no string //车次号
	Station_train_code string //车次
	Start_station_telecode string //始发站编号 
	Start_station_name string //始发站名称
	End_station_telecode string //终点站编号
	End_station_name string //终点站名称
	From_station_telecode string //用户选择的起始站编号
	From_station_name string //用户选择的起始站编号
	To_station_telecode string //用户选择的目的站编号
	To_station_name string //用户选择的目的站编号
	Start_time string //出发时间
	Arrive_time string //到达时间
	Day_difference int //出发到到达跨越的天数
	Train_class_name string 
	Lishi string //经过的时间
	CanWebBuy string 
	LishiValue string 
	Yp_info string 
	Control_train_day string 
	Start_train_date string //出发日期
	Seat_feature string 
	Yp_ex string 
	Train_seat_feature string 
	Seat_types string
	Location_code string 
	From_station_no string 
	To_station_no string 
	Control_day int 
	Sale_time string //开售时间
	Is_support_card string
	Controlled_train_flag string 
	Controlled_train_message string 
	Gg_num string 
	Gr_num string 
	Qt_num string 
	Rw_num string //软卧
	Rz_num string //软座	
	Tz_num string 
	Wz_num string //无座
	Yb_num string 
	Yw_num string //硬卧
	Yz_num string //硬座
	Ze_num string //二等座
	Zy_num string //一等座
	Swz_num string //商务座
}

type cheCi struct {
	QueryLeftNewDTO cheCiInfo
	SecretStr string 
	ButtonTextInfo string 
}

type cheCiData struct {
	ValidateMessageShowId string 
	Status bool
	Httpstatus int 
	Data []cheCi 
	Messages []string 
	ValidateMessages interface{}
}

var (
	queryUrlTemplate = "https://kyfw.12306.cn/otn/leftTicket/queryT?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=ADULT"

	//需要查询的日期
	date = []string{"2016-02-04","2016-02-05","2016-02-06","2016-01-15"}
    //出发站	
	from_station = "BJP"
	//目的站
	to_station = "CQW"
	//报警的信息
	alertMap = map[string]map[string]int{
		"20160206" : map[string]int{
			"Z49" : 1,
			"G307" : 1,
			"G309" : 1,
			"Z3" : 1,
		},
		"20160205" : map[string]int{
			"Z49" : 1,
			"G307" : 1,
			"G309" : 1,
			"Z3" : 1,
		},
		"20160204" : map[string]int{
			"Z49" : 1,
			"G307" : 1,
			"G309" : 1,
			"Z3" : 1,
		},
		"20160115" : map[string]int{
			"Z49" : 1,
			"G307" : 1,
			"G309" : 1,
			"Z3" : 1,
		},
	}
	//条件值
	notIncludeValueMap = map[string]int{
		"--" : 1,
		"无" : 1, 
	}
	//接受报警的手机号
	alertPhone = "15901116997"
	//发送报警的url
	alertUrl = "http://m.cncn.com/account/getcode?mobiletel=%s&idtype=0&inajax=1"
	//上一次发送报警的时间 
	alertPrevTimestamp = int64(0)
	//发送报警的时间间隔(s)
	alertTimePadding = int64(60)

)

func init(){

}

// 获取查询的二进制数据
func getResponse(queryUrl string) []byte{
	client := &http.Client{}
	rep,_:=http.NewRequest("GET",queryUrl,strings.NewReader(""))
	rep.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.10 Safari/537.36");
	rep.Header.Set("X-Requested-With","XMLHttpRequest")

	resp,err := client.Do(rep)
	if err != nil{
		fmt.Println("Error!")
		return []byte("")
	}

	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read body error!")
		return []byte("")
	}
	return body
}

func alert(){
	currentTimestamp := time.Now().Unix()
	expandTime := alertPrevTimestamp + alertTimePadding
	if expandTime >= currentTimestamp {
		return 
	}
	alertPrevTimestamp = currentTimestamp
	url := fmt.Sprintf(alertUrl,alertPhone)
	client := &http.Client{}
	rep,_ := http.NewRequest("GET",url,strings.NewReader(""))
	rep.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.10 Safari/537.36");
    client.Do(rep)
}
func query(){
	for index1 := 0 ; index1 < len(date); index1++{
		queryUrl := fmt.Sprintf(queryUrlTemplate,date[index1],from_station,to_station);
		fmt.Println("Query for date:",date[index1])
		body := getResponse(queryUrl)
		var data cheCiData
		json.Unmarshal(body,&data)
		
		if data.Status == true && data.Httpstatus == 200 {
			for index2 := 0 ; index2 < len(data.Data); index2++ {
				cc := data.Data[index2].QueryLeftNewDTO
				if v1,ok1 := alertMap[cc.Start_train_date]; ok1 {
					if _,ok2 := v1[cc.Station_train_code] ; ok2{
						_,rw := notIncludeValueMap[cc.Rw_num]; 
						_,rz := notIncludeValueMap[cc.Rz_num];
						_,yw := notIncludeValueMap[cc.Yw_num]; 
						_,yz := notIncludeValueMap[cc.Yz_num];
						_,ze := notIncludeValueMap[cc.Ze_num]; 
						if !rw || !rz || !yw || !yz || !ze {
							//短信通知
							alert()
							fmt.Println("Alert for date:",date[index1])
						} 
					} 
				}
				
			}
		}
		// sleep 5s
		time.Sleep(5 * time.Second)
	}
}

func main(){
	for {
		query()
	}
}