package po

// other persistent objects for workflow

// UserStateData user stored data for a state
type UserStateData struct {
	ID        int `json:"id" gorm:"primaryKey;autoIncrement"`
	StateID   int
	StoreData string
	Reference string
}
