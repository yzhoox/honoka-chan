package secretbox

import (
	"fmt"
	secretboxapischema "honoka-chan/internal/schema/api/secretbox"
	secretboxschema "honoka-chan/internal/schema/secretbox"
	"honoka-chan/internal/session"
	"math/rand"
	"slices"
)

const (
	rarityN   = 1
	rarityR   = 2
	raritySR  = 3
	rarityUR  = 4
	raritySSR = 5

	urPromotionAnimationRate = 10
)

type mainQuery interface {
	findUnitIDs(memberCategory int, supportOnly bool, rarity int) ([]int, error)
}

type sessionMainQuery struct {
	ss *session.Session
}

type sessionlessMainQuery struct{}

func (q sessionMainQuery) findUnitIDs(memberCategory int, supportOnly bool, rarity int) ([]int, error) {
	return queryUnitIDs(q.ss, memberCategory, supportOnly, rarity)
}

func (sessionlessMainQuery) findUnitIDs(memberCategory int, supportOnly bool, rarity int) ([]int, error) {
	return queryUnitIDs(nil, memberCategory, supportOnly, rarity)
}

type poolRule struct {
	memberCategory  int
	supportOnly     bool
	allowedRarities []int
}

func buildUnitLineUp(memberCategory int, page secretboxapischema.SecretBoxPage, button secretboxapischema.SecretBoxButton) []secretboxschema.UnitLineUp {
	rarities := []int{rarityUR, raritySSR, raritySR, rarityR, rarityN}
	result := make([]secretboxschema.UnitLineUp, 0, len(rarities))
	for _, rarity := range rarities {
		unitIDs, err := findUnitIDsByRule(sessionlessMainQuery{}, memberCategory, page, button, rarity)
		if err != nil || len(unitIDs) == 0 {
			continue
		}
		result = append(result, secretboxschema.UnitLineUp{
			Rarity:  rarity,
			UnitIDs: unitIDs,
		})
	}
	return result
}

func findUnitIDsByRule(q mainQuery, memberCategory int, page secretboxapischema.SecretBoxPage, button secretboxapischema.SecretBoxButton, rarity int) ([]int, error) {
	rule := buildPoolRule(memberCategory, page, button)
	if len(rule.allowedRarities) > 0 && !slices.Contains(rule.allowedRarities, rarity) {
		return []int{}, nil
	}
	return q.findUnitIDs(rule.memberCategory, rule.supportOnly, rarity)
}

func randomUnitID(ss *session.Session, rng *rand.Rand, memberCategory int, page secretboxapischema.SecretBoxPage, button secretboxapischema.SecretBoxButton, rarity int) (int, error) {
	unitIDs, err := findUnitIDsByRule(sessionMainQuery{ss: ss}, memberCategory, page, button, rarity)
	if err != nil {
		return 0, err
	}
	if len(unitIDs) == 0 {
		return 0, fmt.Errorf("no units available for secretbox=%d rarity=%d", page.SecretBoxInfo.SecretBoxID, rarity)
	}
	return unitIDs[rng.Intn(len(unitIDs))], nil
}

func buildPoolRule(memberCategory int, page secretboxapischema.SecretBoxPage, button secretboxapischema.SecretBoxButton) poolRule {
	rule := poolRule{
		memberCategory:  memberCategory,
		allowedRarities: []int{rarityR, raritySR, raritySSR, rarityUR},
	}

	if isThanksFestivalBox(page) {
		rule.allowedRarities = []int{raritySSR, rarityUR}
	}

	switch page.SecretBoxInfo.SecretBoxID {
	case 1, 61:
		rule.memberCategory = 0
		rule.allowedRarities = []int{rarityN}
		return rule
	case 23:
		rule.memberCategory = 1
		rule.supportOnly = true
		rule.allowedRarities = []int{rarityN, rarityR, raritySR, rarityUR, raritySSR}
		return rule
	}

	switch button.SecretBoxButtonType {
	case 6:
		rule.memberCategory = 0
		rule.allowedRarities = []int{rarityN}
	case 9:
		rule.allowedRarities = []int{raritySR, rarityUR}
	case 10:
		rule.allowedRarities = []int{raritySSR, rarityUR}
	case 11:
		rule.allowedRarities = []int{rarityUR}
	case 14, 15:
		rule.allowedRarities = []int{raritySR, raritySSR, rarityUR}
	}

	return rule
}

func isThanksFestivalBox(page secretboxapischema.SecretBoxPage) bool {
	switch page.SecretBoxInfo.SecretBoxID {
	case 1718, 1719, 1720, 1721:
		return true
	default:
		return false
	}
}

func pickRarity(rng *rand.Rand, page secretboxapischema.SecretBoxPage, button secretboxapischema.SecretBoxButton, guaranteed bool) int {
	rule := buildPoolRule(0, page, button)
	if len(rule.allowedRarities) == 1 {
		return rule.allowedRarities[0]
	}

	weights := map[int]int{
		rarityR:   70,
		raritySR:  20,
		raritySSR: 7,
		rarityUR:  3,
	}

	if isThanksFestivalBox(page) {
		weights = map[int]int{
			raritySSR: 70,
			rarityUR:  30,
		}
	}

	switch button.SecretBoxButtonType {
	case 9:
		weights = map[int]int{raritySR: 80, rarityUR: 20}
	case 10:
		weights = map[int]int{raritySSR: 80, rarityUR: 20}
	case 11:
		return rarityUR
	case 14, 15:
		weights = map[int]int{raritySR: 80, raritySSR: 15, rarityUR: 5}
	}

	if guaranteed {
		return pickGuaranteedRarity(rng, page, button)
	}
	return weightedPick(rng, rule.allowedRarities, weights)
}

func pickGuaranteedRarity(rng *rand.Rand, page secretboxapischema.SecretBoxPage, button secretboxapischema.SecretBoxButton) int {
	rule := buildPoolRule(0, page, button)
	weights := map[int]int{raritySR: 20, raritySSR: 7, rarityUR: 3}
	return weightedPick(rng, rule.allowedRarities, weights)
}

func hasGuaranteedUR(button secretboxapischema.SecretBoxButton, cost secretboxapischema.SecretBoxCost) bool {
	if cost.UnitCount <= 1 {
		return false
	}

	switch button.SecretBoxButtonType {
	case 2, 16:
		return true
	default:
		return false
	}
}

func satisfiesGuaranteedRarity(rarity int, guaranteedRarity int) bool {
	switch guaranteedRarity {
	case rarityUR:
		return rarity == rarityUR
	case raritySR:
		return rarity == raritySR || rarity == raritySSR || rarity == rarityUR
	default:
		return rarity == guaranteedRarity
	}
}

func rarityQualityRank(rarity int) int {
	switch rarity {
	case rarityN:
		return 1
	case rarityR:
		return 2
	case raritySR:
		return 3
	case raritySSR:
		return 4
	case rarityUR:
		return 5
	default:
		return 0
	}
}

func isLowerQualityRarity(left int, right int) bool {
	return rarityQualityRank(left) < rarityQualityRank(right)
}

func shouldShowDirectURHit(rng *rand.Rand, page secretboxapischema.SecretBoxPage, button secretboxapischema.SecretBoxButton, rarity int) bool {
	if rarity != rarityUR {
		return false
	}

	rule := buildPoolRule(0, page, button)
	for _, allowedRarity := range rule.allowedRarities {
		if allowedRarity != rarityUR {
			return rng.Intn(100) >= urPromotionAnimationRate
		}
	}

	return true
}

func weightedPick(rng *rand.Rand, allowed []int, weights map[int]int) int {
	filtered := make([]int, 0, len(allowed))
	total := 0
	for _, rarity := range allowed {
		weight := weights[rarity]
		if weight <= 0 {
			continue
		}
		filtered = append(filtered, rarity)
		total += weight
	}
	if len(filtered) == 0 {
		return allowed[len(allowed)-1]
	}

	n := rng.Intn(total) + 1
	acc := 0
	for _, rarity := range filtered {
		acc += weights[rarity]
		if n <= acc {
			return rarity
		}
	}
	return filtered[len(filtered)-1]
}
