package models

type LocalNotification struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	StoreID  string `json:"store_id"`
	ImageURL string `json:"image_url"`
	Pincode  string `json:"pincode"`
}

type UserData struct {
	Id                string `json:"id" gorm:"primaryKey"`
	UserId            string `json:"user_id"`
	Latitude          string `json:"latitude"`
	Longitude         string `json:"longitude"`
	Address           string `json:"address"`
	Pincode           string `json:"pincode"`
	City              string `json:"city"`
	State             string `json:"state"`
	Country           string `json:"country"`
	NotificationToken string `json:"notification_token"`
}
