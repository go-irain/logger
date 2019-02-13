package logger

//LogObj 日志对象
type LogObj struct {
	logid string
	tag   string
	data  interface{}
	json  bool
}

//NewLog 生成logobj对象
func NewLog(logid, tag string) *LogObj {
	return new(LogObj).ID(logid).Tag(tag)
}

//ID 设置id
func (l *LogObj) ID(id string) *LogObj {
	l.logid = id
	l.json = true
	return l
}

//Tag 设置id
func (l *LogObj) Tag(tag string) *LogObj {
	l.tag = tag
	l.json = true
	return l
}

//Data  设置数据对象
func (l *LogObj) Data(d interface{}) *LogObj {
	l.data = d
	return l
}

//JSON 设置日志格式为json
func (l *LogObj) JSON() *LogObj {
	l.json = true
	return l
}

//Logid 获取logid
func (l *LogObj) Logid() string {
	return l.logid
}

//GetTag 返回tag
func (l *LogObj) GetTag() string {
	return l.tag
}

func (l *LogObj) GetData() interface{} {
	return l.data
}

//UnJSON 设置日志非json格式
func (l *LogObj) UnJSON() *LogObj {
	l.json = false
	return l
}

//Log LOG
func (l *LogObj) Log(v ...interface{}) bool {
	return Trace(LOG, l, v...)
}

//Debug DEBUG
func (l *LogObj) Debug(v ...interface{}) bool {
	return Trace(DEBUG, l, v...)
}

//Info INFO
func (l *LogObj) Info(v ...interface{}) bool {
	return Trace(INFO, l, v...)
}

//Warn WARN
func (l *LogObj) Warn(v ...interface{}) bool {
	return Trace(WARN, l, v...)
}

//Error ERROR
func (l *LogObj) Error(v ...interface{}) bool {
	return Trace(ERROR, l, v...)
}

//Fatal FATAL
func (l *LogObj) Fatal(v ...interface{}) bool {
	return Trace(FATAL, l, v...)
}
