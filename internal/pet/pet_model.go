package pet

type Pet struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	OwnerID uint   `json:"ownerId" gorm:"not null;index"`
	Name    string `json:"name"`
	Age     uint   `json:"age"`
}
