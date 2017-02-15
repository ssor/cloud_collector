package parser

var (
	DebugMode = false
)

func print_debug(f func()) {
	if DebugMode {
		if f != nil {
			f()
		}
	}
}
