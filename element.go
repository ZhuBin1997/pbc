package pbc

/*
#cgo LDFLAGS: /usr/local/lib/libpbc.a -lgmp
#include <pbc/pbc.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
)

var ErrUnknownField = errors.New("unchecked element initialized in unknown field")

type Element interface {
	NewFieldElement() Element

	Set0() Element
	Set1() Element
	SetInt32(int32) Element
	SetBig(*big.Int) Element
	Set(Element) Element

	SetFromHash([]byte) Element
	SetBytes([]byte) Element
	SetXBytes([]byte) Element
	SetCompressedBytes([]byte) Element

	SetString(s string, base int) (Element, bool)

	Format(fmt.State, rune)
	Scan(fmt.ScanState, rune) error

	BigInt() *big.Int
	String() string

	BytesLen() int
	Bytes() []byte
	XBytesLen() int
	XBytes() []byte
	CompressedBytesLen() int
	CompressedBytes() []byte

	Len() int
	Item(int) Element
	X() *big.Int
	Y() *big.Int

	Is0() bool
	Is1() bool
	IsSquare() bool
	Sign() int

	Cmp(x Element) int

	Add(x, y Element) Element
	Sub(x, y Element) Element
	Mul(x, y Element) Element
	MulBig(x Element, i *big.Int) Element
	MulInt32(x Element, i int32) Element
	MulZn(x, y Element) Element
	Div(x, y Element) Element
	Double(x Element) Element
	Halve(x Element) Element
	Square(x Element) Element
	Neg(x Element) Element
	Invert(x Element) Element

	PowBig(x Element, i *big.Int) Element
	PowZn(x, i Element) Element
	Pow2Big(x Element, i *big.Int, y Element, j *big.Int) Element
	Pow2Zn(x, i, y, j Element) Element
	Pow3Big(x Element, i *big.Int, y Element, j *big.Int, z Element, k *big.Int) Element
	Pow3Zn(x, i, y, j, z, k Element) Element

	PreparePower() Power

	BruteForceDL(g, h Element) Element
	PollardRhoDL(g, h Element) Element

	Rand() Element

	impl() *elementImpl
}

type elementImpl struct {
	pairing *pairingImpl
	data    *C.struct_element_s
}

type checkedElement struct {
	elementImpl
	fieldPtr  *C.struct_field_s
	isInteger bool
}

type Power interface {
	PowBig(i *big.Int) Element
	PowZn(i Element) Element
}

type powerImpl struct {
	target *elementImpl
	data   *C.struct_element_pp_s
}

type checkedPower struct {
	powerImpl
}

func (pairing *pairingImpl) NewG1() Element                 { return makeChecked(pairing, G1, pairing.data.G1) }
func (pairing *pairingImpl) NewG2() Element                 { return makeChecked(pairing, G2, pairing.data.G2) }
func (pairing *pairingImpl) NewGT() Element                 { return makeChecked(pairing, GT, &pairing.data.GT[0]) }
func (pairing *pairingImpl) NewZr() Element                 { return makeChecked(pairing, Zr, &pairing.data.Zr[0]) }
func (pairing *pairingImpl) NewElement(field Field) Element { return makeElement(pairing, field) }

func clearElement(element *elementImpl) {
	println("clearelement")
	C.element_clear(element.data)
}

func initElement(element *elementImpl, pairing *pairingImpl, initialize bool, field Field) {
	element.data = &C.struct_element_s{}
	element.pairing = pairing
	if initialize {
		switch field {
		case G1:
			C.element_init_G1(element.data, pairing.data)
		case G2:
			C.element_init_G2(element.data, pairing.data)
		case GT:
			C.element_init_GT(element.data, pairing.data)
		case Zr:
			C.element_init_Zr(element.data, pairing.data)
		default:
			panic(ErrUnknownField)
		}
	}
	runtime.SetFinalizer(element, clearElement)
}

func makeElement(pairing *pairingImpl, field Field) *elementImpl {
	element := &elementImpl{}
	initElement(element, pairing, true, field)
	return element
}

func makeChecked(pairing *pairingImpl, field Field, fieldPtr *C.struct_field_s) *checkedElement {
	element := &checkedElement{
		fieldPtr:  fieldPtr,
		isInteger: field == Zr,
	}
	initElement(&element.elementImpl, pairing, true, field)
	return element
}

func clearPower(power *powerImpl) {
	println("clearpower")
	C.element_pp_clear(power.data)
}

func initPower(power *powerImpl, target *elementImpl) {
	power.target = target
	power.data = &C.struct_element_pp_s{}
	C.element_pp_init(power.data, target.data)
	runtime.SetFinalizer(power, clearPower)
}
