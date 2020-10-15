// 給予Websocket伺服器的opcode
const MESSAGE  = 0 // 聊天室訊息
const INVITE = 1 // 邀請加入聊天室
const ANSWER = 2 // 答覆聊天室邀請
const BROADCAST = 3 // 系統廣播

// 訊息框的分類
const SYSTEM = Symbol( "system" );
const USER = Symbol( "user")
const OTHER = Symbol( "other" );