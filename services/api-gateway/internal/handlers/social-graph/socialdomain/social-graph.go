package socialdomain

type CreateRequestHTTPRequest struct {
	ToUserID string `json:"to_user_id"`
}

type CreateRequestHTTPResponse struct {
	FriendRequest
}

type AcceptRequestHTTPResponse struct {
	FriendRequest
}

type RejectRequestHTTPResponse struct {
	FriendRequest
}

type CancelRequestHTTPResponse struct {
	FriendRequest
}
type FriendRequest struct {
	ID            string `json:"id"`
	FromUserID    string `json:"from_user_id"`
	ToUserID      string `json:"to_user_id"`
	Status        string `json:"status"`
	CreatedAtUnix int64  `json:"created_at_unix"`
	UpdatedAtUnix int64  `json:"updated_at_unix"`
}
type ListFriendsHTTPResponse struct {
	Friends []FriendShip `json:"friends"`
}
type FriendShip struct {
	FriendID      string `json:"user_id"`
	CreatedAtUnix int64  `json:"created_at_unix"`
}

type ListIncomingHTTPResponse struct {
	Incoming []FriendRequest `json:"incoming"`
}
type ListOutgoingHTTPResponse struct {
	Outgoing []FriendRequest `json:"outgoing"`
}

type GetStatusHTTPResponse struct {
	Status string
	ReqID  string
}
