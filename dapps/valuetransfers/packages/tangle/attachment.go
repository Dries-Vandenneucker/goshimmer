package tangle

import (
	"fmt"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/payload"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/transaction"
	"github.com/iotaledger/hive.go/byteutils"
	"github.com/iotaledger/hive.go/marshalutil"
	"github.com/iotaledger/hive.go/objectstorage"
	"github.com/iotaledger/hive.go/stringify"
)

// Attachment stores the information which transaction was attached by which payload. We need this to be able to perform
// reverse lookups from transactions to their corresponding payloads, that attach them.
type Attachment struct {
	objectstorage.StorableObjectFlags

	transactionID transaction.ID
	payloadID     payload.ID
}

// NewAttachment creates an attachment object with the given information.
func NewAttachment(transactionID transaction.ID, payloadID payload.ID) *Attachment {
	return &Attachment{
		transactionID: transactionID,
		payloadID:     payloadID,
	}
}

// AttachmentFromBytes unmarshals an Attachment from a sequence of bytes - it either creates a new object or fills the
// optionally provided one with the parsed information.
func AttachmentFromBytes(bytes []byte) (result *Attachment, consumedBytes int, err error) {
	marshalUtil := marshalutil.New(bytes)
	result, err = ParseAttachment(marshalUtil)
	consumedBytes = marshalUtil.ReadOffset()

	return
}

// ParseAttachment is a wrapper for simplified unmarshaling of Attachments from a byte stream using the marshalUtil package.
func ParseAttachment(marshalUtil *marshalutil.MarshalUtil) (result *Attachment, err error) {
	result = &Attachment{}
	if result.transactionID, err = transaction.ParseID(marshalUtil); err != nil {
		err = fmt.Errorf("failed to parse transaction ID in attachment: %w", err)
		return
	}
	if result.payloadID, err = payload.ParseID(marshalUtil); err != nil {
		err = fmt.Errorf("failed to parse value payload ID in attachment: %w", err)
		return
	}

	return
}

// AttachmentFromObjectStorage gets called when we restore an Attachment from the storage - it parses the key bytes and
// returns the new object.
func AttachmentFromObjectStorage(key []byte, data []byte) (result objectstorage.StorableObject, err error) {
	result, _, err = AttachmentFromBytes(byteutils.ConcatBytes(key, data))
	if err != nil {
		err = fmt.Errorf("failed to parse attachment from object storage: %w", err)
	}

	return
}

// TransactionID returns the transaction id of this Attachment.
func (attachment *Attachment) TransactionID() transaction.ID {
	return attachment.transactionID
}

// PayloadID returns the payload id of this Attachment.
func (attachment *Attachment) PayloadID() payload.ID {
	return attachment.payloadID
}

// Bytes marshals the Attachment into a sequence of bytes.
func (attachment *Attachment) Bytes() []byte {
	return attachment.ObjectStorageKey()
}

// String returns a human readable version of the Attachment.
func (attachment *Attachment) String() string {
	return stringify.Struct("Attachment",
		stringify.StructField("transactionId", attachment.TransactionID()),
		stringify.StructField("payloadId", attachment.PayloadID()),
	)
}

// ObjectStorageKey returns the key that is used to store the object in the database.
func (attachment *Attachment) ObjectStorageKey() []byte {
	return byteutils.ConcatBytes(attachment.transactionID.Bytes(), attachment.payloadID.Bytes())
}

// ObjectStorageValue marshals the "content part" of an Attachment to a sequence of bytes. Since all of the information
// for this object are stored in its key, this method does nothing and is only required to conform with the interface.
func (attachment *Attachment) ObjectStorageValue() (data []byte) {
	return
}

// Update is disabled - updates are supposed to happen through the setters (if existing).
func (attachment *Attachment) Update(other objectstorage.StorableObject) {
	panic("update forbidden")
}

// Interface contract: make compiler warn if the interface is not implemented correctly.
var _ objectstorage.StorableObject = &Attachment{}

// AttachmentLength holds the length of a marshaled Attachment in bytes.
const AttachmentLength = transaction.IDLength + payload.IDLength

// region CachedAttachment /////////////////////////////////////////////////////////////////////////////////////////////

// CachedAttachment is a wrapper for the generic CachedObject returned by the objectstorage, that overrides the accessor
// methods, with a type-casted one.
type CachedAttachment struct {
	objectstorage.CachedObject
}

// Retain marks this CachedObject to still be in use by the program.
func (cachedAttachment *CachedAttachment) Retain() *CachedAttachment {
	return &CachedAttachment{cachedAttachment.CachedObject.Retain()}
}

// Unwrap is the type-casted equivalent of Get. It returns nil if the object does not exist.
func (cachedAttachment *CachedAttachment) Unwrap() *Attachment {
	untypedObject := cachedAttachment.Get()
	if untypedObject == nil {
		return nil
	}

	typedObject := untypedObject.(*Attachment)
	if typedObject == nil || typedObject.IsDeleted() {
		return nil
	}

	return typedObject
}

// Consume unwraps the CachedObject and passes a type-casted version to the consumer (if the object is not empty - it
// exists). It automatically releases the object when the consumer finishes.
func (cachedAttachment *CachedAttachment) Consume(consumer func(attachment *Attachment)) (consumed bool) {
	return cachedAttachment.CachedObject.Consume(func(object objectstorage.StorableObject) {
		consumer(object.(*Attachment))
	})
}

// CachedAttachments represents a collection of CachedAttachments.
type CachedAttachments []*CachedAttachment

// Consume iterates over the CachedObjects, unwraps them and passes a type-casted version to the consumer (if the object
// is not empty - it exists). It automatically releases the object when the consumer finishes. It returns true, if at
// least one object was consumed.
func (cachedAttachments CachedAttachments) Consume(consumer func(attachment *Attachment)) (consumed bool) {
	for _, cachedAttachment := range cachedAttachments {
		consumed = cachedAttachment.Consume(func(output *Attachment) {
			consumer(output)
		}) || consumed
	}

	return
}

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////
