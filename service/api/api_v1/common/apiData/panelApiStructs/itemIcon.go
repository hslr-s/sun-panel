package adminApiStructs

import "sun-panel/models"

type ItemIconEditRequest struct {
	models.ItemIcon
	IconJson string
}
