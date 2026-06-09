package c2pa

// This file defines the public string-typed enumerations used across the
// package. The runtime mapping to the underlying libc2pa_c integer enums
// lives privately in native.go / native_nocgo.go, so callers never need to
// know — or stay in sync with — the upstream enum ordinals.

// SigningAlg names a signing algorithm. The string values match the algorithm
// names accepted by C2paSignerInfo (e.g. "ps256", "ed25519").
type SigningAlg string

const (
	SigningAlgEs256   SigningAlg = "es256"
	SigningAlgEs384   SigningAlg = "es384"
	SigningAlgEs512   SigningAlg = "es512"
	SigningAlgPs256   SigningAlg = "ps256"
	SigningAlgPs384   SigningAlg = "ps384"
	SigningAlgPs512   SigningAlg = "ps512"
	SigningAlgEd25519 SigningAlg = "ed25519"
)

// BuilderIntent names the kind of manifest a Builder will produce.
type BuilderIntent string

const (
	// IntentCreate is a new digital creation; the manifest must not have a
	// parent ingredient.
	IntentCreate BuilderIntent = "create"
	// IntentEdit is an edit of a pre-existing parent asset; the manifest
	// must have a parent ingredient.
	IntentEdit BuilderIntent = "edit"
	// IntentUpdate is a restricted edit for non-editorial changes; there
	// must be exactly one ingredient (as a parent).
	IntentUpdate BuilderIntent = "update"
)

// DigitalSourceType names the source of a digitally created asset.
type DigitalSourceType string

const (
	SourceEmpty                                DigitalSourceType = "empty"
	SourceTrainedAlgorithmicData               DigitalSourceType = "trainedAlgorithmicData"
	SourceDigitalCapture                       DigitalSourceType = "digitalCapture"
	SourceComputationalCapture                 DigitalSourceType = "computationalCapture"
	SourceNegativeFilm                         DigitalSourceType = "negativeFilm"
	SourcePositiveFilm                         DigitalSourceType = "positiveFilm"
	SourcePrint                                DigitalSourceType = "print"
	SourceHumanEdits                           DigitalSourceType = "humanEdits"
	SourceCompositeWithTrainedAlgorithmicMedia DigitalSourceType = "compositeWithTrainedAlgorithmicMedia"
	SourceAlgorithmicallyEnhanced              DigitalSourceType = "algorithmicallyEnhanced"
	SourceDigitalCreation                      DigitalSourceType = "digitalCreation"
	SourceDataDrivenMedia                      DigitalSourceType = "dataDrivenMedia"
	SourceTrainedAlgorithmicMedia              DigitalSourceType = "trainedAlgorithmicMedia"
	SourceAlgorithmicMedia                     DigitalSourceType = "algorithmicMedia"
	SourceScreenCapture                        DigitalSourceType = "screenCapture"
	SourceVirtualRecording                     DigitalSourceType = "virtualRecording"
	SourceComposite                            DigitalSourceType = "composite"
	SourceCompositeCapture                     DigitalSourceType = "compositeCapture"
	SourceCompositeSynthetic                   DigitalSourceType = "compositeSynthetic"
)

// HashType names the hash binding type a Builder will produce for a given
// format in the embeddable signing workflow.
type HashType string

const (
	// HashTypeData uses placeholder + exclusions + hash + sign (JPEG, PNG, …).
	HashTypeData HashType = "data"
	// HashTypeBmff uses placeholder + hash + sign (MP4, AVIF, HEIF/HEIC).
	HashTypeBmff HashType = "bmff"
	// HashTypeBox uses hash + sign with no placeholder needed.
	HashTypeBox HashType = "box"
)

// ProgressPhase identifies the operation stage reported by a progress callback.
type ProgressPhase string

const (
	ProgressReading                ProgressPhase = "reading"
	ProgressVerifyingManifest      ProgressPhase = "verifyingManifest"
	ProgressVerifyingSignature     ProgressPhase = "verifyingSignature"
	ProgressVerifyingIngredient    ProgressPhase = "verifyingIngredient"
	ProgressVerifyingAssetHash     ProgressPhase = "verifyingAssetHash"
	ProgressAddingIngredient       ProgressPhase = "addingIngredient"
	ProgressThumbnail              ProgressPhase = "thumbnail"
	ProgressHashing                ProgressPhase = "hashing"
	ProgressSigning                ProgressPhase = "signing"
	ProgressEmbedding              ProgressPhase = "embedding"
	ProgressFetchingRemoteManifest ProgressPhase = "fetchingRemoteManifest"
	ProgressWriting                ProgressPhase = "writing"
	ProgressFetchingOCSP           ProgressPhase = "fetchingOCSP"
	ProgressFetchingTimestamp      ProgressPhase = "fetchingTimestamp"
)
