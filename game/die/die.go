package die

func Render(value int) (output string) {
	switch value {
	case 1:
		output = "⚀"
	case 2:
		output = "⚁"
	case 3:
		output = "⚂"
	case 4:
		output = "⚃"
	case 5:
		output = "⚄"
	case 6:
		output = "⚅"
	default:
		panic("A die value is between 1 and 6")
	}
	return
}
