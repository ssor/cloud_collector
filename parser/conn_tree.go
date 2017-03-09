package parser

// type ConnectionTree struct {
// 	Tree map[string]*ActiveInternetConnection
// }

func IsConnectingToMongo(port interface{}) bool {
	if port == nil {
		return false
	}
	return port == "27017"
}

type ConnectionTree struct {
	Connections map[string]ActiveInternetConnectionArray
	Predictor   func(para interface{}) bool
}

func NewConnectionTree(predictor func(interface{}) bool) *ConnectionTree {
	return &ConnectionTree{
		Connections: make(map[string]ActiveInternetConnectionArray),
		Predictor:   predictor,
	}
}

func (mct *ConnectionTree) SortToTree(conns ActiveInternetConnectionArray) *ConnectionTree {
	for _, conn := range conns {
		if IsConnectingToMongo(conn.ForeignPort) {
			mct.addConn(conn.ProgramName, conn)
		}
	}

	return mct
}

func (mct *ConnectionTree) addConn(key string, conn *ActiveInternetConnection) {
	if len(conn.ProgramName) <= 0 {
		return
	}

	if list, ok := mct.Connections[key]; ok {
		mct.Connections[key] = append(list, conn)
	} else {
		mct.Connections[key] = ActiveInternetConnectionArray{conn}
	}
}

func (mct *ConnectionTree) ConnStatistics() map[string]int {
	result := make(map[string]int)
	for key, conns := range mct.Connections {
		result[key] = len(conns)
	}
	return result
}
