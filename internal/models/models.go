package models

/**
 * Data Structure to store the user request
 * UserID ->  User wants the ProductID
 * ProductID -> the particular product which UserID taken
 */
type BuyRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}

/**
 * Data Structure to store the global response
 * Success -> we are telling here success or not
 * Message -> the particular message you want to send to user
 * Data -> its optional like you want to send the data to user or not
 */
type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

/**
 * Data Structure which runs in background to store data in our db
 * UserID -> particular user which bought product
 * ProductID -> the product which bought buy the UserID
 * Status -> is that stored in db or not
 * Timestamp -> the particular timestamp when we storing product in db
 */
type OrderMessage struct {
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}
