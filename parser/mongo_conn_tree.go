package parser

// type ConnectionTree struct {
// 	Tree map[string]*ActiveInternetConnection
// }

func IsConnectingToMongo(port string) bool {
	return port == "27017"
}

type MongoConnectionTree map[string]ActiveInternetConnectionArray

func New_MongoConnectionTree() MongoConnectionTree {
	return make(MongoConnectionTree)
}

func (mct MongoConnectionTree) SortToTree(conns ActiveInternetConnectionArray) MongoConnectionTree {
	for _, conn := range conns {
		if IsConnectingToMongo(conn.ForeignPort) {
			mct.addConn(conn.ProgramName, conn)
		}
	}

	return mct
}

func (mct MongoConnectionTree) addConn(key string, conn *ActiveInternetConnection) {
	if len(conn.ProgramName) <= 0 {
		return
	}

	if list, ok := mct[key]; ok {
		mct[key] = append(list, conn)
	} else {
		mct[key] = ActiveInternetConnectionArray{conn}
	}
}

func (mct MongoConnectionTree) ConnStatistics() map[string]int {
	result := make(map[string]int)
	for key, conns := range mct {
		result[key] = len(conns)
	}
	return result
}
