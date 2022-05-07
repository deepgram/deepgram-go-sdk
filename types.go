package deepgram

type InvitationOptions struct {
  Email string `json:"email"`
  Scope string `json:"scope"`
};

type InvitationList struct{
 	Invites []InvitationOptions `json:"invites"`
};

type Message struct {
	Message string `json:"message"`
}