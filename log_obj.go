package logger

//LogObj 日志对象
type LogObj struct {
	logid string
	tag   string
	json  bool
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

//JSON 设置日志格式为json
func (l *LogObj) JSON() *LogObj {
	l.json = true
	return l
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
