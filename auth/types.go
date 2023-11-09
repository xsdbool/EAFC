package auth

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}

type Pid struct {
	ExternalRefType           string `json:"externalRefType"`
	ExternalRefValue          string `json:"externalRefValue"`
	PidId                     int64  `json:"pidId"`
	Email                     string `json:"email"`
	EmailStatus               string `json:"emailStatus"`
	Strength                  string `json:"strength"`
	Dob                       string `json:"dob"`
	Country                   string `json:"country"`
	Language                  string `json:"language"`
	Locale                    string `json:"locale"`
	Status                    string `json:"status"`
	ReasonCode                string `json:"reasonCode"`
	TosVersion                string `json:"tosVersion"`
	ParentalEmail             string `json:"parentalEmail"`
	ThirdPartyOptin           string `json:"thirdPartyOptin"`
	GlobalOptin               string `json:"globalOptin"`
	DateCreated               string `json:"dateCreated"`
	DateModified              string `json:"dateModified"`
	LastAuthDate              string `json:"lastAuthDate"`
	RegistrationSource        string `json:"registrationSource"`
	AuthenticationSource      string `json:"authenticationSource"`
	ShowEmail                 string `json:"showEmail"`
	DiscoverableEmail         string `json:"discoverableEmail"`
	AnonymousPid              string `json:"anonymousPid"`
	UnderagePid               string `json:"underagePid"`
	DefaultBillingAddressUri  string `json:"defaultBillingAddressUri"`
	DefaultShippingAddressUri string `json:"defaultShippingAddressUri"`
	PasswordSignature         int    `json:"passwordSignature"`
}

type IdentityMeResponse struct {
	Pid Pid `json:"pid"`
}

type Club struct {
	Year           string           `json:"year"`
	AssetId        int              `json:"assetId"`
	TeamId         int              `json:"teamId"`
	LastAccessTime int64            `json:"lastAccessTime"`
	Platform       string           `json:"platform"`
	ClubName       string           `json:"clubName"`
	ClubAbbr       string           `json:"clubAbbr"`
	Established    int64            `json:"established"`
	DivisionOnline int              `json:"divisionOnline"`
	BadgeId        int              `json:"badgeId"`
	SkuAccessList  map[string]int64 `json:"skuAccessList"`
	ActiveHomeKit  int              `json:"activeHomeKit"`
	ActiveCaptain  int              `json:"activeCaptain"`
}

type Persona struct {
	PersonaId     int64  `json:"personaId"`
	PersonaName   string `json:"personaName"`
	ReturningUser int    `json:"returningUser"`
	OnlineAccess  bool   `json:"onlineAccess"`
	Trial         bool   `json:"trial"`
	UserClubList  []Club `json:"userClubList"`
	TrialFree     bool   `json:"trialFree"`
}

type UserAccountInfo struct {
	Personas []Persona `json:"personas"`
}

type AccountInfoResponse struct {
	UserAccountInfo UserAccountInfo `json:"userAccountInfo"`
	NucEnabled      bool            `json:"nucEnabled"`
}

type UTASAuthIdentification struct {
	AuthCode    string `json:"authCode"`
	RedirectUrl string `json:"redirectUrl"`
}

type AccessCodeResponse struct {
	Code string `json:"code"`
}

type UTASAuthPayload struct {
	ClientVersion    int                    `json:"clientVersion"`
	Ds               string                 `json:"ds"`
	GameSku          string                 `json:"gameSku"`
	Identification   UTASAuthIdentification `json:"identification"`
	IsReadOnly       bool                   `json:"isReadOnly"`
	Locale           string                 `json:"locale"`
	Method           string                 `json:"method"`
	NucleusPersonaId int64                  `json:"nucleusPersonaId"`
	PriorityLevel    int                    `json:"priorityLevel"`
	Sku              string                 `json:"sku"`
}
type UTASAuthResponse struct {
	Protocol      string `json:"protocol"`
	IpPort        string `json:"ipPort"`
	SID           string `json:"sid"`
	PhishingToken string `json:"phishingToken"`
}
