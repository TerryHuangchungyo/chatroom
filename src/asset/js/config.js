// 給予Websocket伺服器的opcode
const SEND   = 0 // 傳送訊息到聊天室
const REPLY  = 1 // 聊天室訊息回覆
const INVITE = 2 // 邀請加入聊天室
const ANSWER = 3 // 答覆聊天室邀請

// 訊息框的分類
const BROADCAST = Symbol( "system" );
const USER = Symbol( "user")
const OTHER = Symbol( "other" );