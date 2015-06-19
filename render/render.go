package render

const (
	Black   = "black"
	Red     = "red"
	Green   = "green"
	Yellow  = "yellow"
	Blue    = "blue"
	Magenta = "magenta"
	Cyan    = "cyan"
	Gray    = "gray"
)

type Renderer func(string) (string, error)

type Markupper interface {
	StartColour(string) interface{}
	EndColour() interface{}
	StartBold() interface{}
	EndBold() interface{}
}

type Context struct{}

// @see http://en.wikipedia.org/wiki/ANSI_escape_code#Colours
func ValidColours() []string {
	return []string{
		Black,
		Red,
		Green,
		Yellow,
		Blue,
		Magenta,
		Cyan,
		Gray,
	}
}

func IsValidColour(c string) bool {
	for _, validColour := range ValidColours() {
		if validColour == c {
			return true
		}
	}
	return false
}

func AttachTemplateFuncs(to map[string]interface{}, m Markupper) map[string]interface{} {
	to["c"] = func(colour string) interface{} {
		if !IsValidColour(colour) {
			panic(colour + " is not a valid colour")
		}
		return m.StartColour(colour)
	}
	to["_c"] = func() interface{} {
		return m.EndColour()
	}
	to["b"] = func() interface{} {
		return m.StartBold()
	}
	to["_b"] = func() interface{} {
		return m.EndBold()
	}
	// The l functions were removed but remain here for backwards compatability.
	to["l"] = func() interface{} {
		return ""
	}
	to["_l"] = func() interface{} {
		return ""
	}
	return to
}

func RenderTemplates(tmpls []string, renderer Renderer) (ret []string, err error) {
	ret = make([]string, len(tmpls))
	for i, t := range tmpls {
		if ret[i], err = renderer(t); err != nil {
			return
		}
	}
	return
}
