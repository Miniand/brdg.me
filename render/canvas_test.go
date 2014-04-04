package render

func ExampleCanvasRenderTagPositioning() {
	c := &Canvas{}
	c.Draw(-1, -3, `{{b}}1{{_b}}`)
	c.Draw(1, -3, `{{c "green"}}2{{_c}}`)
	c.Draw(0, -2, `3`)
	// Output:
	// {{b}}1{{_b}} {{c "green"}}2{{_c}}
	//  3
}
