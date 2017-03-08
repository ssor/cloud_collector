package conn_parser

var (
	DebugMode = false
	// DebugMode = true
)

func print_debug(f func()) {
	if DebugMode {
		if f != nil {
			f()
		}
	}
}
