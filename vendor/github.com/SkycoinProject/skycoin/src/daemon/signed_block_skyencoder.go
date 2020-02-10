// Code generated by github.com/SkycoinProject/skyencoder. DO NOT EDIT.

package daemon

import (
	"errors"
	"math"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
	"github.com/SkycoinProject/skycoin/src/coin"
)

// encodeSizeSignedBlock computes the size of an encoded object of type SignedBlock
func encodeSizeSignedBlock(obj *coin.SignedBlock) uint64 {
	i0 := uint64(0)

	// obj.Block.Head.Version
	i0 += 4

	// obj.Block.Head.Time
	i0 += 8

	// obj.Block.Head.BkSeq
	i0 += 8

	// obj.Block.Head.Fee
	i0 += 8

	// obj.Block.Head.PrevHash
	i0 += 32

	// obj.Block.Head.BodyHash
	i0 += 32

	// obj.Block.Head.UxHash
	i0 += 32

	// obj.Block.Body.Transactions
	i0 += 4
	for _, x1 := range obj.Block.Body.Transactions {
		i1 := uint64(0)

		// x1.Length
		i1 += 4

		// x1.Type
		i1++

		// x1.InnerHash
		i1 += 32

		// x1.Sigs
		i1 += 4
		{
			i2 := uint64(0)

			// x2
			i2 += 65

			i1 += uint64(len(x1.Sigs)) * i2
		}

		// x1.In
		i1 += 4
		{
			i2 := uint64(0)

			// x2
			i2 += 32

			i1 += uint64(len(x1.In)) * i2
		}

		// x1.Out
		i1 += 4
		{
			i2 := uint64(0)

			// x2.Address.Version
			i2++

			// x2.Address.Key
			i2 += 20

			// x2.Coins
			i2 += 8

			// x2.Hours
			i2 += 8

			i1 += uint64(len(x1.Out)) * i2
		}

		i0 += i1
	}

	// obj.Sig
	i0 += 65

	return i0
}

// encodeSignedBlock encodes an object of type SignedBlock to a buffer allocated to the exact size
// required to encode the object.
func encodeSignedBlock(obj *coin.SignedBlock) ([]byte, error) {
	n := encodeSizeSignedBlock(obj)
	buf := make([]byte, n)

	if err := encodeSignedBlockToBuffer(buf, obj); err != nil {
		return nil, err
	}

	return buf, nil
}

// encodeSignedBlockToBuffer encodes an object of type SignedBlock to a []byte buffer.
// The buffer must be large enough to encode the object, otherwise an error is returned.
func encodeSignedBlockToBuffer(buf []byte, obj *coin.SignedBlock) error {
	if uint64(len(buf)) < encodeSizeSignedBlock(obj) {
		return encoder.ErrBufferUnderflow
	}

	e := &encoder.Encoder{
		Buffer: buf[:],
	}

	// obj.Block.Head.Version
	e.Uint32(obj.Block.Head.Version)

	// obj.Block.Head.Time
	e.Uint64(obj.Block.Head.Time)

	// obj.Block.Head.BkSeq
	e.Uint64(obj.Block.Head.BkSeq)

	// obj.Block.Head.Fee
	e.Uint64(obj.Block.Head.Fee)

	// obj.Block.Head.PrevHash
	e.CopyBytes(obj.Block.Head.PrevHash[:])

	// obj.Block.Head.BodyHash
	e.CopyBytes(obj.Block.Head.BodyHash[:])

	// obj.Block.Head.UxHash
	e.CopyBytes(obj.Block.Head.UxHash[:])

	// obj.Block.Body.Transactions maxlen check
	if len(obj.Block.Body.Transactions) > 65535 {
		return encoder.ErrMaxLenExceeded
	}

	// obj.Block.Body.Transactions length check
	if uint64(len(obj.Block.Body.Transactions)) > math.MaxUint32 {
		return errors.New("obj.Block.Body.Transactions length exceeds math.MaxUint32")
	}

	// obj.Block.Body.Transactions length
	e.Uint32(uint32(len(obj.Block.Body.Transactions)))

	// obj.Block.Body.Transactions
	for _, x := range obj.Block.Body.Transactions {

		// x.Length
		e.Uint32(x.Length)

		// x.Type
		e.Uint8(x.Type)

		// x.InnerHash
		e.CopyBytes(x.InnerHash[:])

		// x.Sigs maxlen check
		if len(x.Sigs) > 65535 {
			return encoder.ErrMaxLenExceeded
		}

		// x.Sigs length check
		if uint64(len(x.Sigs)) > math.MaxUint32 {
			return errors.New("x.Sigs length exceeds math.MaxUint32")
		}

		// x.Sigs length
		e.Uint32(uint32(len(x.Sigs)))

		// x.Sigs
		for _, x := range x.Sigs {

			// x
			e.CopyBytes(x[:])

		}

		// x.In maxlen check
		if len(x.In) > 65535 {
			return encoder.ErrMaxLenExceeded
		}

		// x.In length check
		if uint64(len(x.In)) > math.MaxUint32 {
			return errors.New("x.In length exceeds math.MaxUint32")
		}

		// x.In length
		e.Uint32(uint32(len(x.In)))

		// x.In
		for _, x := range x.In {

			// x
			e.CopyBytes(x[:])

		}

		// x.Out maxlen check
		if len(x.Out) > 65535 {
			return encoder.ErrMaxLenExceeded
		}

		// x.Out length check
		if uint64(len(x.Out)) > math.MaxUint32 {
			return errors.New("x.Out length exceeds math.MaxUint32")
		}

		// x.Out length
		e.Uint32(uint32(len(x.Out)))

		// x.Out
		for _, x := range x.Out {

			// x.Address.Version
			e.Uint8(x.Address.Version)

			// x.Address.Key
			e.CopyBytes(x.Address.Key[:])

			// x.Coins
			e.Uint64(x.Coins)

			// x.Hours
			e.Uint64(x.Hours)

		}

	}

	// obj.Sig
	e.CopyBytes(obj.Sig[:])

	return nil
}

// decodeSignedBlock decodes an object of type SignedBlock from a buffer.
// Returns the number of bytes used from the buffer to decode the object.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
func decodeSignedBlock(buf []byte, obj *coin.SignedBlock) (uint64, error) {
	d := &encoder.Decoder{
		Buffer: buf[:],
	}

	{
		// obj.Block.Head.Version
		i, err := d.Uint32()
		if err != nil {
			return 0, err
		}
		obj.Block.Head.Version = i
	}

	{
		// obj.Block.Head.Time
		i, err := d.Uint64()
		if err != nil {
			return 0, err
		}
		obj.Block.Head.Time = i
	}

	{
		// obj.Block.Head.BkSeq
		i, err := d.Uint64()
		if err != nil {
			return 0, err
		}
		obj.Block.Head.BkSeq = i
	}

	{
		// obj.Block.Head.Fee
		i, err := d.Uint64()
		if err != nil {
			return 0, err
		}
		obj.Block.Head.Fee = i
	}

	{
		// obj.Block.Head.PrevHash
		if len(d.Buffer) < len(obj.Block.Head.PrevHash) {
			return 0, encoder.ErrBufferUnderflow
		}
		copy(obj.Block.Head.PrevHash[:], d.Buffer[:len(obj.Block.Head.PrevHash)])
		d.Buffer = d.Buffer[len(obj.Block.Head.PrevHash):]
	}

	{
		// obj.Block.Head.BodyHash
		if len(d.Buffer) < len(obj.Block.Head.BodyHash) {
			return 0, encoder.ErrBufferUnderflow
		}
		copy(obj.Block.Head.BodyHash[:], d.Buffer[:len(obj.Block.Head.BodyHash)])
		d.Buffer = d.Buffer[len(obj.Block.Head.BodyHash):]
	}

	{
		// obj.Block.Head.UxHash
		if len(d.Buffer) < len(obj.Block.Head.UxHash) {
			return 0, encoder.ErrBufferUnderflow
		}
		copy(obj.Block.Head.UxHash[:], d.Buffer[:len(obj.Block.Head.UxHash)])
		d.Buffer = d.Buffer[len(obj.Block.Head.UxHash):]
	}

	{
		// obj.Block.Body.Transactions

		ul, err := d.Uint32()
		if err != nil {
			return 0, err
		}

		length := int(ul)
		if length < 0 || length > len(d.Buffer) {
			return 0, encoder.ErrBufferUnderflow
		}

		if length > 65535 {
			return 0, encoder.ErrMaxLenExceeded
		}

		if length != 0 {
			obj.Block.Body.Transactions = make([]coin.Transaction, length)

			for z3 := range obj.Block.Body.Transactions {
				{
					// obj.Block.Body.Transactions[z3].Length
					i, err := d.Uint32()
					if err != nil {
						return 0, err
					}
					obj.Block.Body.Transactions[z3].Length = i
				}

				{
					// obj.Block.Body.Transactions[z3].Type
					i, err := d.Uint8()
					if err != nil {
						return 0, err
					}
					obj.Block.Body.Transactions[z3].Type = i
				}

				{
					// obj.Block.Body.Transactions[z3].InnerHash
					if len(d.Buffer) < len(obj.Block.Body.Transactions[z3].InnerHash) {
						return 0, encoder.ErrBufferUnderflow
					}
					copy(obj.Block.Body.Transactions[z3].InnerHash[:], d.Buffer[:len(obj.Block.Body.Transactions[z3].InnerHash)])
					d.Buffer = d.Buffer[len(obj.Block.Body.Transactions[z3].InnerHash):]
				}

				{
					// obj.Block.Body.Transactions[z3].Sigs

					ul, err := d.Uint32()
					if err != nil {
						return 0, err
					}

					length := int(ul)
					if length < 0 || length > len(d.Buffer) {
						return 0, encoder.ErrBufferUnderflow
					}

					if length > 65535 {
						return 0, encoder.ErrMaxLenExceeded
					}

					if length != 0 {
						obj.Block.Body.Transactions[z3].Sigs = make([]cipher.Sig, length)

						for z5 := range obj.Block.Body.Transactions[z3].Sigs {
							{
								// obj.Block.Body.Transactions[z3].Sigs[z5]
								if len(d.Buffer) < len(obj.Block.Body.Transactions[z3].Sigs[z5]) {
									return 0, encoder.ErrBufferUnderflow
								}
								copy(obj.Block.Body.Transactions[z3].Sigs[z5][:], d.Buffer[:len(obj.Block.Body.Transactions[z3].Sigs[z5])])
								d.Buffer = d.Buffer[len(obj.Block.Body.Transactions[z3].Sigs[z5]):]
							}

						}
					}
				}

				{
					// obj.Block.Body.Transactions[z3].In

					ul, err := d.Uint32()
					if err != nil {
						return 0, err
					}

					length := int(ul)
					if length < 0 || length > len(d.Buffer) {
						return 0, encoder.ErrBufferUnderflow
					}

					if length > 65535 {
						return 0, encoder.ErrMaxLenExceeded
					}

					if length != 0 {
						obj.Block.Body.Transactions[z3].In = make([]cipher.SHA256, length)

						for z5 := range obj.Block.Body.Transactions[z3].In {
							{
								// obj.Block.Body.Transactions[z3].In[z5]
								if len(d.Buffer) < len(obj.Block.Body.Transactions[z3].In[z5]) {
									return 0, encoder.ErrBufferUnderflow
								}
								copy(obj.Block.Body.Transactions[z3].In[z5][:], d.Buffer[:len(obj.Block.Body.Transactions[z3].In[z5])])
								d.Buffer = d.Buffer[len(obj.Block.Body.Transactions[z3].In[z5]):]
							}

						}
					}
				}

				{
					// obj.Block.Body.Transactions[z3].Out

					ul, err := d.Uint32()
					if err != nil {
						return 0, err
					}

					length := int(ul)
					if length < 0 || length > len(d.Buffer) {
						return 0, encoder.ErrBufferUnderflow
					}

					if length > 65535 {
						return 0, encoder.ErrMaxLenExceeded
					}

					if length != 0 {
						obj.Block.Body.Transactions[z3].Out = make([]coin.TransactionOutput, length)

						for z5 := range obj.Block.Body.Transactions[z3].Out {
							{
								// obj.Block.Body.Transactions[z3].Out[z5].Address.Version
								i, err := d.Uint8()
								if err != nil {
									return 0, err
								}
								obj.Block.Body.Transactions[z3].Out[z5].Address.Version = i
							}

							{
								// obj.Block.Body.Transactions[z3].Out[z5].Address.Key
								if len(d.Buffer) < len(obj.Block.Body.Transactions[z3].Out[z5].Address.Key) {
									return 0, encoder.ErrBufferUnderflow
								}
								copy(obj.Block.Body.Transactions[z3].Out[z5].Address.Key[:], d.Buffer[:len(obj.Block.Body.Transactions[z3].Out[z5].Address.Key)])
								d.Buffer = d.Buffer[len(obj.Block.Body.Transactions[z3].Out[z5].Address.Key):]
							}

							{
								// obj.Block.Body.Transactions[z3].Out[z5].Coins
								i, err := d.Uint64()
								if err != nil {
									return 0, err
								}
								obj.Block.Body.Transactions[z3].Out[z5].Coins = i
							}

							{
								// obj.Block.Body.Transactions[z3].Out[z5].Hours
								i, err := d.Uint64()
								if err != nil {
									return 0, err
								}
								obj.Block.Body.Transactions[z3].Out[z5].Hours = i
							}

						}
					}
				}
			}
		}
	}

	{
		// obj.Sig
		if len(d.Buffer) < len(obj.Sig) {
			return 0, encoder.ErrBufferUnderflow
		}
		copy(obj.Sig[:], d.Buffer[:len(obj.Sig)])
		d.Buffer = d.Buffer[len(obj.Sig):]
	}

	return uint64(len(buf) - len(d.Buffer)), nil
}

// decodeSignedBlockExact decodes an object of type SignedBlock from a buffer.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
// If the buffer is longer than required to decode the object, returns encoder.ErrRemainingBytes.
func decodeSignedBlockExact(buf []byte, obj *coin.SignedBlock) error {
	if n, err := decodeSignedBlock(buf, obj); err != nil {
		return err
	} else if n != uint64(len(buf)) {
		return encoder.ErrRemainingBytes
	}

	return nil
}
