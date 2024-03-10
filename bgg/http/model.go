package http

type CollectionXML struct {
	TotalItems int                 `xml:"totalitems,attr"`
	TermsOfUse string              `xml:"termsofuse,attr"`
	PubDate    string              `xml:"pubdate,attr"`
	Items      []CollectionItemXML `xml:"item"`
}

type CollectionItemXML struct {
	ObjectType    string              `xml:"objecttype,attr"`
	ObjectID      string              `xml:"objectid,attr"`
	Subtype       string              `xml:"subtype,attr"`
	CollID        string              `xml:"collid,attr"`
	Name          string              `xml:"name"`
	YearPublished int                 `xml:"yearpublished"`
	Image         string              `xml:"image"`
	Thumbnail     string              `xml:"thumbnail"`
	Status        CollectionStatusXML `xml:"status"`
	NumPlays      int                 `xml:"numplays"`
}

type CollectionStatusXML struct {
	Own          int    `xml:"own,attr"`
	PrevOwned    int    `xml:"prevowned,attr"`
	ForTrade     int    `xml:"fortrade,attr"`
	Want         int    `xml:"want,attr"`
	WantToPlay   int    `xml:"wanttoplay,attr"`
	WantToBuy    int    `xml:"wanttobuy,attr"`
	Wishlist     int    `xml:"wishlist,attr"`
	Preordered   int    `xml:"preordered,attr"`
	LastModified string `xml:"lastmodified,attr"`
}

type ThingsXML struct {
	TotalItems int            `xml:"totalitems,attr"`
	TermsOfUse string         `xml:"termsofuse,attr"`
	PubDate    string         `xml:"pubdate,attr"`
	Items      []ThingItemXML `xml:"item"`
}

type ThingItemXML struct {
	Type          string               `xml:"type,attr"`
	ID            int                  `xml:"id,attr"`
	Thumbnail     string               `xml:"thumbnail"`
	Image         string               `xml:"image"`
	Name          []ThingNameXML       `xml:"name"`
	Description   string               `xml:"description"`
	YearPublished ThingYearXML         `xml:"yearpublished"`
	MinPlayers    ThingMinMaxXML       `xml:"minplayers"`
	MaxPlayers    ThingMinMaxXML       `xml:"maxplayers"`
	PlayingTime   ThingMinMaxXML       `xml:"playingtime"`
	MinPlayTime   ThingMinMaxXML       `xml:"minplaytime"`
	MaxPlayTime   ThingMinMaxXML       `xml:"maxplaytime"`
	MinAge        ThingMinMaxXML       `xml:"minage"`
	Poll          []ThingPollXML       `xml:"poll"`
	Statistics    []ThingStatisticsXML `xml:"statistics"`
}

type ThingPollXML struct {
	Name       string            `xml:"name,attr"`
	Title      string            `xml:"title,attr"`
	TotalVotes string            `xml:"totalvotes,attr"`
	Results    []ThingResultsXML `xml:"results"`
}

type ThingResultsXML struct {
	NumPlayers string           `xml:"numplayers,attr"`
	Results    []ThingResultXML `xml:"result"`
}

type ThingStatisticsXML struct {
	Ratings []ThingRatingsXML `xml:"ratings"`
}

type ThingRatingsXML struct {
	AverageWeight AverageWeightXML `xml:"averageweight"`
}

func (t *ThingResultsXML) GetBestNumVotes() int {
	for _, result := range t.Results {
		if result.Value == "Best" {
			return result.NumVotes
		}
	}

	return -1
}

type ThingResultXML struct {
	Value    string `xml:"value,attr"`
	NumVotes int    `xml:"numvotes,attr"`
}

type ThingNameXML struct {
	Type      string `xml:"type,attr"`
	SortIndex int    `xml:"sortindex,attr"`
	Value     string `xml:"value,attr"`
}

type ThingYearXML struct {
	Value int `xml:"value,attr"`
}

type ThingMinMaxXML struct {
	Value int `xml:"value,attr"`
}

type AverageWeightXML struct {
	Value float64 `xml:"value,attr"`
}
