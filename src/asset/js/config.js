// 給予Websocket伺服器的opcode
const MESSAGE  = 0 // 聊天室訊息
const INVITE = 1 // 邀請加入聊天室
const ANSWER = 2 // 答覆聊天室邀請
const BROADCAST = 3 // 系統廣播
const USER_ONLINE = 4 // 使用者上線
const USER_OFFLINE = 5 // 使用者下線
const ADD_USER = 6 // 新使用者加入聊天室

// 訊息框的分類
const SYSTEM = Symbol( "system" );
const USER = Symbol( "user")
const OTHER = Symbol( "other" );

// 正則表達式
var urlPattern = /(((ht|f)tp(s?))\:\/\/)?(www.|[a-zA-Z].)[a-zA-Z0-9\-\.]+\.(com|edu|gov|mil|net|org|biz|info|name|museum|us|ca|uk)(\:[0-9]+)*(\/($|[a-zA-Z0-9\.\,\;\?\'\\\+&amp;%\$#\=~_\-]+))*/g;