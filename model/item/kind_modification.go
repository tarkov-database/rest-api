package item

const (
	// KindModification represents the kind of Modification
	KindModification Kind = "modification"

	// KindModificationBarrel represents the kind of Barrel
	KindModificationBarrel Kind = "modificationBarrel"

	// KindModificationBipod represents the kind of Bipod
	KindModificationBipod Kind = "modificationBipod"

	// KindModificationCharge represents the kind of Charge
	KindModificationCharge Kind = "modificationCharge"

	// KindModificationDevice represents the kind of Device
	KindModificationDevice Kind = "modificationDevice"

	// KindModificationForegrip represents the kind of Foregrip
	KindModificationForegrip Kind = "modificationForegrip"

	// KindModificationGasblock represents the kind of GasBlock
	KindModificationGasblock Kind = "modificationGasblock"

	// KindModificationHandguard represents the kind of Handguard
	KindModificationHandguard Kind = "modificationHandguard"

	// KindModificationLauncher represents the kind of Launcher
	KindModificationLauncher Kind = "modificationLauncher"

	// KindModificationMount represents the kind of Mount
	KindModificationMount Kind = "modificationMount"

	// KindModificationMuzzle represents the kind of Muzzle
	KindModificationMuzzle Kind = "modificationMuzzle"

	// KindModificationGoggles represents the kind of Goggles
	KindModificationGoggles Kind = "modificationGoggles"

	// KindModificationPistolgrip represents the kind of PistolGrip
	KindModificationPistolgrip Kind = "modificationPistolgrip"

	// KindModificationReceiver represents the kind of Receiver
	KindModificationReceiver Kind = "modificationReceiver"

	// KindModificationSight represents the kind of Sight
	KindModificationSight Kind = "modificationSight"

	// KindModificationSightSpecial represents the kind of SightSpecial
	KindModificationSightSpecial Kind = "modificationSightSpecial"

	// KindModificationStock represents the kind of Stock
	KindModificationStock Kind = "modificationStock"
)

// Modification represents the basic data of modification item
type Modification struct {
	Item `bson:",inline"`

	ErgonomicsFloat float64      `json:"ergonomicsFP" bson:"ergonomicsFP"`
	Ergonomics      int64        `json:"ergonomics" bson:"ergonomics"` // Deprecated
	Accuracy        float64      `json:"accuracy" bson:"accuracy"`
	Recoil          float64      `json:"recoil" bson:"recoil"`
	RaidModdable    int64        `json:"raidModdable" bson:"raidModdable"`
	GridModifier    GridModifier `json:"gridModifier" bson:"gridModifier"`
	Slots           Slots        `json:"slots" bson:"slots"`
	Compatibility   List         `json:"compatibility" bson:"compatibility"`
	Conflicts       List         `json:"conflicts" bson:"conflicts"`
}

// Weapon modifications //

// Barrel describes the entity of an barrel item
type Barrel struct {
	Modification `bson:",inline"`

	Length     float64 `json:"length" bson:"length"`
	Velocity   float64 `json:"velocity" bson:"velocity"`
	Suppressor bool    `json:"suppressor" bson:"suppressor"`
}

// Bipod describes the entity of an bipod item
type Bipod struct {
	Modification `bson:",inline"`
}

// Charge describes the entity of an charging handle item
type Charge struct {
	Modification `bson:",inline"`
}

// Device describes the entity of an tactical device item
type Device struct {
	Modification `bson:",inline"`

	Type  string   `json:"type" bson:"type"`
	Modes []string `json:"modes" bson:"modes"`
}

// Foregrip describes the entity of an foregrip item
type Foregrip struct {
	Modification `bson:",inline"`
}

// GasBlock describes the entity of an gas block item
type GasBlock struct {
	Modification `bson:",inline"`
}

// Handguard describes the entity of an handguard item
type Handguard struct {
	Modification `bson:",inline"`
}

// Launcher describes the entity of an launcher item
type Launcher struct {
	Modification `bson:",inline"`

	Caliber string `json:"caliber" bson:"caliber"`
}

// Mount describes the entity of an mount item
type Mount struct {
	Modification `bson:",inline"`
}

// Muzzle describes the entity of an muzzle item
type Muzzle struct {
	Modification `bson:",inline"`

	Type     string  `json:"type" bson:"type"`
	Velocity float64 `json:"velocity" bson:"velocity"`
}

// PistolGrip describes the entity of an pistol grip item
type PistolGrip struct {
	Modification `bson:",inline"`
}

// Receiver describes the entity of an receiver item
type Receiver struct {
	Modification `bson:",inline"`

	Velocity float64 `json:"velocity" bson:"velocity"`
}

// Sight describes the entity of an sight item
type Sight struct {
	Modification `bson:",inline"`

	Type          string   `json:"type" bson:"type"`
	Magnification []string `json:"magnification" bson:"magnification"`
	VariableZoom  bool     `json:"variableZoom" bson:"variableZoom"`
	ZeroDistances []int64  `json:"zeroDistances" bson:"zeroDistances"`
}

// SightSpecial describes the entity of an special sights item
type SightSpecial struct {
	Sight        `bson:",inline"`
	OpticSpecial `bson:",inline"`
}

// Stock describes the entity of an stock item
type Stock struct {
	Modification `bson:",inline"`

	FoldRectractable bool `json:"foldRectractable" bson:"foldRectractable"`
}

// Gear modifications //

// Goggles describes the entity of an goggles item
type Goggles struct {
	Modification `bson:",inline"`

	Type string `json:"type" bson:"type"`

	OpticSpecial `bson:",inline"`
}

// Properties //

// GridModifier represents the properties of an grid modifier
type GridModifier struct {
	Height int64 `json:"height" bson:"height"`
	Width  int64 `json:"width" bson:"width"`
}

// OpticSpecial represents the properties of GogglesSpecial and SightSpecial
type OpticSpecial struct {
	Modes []string `json:"modes" bson:"modes"`
	Color RGBA     `json:"color" bson:"color"`
	Noise string   `json:"noise" bson:"noise"`
}
