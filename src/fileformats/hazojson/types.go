package hazojson

type Photo struct {
	Id     string `json:"id"`
	Url    string `json:"url"`
	HubUrl string `json:"hubUrl"`
	Label  string `json:"label"`
}

type State struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	NameEN      string  `json:"nameEN"`
	NameCN      string  `json:"nameCN"`
	Photos      []Photo `json:"photos"`
	Description string  `json:"description"`
	Color       string  `json:"color,omitempty"`
}

type Book struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type BookInfo struct {
	Fasc   string `json:"fasc"`
	Page   string `json:"page"`
	Detail string `json:"detail"`
}

type Descriptions struct {
	DescriptorId string   `json:"descriptorId"`
	StatesIds    []string `json:"statesIds"`
}

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type Taxon struct {
	Id                string                 `json:"id"`
	ParentId          string                 `json:"parentId,omitempty"`
	Name              string                 `json:"name"`
	NameEN            string                 `json:"nameEN"`
	NameCN            string                 `json:"nameCN"`
	VernacularName    string                 `json:"vernacularName"`
	Detail            string                 `json:"detail"`
	Children          []string               `json:"children"`
	Photos            []Photo                `json:"photos"`
	Descriptions      []Descriptions         `json:"descriptions"`
	SpecimenLocations []Location             `json:"specimenLocations"`
	Author            string                 `json:"author"`
	VernacularName2   string                 `json:"vernacularName2,omitempty"`
	Name2             string                 `json:"name2,omitempty"`
	Meaning           string                 `json:"meaning,omitempty"`
	HerbariumPicture  string                 `json:"herbariumpicture,omitempty"`
	Website           string                 `json:"website,omitempty"`
	NoHerbier         string                 `json:"noHerbier,omitempty"`
	Fasc              string                 `json:"fasc,omitempty"`
	Page              string                 `json:"page,omitempty"`
	BookInfoByIds     map[string]BookInfo    `json:"bookInfobyids,omitempty"`
	Extra             map[string]interface{} `json:"extra,omitempty"`
}

func NewTaxon(id string, parentId string, name string, vernacularName string, chineseName string, website string) *Taxon {
	return &Taxon{
		Id:                id,
		ParentId:          parentId,
		Name:              name,
		VernacularName:    vernacularName,
		NameCN:            chineseName,
		Website:           website,
		Photos:            make([]Photo, 0),
		Children:          make([]string, 0),
		Descriptions:      make([]Descriptions, 0),
		SpecimenLocations: make([]Location, 0),
	}
}

type CharacterPreset string
type CharacterType string

const (
	CharacterPresetMap       CharacterPreset = "map"
	CharacterPresetFlowering CharacterPreset = "flowering"
	CharacterPresetFamily    CharacterPreset = "family"

	CharacterTypeRange    CharacterType = "range"
	CharacterTypeDiscrete CharacterType = "discrete"
)

type Character struct {
	Id                    string          `json:"id"`
	ParentId              string          `json:"parentId,omitempty"`
	Name                  string          `json:"name"`
	NameEN                string          `json:"nameEN"`
	NameCN                string          `json:"nameCN"`
	VernacularName        string          `json:"vernacularName"`
	Detail                string          `json:"detail"`
	Children              []string        `json:"children"`
	Photos                []Photo         `json:"photos"`
	Color                 string          `json:"color"`
	Preset                CharacterPreset `json:"preset"`
	Type                  CharacterType   `json:"chracterType"`
	States                []string        `json:"states"`
	InherentStateId       string          `json:"inherentStateId"`
	InapplicableStatesIds []string        `json:"inapplicableStatesIds"`
	RequiredStatesIds     []string        `json:"requiredStatesIds"`
	MapFile               string          `json:"mapFile"`
	Min                   int             `json:"min"`
	Max                   int             `json:"max"`
	Unit                  string          `json:"unit"`
}

type Dataset struct {
	Id         string       `json:"id"`
	Taxons     []*Taxon     `json:"taxons"`
	Characters []*Character `json:"characters"`
	States     []*State     `json:"states"`
	Books      []*Book      `json:"books"`
}
