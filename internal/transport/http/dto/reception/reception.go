package reception

type CreateReceptionRequest struct {
    PvzID string `json:"pvzId"`
}

type ReceptionResponse struct {
    ID       string `json:"id"`
    PvzID    string `json:"pvzId"`
    DateTime string `json:"dateTime"`
    Status   string `json:"status"`
}