package item

import (
	"encoding/json"

	"github.com/tarkov-database/rest-api/model"
)

// Entity defines the methods of an item entity
type Entity interface {
	GetID() objectID
	SetID(objectID)

	GetKind() Kind
	SetKind(Kind)

	GetModified() timestamp
	SetModified(timestamp)

	Validate() error
}

// A Kind represents an item type
type Kind string

// IsValid checks if a kind is valid
func (k Kind) IsValid() bool {
	if _, err := k.GetEntity(); err != nil {
		return false
	}

	return true
}

// IsEmpty checks if a kind is empty
func (k Kind) IsEmpty() bool {
	return k == Kind("")
}

// String returns the string representation of the kind
func (k Kind) String() string {
	return string(k)
}

// MarshalJSON implements the JSON marshaler
func (k *Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

// UnmarshalJSON implements the JSON unmarshaler
func (k *Kind) UnmarshalJSON(b []byte) error {
	var kind string

	if err := json.Unmarshal(b, &kind); err != nil {
		return err
	}

	*k = Kind(kind)

	return nil
}

// GetEntity returns the associated entity of the kind
func (k Kind) GetEntity() (Entity, error) {
	var e Entity

	switch k {
	case KindAmmunition:
		e = &Ammunition{}
	case KindArmor:
		e = &Armor{}
	case KindBackpack:
		e = &Backpack{}
	case KindBarter:
		e = &Barter{}
	case KindClothing:
		e = &Clothing{}
	case KindCommon:
		e = &Item{}
	case KindContainer:
		e = &Container{}
	case KindFirearm:
		e = &Firearm{}
	case KindFood:
		e = &Food{}
	case KindGrenade:
		e = &Grenade{}
	case KindHeadphone:
		e = &Headphone{}
	case KindKey:
		e = &Key{}
	case KindMagazine:
		e = &Magazine{}
	case KindMap:
		e = &Map{}
	case KindMedical:
		e = &Medical{}
	case KindMelee:
		e = &Melee{}
	case KindModification:
		e = &Modification{}
	case KindModificationBarrel:
		e = &Barrel{}
	case KindModificationBipod:
		e = &Bipod{}
	case KindModificationCharge:
		e = &Charge{}
	case KindModificationDevice:
		e = &Device{}
	case KindModificationForegrip:
		e = &Foregrip{}
	case KindModificationGasblock:
		e = &GasBlock{}
	case KindModificationGoggles:
		e = &Goggles{}
	case KindModificationHandguard:
		e = &Handguard{}
	case KindModificationLauncher:
		e = &Launcher{}
	case KindModificationMount:
		e = &Mount{}
	case KindModificationMuzzle:
		e = &Muzzle{}
	case KindModificationPistolgrip:
		e = &PistolGrip{}
	case KindModificationReceiver:
		e = &Receiver{}
	case KindModificationSight:
		e = &Sight{}
	case KindModificationSightSpecial:
		e = &SightSpecial{}
	case KindModificationStock:
		e = &Stock{}
	case KindMoney:
		e = &Money{}
	case KindTacticalrig:
		e = &TacticalRig{}
	default:
		return e, model.ErrInvalidKind
	}

	return e, nil
}

// KindList holds all kinds
var KindList = [...]Kind{
	KindAmmunition,
	KindArmor,
	KindBackpack,
	KindBarter,
	KindClothing,
	KindCommon,
	KindContainer,
	KindFirearm,
	KindFood,
	KindGrenade,
	KindHeadphone,
	KindKey,
	KindMagazine,
	KindMap,
	KindMedical,
	KindMelee,
	KindModification,
	KindModificationBarrel,
	KindModificationBipod,
	KindModificationCharge,
	KindModificationDevice,
	KindModificationForegrip,
	KindModificationGasblock,
	KindModificationGoggles,
	KindModificationHandguard,
	KindModificationLauncher,
	KindModificationMount,
	KindModificationMuzzle,
	KindModificationPistolgrip,
	KindModificationReceiver,
	KindModificationSight,
	KindModificationSightSpecial,
	KindModificationStock,
	KindMoney,
	KindTacticalrig,
}
