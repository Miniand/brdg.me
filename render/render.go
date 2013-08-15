package render

type Markupper interface {
	StartColour(string) interface{}
	EndColour() interface{}
	StartBold() interface{}
	EndBold() interface{}
	StartLarge() interface{}
	EndLarge() interface{}
}

type Context struct{}

// @see http://en.wikipedia.org/wiki/ANSI_escape_code#Colours
func ValidColours() []string {
	return []string{
		"black",
		"red",
		"green",
		"yellow",
		"blue",
		"magenta",
		"cyan",
		"gray",
	}
}

func AttachTemplateFuncs(to map[string]interface{}, m Markupper) map[string]interface{} {
	to["c"] = func(colour string) interface{} {
		found := false
		for _, validColour := range ValidColours() {
			if validColour == colour {
				found = true
				break
			}
		}
		if !found {
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
	to["l"] = func() interface{} {
		return m.StartLarge()
	}
	to["_l"] = func() interface{} {
		return m.EndLarge()
	}
	return to
}
