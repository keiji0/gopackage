package internal

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/keiji0/sandbox/201904_p2p/net"
	"github.com/pkg/errors"
)

// Serialize は値をプロトコルに応じたデータに変換してwに書き込みます
func Serialize(w io.Writer, i interface{}) error {
	switch v := i.(type) {
	case net.VarUint:
		return serializeVarUint(w, v)
	case string:
		return serializeString(w, v)
	default:
		return binary.Write(w, net.ByteOrder, v)
	}
}

// Deserialize はrからローカル環境で利用できるデータに変換してvに読み込みます
func Deserialize(r io.Reader, i interface{}) error {
	switch p := i.(type) {
	case *net.VarUint:
		return deserializeVarUint(r, p)
	case *string:
		return deserializeString(r, p)
	default:
		if err := binary.Read(r, net.ByteOrder, p); err != nil {
			return errors.Wrapf(err, "読み込みに失敗しました: %T", p)
		}
		return nil
	}
}

// BulkSerialize は一括でシリアライズをします
func BulkSerialize(w io.Writer, args ...interface{}) error {
	for _, i := range args {
		if err := Serialize(w, i); err != nil {
			return err
		}
	}
	return nil
}

// BulkDeserialize は一括でデシリアライズします
func BulkDeserialize(r io.Reader, args ...interface{}) error {
	for _, i := range args {
		if err := Deserialize(r, i); err != nil {
			return err
		}
	}
	return nil
}

// serializeVarUint は可変長数値をシリアライズする
func serializeVarUint(w io.Writer, v net.VarUint) error {
	if v <= net.VarUint8Max {
		return binary.Write(w, net.ByteOrder, byte(v))
	}
	if v <= math.MaxUint16 {
		if err := binary.Write(w, net.ByteOrder, net.VarUint16Tag); err != nil {
			return err
		}
		return binary.Write(w, net.ByteOrder, uint16(v))
	}
	if v <= math.MaxUint32 {
		if err := binary.Write(w, net.ByteOrder, net.VarUint32Tag); err != nil {
			return err
		}
		return binary.Write(w, net.ByteOrder, uint32(v))
	}
	if err := binary.Write(w, net.ByteOrder, net.VarUint64Tag); err != nil {
		return err
	}
	return binary.Write(w, net.ByteOrder, uint64(v))
}

// deserializeVarUint は可変長数値をデシリアライズする
func deserializeVarUint(r io.Reader, p *net.VarUint) error {
	var length uint8
	if err := binary.Read(r, net.ByteOrder, &length); err != nil {
		return errors.Wrap(err, "VarUintの長さの読み込みに失敗しました")
	}

	switch length {
	case net.VarUint16Tag:
		var v uint16
		if err := binary.Read(r, net.ByteOrder, &v); err != nil {
			return errors.Wrap(err, "VarUint16の読み込みに失敗しました")
		}
		*p = net.VarUint(v)

	case net.VarUint32Tag:
		var v uint32
		if err := binary.Read(r, net.ByteOrder, &v); err != nil {
			return errors.Wrap(err, "VarUint32の読み込みに失敗しました")
		}
		*p = net.VarUint(v)

	case net.VarUint64Tag:
		var v uint64
		if err := binary.Read(r, net.ByteOrder, &v); err != nil {
			return errors.Wrap(err, "VarUint64の読み込みに失敗しました")
		}
		*p = net.VarUint(v)

	default:
		*p = net.VarUint(length)
	}

	return nil
}

// serializeString は可変長の文字列をシリアライズします
func serializeString(w io.Writer, v string) error {
	if err := serializeVarUint(w, net.VarUint(len(v))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(v)); err != nil {
		return errors.Wrap(err, "Stringの書き込みに失敗しました")
	}
	return nil
}

// serializeString は可変長の文字列をデシリアライズします
func deserializeString(r io.Reader, p *string) error {
	*p = ""
	var len net.VarUint
	if err := deserializeVarUint(r, &len); err != nil {
		return err
	}
	if net.MaxStringLength < len {
		return errors.Errorf("Stringが読み込める最大長を超えました: length=%d", len)
	}
	if len == 0 {
		// 0は正常、さっさと終わらす
		return nil
	}

	buf := make([]byte, len)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return errors.Wrapf(err, "StringのBodyの読み込みに失敗しました: length=%d", len)
	}
	if net.VarUint(n) != len {
		return errors.Wrapf(err, "Stringの長さと読み込んだ長さが一致しません: length=%d, read_length=%d", len, n)
	}

	*p = string(buf)
	return nil
}
