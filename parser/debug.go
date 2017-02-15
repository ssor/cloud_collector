package parser

var (
	DebugMode = true
)

func print_debug(f func()) {
	if DebugMode {
		if f != nil {
			f()
		}
	}
}
