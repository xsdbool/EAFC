package api

type Currency struct {
	Name       string `json:"name"`
	Funds      int    `json:"funds"`
	FinalFunds int    `json:"finalFunds"`
}

type ItemData struct {
	ID                 int      `json:"id"`
	Timestamp          int      `json:"timestamp"`
	Formation          string   `json:"formation"`
	Untradeable        bool     `json:"untradeable"`
	AssetId            int      `json:"assetId"`
	Rating             int      `json:"rating"`
	ItemType           string   `json:"itemType"`
	ResourceId         int      `json:"resourceId"`
	Owners             int      `json:"owners"`
	DiscardValue       int      `json:"discardValue"`
	ItemState          string   `json:"itemState"`
	CardSubtypeID      int      `json:"cardsubtypeid"`
	LastSalePrice      int      `json:"lastSalePrice"`
	InjuryType         string   `json:"injuryType"`
	InjuryGames        int      `json:"injuryGames"`
	PreferredPosition  string   `json:"preferredPosition"`
	Contract           int      `json:"contract"`
	TeamID             int      `json:"teamid"`
	RareFlag           int      `json:"rareflag"`
	PlayStyle          int      `json:"playStyle"`
	LeagueID           int      `json:"leagueId"`
	Assists            int      `json:"assists"`
	LifetimeAssists    int      `json:"lifetimeAssists"`
	LoyaltyBonus       int      `json:"loyaltyBonus"`
	Pile               int      `json:"pile"`
	Nation             int      `json:"nation"`
	ResourceGameYear   int      `json:"resourceGameYear"`
	AttributeArray     []int    `json:"attributeArray"`
	StatsArray         []int    `json:"statsArray"`
	LifetimeStatsArray []int    `json:"lifetimeStatsArray"`
	SkillMoves         int      `json:"skillmoves"`
	WeakFootAbility    int      `json:"weakfootabilitytypecode"`
	AttackingWorkRate  int      `json:"attackingworkrate"`
	DefensiveWorkRate  int      `json:"defensiveworkrate"`
	PreferredFoot      int      `json:"preferredfoot"`
	PossiblePositions  []string `json:"possiblePositions"`
	Gender             int      `json:"gender"`
	BaseTraits         []int    `json:"baseTraits"`
	IconTraits         []int    `json:"iconTraits"`
}

type AuctionInfo struct {
	TradeId           int      `json:"tradeId"`
	ItemData          ItemData `json:"itemData"`
	TradeState        string   `json:"tradeState"`
	BuyNowPrice       int      `json:"buyNowPrice"`
	CurrentBid        int      `json:"currentBid"`
	Offers            int      `json:"offers"`
	Watched           bool     `json:"watched"`
	BidState          string   `json:"bidState"`
	StartingBid       int      `json:"startingBid"`
	ConfidenceValue   int      `json:"confidenceValue"`
	Expires           int      `json:"expires"`
	SellerName        string   `json:"sellerName"`
	SellerEstablished int      `json:"sellerEstablished"`
	SellerId          int      `json:"sellerId"`
	TradeOwner        bool     `json:"tradeOwner"`
	TradeIdStr        string   `json:"tradeIdStr"`
}

type TransfermarketResponse struct {
	AuctionInfo []AuctionInfo `json:"auctionInfo"`
	BidTokens   struct{}      `json:"bidTokens"`
}

type BidResponse struct {
	Credits     int           `json:"credits"`
	AuctionInfo []AuctionInfo `json:"auctionInfo"`
	BidTokens   struct{}      `json:"bidTokens"`
	Currencies  []Currency    `json:"currencies"`
}

type WatchlistResponse struct {
	Total       int           `json:"total"`
	Credits     int           `json:"credits"`
	AuctionInfo []AuctionInfo `json:"auctionInfo"`
}

type TradeIdList struct {
	Id    int    `json:"id"`
	IdStr string `json:"idStr"`
}

type ObjectiveProgressList struct {
	ObjectiveId   int `json:"objectiveId"`
	State         int `json:"state"`
	ProgressCount int `json:"progressCount"`
}

type ScmpGroupProgressList struct {
	GroupId               int                     `json:"groupId"`
	State                 int                     `json:"state"`
	ObjectiveProgressList []ObjectiveProgressList `json:"objectiveProgressList"`
	GroupType             int                     `json:"groupType"`
}

type LearningGroupProgressList struct {
	CategoryId            int                     `json:"categoryId"`
	ScmpGroupProgressList []ScmpGroupProgressList `json:"scmpGroupProgressList"`
}

type ProgressOnAcademyObjectives struct {
	CompletedObjectivesCount int  `json:"completedObjectivesCount"`
	HasReachedLevel          bool `json:"hasReachedLevel"`
}

type DynamicObjectivesUpdates struct {
	NeedsGroupsRefresh          bool                        `json:"needsGroupsRefresh"`
	LearningGroupProgressList   []LearningGroupProgressList `json:"learningGroupProgressList"`
	NeedsAutoClaim              bool                        `json:"needsAutoClaim"`
	NeedsMilestonesAutoClaim    bool                        `json:"needsMilestonesAutoClaim"`
	ProgressOnAcademyObjectives ProgressOnAcademyObjectives `json:"progressOnAcademyObjectives"`
}

type PileItemResult struct {
	ID      int    `json:"id"`
	Pile    string `json:"pile"`
	Success bool   `json:"success"`
}

type ItemResponse struct {
	ItemData []PileItemResult `json:"itemData"`
}

type RelistResponse struct {
	TradeIdList              []TradeIdList            `json:"tradeIdList"`
	DynamicObjectivesUpdates DynamicObjectivesUpdates `json:"dynamicObjectivesUpdates"`
}

type AuctionHouseResponse struct {
	Id                       int64                    `json:"id"`
	IdStr                    string                   `json:"idStr"`
	DynamicObjectivesUpdates DynamicObjectivesUpdates `json:"dynamicObjectivesUpdates"`
}
