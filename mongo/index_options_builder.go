package mongo

import "github.com/mongodb/mongo-go-driver/bson"

// IndexOptionsBuilder constructs a BSON document for index options
type IndexOptionsBuilder struct {
	document *bson.Document
}

// NewIndexOptionsBuilder creates a new instance of IndexOptionsBuilder
func NewIndexOptionsBuilder() *IndexOptionsBuilder {
	var b IndexOptionsBuilder
	b.document = bson.NewDocument()
	return &b
}

// Background sets the background option
func (iob *IndexOptionsBuilder) Background(background bool) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.Boolean("background", background))
	return iob
}

// ExpireAfter sets the expireAfter option
func (iob *IndexOptionsBuilder) ExpireAfter(expireAfter int32) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.Int32("expireAfter", expireAfter))
	return iob
}

// Name sets the name option
func (iob *IndexOptionsBuilder) Name(name string) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.String("name", name))
	return iob
}

// Sparse sets the sparse option
func (iob *IndexOptionsBuilder) Sparse(sparse bool) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.Boolean("sparse", sparse))
	return iob
}

// StorageEngine sets the storageEngine option
func (iob *IndexOptionsBuilder) StorageEngine(storageEngine string) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.String("storageEngine", storageEngine))
	return iob
}

// Unique sets the unique option
func (iob *IndexOptionsBuilder) Unique(unique bool) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.Boolean("unique", unique))
	return iob
}

// Version sets the verison option
func (iob *IndexOptionsBuilder) Version(version int32) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.Int32("version", version))
	return iob
}

// DefaultLanguage sets the defaultLanguage option
func (iob *IndexOptionsBuilder) DefaultLanguage(defaultLanguage string) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.String("defaultLanguage", defaultLanguage))
	return iob
}

// LanguageOverride sets the languageOverride option
func (iob *IndexOptionsBuilder) LanguageOverride(languageOverride string) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.String("languageOverride", languageOverride))
	return iob
}

// TextVersion sets the textVersion option
func (iob *IndexOptionsBuilder) TextVersion(textVersion int32) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.Int32("textVersion", textVersion))
	return iob
}

// Weights sets the weights option
func (iob *IndexOptionsBuilder) Weights(weights *bson.Document) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.SubDocument("weights", weights))
	return iob
}

// SphereVersion sets the sphereVersion option
func (iob *IndexOptionsBuilder) SphereVersion(sphereVersion int32) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.Int32("sphereVersion", sphereVersion))
	return iob
}

// Bits sets the bits option
func (iob *IndexOptionsBuilder) Bits(bits int32) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.Int32("bits", bits))
	return iob
}

// Max sets the max option
func (iob *IndexOptionsBuilder) Max(max float64) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.Double("max", max))
	return iob
}

// Min sets the min option
func (iob *IndexOptionsBuilder) Min(min float64) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.Double("min", min))
	return iob
}

// BucketSize sets the bucketSize option
func (iob *IndexOptionsBuilder) BucketSize(bucketSize int32) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.Int32("bucketSize", bucketSize))
	return iob
}

// PartialFilterExpression sets the partialFilterExpression option
func (iob *IndexOptionsBuilder) PartialFilterExpression(partialFilterExpression *bson.Document) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.SubDocument("partialFilterExpression", partialFilterExpression))
	return iob
}

// Collation sets the collation option
func (iob *IndexOptionsBuilder) Collation(collation *bson.Document) *IndexOptionsBuilder {
	iob.document.Append(bson.EC.SubDocument("collation", collation))
	return iob
}

// Build returns the BSON document from the builder
func (iob *IndexOptionsBuilder) Build() *bson.Document {
	return iob.document
}
