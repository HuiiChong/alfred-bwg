package main

import (
	"github.com/deanishe/awgo"
	"github.com/valyala/fasthttp"
	"github.com/json-iterator/go"
	"fmt"
	"time"
	"os"
)

type Result struct {
	Vm_type 					string		`json:"vm_type"`
	Hostname 					string		`json:"hostname"`
	Node_ip 					string		`json:"node_ip"`
	Node_alias 					string		`json:"node_alias"`
	Node_location 				string		`json:"node_location"`
	Node_location_id 			string		`json:"node_location_id"`
	Node_datacenter 			string		`json:"node_datacenter"`
	Location_ipv6_ready 		bool		`json:"location_ipv6_ready"`
	Plan 						string		`json:"plan"`
	Plan_monthly_data 			int64		`json:"plan_monthly_data"`
	Monthly_data_multiplier 	int			`json:"monthly_data_multiplier"`
	Plan_disk 					int64		`json:"plan_disk"`
	Plan_ram 					int64		`json:"plan_ram"`
	Plan_swap 					int			`json:"plan_swap"`
	Plan_max_ipv6s 				int			`json:"plan_max_ipv6s"`
	Os 							string		`json:"os"`
	Email 						string		`json:"email"`
	Data_counter 				int64		`json:"data_counter"`
	Data_next_reset 			int64		`json:"data_next_reset"`
	Ip_addresses 				[]string	`json:"ip_addresses"`
	Rdns_api_available 			bool		`json:"rdns_api_available"`
	Suspended 					bool		`json:"suspended"`
	Ptr 						string		`json:"-"`
	Error 						int			`json:"error"`
}

var wf *aw.Workflow
var query string

func init() {
	wf = aw.New()
}

func run() {
	id := os.Args[1]
	key := os.Args[2]
	query = ""
	url := "https://api.64clouds.com/v1/getServiceInfo?veid=" + id + "&api_key=" + key
	_, res, err := fasthttp.Get(nil, url)
	if err != nil {
		wf.NewItem("获取信息失败！").Subtitle("抱歉, 请更新版本或者配置正确!").Valid(false).Icon(aw.IconError)
		wf.SendFeedback()
	}
	//解析json
	var result = &Result{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if json.Unmarshal(res, result) != nil {
		wf.NewItem("获取信息失败！").Subtitle("抱歉, 请更新版本或者配置正确!").Valid(false).Icon(aw.IconError)
		wf.SendFeedback()
	}
	//
	tm := time.Unix(result.Data_next_reset, 0)
	sub := fmt.Sprintf("本月已使用: %.2f GB (%d GB, %v, %v)",
		float64(result.Data_counter)/(1024.0*1024.0*1024.8),
		result.Plan_monthly_data/(1024*1024*1024),tm.Format("2006-01-02"),
		result.Node_datacenter)
	wf.NewItem(result.Hostname).Subtitle(sub).Icon(aw.IconInfo).Arg(result.Hostname+"\n"+sub).Valid(true)
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
