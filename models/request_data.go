package models

type Request struct {
    Type string `json:"type"`
    Data string `json:"data"`
}

type GetValue struct {
    Key string `json:"key"`
}

type SetValue struct {
    Key string `json:"key"`
    Value string `json:"value"`
}

type Message struct {
    From string `json:"from"`
    To string `json:"to"`
    Text string `json:"text"`
}

type Answer struct {
    Error error
    Data string
}

//func (m NewMessage) GetUserChatKey() string {
//    arr:=[]string{m.From,m.To}
//    sort.Strings(arr)
//    return strings.Join(arr,"_")
//}