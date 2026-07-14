// Code generated from JSON Schema using quicktype. DO NOT EDIT.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    manifestDefinition, err := UnmarshalManifestDefinition(bytes)
//    bytes, err = manifestDefinition.Marshal()
//
//    manifestStore, err := UnmarshalManifestStore(bytes)
//    bytes, err = manifestStore.Marshal()
//
//    settings, err := UnmarshalSettings(bytes)
//    bytes, err = settings.Marshal()

package schema

import "bytes"
import "errors"

import "encoding/json"

func UnmarshalManifestDefinition(data []byte) (ManifestDefinition, error) {
	var r ManifestDefinition
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ManifestDefinition) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalManifestStore(data []byte) (ManifestStore, error) {
	var r ManifestStore
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ManifestStore) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalSettings(data []byte) (Settings, error) {
	var r Settings
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Settings) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Use a ManifestDefinition to define a manifest and to build a `ManifestStore`.
// A manifest is a collection of ingredients and assertions
// used to define a claim that can be signed and embedded into a file.
type ManifestDefinition struct {
	// A list of assertions
	Assertions []AssertionElement `json:"assertions,omitempty"`
	// Software that generated the claim, as a list of [`ClaimGeneratorInfo`].
	//
	// In JSON, when this key is **omitted** (or the array is empty), the value is
	// resolved at claim-building time: settings.builder.claim_generator_info if set, otherwise
	// on the active [`Context`] is used if set, otherwise [`ClaimGeneratorInfo::default`].
	// A non-empty list in the definition is used as given.
	ClaimGeneratorInfo []ClaimGeneratorInfoElement `json:"claim_generator_info,omitempty"`
	// The version of the claim.  Defaults to 2.
	ClaimVersion *int64 `json:"claim_version"`
	// The format of the source file as a MIME type.
	Format *string `json:"format,omitempty"`
	// Hash algorithm used for asset hashing (DataHash, BmffHash) and assertion hashing
	// in the claim.  Defaults to `"sha256"` when not set.
	//
	// Valid values: `"sha256"`, `"sha384"`, `"sha512"`.
	//
	// This sets the claim-level `alg` field and is the default used by
	// [`Builder::update_hash_from_stream`].  It can be overridden on individual
	// hard binding assertions (e.g. a pre-constructed `DataHash`) by setting
	// the assertion's own `alg` field before adding it to the builder.
	//
	// Named `hash_alg` (rather than `alg`) to avoid collision with the signer's
	// signature algorithm, which uses the same key in some combined JSON configurations.
	HashAlg *string `json:"hash_alg"`
	// A List of ingredients
	Ingredients []IngredientElement `json:"ingredients,omitempty"`
	// Instance ID from `xmpMM:InstanceID` in XMP metadata.
	InstanceID *string `json:"instance_id,omitempty"`
	// Allows you to pre-define the manifest label, which must be unique.
	// Not intended for general use.  If not set, it will be assigned automatically.
	Label *string `json:"label"`
	// Optional manifest metadata. This will be deprecated in the future; not recommended to use.
	Metadata []MetadatumElement `json:"metadata"`
	// JUMBF URIs of assertions to redact from ingredient manifests.
	//
	// Each URI has the form
	// `self#jumbf=/c2pa/<manifest_label>/c2pa.assertions/<assertion_label>`.
	// Use a [`Reader`](crate::Reader) to discover the manifest label.
	// See the [redaction
	// guide](https://github.com/contentauth/c2pa-rs/blob/main/docs/redaction.md)
	// for details.
	Redactions []string `json:"redactions"`
	// An optional ResourceRef to a thumbnail image that represents the asset that was signed.
	// Must be available when the manifest is signed.
	Thumbnail *ThumbnailClass `json:"thumbnail"`
	// A human-readable title, generally source filename.
	Title *string `json:"title"`
	// Optional prefix added to the generated Manifest Label
	// This is typically a reverse domain name.
	Vendor *string `json:"vendor"`
}

// Defines an assertion that consists of a label that can be either
// a C2PA-defined assertion label or a custom label in reverse domain format.
type AssertionElement struct {
	// True if this assertion is attributed to the signer (defaults to false)
	Created *bool `json:"created,omitempty"`
	// The assertion data
	Data interface{} `json:"data"`
	// The kind of assertion data, either Cbor or Json (defaults to Cbor)
	Kind *KindEnum `json:"kind"`
	// An assertion label in reverse domain format
	Label string `json:"label"`
}

// Description of the claim generator, or the software used in generating the claim.
//
// This structure is also used for actions softwareAgent
type ClaimGeneratorInfoElement struct {
	// hashed URI to the icon (either embedded or remote)
	Icon *IconClass `json:"icon"`
	// A human readable string naming the claim_generator
	Name string `json:"name"`
	// A human readable string of the OS the claim generator is running on.
	// CrJSON schema uses `operating_system`; C2PA CBOR may use
	// `schema.org.SoftwareApplication.operatingSystem`.
	OperatingSystem *string `json:"operating_system"`
	// A human readable string of the product's version
	Version *string `json:"version"`
}

// A reference to a resource to be used in JSON serialization.
//
// The underlying data can be read as a stream via
// [`Reader::resource_to_stream`][crate::Reader::resource_to_stream].
//
// A `HashedUri` provides a reference to content available within the same
// manifest store.
//
// This is described in [URI References in the C2PA Technical
//
// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_uri_references).
type IconClass struct {
	// The algorithm used to hash the resource (if applicable).
	//
	// A string identifying the cryptographic hash algorithm used to compute
	// the hash
	Alg *string `json:"alg"`
	// More detailed data types as defined in the C2PA spec.
	DataTypes []DataTypeElement `json:"data_types"`
	// The mime type of the referenced resource.
	Format *string `json:"format,omitempty"`
	// The hash of the resource (if applicable).
	//
	// Byte string containing the hash value
	Hash *HashUnion `json:"hash"`
	// A URI that identifies the resource as referenced from the manifest.
	//
	// This may be a JUMBF URI, a file path, a URL or any other string.
	// Relative JUMBF URIs will be resolved with the manifest label.
	// Relative file paths will be resolved with the base path if provided.
	Identifier *string `json:"identifier,omitempty"`
	// JUMBF URI reference
	URL *string `json:"url,omitempty"`
}

type DataTypeElement struct {
	Type    string  `json:"type"`
	Version *string `json:"version"`
}

// An `Ingredient` is any external asset that has been used in the creation of an asset.
type IngredientElement struct {
	// The active manifest label (if one exists).
	//
	// If this ingredient has a [`ManifestStore`],
	// this will hold the label of the active [`Manifest`].
	//
	// [`Manifest`]: crate::Manifest
	// [`ManifestStore`]: crate::ManifestStore
	ActiveManifest *string `json:"active_manifest"`
	// A reference to the actual data of the ingredient.
	Data *ThumbnailClass `json:"data"`
	// Additional information about the data's type to the ingredient V2 structure.
	DataTypes []DataTypeElement `json:"data_types"`
	// Additional description of the ingredient.
	Description *string `json:"description"`
	// Document ID from `xmpMM:DocumentID` in XMP metadata.
	DocumentID *string `json:"document_id"`
	// The format of the source file as a MIME type.
	Format *string `json:"format"`
	// An optional hash of the asset to prevent duplicates.
	Hash *string `json:"hash"`
	// URI to an informational page about the ingredient or its data.
	InformationalURI *string `json:"informational_URI"`
	// Instance ID from `xmpMM:InstanceID` in XMP metadata.
	InstanceID *string `json:"instance_id"`
	// The ingredient's label as assigned in the manifest.
	Label *string `json:"label"`
	// A [`ManifestStore`] from the source asset extracted as a binary C2PA blob.
	//
	// [`ManifestStore`]: crate::ManifestStore
	ManifestData *ThumbnailClass `json:"manifest_data"`
	// Any additional [`Metadata`] as defined in the C2PA spec.
	//
	// [`Metadata`]: crate::Metadata
	Metadata      *MetadatumElement `json:"metadata"`
	OcspResponses []ThumbnailClass  `json:"ocsp_responses"`
	// URI from `dcterms:provenance` in XMP metadata.
	Provenance *string `json:"provenance"`
	// Set to `ParentOf` if this is the parent ingredient.
	//
	// There can only be one parent ingredient in the ingredients.
	Relationship *Relationship `json:"relationship,omitempty"`
	// A thumbnail image capturing the visual state at the time of import.
	//
	// A tuple of thumbnail MIME format (for example `image/jpeg`) and binary bits of the image.
	Thumbnail *ThumbnailClass `json:"thumbnail"`
	// A human-readable title, generally source filename.
	Title *string `json:"title"`
	// Validation results (Ingredient.V3)
	ValidationResults *ValidationResultsClass `json:"validation_results"`
	// Validation status (Ingredient v1 & v2)
	ValidationStatus []ValidationStatusElement `json:"validation_status"`
}

// A reference to a resource to be used in JSON serialization.
//
// The underlying data can be read as a stream via
// [`Reader::resource_to_stream`][crate::Reader::resource_to_stream].
type ThumbnailClass struct {
	// The algorithm used to hash the resource (if applicable).
	Alg *string `json:"alg"`
	// More detailed data types as defined in the C2PA spec.
	DataTypes []DataTypeElement `json:"data_types"`
	// The mime type of the referenced resource.
	Format string `json:"format"`
	// The hash of the resource (if applicable).
	Hash *string `json:"hash"`
	// A URI that identifies the resource as referenced from the manifest.
	//
	// This may be a JUMBF URI, a file path, a URL or any other string.
	// Relative JUMBF URIs will be resolved with the manifest label.
	// Relative file paths will be resolved with the base path if provided.
	Identifier string `json:"identifier"`
}

// A region of interest within an asset describing the change.
//
// This struct can be used from [`Action::changes`][crate::assertions::Action::changes],
//
// [`AssertionMetadata::region_of_interest`][crate::assertions::AssertionMetadata::region_of_interest],
// or
// [`SoftBindingScope::region`][crate::assertions::soft_binding::SoftBindingScope::region].
type RegionOfInterestClass struct {
	// A free-text string.
	Description *string `json:"description"`
	// A free-text string representing a machine-readable, unique to this assertion, identifier
	// for the region.
	Identifier *string `json:"identifier"`
	// Additional information about the asset.
	Metadata *MetadatumElement `json:"metadata"`
	// A free-text string representing a human-readable name for the region which might be used
	// in a user interface.
	Name *string `json:"name"`
	// A range describing the region of interest for the specific asset.
	Region []RegionElement `json:"region"`
	// A value from our controlled vocabulary or an entity-specific value (e.g.,
	// com.litware.coolArea) that represents
	// the role of a region among other regions.
	Role *RoleEnum `json:"role"`
	// A value from a controlled vocabulary such as
	// <https://cv.iptc.org/newscodes/imageregiontype/> or an entity-specific
	// value (e.g., com.litware.newType) that represents the type of thing(s) depicted by a
	// region.
	//
	// Note this field serializes/deserializes into the name `type`.
	Type *string `json:"type"`
}

// The AssertionMetadata structure can be used as part of other assertions or on its own to
// reference others
type MetadatumElement struct {
	DataSource       *DataSourceClass               `json:"dataSource"`
	DateTime         *string                        `json:"dateTime"`
	Localizations    []map[string]map[string]string `json:"localizations"`
	Reference        *ReferenceElement              `json:"reference"`
	RegionOfInterest *RegionOfInterestClass         `json:"regionOfInterest"`
	ReviewRatings    []ReviewRatingElement          `json:"reviewRatings"`
}

// A spatial, temporal, frame, or textual range describing the region of interest.
type RegionElement struct {
	// A frame range.
	Frame *FrameClass `json:"frame"`
	// A item identifier.
	Item *ItemClass `json:"item"`
	// A spatial range.
	Shape *ShapeClass `json:"shape"`
	// A textual range.
	Text *TextClass `json:"text"`
	// A temporal range.
	Time *TimeClass `json:"time"`
	// The type of range of interest.
	Type RegionType `json:"type"`
}

// A frame range representing starting and ending frames or pages.
//
// If both `start` and `end` are missing, the frame will span the entire asset.
type FrameClass struct {
	// The end of the frame inclusive or the end of the asset if not present.
	End *int64 `json:"end"`
	// The start of the frame or the end of the asset if not present.
	//
	// The first frame/page starts at 0.
	Start *int64 `json:"start"`
}

// Description of the boundaries of an identified range.
type ItemClass struct {
	// The container-specific term used to identify items, such as "track_id" for MP4 or
	// "item_ID" for HEIF.
	Identifier string `json:"identifier"`
	// The value of the identifier, e.g. a value of "2" for an identifier of "track_id" would
	// imply track 2 of the asset.
	Value string `json:"value"`
}

// A spatial range representing rectangle, circle, or a polygon.
type ShapeClass struct {
	// The height of a rectnagle.
	//
	// This field can be ignored for circles and polygons.
	Height *float64 `json:"height"`
	// If the range is inside the shape.
	//
	// The default value is true.
	Inside *bool `json:"inside"`
	// THe origin of the coordinate in the shape.
	Origin PurpleOrigin `json:"origin"`
	// The type of shape.
	Type ShapeType `json:"type"`
	// The type of unit for the shape range.
	Unit Unit `json:"unit"`
	// The vertices of the polygon.
	//
	// This field can be ignored for rectangles and circles.
	Vertices []PurpleOrigin `json:"vertices"`
	// The width for rectangles or diameter for circles.
	//
	// This field can be ignored for polygons.
	Width *float64 `json:"width"`
}

// THe origin of the coordinate in the shape.
//
// An x, y coordinate used for specifying vertices in polygons.
type PurpleOrigin struct {
	// The coordinate along the x-axis.
	X float64 `json:"x"`
	// The coordinate along the y-axis.
	Y float64 `json:"y"`
}

// A textual range representing multiple (possibly discontinuous) ranges of text.
type TextClass struct {
	// The ranges of text to select.
	Selectors []SelectorElement `json:"selectors"`
}

// One or two [`TextSelector`] identifiying the range to select.
type SelectorElement struct {
	// The end of the text range.
	End *EndClass `json:"end"`
	// The start (or entire) text range.
	Selector EndClass `json:"selector"`
}

// Selects a range of text via a fragment identifier.
//
// This is modeled after the W3C Web Annotation selector model.
//
// The start (or entire) text range.
type EndClass struct {
	// The end character offset or the end of the fragment if not present.
	End *int64 `json:"end"`
	// Fragment identifier as per RFC3023 (XML) or ISO 32000-2 (PDF), Annex O.
	Fragment string `json:"fragment"`
	// The start character offset or the start of the fragment if not present.
	Start *int64 `json:"start"`
}

// A temporal range representing a starting time to an ending time.
type TimeClass struct {
	// The end time or the end of the asset if not present.
	End *string `json:"end"`
	// The start time or the start of the asset if not present.
	Start *string `json:"start"`
	// The type of time.
	Type *TimeType `json:"type,omitempty"`
}

// A description of the source for assertion data
type DataSourceClass struct {
	// A list of [`Actor`]s associated with this source.
	Actors []ActorElement `json:"actors"`
	// A human-readable string giving details about the source of the assertion data.
	Details *string `json:"details"`
	// A value from among the enumerated list indicating the source of the assertion.
	Type string `json:"type"`
}

// Identifies a person responsible for an action.
type ActorElement struct {
	// List of references to W3C Verifiable Credentials.
	Credentials []ReferenceElement `json:"credentials"`
	// An identifier for a human actor, used when the "type" is `humanEntry.identified`.
	Identifier *string `json:"identifier"`
}

// A `HashedUri` provides a reference to content available within the same
// manifest store.
//
// This is described in [URI References in the C2PA Technical
//
// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_uri_references).
type ReferenceElement struct {
	// A string identifying the cryptographic hash algorithm used to compute
	// the hash
	Alg *string `json:"alg"`
	// Byte string containing the hash value
	Hash []int64 `json:"hash"`
	// JUMBF URI reference
	URL string `json:"url"`
}

// A rating on an Assertion.
//
// See [C2PA Specification - Review
// Ratings](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_review_ratings).
type ReviewRatingElement struct {
	Code        *string `json:"code"`
	Explanation string  `json:"explanation"`
	Value       int64   `json:"value"`
}

// A map of validation results for a manifest store.
//
// The map contains the validation results for the active manifest and any ingredient
// deltas.
// It is normal for there to be many
type ValidationResultsClass struct {
	// Validation status codes for the ingredient's active manifest. Present if ingredient is a
	// C2PA
	// asset. Not present if the ingredient is not a C2PA asset.
	ActiveManifest *ActiveManifestClass `json:"activeManifest"`
	// List of any changes/deltas between the current and previous validation results for each
	// ingredient's
	// manifest. Present if the the ingredient is a C2PA asset.
	IngredientDeltas []IngredientDeltaElement `json:"ingredientDeltas"`
	// Time when the validation was performed (RFC 3339 date-time). Used only for document-level
	// validationInfo; not serialized in validationResults (e.g. ingredient assertions).
	ValidationTime *string `json:"validationTime"`
}

// Contains a set of success, informational, and failure validation status codes.
//
// Validation results for the ingredient's active manifest
type ActiveManifestClass struct {
	// An array of validation failure codes. May be empty.
	Failure []ValidationStatusElement `json:"failure"`
	// An array of validation informational codes. May be empty.
	Informational []ValidationStatusElement `json:"informational"`
	// An array of validation success codes. May be empty.
	Success []ValidationStatusElement `json:"success"`
}

// A `ValidationStatus` struct describes the validation status of a
// specific part of a manifest.
//
// See [Existing Manifests - C2PA Technical
// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_existing_manifests).
type ValidationStatusElement struct {
	Code        string  `json:"code"`
	Explanation *string `json:"explanation"`
	Success     *bool   `json:"success"`
	URL         *string `json:"url"`
}

// Represents any changes or deltas between the current and previous validation results for
// an ingredient's manifest.
type IngredientDeltaElement struct {
	// JUMBF URI reference to the ingredient assertion
	IngredientAssertionURI string `json:"ingredientAssertionURI"`
	// Validation results for the ingredient's active manifest
	ValidationDeltas ActiveManifestClass `json:"validationDeltas"`
}

// Use a Reader to read and validate a manifest store.
type ManifestStore struct {
	// A label for the active (most recent) manifest in the store
	ActiveManifest *string `json:"active_manifest"`
	// A HashMap of Manifests
	Manifests map[string]ManifestValue `json:"manifests,omitempty"`
	// ValidationStatus generated when loading the ManifestStore from an asset
	ValidationResults *IngredientValidationResults `json:"validation_results"`
	// The validation state of the manifest store
	ValidationState *ValidationStateEnum `json:"validation_state"`
	// ValidationStatus generated when loading the ManifestStore from an asset
	ValidationStatus []FailureElement `json:"validation_status"`
}

// A Manifest represents all the information in a c2pa manifest
type ManifestValue struct {
	// A list of assertions
	Assertions []AssertionClass `json:"assertions,omitempty"`
	// A User Agent formatted string identifying the software/hardware/system produced this
	// claim
	// Spaces are not allowed in names, versions can be specified with product/1.0 syntax.
	ClaimGenerator *string `json:"claim_generator"`
	// A list of claim generator info data identifying the software/hardware/system produced
	// this claim.
	ClaimGeneratorInfo []ClaimGeneratorInfoClass `json:"claim_generator_info"`
	// The version of the claim, parsed from the claim label.
	//
	// For example:
	// - `c2pa.claim.v2` -> 2
	// - `c2pa.claim` -> 1
	ClaimVersion *int64 `json:"claim_version"`
	// A List of verified credentials
	Credentials []interface{} `json:"credentials"`
	// The format of the source file as a MIME type.
	Format *string `json:"format"`
	// A List of ingredients
	Ingredients []IngredientClass `json:"ingredients,omitempty"`
	// Instance ID from `xmpMM:InstanceID` in XMP metadata.
	InstanceID *string `json:"instance_id,omitempty"`
	Label      *string `json:"label"`
	// A list of user metadata for this claim.
	Metadata []MetadataElement `json:"metadata"`
	// JUMBF URIs of assertions that were redacted by this manifest.
	//
	// Each entry has the form
	// `self#jumbf=/c2pa/<manifest_label>/c2pa.assertions/<assertion_label>`
	// and corresponds to an assertion that was intentionally removed from an
	// ingredient manifest in the claim chain.
	Redactions []string `json:"redactions"`
	// Signature data (only used for reporting)
	SignatureInfo *SignatureInfoClass `json:"signature_info"`
	Thumbnail     *DataClass          `json:"thumbnail"`
	// A human-readable title, generally source filename.
	Title *string `json:"title"`
	// Optional prefix added to the generated Manifest label.
	// This is typically an internet domain name for the vendor (i.e. `adobe`).
	Vendor *string `json:"vendor"`
}

// A labeled container for an Assertion value in a Manifest
type AssertionClass struct {
	// True if this assertion is attributed to the signer
	// This maps to a created vs a gathered assertion. (defaults to false)
	Created *bool `json:"created,omitempty"`
	// The data of the assertion as Value
	Data interface{} `json:"data"`
	// There can be more than one assertion for any label
	Instance *int64 `json:"instance"`
	// The [ManifestAssertionKind] for this assertion (as stored in c2pa content)
	Kind *KindEnum `json:"kind"`
	// An assertion label in reverse domain format
	Label string `json:"label"`
}

// Description of the claim generator, or the software used in generating the claim.
//
// This structure is also used for actions softwareAgent
type ClaimGeneratorInfoClass struct {
	// hashed URI to the icon (either embedded or remote)
	Icon *ClaimGeneratorInfoIcon `json:"icon"`
	// A human readable string naming the claim_generator
	Name string `json:"name"`
	// A human readable string of the OS the claim generator is running on.
	// CrJSON schema uses `operating_system`; C2PA CBOR may use
	// `schema.org.SoftwareApplication.operatingSystem`.
	OperatingSystem *string `json:"operating_system"`
	// A human readable string of the product's version
	Version *string `json:"version"`
}

// A reference to a resource to be used in JSON serialization.
//
// The underlying data can be read as a stream via
// [`Reader::resource_to_stream`][crate::Reader::resource_to_stream].
//
// A `HashedUri` provides a reference to content available within the same
// manifest store.
//
// This is described in [URI References in the C2PA Technical
//
// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_uri_references).
type ClaimGeneratorInfoIcon struct {
	// The algorithm used to hash the resource (if applicable).
	//
	// A string identifying the cryptographic hash algorithm used to compute
	// the hash
	Alg *string `json:"alg"`
	// More detailed data types as defined in the C2PA spec.
	DataTypes []DataTypeClass `json:"data_types"`
	// The mime type of the referenced resource.
	Format *string `json:"format,omitempty"`
	// The hash of the resource (if applicable).
	//
	// Byte string containing the hash value
	Hash *HashUnion `json:"hash"`
	// A URI that identifies the resource as referenced from the manifest.
	//
	// This may be a JUMBF URI, a file path, a URL or any other string.
	// Relative JUMBF URIs will be resolved with the manifest label.
	// Relative file paths will be resolved with the base path if provided.
	Identifier *string `json:"identifier,omitempty"`
	// JUMBF URI reference
	URL *string `json:"url,omitempty"`
}

type DataTypeClass struct {
	Type    string  `json:"type"`
	Version *string `json:"version"`
}

// An `Ingredient` is any external asset that has been used in the creation of an asset.
type IngredientClass struct {
	// The active manifest label (if one exists).
	//
	// If this ingredient has a [`ManifestStore`],
	// this will hold the label of the active [`Manifest`].
	//
	// [`Manifest`]: crate::Manifest
	// [`ManifestStore`]: crate::ManifestStore
	ActiveManifest *string `json:"active_manifest"`
	// A reference to the actual data of the ingredient.
	Data *DataClass `json:"data"`
	// Additional information about the data's type to the ingredient V2 structure.
	DataTypes []DataTypeClass `json:"data_types"`
	// Additional description of the ingredient.
	Description *string `json:"description"`
	// Document ID from `xmpMM:DocumentID` in XMP metadata.
	DocumentID *string `json:"document_id"`
	// The format of the source file as a MIME type.
	Format *string `json:"format"`
	// An optional hash of the asset to prevent duplicates.
	Hash *string `json:"hash"`
	// URI to an informational page about the ingredient or its data.
	InformationalURI *string `json:"informational_URI"`
	// Instance ID from `xmpMM:InstanceID` in XMP metadata.
	InstanceID *string `json:"instance_id"`
	// The ingredient's label as assigned in the manifest.
	Label *string `json:"label"`
	// A [`ManifestStore`] from the source asset extracted as a binary C2PA blob.
	//
	// [`ManifestStore`]: crate::ManifestStore
	ManifestData *DataClass `json:"manifest_data"`
	// Any additional [`Metadata`] as defined in the C2PA spec.
	//
	// [`Metadata`]: crate::Metadata
	Metadata      *MetadataElement `json:"metadata"`
	OcspResponses []DataClass      `json:"ocsp_responses"`
	// URI from `dcterms:provenance` in XMP metadata.
	Provenance *string `json:"provenance"`
	// Set to `ParentOf` if this is the parent ingredient.
	//
	// There can only be one parent ingredient in the ingredients.
	Relationship *Relationship `json:"relationship,omitempty"`
	// A thumbnail image capturing the visual state at the time of import.
	//
	// A tuple of thumbnail MIME format (for example `image/jpeg`) and binary bits of the image.
	Thumbnail *DataClass `json:"thumbnail"`
	// A human-readable title, generally source filename.
	Title *string `json:"title"`
	// Validation results (Ingredient.V3)
	ValidationResults *IngredientValidationResults `json:"validation_results"`
	// Validation status (Ingredient v1 & v2)
	ValidationStatus []FailureElement `json:"validation_status"`
}

// A reference to a resource to be used in JSON serialization.
//
// The underlying data can be read as a stream via
// [`Reader::resource_to_stream`][crate::Reader::resource_to_stream].
type DataClass struct {
	// The algorithm used to hash the resource (if applicable).
	Alg *string `json:"alg"`
	// More detailed data types as defined in the C2PA spec.
	DataTypes []DataTypeClass `json:"data_types"`
	// The mime type of the referenced resource.
	Format string `json:"format"`
	// The hash of the resource (if applicable).
	Hash *string `json:"hash"`
	// A URI that identifies the resource as referenced from the manifest.
	//
	// This may be a JUMBF URI, a file path, a URL or any other string.
	// Relative JUMBF URIs will be resolved with the manifest label.
	// Relative file paths will be resolved with the base path if provided.
	Identifier string `json:"identifier"`
}

// A region of interest within an asset describing the change.
//
// This struct can be used from [`Action::changes`][crate::assertions::Action::changes],
//
// [`AssertionMetadata::region_of_interest`][crate::assertions::AssertionMetadata::region_of_interest],
// or
// [`SoftBindingScope::region`][crate::assertions::soft_binding::SoftBindingScope::region].
type MetadatumRegionOfInterest struct {
	// A free-text string.
	Description *string `json:"description"`
	// A free-text string representing a machine-readable, unique to this assertion, identifier
	// for the region.
	Identifier *string `json:"identifier"`
	// Additional information about the asset.
	Metadata *MetadataElement `json:"metadata"`
	// A free-text string representing a human-readable name for the region which might be used
	// in a user interface.
	Name *string `json:"name"`
	// A range describing the region of interest for the specific asset.
	Region []RegionClass `json:"region"`
	// A value from our controlled vocabulary or an entity-specific value (e.g.,
	// com.litware.coolArea) that represents
	// the role of a region among other regions.
	Role *RoleEnum `json:"role"`
	// A value from a controlled vocabulary such as
	// <https://cv.iptc.org/newscodes/imageregiontype/> or an entity-specific
	// value (e.g., com.litware.newType) that represents the type of thing(s) depicted by a
	// region.
	//
	// Note this field serializes/deserializes into the name `type`.
	Type *string `json:"type"`
}

// The AssertionMetadata structure can be used as part of other assertions or on its own to
// reference others
type MetadataElement struct {
	DataSource       *MetadatumDataSource           `json:"dataSource"`
	DateTime         *string                        `json:"dateTime"`
	Localizations    []map[string]map[string]string `json:"localizations"`
	Reference        *CredentialElement             `json:"reference"`
	RegionOfInterest *MetadatumRegionOfInterest     `json:"regionOfInterest"`
	ReviewRatings    []ReviewRatingClass            `json:"reviewRatings"`
}

// A spatial, temporal, frame, or textual range describing the region of interest.
type RegionClass struct {
	// A frame range.
	Frame *RegionFrame `json:"frame"`
	// A item identifier.
	Item *RegionItem `json:"item"`
	// A spatial range.
	Shape *RegionShape `json:"shape"`
	// A textual range.
	Text *RegionText `json:"text"`
	// A temporal range.
	Time *RegionTime `json:"time"`
	// The type of range of interest.
	Type RegionType `json:"type"`
}

// A frame range representing starting and ending frames or pages.
//
// If both `start` and `end` are missing, the frame will span the entire asset.
type RegionFrame struct {
	// The end of the frame inclusive or the end of the asset if not present.
	End *int64 `json:"end"`
	// The start of the frame or the end of the asset if not present.
	//
	// The first frame/page starts at 0.
	Start *int64 `json:"start"`
}

// Description of the boundaries of an identified range.
type RegionItem struct {
	// The container-specific term used to identify items, such as "track_id" for MP4 or
	// "item_ID" for HEIF.
	Identifier string `json:"identifier"`
	// The value of the identifier, e.g. a value of "2" for an identifier of "track_id" would
	// imply track 2 of the asset.
	Value string `json:"value"`
}

// A spatial range representing rectangle, circle, or a polygon.
type RegionShape struct {
	// The height of a rectnagle.
	//
	// This field can be ignored for circles and polygons.
	Height *float64 `json:"height"`
	// If the range is inside the shape.
	//
	// The default value is true.
	Inside *bool `json:"inside"`
	// THe origin of the coordinate in the shape.
	Origin FluffyOrigin `json:"origin"`
	// The type of shape.
	Type ShapeType `json:"type"`
	// The type of unit for the shape range.
	Unit Unit `json:"unit"`
	// The vertices of the polygon.
	//
	// This field can be ignored for rectangles and circles.
	Vertices []FluffyOrigin `json:"vertices"`
	// The width for rectangles or diameter for circles.
	//
	// This field can be ignored for polygons.
	Width *float64 `json:"width"`
}

// THe origin of the coordinate in the shape.
//
// An x, y coordinate used for specifying vertices in polygons.
type FluffyOrigin struct {
	// The coordinate along the x-axis.
	X float64 `json:"x"`
	// The coordinate along the y-axis.
	Y float64 `json:"y"`
}

// A textual range representing multiple (possibly discontinuous) ranges of text.
type RegionText struct {
	// The ranges of text to select.
	Selectors []TextSelector `json:"selectors"`
}

// One or two [`TextSelector`] identifiying the range to select.
type TextSelector struct {
	// The end of the text range.
	End *SelectorSelector `json:"end"`
	// The start (or entire) text range.
	Selector SelectorSelector `json:"selector"`
}

// Selects a range of text via a fragment identifier.
//
// This is modeled after the W3C Web Annotation selector model.
//
// The start (or entire) text range.
type SelectorSelector struct {
	// The end character offset or the end of the fragment if not present.
	End *int64 `json:"end"`
	// Fragment identifier as per RFC3023 (XML) or ISO 32000-2 (PDF), Annex O.
	Fragment string `json:"fragment"`
	// The start character offset or the start of the fragment if not present.
	Start *int64 `json:"start"`
}

// A temporal range representing a starting time to an ending time.
type RegionTime struct {
	// The end time or the end of the asset if not present.
	End *string `json:"end"`
	// The start time or the start of the asset if not present.
	Start *string `json:"start"`
	// The type of time.
	Type *TimeType `json:"type,omitempty"`
}

// A description of the source for assertion data
type MetadatumDataSource struct {
	// A list of [`Actor`]s associated with this source.
	Actors []ActorClass `json:"actors"`
	// A human-readable string giving details about the source of the assertion data.
	Details *string `json:"details"`
	// A value from among the enumerated list indicating the source of the assertion.
	Type string `json:"type"`
}

// Identifies a person responsible for an action.
type ActorClass struct {
	// List of references to W3C Verifiable Credentials.
	Credentials []CredentialElement `json:"credentials"`
	// An identifier for a human actor, used when the "type" is `humanEntry.identified`.
	Identifier *string `json:"identifier"`
}

// A `HashedUri` provides a reference to content available within the same
// manifest store.
//
// This is described in [URI References in the C2PA Technical
//
// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_uri_references).
type CredentialElement struct {
	// A string identifying the cryptographic hash algorithm used to compute
	// the hash
	Alg *string `json:"alg"`
	// Byte string containing the hash value
	Hash []int64 `json:"hash"`
	// JUMBF URI reference
	URL string `json:"url"`
}

// A rating on an Assertion.
//
// See [C2PA Specification - Review
// Ratings](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_review_ratings).
type ReviewRatingClass struct {
	Code        *string `json:"code"`
	Explanation string  `json:"explanation"`
	Value       int64   `json:"value"`
}

// A map of validation results for a manifest store.
//
// The map contains the validation results for the active manifest and any ingredient
// deltas.
// It is normal for there to be many
type IngredientValidationResults struct {
	// Validation status codes for the ingredient's active manifest. Present if ingredient is a
	// C2PA
	// asset. Not present if the ingredient is not a C2PA asset.
	ActiveManifest *ValidationDeltasClass `json:"activeManifest"`
	// List of any changes/deltas between the current and previous validation results for each
	// ingredient's
	// manifest. Present if the the ingredient is a C2PA asset.
	IngredientDeltas []IngredientDeltaClass `json:"ingredientDeltas"`
	// Time when the validation was performed (RFC 3339 date-time). Used only for document-level
	// validationInfo; not serialized in validationResults (e.g. ingredient assertions).
	ValidationTime *string `json:"validationTime"`
}

// Contains a set of success, informational, and failure validation status codes.
//
// Validation results for the ingredient's active manifest
type ValidationDeltasClass struct {
	// An array of validation failure codes. May be empty.
	Failure []FailureElement `json:"failure"`
	// An array of validation informational codes. May be empty.
	Informational []FailureElement `json:"informational"`
	// An array of validation success codes. May be empty.
	Success []FailureElement `json:"success"`
}

// A `ValidationStatus` struct describes the validation status of a
// specific part of a manifest.
//
// See [Existing Manifests - C2PA Technical
// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_existing_manifests).
type FailureElement struct {
	Code        string  `json:"code"`
	Explanation *string `json:"explanation"`
	Success     *bool   `json:"success"`
	URL         *string `json:"url"`
}

// Represents any changes or deltas between the current and previous validation results for
// an ingredient's manifest.
type IngredientDeltaClass struct {
	// JUMBF URI reference to the ingredient assertion
	IngredientAssertionURI string `json:"ingredientAssertionURI"`
	// Validation results for the ingredient's active manifest
	ValidationDeltas ValidationDeltasClass `json:"validationDeltas"`
}

// Holds information about a signature
type SignatureInfoClass struct {
	// Human-readable issuing authority for this signature.
	Alg *AlgEnum `json:"alg"`
	// The serial number of the certificate.
	CERTSerialNumber *string `json:"cert_serial_number"`
	// Human-readable for common name of this certificate.
	CommonName *string `json:"common_name"`
	// Human-readable issuing authority for this signature.
	Issuer *string `json:"issuer"`
	// Revocation status of the certificate.
	RevocationStatus *bool `json:"revocation_status"`
	// The time the signature was created.
	Time *string `json:"time"`
}

// Settings for configuring all aspects of c2pa-rs.
//
// [Settings::default] will be set thread-locally by default. Any settings set via
// [Settings::from_toml] or [Settings::from_file] will also be thread-local.
type Settings struct {
	// Settings for configuring the [`Builder`].
	//
	// [`Builder`]: crate::Builder
	Builder *Builder `json:"builder,omitempty"`
	// Settings for configuring the CAWG trust lists.
	CawgTrust *CawgTrust `json:"cawg_trust,omitempty"`
	// Settings for configuring the CAWG x509 signer, accessible via [`Settings::signer`].
	CawgX509Signer *CawgX509SignerClass `json:"cawg_x509_signer"`
	// Settings for configuring core features.
	Core *Core `json:"core,omitempty"`
	// Settings for configuring the base C2PA signer, accessible via [`Settings::signer`].
	Signer *CawgX509SignerClass `json:"signer"`
	// Settings for configuring the C2PA trust lists.
	Trust *CawgTrust `json:"trust,omitempty"`
	// Settings for configuring verification.
	Verify *Verify `json:"verify,omitempty"`
	// Version of the configuration.
	Version *int64 `json:"version,omitempty"`
}

// Settings for configuring the [`Builder`].
//
// [`Builder`]: crate::Builder
//
// Settings for the [Builder][crate::Builder].
type Builder struct {
	// Settings for configuring fields in an [Actions][crate::assertions::Actions] assertion.
	//
	// For more information on the reasoning behind this field see [ActionsSettings].
	Actions Actions `json:"actions"`
	// Settings for configuring auto-generation of the [`TimeStamp`] assertion.
	//
	// [`TimeStamp`]: crate::assertions::TimeStamp
	AutoTimestampAssertion AutoTimestampAssertion `json:"auto_timestamp_assertion"`
	// Whether to create [`CertificateStatus`] assertions for manifests to store certificate
	// revocation
	// status. The assertion can be fetched for the active manifest or for all manifests
	// (including
	// ingredients).
	//
	// The default is to not fetch them at all.
	//
	// For more information, see [Certificate status assertion - C2PA Technical
	// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#certificate_status_assertion).
	//
	// [`CertificateStatus`]: crate::assertions::CertificateStatus
	CertificateStatusFetch *CertificateStatusFetchEnum `json:"certificate_status_fetch"`
	// Whether to only use [`CertificateStatus`] assertions to check certificate revocation
	// status. If there
	// is a stapled OCSP in the COSE claim of the manifest, it will be ignored. If
	// [`Verify::ocsp_fetch`] is
	// enabled, it will also be ignored.
	//
	// The default value is false.
	//
	// [`CertificateStatus`]: crate::assertions::CertificateStatus
	// [`Verify::ocsp_fetch`]: crate::settings::Verify::ocsp_fetch
	CertificateStatusShouldOverride *bool `json:"certificate_status_should_override"`
	// When set, used as [`ClaimGeneratorInfo`] when
	// [`ManifestDefinition::claim_generator_info`](crate::builder::ManifestDefinition) is empty
	// (e.g. key omitted in JSON or an empty array). If `None` or when the definition lists at
	// least one generator, that path does not use this value.
	ClaimGeneratorInfo *SoftwareAgentClass `json:"claim_generator_info"`
	// Assertions with a base label included in this list will be automatically marked as a
	// created assertion.
	// Assertions not in this list will be automatically marked as gathered.
	//
	// Note that the label should be a **base label**, not including the assertion version nor
	// instance.
	//
	// See more information on the difference between created vs gathered assertions in the spec
	// here:
	// [fields - C2PA Technical
	// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_fields)
	CreatedAssertionLabels []string `json:"created_assertion_labels"`
	// Whether to generate a C2PA archive (instead of zip) when writing the manifest builder.
	// Now always defaults to true - the ability to disable it will be removed in the future.
	GenerateC2PaArchive *bool `json:"generate_c2pa_archive"`
	// The default [`BuilderIntent`] for the [`Builder`].
	//
	// See [`BuilderIntent`] for more information.
	//
	// [`BuilderIntent`]: crate::BuilderIntent
	// [`Builder`]: crate::Builder
	Intent *Intent `json:"intent"`
	// When `true`, use [`BoxHash`] instead of [`crate::assertions::DataHash`] for formats
	// that support it (JPEG, PNG, GIF, etc.) when no explicit hard binding assertion has
	// been set.
	//
	// Formats that support `BoxHash` can embed the C2PA manifest as a new chunk/segment
	// without shifting existing byte offsets, so a placeholder is never required.
	// Setting this to `true` enables the direct workflow (`Builder::sign_embeddable`
	// Mode 2) for those formats and makes `Builder::needs_placeholder` return `false`.
	//
	// Defaults to `false` to preserve existing behaviour until `BoxHash` support is
	// more widely tested.  Set to `true` (or configure it per-[`Context`]) whenever
	// you are ready to prefer box-based hashing for supported formats.
	//
	// [`BoxHash`]: crate::assertions::BoxHash
	// [`Context`]: crate::Context
	PreferBoxHash bool `json:"prefer_box_hash"`
	// Various settings for configuring automatic thumbnail generation.
	Thumbnail Thumbnail `json:"thumbnail"`
	// The name of the vendor creating the content credential.
	Vendor *string `json:"vendor"`
}

// Settings for configuring fields in an [Actions][crate::assertions::Actions] assertion.
//
// For more information on the reasoning behind this field see [ActionsSettings].
//
// Settings for configuring the "base" [Actions][crate::assertions::Actions] assertion.
//
// The reason this setting exists only for an [Actions][crate::assertions::Actions]
// assertion
// is because of its mandations and reusable fields.
type Actions struct {
	// Whether or not to set the
	// [Actions::all_actions_included][crate::assertions::Actions::all_actions_included]
	// field.
	AllActionsIncluded *bool `json:"all_actions_included"`
	// Whether to automatically generate a c2pa.created [Action] assertion or error that it
	// doesn't already exist.
	//
	// For more information about the mandatory conditions for a c2pa.created action assertion,
	// see the
	// [C2PA Technical
	// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_mandatory_presence_of_at_least_one_actions_assertion).
	AutoCreatedAction AutoCreatedAction `json:"auto_created_action"`
	// Whether to automatically generate a c2pa.opened [Action] assertion or error that it
	// doesn't already exist.
	//
	// For more information about the mandatory conditions for a c2pa.opened action assertion,
	// see the
	// [C2PA Technical
	// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_mandatory_presence_of_at_least_one_actions_assertion).
	AutoOpenedAction AutoCreatedAction `json:"auto_opened_action"`
	// Whether to automatically generate a c2pa.placed [Action] assertion or error that it
	// doesn't already exist.
	//
	// For more information about the mandatory conditions for a c2pa.placed action assertion,
	// see
	// [Relationship - C2PA Technical
	// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_relationship)
	AutoPlacedAction AutoCreatedAction `json:"auto_placed_action"`
	// Templates to be added to the [Actions::templates][crate::assertions::Actions::templates]
	// field.
	Templates []TemplateElement `json:"templates"`
}

// Whether to automatically generate a c2pa.created [Action] assertion or error that it
// doesn't already exist.
//
// For more information about the mandatory conditions for a c2pa.created action assertion,
// see the
// [C2PA Technical
// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_mandatory_presence_of_at_least_one_actions_assertion).
//
// Settings for the auto actions (e.g. created, opened, placed).
//
// Whether to automatically generate a c2pa.opened [Action] assertion or error that it
// doesn't already exist.
//
// For more information about the mandatory conditions for a c2pa.opened action assertion,
// see the
// [C2PA Technical
// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_mandatory_presence_of_at_least_one_actions_assertion).
//
// Whether to automatically generate a c2pa.placed [Action] assertion or error that it
// doesn't already exist.
//
// For more information about the mandatory conditions for a c2pa.placed action assertion,
// see
// [Relationship - C2PA Technical
// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_relationship)
type AutoCreatedAction struct {
	// Whether to enable this auto action or not.
	Enabled bool `json:"enabled"`
	// The default source type for the auto action.
	SourceType *string `json:"source_type"`
}

// Settings for an action template.
type TemplateElement struct {
	// The label associated with this action. See
	// ([c2pa_action][crate::assertions::actions::c2pa_action]).
	Action string `json:"action"`
	// Description of the template.
	Description *string `json:"description"`
	// Reference to an icon.
	Icon *TemplateIcon `json:"icon"`
	// The software agent that performed the action.
	SoftwareAgent *SoftwareAgentClass `json:"software_agent"`
	// 0-based index into the softwareAgents array
	SoftwareAgentIndex *int64 `json:"software_agent_index"`
	// One of the defined URI values at `<https://cv.iptc.org/newscodes/digitalsourcetype/>`
	SourceType *string `json:"source_type"`
	// Additional parameters for the template
	TemplateParameters map[string]interface{} `json:"template_parameters"`
}

// A reference to a resource to be used in JSON serialization.
//
// The underlying data can be read as a stream via
// [`Reader::resource_to_stream`][crate::Reader::resource_to_stream].
type TemplateIcon struct {
	// The algorithm used to hash the resource (if applicable).
	Alg *string `json:"alg"`
	// More detailed data types as defined in the C2PA spec.
	DataTypes []IconDataType `json:"data_types"`
	// The mime type of the referenced resource.
	Format string `json:"format"`
	// The hash of the resource (if applicable).
	Hash *string `json:"hash"`
	// A URI that identifies the resource as referenced from the manifest.
	//
	// This may be a JUMBF URI, a file path, a URL or any other string.
	// Relative JUMBF URIs will be resolved with the manifest label.
	// Relative file paths will be resolved with the base path if provided.
	Identifier string `json:"identifier"`
}

type IconDataType struct {
	Type    string  `json:"type"`
	Version *string `json:"version"`
}

// Settings for the claim generator info.
type SoftwareAgentClass struct {
	// Reference to an icon.
	Icon *TemplateIcon `json:"icon"`
	// A human readable string naming the claim_generator.
	Name string `json:"name"`
	// Settings for the claim generator info's operating system field.
	OperatingSystem *string `json:"operating_system"`
	// A human readable string of the product's version.
	Version *string `json:"version"`
}

// Settings for configuring auto-generation of the [`TimeStamp`] assertion.
//
// [`TimeStamp`]: crate::assertions::TimeStamp
//
// Settings for configuring auto-generation of the [`TimeStamp`] assertion.
//
// Useful when a manifest was signed offline and you want to attach a trusted timestamp to
// it later.
//
// [`TimeStamp`]: crate::assertions::TimeStamp
type AutoTimestampAssertion struct {
	// Whether to auto-generate a [`TimeStamp`] assertion for the
	// [`TimeStampSettings::fetch_scope`].
	//
	// Note that for this setting to take effect, a timestamping authority URL must be set in
	// the
	// [`Signer::time_authority_url`]. If the signer is acquired from settings via
	// [`Settings::signer`],
	// the URL can be set in [`SignerSettings`].
	//
	// The default value is false.
	//
	// [`TimeStamp`]: crate::assertions::TimeStamp
	// [`Signer::time_authority_url`]: crate::Signer::time_authority_url
	// [`Settings::signer`]: crate::settings::signer
	// [`SignerSettings`]: crate::settings::signer::SignerSettings
	Enabled bool `json:"enabled"`
	// Which manifests to fetch timestamps for.
	//
	// The default value is [`TimeStampFetchScope::All`].
	FetchScope FetchScope `json:"fetch_scope"`
	// Whether to skip fetching timestamps for manifests that already have one.
	//
	// This setting will account for both existing [`TimeStamp`] assertions and timestamps
	// embedded
	// in the claim.
	//
	// The default value is true.
	//
	// [`TimeStamp`]: crate::assertions::TimeStamp
	SkipExisting bool `json:"skip_existing"`
}

// This is a new digital creation, a DigitalSourceType is required.
//
// The Manifest must not have have a parent ingredient.
// A `c2pa.created` action will be added if not provided.
type IntentClass struct {
	Create string `json:"create"`
}

// Various settings for configuring automatic thumbnail generation.
//
// Settings for controlling automatic thumbnail generation.
type Thumbnail struct {
	// Whether or not to automatically generate thumbnails.
	//
	// The default value is true.
	//
	// <div class="warning">
	// This setting is only applicable if the crate is compiled with the `add_thumbnails`
	// feature.
	// </div>
	Enabled bool `json:"enabled"`
	// Format of the thumbnail.
	//
	// If this field isn't specified, the thumbnail format will correspond to the
	// input format.
	//
	// The default value is None.
	Format *FormatEnum `json:"format"`
	// Whether to ignore thumbnail generation errors.
	//
	// This may occur, for instance, if the thumbnail media type or color layout isn't
	// supported.
	//
	// The default value is true.
	IgnoreErrors bool `json:"ignore_errors"`
	// The size of the longest edge of the thumbnail.
	//
	// This function will resize the input to preserve aspect ratio.
	//
	// The default value is 1024.
	LongEdge int64 `json:"long_edge"`
	// Whether or not to prefer a smaller sized media format for the thumbnail.
	//
	// Note that [ThumbnailSettings::format] takes precedence over this field. In addition,
	// if the output format is unsupported, it will default to the smallest format regardless
	// of the value of this field.
	//
	// For instance, if the source input type is a PNG, but it doesn't have an alpha channel,
	// the image will be converted to a JPEG of smaller size.
	//
	// The default value is true.
	PreferSmallestFormat bool `json:"prefer_smallest_format"`
	// The output quality of the thumbnail.
	//
	// This setting contains sensible defaults for things like quality, compression, and
	// algorithms for various formats.
	//
	// The default value is [`ThumbnailQuality::Medium`].
	Quality Quality `json:"quality"`
}

// Settings for configuring the CAWG trust lists.
//
// Settings to configure the trust list.
//
// Settings for configuring the C2PA trust lists.
type CawgTrust struct {
	// List of explicitly allowed certificates as a PEM bundle.
	AllowedList *string `json:"allowed_list"`
	// List of default trust anchor root certificates as a PEM bundle.
	//
	// Normally this option contains the official C2PA-recognized trust anchors found here:
	// <https://github.com/c2pa-org/conformance-public/tree/main/trust-list>
	TrustAnchors *string `json:"trust_anchors"`
	// List of allowed extended key usage (EKU) object identifiers (OID) that
	// certificates must have.
	TrustConfig *string `json:"trust_config"`
	// List of additional user-provided trust anchor root certificates as a PEM bundle.
	UserAnchors *string `json:"user_anchors"`
	// Whether to verify certificates against the trust lists specified in [`Trust`]. This
	// option is ONLY applicable to CAWG.
	//
	// The default value is true.
	//
	// <div class="warning">
	// Verifying trust is REQUIRED by the CAWG spec. This option should only be used for
	// development or testing.
	// </div>
	VerifyTrustList *bool `json:"verify_trust_list,omitempty"`
}

// A signer configured locally.
//
// A signer configured remotely.
type CawgX509SignerClass struct {
	Local  *Local  `json:"local,omitempty"`
	Remote *Remote `json:"remote,omitempty"`
}

type Local struct {
	// Algorithm to use for signing.
	Alg AlgEnum `json:"alg"`
	// Private key used for signing (PEM format).
	PrivateKey string `json:"private_key"`
	// Referenced assertions for CAWG identity signing (optional).
	ReferencedAssertions []string `json:"referenced_assertions"`
	// Roles for CAWG identity signing (optional).
	Roles []string `json:"roles"`
	// Certificate used for signing (PEM format).
	SignCERT string `json:"sign_cert"`
	// Time stamp authority URL for signing.
	TsaURL *string `json:"tsa_url"`
}

type Remote struct {
	// Algorithm to use for signing.
	Alg AlgEnum `json:"alg"`
	// Referenced assertions for CAWG identity signing (optional).
	ReferencedAssertions []string `json:"referenced_assertions"`
	// Roles for CAWG identity signing (optional).
	Roles []string `json:"roles"`
	// Certificate used for signing (PEM format).
	SignCERT string `json:"sign_cert"`
	// Time stamp authority URL for signing.
	TsaURL *string `json:"tsa_url"`
	// URL that the signer will use for signing.
	// A POST request with a byte-stream will be sent to this URL.
	URL string `json:"url"`
}

// Settings for configuring core features.
//
// Settings to configure core features.
type Core struct {
	// <div class="warning">
	// The CAWG identity assertion does not currently respect this setting.
	// See [Issue #1645](https://github.com/contentauth/c2pa-rs/issues/1645).
	// </div>
	//
	// List of host patterns that are allowed for network requests.
	//
	// Each pattern may include:
	// - A scheme (e.g. `https://` or `http://`)
	// - A hostname or IP address (e.g. `contentauthenticity.org` or `192.0.2.1`)
	// - The hostname may contain a single leading wildcard (e.g. `*.contentauthenticity.org`)
	// - An optional port (e.g. `contentauthenticity.org:443` or `192.0.2.1:8080`)
	//
	// Matching is case-insensitive. A wildcard pattern such as `*.contentauthenticity.org`
	// matches
	// `sub.contentauthenticity.org`, but does not match `contentauthenticity.org` or
	// `fakecontentauthenticity.org`.
	// If a scheme is present in the pattern, only URIs using the same scheme are considered a
	// match. If the scheme
	// is omitted, any scheme is allowed as long as the host matches.
	//
	// The behavior is as follows:
	// - `None` (default) no filtering enabled.
	// - `Some(vec)` where `vec` is empty, all traffic is blocked.
	// - `Some(vec)` with at least one pattern, filtering enabled for only those patterns.
	//
	// # Examples
	//
	// Pattern: `*.contentauthenticity.org`
	// - Does match:
	// - `https://sub.contentauthenticity.org`
	// - `http://api.contentauthenticity.org`
	// - Does **not** match:
	// - `https://contentauthenticity.org` (no subdomain)
	// - `https://sub.fakecontentauthenticity.org` (different host)
	//
	// Pattern: `http://192.0.2.1:8080`
	// - Does match:
	// - `http://192.0.2.1:8080`
	// - Does **not** match:
	// - `https://192.0.2.1:8080` (scheme mismatch)
	// - `http://192.0.2.1` (port omitted)
	// - `http://192.0.2.2:8080` (different IP address)
	//
	// These settings are applied by the SDK's HTTP resolvers to restrict network requests.
	// When network requests occur depends on the operations being performed (reading manifests,
	// validating credentials, timestamping, etc.).
	AllowedNetworkHosts []string `json:"allowed_network_hosts"`
	// Maximum amount of data in megabytes that will be loaded into memory before
	// being stored in temporary files on the disk.
	//
	// This option defaults to 512MB and can result in noticeable performance improvements.
	BackingStoreMemoryThresholdInMB *int64 `json:"backing_store_memory_threshold_in_mb,omitempty"`
	// Whether to decode CAWG [`IdentityAssertion`]s during reading in the [`Reader`].
	//
	// This option defaults to true.
	//
	// [`IdentityAssertion`]: crate::identity::IdentityAssertion
	// [`Reader`]: crate::Reader
	DecodeIdentityAssertions *bool `json:"decode_identity_assertions,omitempty"`
	// Maximum size in megabytes of a Brotli-decompressed JUMBF manifest.
	// Limits memory consumption from decompression bomb attacks.
	//
	// The default is 32 MB.
	MaxDecompressedManifestSizeInMB *int64 `json:"max_decompressed_manifest_size_in_mb,omitempty"`
	// Size of the [`BmffHash`] merkle tree chunks in kilobytes.
	//
	// This option is associated with the [`MerkleMap::fixed_block_size`] field.
	//
	// See more information in the spec here:
	// [bmff_based_hash - C2PA Technical
	// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_bmff_based_hash)
	//
	// [`MerkleMap::fixed_block_size`]: crate::assertions::MerkleMap::fixed_block_size
	// [`BmffHash`]: crate::assertions::BmffHash
	MerkleTreeChunkSizeInKB *int64 `json:"merkle_tree_chunk_size_in_kb"`
	// Maximum number of proof hashes stored in UUID merkle boxes when  generating a
	// [`BmffHash`] merkle tree.  This
	// determines the Merkle tree row stored in the manifest and thus the number of proof hashes
	// that need to be
	// provided during validation. The value may be 0 to store just leaf node hashes (no UUID
	// boxes are generated in this case).
	//
	// This option defaults to 5.
	//
	// See more information in the spec here:
	// [bmff_based_hash - C2PA Technical
	// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_bmff_based_hash)
	//
	// [`BmffHash`]: crate::assertions::BmffHash
	MerkleTreeMaxProofs *int64 `json:"merkle_tree_max_proofs,omitempty"`
	// Whether to prefer compressing manifests. This can reduce the size of the manifest.
	// Compressed manifest
	// are not always possible and will default back to uncompressed if the manifest contains
	// features
	// that are not compatible with compression.
	//
	// The default value is false.
	//
	// See more information in the spec here:
	// [Compressed manifests - C2PA Technical
	// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_compressed_boxes)
	PreferCompressManifests *bool `json:"prefer_compress_manifests,omitempty"`
}

// Settings for configuring verification.
//
// Settings to configure the verification process.
type Verify struct {
	// Whether to fetch the certificates OCSP status during validation.
	//
	// Revocation status is checked in the following order:
	// 1. The OCSP staple stored in the COSE claim of the manifest
	// 2. Otherwise if `ocsp_fetch` is enabled, it fetches a new OCSP status
	// 3. Otherwise if `ocsp_fetch` is disabled, it checks `CertificateStatus` assertions
	//
	// The default value is false.
	OcspFetch *bool `json:"ocsp_fetch,omitempty"`
	// Whether to fetch remote manifests in the following scenarios:
	// - Constructing a [`Reader`]
	// - Adding an [`Ingredient`] to the [`Builder`]
	//
	// The default value is true.
	//
	// <div class="warning">
	// This setting is only applicable if the crate is compiled with the
	// `fetch_remote_manifests` feature.
	// </div>
	//
	// [`Reader`]: crate::Reader
	// [`Ingredient`]: crate::Ingredient
	// [`Builder`]: crate::Builder
	RemoteManifestFetch *bool `json:"remote_manifest_fetch,omitempty"`
	// Whether to skip ingredient conflict resolution when multiple ingredients have the same
	// manifest identifier. This settings is only applicable for C2PA v2 validation.
	//
	// The default value is false.
	//
	// See more information in the spec here:
	// [versioning_manifests_due_to_conflicts - C2PA Technical
	// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_versioning_manifests_due_to_conflicts)
	SkipIngredientConflictResolution *bool `json:"skip_ingredient_conflict_resolution,omitempty"`
	// Whether to do strictly C2PA v1 validation or otherwise the latest validation.
	//
	// The default value is false.
	StrictV1Validation *bool `json:"strict_v1_validation,omitempty"`
	// Whether to verify the manifest after reading in the [`Reader`].
	//
	// The default value is true.
	//
	// <div class="warning">
	// Disabling validation can improve reading performance, BUT it carries the risk of reading
	// an invalid
	// manifest.
	// </div>
	//
	// [`Reader`]: crate::Reader
	VerifyAfterReading *bool `json:"verify_after_reading,omitempty"`
	// Whether to verify the manifest after signing in the [`Builder`].
	//
	// The default value is false.
	//
	// In the future, this setting will default to true.
	//
	// <div class="warning">
	// Disabling validation can improve signing performance, BUT it carries the risk of signing
	// an invalid
	// manifest.
	// </div>
	//
	// [`Builder`]: crate::Builder
	VerifyAfterSign *bool `json:"verify_after_sign,omitempty"`
	// Whether to include asset hash validation when verifying after signing.
	//
	// The default value is false.
	//
	// Has no effect when [`Verify::verify_after_sign`] is false.
	VerifyAfterSignHash *bool `json:"verify_after_sign_hash,omitempty"`
	// Whether to verify the timestamp certificates against the trust lists specified in
	// [`Trust`].
	//
	// The default value is true.
	//
	// <div class="warning">
	// Verifying timestamp trust is REQUIRED by the C2PA spec. This option should only be used
	// for development or testing.
	// </div>
	VerifyTimestampTrust *bool `json:"verify_timestamp_trust,omitempty"`
	// Whether to verify certificates against the trust lists specified in [`Trust`]. To
	// configure
	// timestamp certificate verification, see [`Verify::verify_timestamp_trust`].
	//
	// The default value is true.
	//
	// <div class="warning">
	// Verifying trust is REQUIRED by the C2PA spec. This option should only be used for
	// development or testing.
	// </div>
	VerifyTrust *bool `json:"verify_trust,omitempty"`
}

// Assertions in C2PA can be stored in several formats
type KindEnum string

const (
	Binary KindEnum = "Binary"
	Cbor   KindEnum = "Cbor"
	JSON   KindEnum = "Json"
	URI    KindEnum = "Uri"
)

// The type of shape.
//
// The type of shape for the range.
//
// A rectangle.
//
// A circle.
//
// A polygon.
type ShapeType string

const (
	Circle    ShapeType = "circle"
	Polygon   ShapeType = "polygon"
	Rectangle ShapeType = "rectangle"
)

// The type of unit for the shape range.
//
// The type of unit for the range.
//
// Use pixels.
//
// Use percentage.
type Unit string

const (
	Percent Unit = "percent"
	Pixel   Unit = "pixel"
)

// The type of time.
//
// Times are described using Normal Play Time (npt) as described in RFC 2326.
type TimeType string

const (
	Npt TimeType = "npt"
)

// The type of range of interest.
//
// The type of range for the region of interest.
//
// A spatial range, see [`Shape`] for more details.
//
// A temporal range, see [`Time`] for more details.
//
// A spatial range, see [`Frame`] for more details.
//
// A textual range, see [`Text`] for more details.
//
// A range identified by a specific identifier and value, see [`Item`] for more details.
type RegionType string

const (
	Frame      RegionType = "frame"
	Identified RegionType = "identified"
	Spatial    RegionType = "spatial"
	Temporal   RegionType = "temporal"
	Textual    RegionType = "textual"
)

// Arbitrary area worth identifying.
//
// This area is all that is left after a crop action.
//
// This area has had edits applied to it.
//
// The area where an ingredient was placed/added.
//
// Something in this area was redacted.
//
// Area specific to a subject (human or not).
//
// A range of information was removed/deleted.
//
// Styling was applied to this area.
//
// Invisible watermarking was applied to this area for the purpose of soft binding.
type RoleEnum string

const (
	C2PaAreaOfInterest RoleEnum = "c2pa.areaOfInterest"
	C2PaCropped        RoleEnum = "c2pa.cropped"
	C2PaDeleted        RoleEnum = "c2pa.deleted"
	C2PaEdited         RoleEnum = "c2pa.edited"
	C2PaPlaced         RoleEnum = "c2pa.placed"
	C2PaRedacted       RoleEnum = "c2pa.redacted"
	C2PaStyled         RoleEnum = "c2pa.styled"
	C2PaSubjectArea    RoleEnum = "c2pa.subjectArea"
	C2PaWatermarked    RoleEnum = "c2pa.watermarked"
)

// Set to `ParentOf` if this is the parent ingredient.
//
// There can only be one parent ingredient in the ingredients.
//
// The relationship of the ingredient to the current asset.
//
// The current asset is derived from this ingredient.
//
// The current asset is a part of this ingredient.
//
// The ingredient was used as an input to a computational process to create or modify the
// asset.
type Relationship string

const (
	ComponentOf Relationship = "componentOf"
	InputTo     Relationship = "inputTo"
	ParentOf    Relationship = "parentOf"
)

// ECDSA with SHA-256
//
// # ECDSA with SHA-384
//
// # ECDSA with SHA-512
//
// # RSASSA-PSS using SHA-256 and MGF1 with SHA-256
//
// # RSASSA-PSS using SHA-384 and MGF1 with SHA-384
//
// # RSASSA-PSS using SHA-512 and MGF1 with SHA-512
//
// Edwards-Curve DSA (Ed25519 instance only)
//
// Algorithm to use for signing.
//
// Describes the digital signature algorithms allowed by the C2PA spec.
//
// Per [§13.2, “Digital Signatures”]:
//
// > All digital signatures applied as per the technical requirements of this
// > specification shall be generated using one of the digital signature
// > algorithms and key types listed as described in this section.
//
// [§13.2, “Digital Signatures”]:
// https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_digital_signatures
type AlgEnum string

const (
	Ed25519 AlgEnum = "Ed25519"
	Es256   AlgEnum = "Es256"
	Es384   AlgEnum = "Es384"
	Es512   AlgEnum = "Es512"
	Ps256   AlgEnum = "Ps256"
	Ps384   AlgEnum = "Ps384"
	Ps512   AlgEnum = "Ps512"
)

// The manifest store fails to meet ValidationState::WellFormed requirements, meaning it
// cannot
// even be parsed or its basic structure is non-compliant.
//
// This case may also occur if validation is disabled in the SDK.
//
// The manifest store is well-formed and the cryptographic integrity checks succeed.
//
// See [Valid Manifest - C2PA Technical
// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_valid_manifest).
//
// The manifest store is valid and signed by a certificate that chains up to a trusted root
// or known
// authority in the trust list.
//
// See [Trusted Manifest - C2PA Technical
// Specification](https://spec.c2pa.org/specifications/specifications/2.3/specs/C2PA_Specification.html#_trusted_manifest).
type ValidationStateEnum string

const (
	Invalid ValidationStateEnum = "Invalid"
	Trusted ValidationStateEnum = "Trusted"
	Valid   ValidationStateEnum = "Valid"
)

// Which manifests to fetch timestamps for.
//
// The default value is [`TimeStampFetchScope::All`].
//
// The scope of manifests to fetch timestamps for.
//
// See [`TimeStampSettings`] for more information.
//
// Fetch timestamps for only the parent manifest.
//
// Fetch timestmaps for all manifests in the manifest store.
type FetchScope string

const (
	FetchScopeAll FetchScope = "all"
	Parent        FetchScope = "parent"
)

// Fetch OCSP for all manifests.
//
// Fetch OCSP for the active manifest only.
type CertificateStatusFetchEnum string

const (
	Active            CertificateStatusFetchEnum = "active"
	SettingsSchemaAll CertificateStatusFetchEnum = "all"
)

// This is an edit of a pre-existing parent asset.
//
// The Manifest must have a parent ingredient.
// A parent ingredient will be generated from the source stream if not otherwise provided.
// A `c2pa.opened action will be tied to the parent ingredient.
//
// A restricted version of Edit for non-editorial changes.
//
// There must be only one ingredient, as a parent.
// No changes can be made to the hashed content of the parent.
// There are additional restrictions on the types of changes that can be made.
type IntentEnum string

const (
	Edit   IntentEnum = "edit"
	Update IntentEnum = "update"
)

// An image in PNG format.
//
// An image in JPEG format.
//
// An image in GIF format.
//
// An image in WEBP format.
//
// An image in TIFF format.
type FormatEnum string

const (
	GIF  FormatEnum = "gif"
	JPEG FormatEnum = "jpeg"
	PNG  FormatEnum = "png"
	Tiff FormatEnum = "tiff"
	Webp FormatEnum = "webp"
)

// The output quality of the thumbnail.
//
// This setting contains sensible defaults for things like quality, compression, and
// algorithms for various formats.
//
// The default value is [`ThumbnailQuality::Medium`].
//
// Quality of the thumbnail.
//
// Low quality.
//
// Medium quality.
//
// High quality.
type Quality string

const (
	High   Quality = "high"
	Low    Quality = "low"
	Medium Quality = "medium"
)

type HashUnion struct {
	IntegerArray []int64
	String       *string
}

func (x *HashUnion) UnmarshalJSON(data []byte) error {
	x.IntegerArray = nil
	object, err := unmarshalUnion(data, nil, nil, nil, &x.String, true, &x.IntegerArray, false, nil, false, nil, false, nil, true)
	if err != nil {
		return err
	}
	if object {
	}
	return nil
}

func (x *HashUnion) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, x.String, x.IntegerArray != nil, x.IntegerArray, false, nil, false, nil, false, nil, true)
}

// The default [`BuilderIntent`] for the [`Builder`].
//
// See [`BuilderIntent`] for more information.
//
// [`BuilderIntent`]: crate::BuilderIntent
// [`Builder`]: crate::Builder
type Intent struct {
	Enum        *IntentEnum
	IntentClass *IntentClass
}

func (x *Intent) UnmarshalJSON(data []byte) error {
	x.IntentClass = nil
	x.Enum = nil
	var c IntentClass
	object, err := unmarshalUnion(data, nil, nil, nil, nil, false, nil, true, &c, false, nil, true, &x.Enum, true)
	if err != nil {
		return err
	}
	if object {
		x.IntentClass = &c
	}
	return nil
}

func (x *Intent) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, false, nil, x.IntentClass != nil, x.IntentClass, false, nil, x.Enum != nil, x.Enum, true)
}

func unmarshalUnion(data []byte, pi **int64, pf **float64, pb **bool, ps **string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) (bool, error) {
	if pi != nil {
		*pi = nil
	}
	if pf != nil {
		*pf = nil
	}
	if pb != nil {
		*pb = nil
	}
	if ps != nil {
		*ps = nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return false, err
	}

	switch v := tok.(type) {
	case json.Number:
		if pi != nil {
			i, err := v.Int64()
			if err == nil {
				*pi = &i
				return false, nil
			}
		}
		if pf != nil {
			f, err := v.Float64()
			if err == nil {
				*pf = &f
				return false, nil
			}
			return false, errors.New("Unparsable number")
		}
		return false, errors.New("Union does not contain number")
	case float64:
		return false, errors.New("Decoder should not return float64")
	case bool:
		if pb != nil {
			*pb = &v
			return false, nil
		}
		return false, errors.New("Union does not contain bool")
	case string:
		if haveEnum {
			return false, json.Unmarshal(data, pe)
		}
		if ps != nil {
			*ps = &v
			return false, nil
		}
		return false, errors.New("Union does not contain string")
	case nil:
		if nullable {
			return false, nil
		}
		return false, errors.New("Union does not contain null")
	case json.Delim:
		if v == '{' {
			if haveObject {
				return true, json.Unmarshal(data, pc)
			}
			if haveMap {
				return false, json.Unmarshal(data, pm)
			}
			return false, errors.New("Union does not contain object")
		}
		if v == '[' {
			if haveArray {
				return false, json.Unmarshal(data, pa)
			}
			return false, errors.New("Union does not contain array")
		}
		return false, errors.New("Cannot handle delimiter")
	}
	return false, errors.New("Cannot unmarshal union")
}

func marshalUnion(pi *int64, pf *float64, pb *bool, ps *string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) ([]byte, error) {
	if pi != nil {
		return json.Marshal(*pi)
	}
	if pf != nil {
		return json.Marshal(*pf)
	}
	if pb != nil {
		return json.Marshal(*pb)
	}
	if ps != nil {
		return json.Marshal(*ps)
	}
	if haveArray {
		return json.Marshal(pa)
	}
	if haveObject {
		return json.Marshal(pc)
	}
	if haveMap {
		return json.Marshal(pm)
	}
	if haveEnum {
		return json.Marshal(pe)
	}
	if nullable {
		return json.Marshal(nil)
	}
	return nil, errors.New("Union must not be null")
}
