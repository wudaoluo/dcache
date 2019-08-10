package cache

// Command（命令）
//  命令模式（Command Pattern）是一种数据驱动的设计模式，它属于行为型模式。
//  请求以命令的形式包裹在对象中，并传给调用对象。 调用对象寻找可以处理该命令
//  的合适的对象，并把该命令传给相应的对象，该对象执行命令。
//

// 命令接口
//type Command interface {
//	Do(args interface{}) (interface{}, error)
//}
//
//// 上下文
//type CmdContext struct {
//	CmdType byte
//	Args    internal.Data
//}

// 命令管理者
type CommandHandler struct {
	CmdMap map[byte]cache
}

//// 处理命令
//func (ch *CommandHandler) Handle(ctx *CmdContext) (interface{}, error) {
//	if ctx == nil {
//		return nil, errors.New("")
//	}
//	cmd, ok := ch.CmdMap[ctx.CmdType]
//	if ok {
//		return cmd.Do(ctx.Args)
//	}
//	return nil, errors.New("invalid Command ")
//}

// 注册命令
func (ch *CommandHandler) Register(cmdType byte, cmd cache) {
	ch.CmdMap[cmdType] = cmd
}

//func init() {
//	//postCtx := &CmdContext{CmdType: "post",}
//	//putCtx := &CmdContext{CmdType: "post", Args: " Post"}
//	getCtx := &CmdContext{CmdType: internal.OP_REQ_GET}
//	nullCtx := &CmdContext{CmdType: "null", Args: " Get"}
//	//cmdHandler.Register("post", &PostCommand{})
//	HandlerReq.Register("get", &GetCommand{})
//
//	fmt.Println(cmdHandler.Handle(postCtx))
//	fmt.Println(cmdHandler.Handle(getCtx))
//	fmt.Println(cmdHandler.Handle(nullCtx))
//}
