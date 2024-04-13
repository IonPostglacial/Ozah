package link

import (
	"fmt"
	"strings"
)

type Maker = func(dsName string, ref string) string

func ToTaxons(dsName string) string {
	return fmt.Sprintf("/ds/%s/taxons", dsName)
}

func ToTaxon(dsName string, ref string) string {
	return fmt.Sprintf("/ds/%s/taxons/%s", dsName, ref)
}

func ToCharacters(dsName string) string {
	return fmt.Sprintf("/ds/%s/characters", dsName)
}

func ToCharacter(dsName string, ref string) string {
	return fmt.Sprintf("/ds/%s/characters/%s", dsName, ref)
}

func ToDescriptor(taxonRef string) Maker {
	return func(dsName string, ref string) string {
		return fmt.Sprintf("/ds/%s/taxons/%s?d=%s", dsName, taxonRef, ref)
	}
}

func ToIdentify(dsName string) string {
	return fmt.Sprintf("/ds/%s/identify", dsName)
}

func ToDocument(dsName string, ref string) string {
	switch {
	case strings.HasPrefix(ref, "t"):
		return ToTaxon(dsName, ref)
	case strings.HasPrefix(ref, "c"):
		return ToCharacter(dsName, ref)
	default:
		return ToTaxons(dsName)
	}
}
