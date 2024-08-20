package duguncom

type LeadDetails struct {
	Name               string `json:"name"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	OrganizationTypeID int    `json:"organizationTypeId"`
	AgentID            int    `json:"agentId"`
	NotWedding         bool   `json:"notWedding"`
}

type Lead struct {
	ID          int         `json:"id"`
	ProviderID  int         `json:"providerId"`
	CoupleID    int         `json:"coupleId"`
	LeadDetails LeadDetails `json:"leadDetails"`
	Status      string      `json:"status"`
}
