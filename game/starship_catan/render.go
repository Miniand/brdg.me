package starship_catan

import (
	"fmt"
	"log"
)

func RenderMoney(amount int) string {
	return fmt.Sprintf(`{{b}}{{c "green"}}$%d{{_c}}{{_b}}`, amount)
}

func RenderResource(resource int) string {
	if _, ok := ResourceColours[resource]; !ok {
		log.Fatalf(
			"There is no resource colour for %s (%d)",
			ResourceNames[resource],
			resource,
		)
	}
	return fmt.Sprintf(
		`{{b}}{{c "%s"}}%s{{_c}}{{_b}}`,
		ResourceColours[resource],
		ResourceNames[resource],
	)
}
