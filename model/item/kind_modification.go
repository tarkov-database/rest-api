package item

const (
	KindModification               Kind = "modification"
	KindModificationBarrel         Kind = "modificationBarrel"
	KindModificationBipod          Kind = "modificationBipod"
	KindModificationCharge         Kind = "modificationCharge"
	KindModificationDevice         Kind = "modificationDevice"
	KindModificationForegrip       Kind = "modificationForegrip"
	KindModificationGasblock       Kind = "modificationGasblock"
	KindModificationHandguard      Kind = "modificationHandguard"
	KindModificationLauncher       Kind = "modificationLauncher"
	KindModificationMount          Kind = "modificationMount"
	KindModificationMuzzle         Kind = "modificationMuzzle"
	KindModificationGoggles        Kind = "modificationGoggles"
	KindModificationGogglesSpecial Kind = "modificationGogglesSpecial"
	KindModificationPistolgrip     Kind = "modificationPistolgrip"
	KindModificationReceiver       Kind = "modificationReceiver"
	KindModificationSight          Kind = "modificationSight"
	KindModificationSightSpecial   Kind = "modificationSightSpecial"
	KindModificationStock          Kind = "modificationStock"
)

type Modification struct {
	Item `json:",inline" bson:",inline"`

	Ergonomics    int64        `json:"ergonomics" bson:"ergonomics"`
	RaidModdable  int64        `json:"raidModdable" bson:"raidModdable"`
	GridModifier  GridModifier `json:"gridModifier" bson:"gridModifier"`
	Slots         Slots        `json:"slots" bson:"slots"`
	Compatibility ItemList     `json:"compatibility" bson:"compatibility"`
	Conflicts     ItemList     `json:"conflicts" bson:"conflicts"`
}

// Weapon modifications //

type Barrel struct {
	Modification `json:",inline" bson:",inline"`

	Length     float64 `json:"length" bson:"length"`
	Accuracy   float64 `json:"accuracy" bson:"accuracy"`
	Velocity   float64 `json:"velocity" bson:"velocity"`
	Recoil     float64 `json:"recoil" bson:"recoil"`
	Suppressor bool    `json:"suppressor" bson:"suppressor"`
}

type Bipod struct {
	Modification `json:",inline" bson:",inline"`

	Recoil float64 `json:"recoil" bson:"recoil"`
}

type Charge struct {
	Modification `json:",inline" bson:",inline"`
}

type Device struct {
	Modification `json:",inline" bson:",inline"`

	Type  string   `json:"type" bson:"type"`
	Modes []string `json:"modes" bson:"modes"`
}

type Foregrip struct {
	Modification `json:",inline" bson:",inline"`

	Recoil float64 `json:"recoil" bson:"recoil"`
}

type GasBlock struct {
	Modification `json:",inline" bson:",inline"`
}

type Handguard struct {
	Modification `json:",inline" bson:",inline"`

	Recoil float64 `json:"recoil" bson:"recoil"`
}

type Launcher struct {
	Modification `json:",inline" bson:",inline"`

	Recoil  float64 `json:"recoil" bson:"recoil"`
	Caliber string  `json:"caliber" bson:"caliber"`
}

type Mount struct {
	Modification `json:",inline" bson:",inline"`
}

type Muzzle struct {
	Modification `json:",inline" bson:",inline"`

	Type     string  `json:"type" bson:"type"`
	Accuracy float64 `json:"accuracy" bson:"accuracy"`
	Velocity float64 `json:"velocity" bson:"velocity"`
	Recoil   float64 `json:"recoil" bson:"recoil"`
}

type PistolGrip struct {
	Modification `json:",inline" bson:",inline"`

	Recoil float64 `json:"recoil" bson:"recoil"`
}

type Receiver struct {
	Modification `json:",inline" bson:",inline"`
}

type Sight struct {
	Modification `json:",inline" bson:",inline"`

	Type          string   `json:"type" bson:"type"`
	Magnification []string `json:"magnification" bson:"magnification"`
	VariableZoom  bool     `json:"variableZoom" bson:"variableZoom"`
	ZeroDistances []int64  `json:"zeroDistances" bson:"zeroDistances"`
}

type SightSpecial struct {
	Sight        `json:",inline" bson:",inline"`
	OpticSpecial `json:",inline" bson:",inline"`
}

type Stock struct {
	Modification `json:",inline" bson:",inline"`

	Recoil           float64 `json:"recoil" bson:"recoil"`
	FoldRectractable bool    `json:"foldRectractable" bson:"foldRectractable"`
}

// Gear modifications //

type Goggles struct {
	Modification `json:",inline" bson:",inline"`

	Type string `json:"type" bson:"type"`
}

type GogglesSpecial struct {
	Goggles      `json:",inline" bson:",inline"`
	OpticSpecial `json:",inline" bson:",inline"`
}

// Properties //

type GridModifier struct {
	Height int64 `json:"height" bson:"height"`
	Width  int64 `json:"width" bson:"width"`
}

type OpticSpecial struct {
	Modes []string `json:"modes" bson:"modes"`
	Color RGBA     `json:"color" bson:"color"`
	Noise string   `json:"noise" bson:"noise"`
}
