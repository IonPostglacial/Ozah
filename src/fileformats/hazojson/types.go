package hazojson

import "strings"

type Photo struct {
	Id     string `json:"id"`
	Url    string `json:"url"`
	HubUrl string `json:"hubUrl"`
	Label  string `json:"label"`
}

type State struct {
	Id          string   `json:"id"`
	Path        []string `json:"path"`
	Name        string   `json:"name"`
	NameEN      string   `json:"nameEN"`
	NameCN      string   `json:"nameCN"`
	Photos      []Photo  `json:"photos"`
	Description string   `json:"description"`
	Color       string   `json:"color,omitempty"`
}

type Book struct {
	Id    string   `json:"id"`
	Path  []string `json:"path"`
	Label string   `json:"label"`
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

type Measurement struct {
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	CharacterRef string  `json:"character"`
}

type Taxon struct {
	Id                string              `json:"id"`
	Path              []string            `json:"path"`
	ParentId          string              `json:"parentId,omitempty"`
	Name              string              `json:"name"`
	NameEN            string              `json:"nameEN"`
	NameCN            string              `json:"nameCN"`
	VernacularName    string              `json:"vernacularName"`
	Detail            string              `json:"detail"`
	Children          []string            `json:"children"`
	Photos            []Photo             `json:"photos"`
	Descriptions      []Descriptions      `json:"descriptions"`
	SpecimenLocations []Location          `json:"specimenLocations"`
	Author            string              `json:"author"`
	VernacularName2   string              `json:"vernacularName2,omitempty"`
	Name2             string              `json:"name2,omitempty"`
	Meaning           string              `json:"meaning,omitempty"`
	HerbariumPicture  string              `json:"herbariumpicture,omitempty"`
	Website           string              `json:"website,omitempty"`
	NoHerbier         string              `json:"noHerbier,omitempty"`
	Fasc              string              `json:"fasc,omitempty"`
	Page              string              `json:"page,omitempty"`
	Measurements      []Measurement       `json:"measurements"`
	BookInfoByIds     map[string]BookInfo `json:"bookInfobyids,omitempty"`
	Extra             map[string]any      `json:"extra,omitempty"`
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
	Path                  []string        `json:"path"`
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
	Type                  CharacterType   `json:"characterType"`
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

type EncodedDocument interface {
	Id() string
	Path() string
	Name() string
	NameEN() string
	NameCN() string
	NameFR() string
	NameV() string
	NameV2() string
	Description() string
	Photos() []Photo
}

type StateAsDocument struct {
	*State
}

func (s StateAsDocument) Id() string {
	return s.State.Id
}

func (s StateAsDocument) Path() string {
	return strings.Join(s.State.Path, ".")
}

func (s StateAsDocument) Name() string {
	return s.State.Name
}

func (s StateAsDocument) NameEN() string {
	return s.State.NameEN
}

func (s StateAsDocument) NameCN() string {
	return s.State.NameCN
}

func (s StateAsDocument) NameFR() string {
	return ""
}

func (s StateAsDocument) NameV() string {
	return ""
}

func (s StateAsDocument) NameV2() string {
	return ""
}

func (s StateAsDocument) Description() string {
	return s.State.Description
}

func (s StateAsDocument) Photos() []Photo {
	return s.State.Photos
}

type TaxonAsDocument struct {
	*Taxon
}

func (s TaxonAsDocument) Id() string {
	return s.Taxon.Id
}

func (s TaxonAsDocument) Path() string {
	return strings.Join(s.Taxon.Path, ".")
}

func (s TaxonAsDocument) Name() string {
	return s.Taxon.Name
}

func (s TaxonAsDocument) NameEN() string {
	return s.Taxon.NameEN
}

func (s TaxonAsDocument) NameCN() string {
	return s.Taxon.NameCN
}

func (s TaxonAsDocument) NameFR() string {
	return ""
}

func (s TaxonAsDocument) NameV() string {
	return ""
}

func (s TaxonAsDocument) NameV2() string {
	return ""
}

func (s TaxonAsDocument) Description() string {
	return s.Taxon.Detail
}

func (s TaxonAsDocument) Photos() []Photo {
	return s.Taxon.Photos
}

type CharacterAsDocument struct {
	*Character
}

func (s CharacterAsDocument) Id() string {
	return s.Character.Id
}

func (s CharacterAsDocument) Path() string {
	return strings.Join(s.Character.Path, ".")
}

func (s CharacterAsDocument) Name() string {
	return s.Character.Name
}

func (s CharacterAsDocument) NameEN() string {
	return s.Character.NameEN
}

func (s CharacterAsDocument) NameCN() string {
	return s.Character.NameCN
}

func (s CharacterAsDocument) NameFR() string {
	return ""
}

func (s CharacterAsDocument) NameV() string {
	return ""
}

func (s CharacterAsDocument) NameV2() string {
	return ""
}

func (s CharacterAsDocument) Description() string {
	return s.Character.Detail
}

func (s CharacterAsDocument) Photos() []Photo {
	return s.Character.Photos
}

type BookAsDocument struct {
	*Book
}

func (s BookAsDocument) Id() string {
	return s.Book.Id
}

func (s BookAsDocument) Path() string {
	return strings.Join(s.Book.Path, ".")
}

func (s BookAsDocument) Name() string {
	return s.Book.Label
}

func (s BookAsDocument) NameEN() string {
	return ""
}

func (s BookAsDocument) NameCN() string {
	return ""
}

func (s BookAsDocument) NameFR() string {
	return ""
}

func (s BookAsDocument) NameV() string {
	return ""
}

func (s BookAsDocument) NameV2() string {
	return ""
}

func (s BookAsDocument) Description() string {
	return ""
}

func (s BookAsDocument) Photos() []Photo {
	return []Photo{}
}
