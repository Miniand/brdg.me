package starship_catan

type EmptyCard struct {
	UnsortableCard
}

func (e EmptyCard) String() string {
	return `{{b}}{{c "gray"}}Lost Planet{{_c}}{{_b}} (empty space)`
}
