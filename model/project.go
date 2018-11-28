package model

type Project struct {
	Name string			`json:"name"`
	Services []Service  `json:"services"`
	Links    []Link		`json:"links"`

}

type Service struct {
	Name 		string	`json:"name"`
	Type 		string  `json:"type"` // address type, default is ip:port
	Project 	string	`json:"project"`
	Address 	string 	`json:"address"`//
	Port    	int		`json:"port"`
	Protocol  	string	`json:"protocol"`
	Method 		string	`json:"method"`
	ReqPrefix   string	`json:"reqPrefix"`
}

type Link struct {
	Name 	string	`json:"name"`
	ServiceName string	`json:"serviceName"`
	Project string	`json:"project"`
	Query   string  `json:"query"`//file path. http or local
	HttpHeader string `json:"httpHeader"`//file path. http or local. add to every query
	TargetQps  int	`json:"targetQps"`

	Service  Service
}
