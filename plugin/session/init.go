package session

var sesshub SessHub

func init() {
	sesshub = NewSessionHub()
}
//添加session
func HubAdd(sess *Session) {sesshub.Add(sess)}

//删除session
func HubClose(sess *Session) {
	sess.Close()
	sesshub.Del(sess.ID())
}