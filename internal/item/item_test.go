package item

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"testing"
)

var checks = []string{
	"A6cRHH+sfCuWGEZz2Lc5FWDbSfcQLmbaOV6SzgYP",
	"AwAAAADFtIC3/mrBkEsaj5NM0xGVIBFDCAAAAAAAMAYA",
	"AwAAAABLZ4A3RhkBkWMalJ8AEtSYWC1gJWYIAQAAAAAAyhgA",
	"AwAAAADuCYA3RhkBkWMalJ8AEtSYWC1gJmYIAQAAAAAAyhgA",
	"AwAAAACiM4A3rVMBk2tkjwhEwkRYO5cMpwkIAQAAAAAAyIAA",
	"AwAAAAA+uIC3syBDllvs4u4gGyP7LLHEEkssMQAA",
	"AwAAAADtqIC31oBBkWMEBcKAJnqQTAdOLWIIAQAAAAAAyIAA",
	"AwAAAACGEoC36JCAkTsKGoSgBASiIgsA",
	"AwAAAAAL94C3t9hAkysShLxMKmMLAA==",
	"AwAAAAByhIC3A/pBkWMGBYLB+IDMbhnFMWYIAQAAAAAAzoAA",
	"AwAAAACr6IC37xABkWsIqFPqeE0YjJJYUxUxhAAAAAAAAGQMAA==",
	"AwAAAAB1WoC3t9hAkysShLxMKkMLAA==",
	// "AwAAAACaOoA3VJMAkSsQUhYFGGMLAA==",
	// "AwAAAADkl4A3VJMAkSsQUhYFGEMLAA==",
	// "AwAAAAAZzYA3VJMAkSsQUhYFGAMLAA==",
	// "AwAAAACd1YA3VJMAkSsQUhYFGKMLAA==",
	"AwAAAAC754A3ElwAmCtYUlWPjAAAAA==",
	"AwAAAADBk4A3ElwAmCtYUlWxjgAAAA==",
	"AwAAAABM+oC33IBBkWMEA0LBZlmkGKfELb4IAQAAAAAAzoAA",
	"AwAAAABDBIC3syBDllvs4u4gDG3MtcQVE0tsRQAA",
	"AwAAAADT2YC3syBDllvs4u6gz2zcdcUVS0ysRAAA",
	"AwAAAAB0xYA3pNNBkXMIKNMJSiplJhFOghEFhAAAAAAAAGUMAA==",
	"AwAAAADh6oA3p+vCkHsiOIpkJRQgNB8QyRCcxAwhAAAAAABAGQMA",
	"AwAAAACr94A3wNBBkWMIJxRMBMrEIJyELGIIAQAAAAAAyoAA",
	"AwAAAAA2EYC3pGNBk2MaDghSkel4SJFSMnQIAQAAAAAAyhgA",
	"AwAAAABvRoA38QgBk0sap9jjQvFwZDFDCAAAAAAAQAYA",
	"AwAAAABRaYC3DUGBkGMGuXk40BtJSotjghAAAAAAAGAMAA==",
	"AwAAAAAHiIC3NmvBkEsaD4dOwwlNchFDCAAAAAAAUAYA",
	"AwAAAAC2soA31ECBkFOGteE+ViSvNkwIAQAAAAAA0oAA",
	"AwAAAAB0hoA3a1IBk3MeMkhEkisIJhZQhqOLGUIAAAAAAIAzIAA=",
	"AwAAAACxi4C3ZINBkXMEA0KBl9EoUowTAtQDhAAAAAAAAGNAAA==",
	"AwAAAABMd4C3y+qAEmCaB3LONQrZ6stiihAAAAAAAGAJAA==",
	"AwAAAABnloC3z5lBk1saN8zHFMEJCKlMzBACAAAAAACQAQA=",
	"AwAAAACIAIC3t9hAkysShLxMKgMLAA==",
	"AwAAAADjeIA3VJMAkSsQUhYFGIMLAA==",
	"AwAAAADEyIA37wgBk1sap5fBcYmAShxYzBACAAAAAACkAQA=",
	// potentially corrupt items
	//"AwAAAAAZDYC3/mrBkEsaj5NM0xGVIBFDCAAAAAAAMA==",
	//"AwAAAAA0qYA3RhkBkWMalJ8AEtSYWC1gJWYIAQAAAAAAyg==",
	//"AwAAAABknYA3RhkBkWMalJ8AEtSYWC1gJmYIAQAAAAAAyg==",
	//"AwAAAACrOoA3RhkBkWMalJ8AEtSYWC0gJmYIAQAAAAAAyg==",
	//"AwAAAADHT4A3rVMBk2tkjwhEwkRYO5cMpwkIAQAAAAAAyA==",
	//"AwAAAAAicYC3syBDllvs4u4gGyP7LLHEEkssMQ==",
	//"AwAAAAAXDIC31oBBkWMEBcKAJnqQTAdOLWIIAQAAAAAAyA==",
	//"AwAAAAD+8YC36JCAkTsKGoSgBASiIg==",
	//"AwAAAADDxoC3t9hAkysShLxMKmM=",
	//"AwAAAAAxn4C3A/pBkWMGBYLB+IDMbhnFMWYIAQAAAAAAzg==",
	//"AwAAAACAxIC37xABkWsIqFPqeE0YjJJYUxUxhAAAAAAAAGQ=",
	//"AwAAAADYYIC3t9hAkysShLxMKkM=",
	//"AwAAAAA/dIA34pIAkSsQUgAFBGM=",
	//"AwAAAAAk0oA34pIAkSsQUgAFBEM=",
	//"AwAAAAATnoA34pIAkSsQUgAFBAM=",
	//"AwAAAABmoIA34pIAkSsQUgAFBKM=",
	//"AwAAAACSJoA3ElwAmCtYUlWPjAA=",
	//"AwAAAAAk/YA3ElwAmCtYUlWxjgA=",
	//"AwAAAAB00IC33IBBkWMEA0LBZlmkGKfELb4IAQAAAAAAzg==",
	//"AwAAAACIjYC3syBDllvs4u4gDG3MtcQVE0tsRQ==",
	//"AwAAAACIjYC3syBDllvs4u4gDG3MtcQVE0tsRQ==",
	//"AwAAAACIBYC3syBDllvs4u6gz2zcdcUVS0ysRA==",
	//"AwAAAADCYoA3pNNBkXMIKNMJSiplJhFOghEFhAAAAAAAAGU=",
	//"AwAAAACsaIA3p+vCkHsiOIpkJRQgNB8QyRCcxAwhAAAAAABAGQ==",
	//"AwAAAACEUoA3wNBBkWMIJxRMBMrEIJyELGIIAQAAAAAAyg==",
	//"AwAAAABLQIC3pGNBk2MaDghSkel4SJFSMnQIAQAAAAAAyg==",
	//"AwAAAABzIoA38QgBk0sap9jjQvFwZDFDCAAAAAAAQA==",
	//"AwAAAABf1YC3DUGBkGMGuXk40BtJSotjghAAAAAAAGA=",
	//"AwAAAAAljoC3NmvBkEsaD4dOwwlNchFDCAAAAAAAUA==",
	//"AwAAAABQjoA31ECBkFOGteE+ViSvNkwIAQAAAAAA0g==",
	//"AwAAAAAqpYA3a1IBk3MeMkhEkisIJhZQhqOLGUIAAAAAAIAz",
	//"AwAAAABDxoC3ZINBkXMEA0KBl9EoUowTAtQDhAAAAAAAAGM=",
	//"AwAAAADQCIC3y+qAEmCaB3LONQrZ6stiihAAAAAAAGA=",
	//"AwAAAAB4aIC3z5lBk1saN8zHFMEJCKlMzBACAAAAAACQ",
	//"AwAAAADvLIC3t9hAkysShLxMKgM=",
	//"AwAAAAB9BoA34pIAkSsQUgAFBIM=",
	//"AwAAAAAfmYA37wgBk1sap5fBcYmAShxYzBACAAAAAACk",
}

func TestDecryptSerial(t *testing.T) {
	for _, check := range checks {
		bs, err := base64.StdEncoding.DecodeString(check)
		if err != nil {
			t.Fatal(err)
		}
		item, err := DecryptSerial(bs)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(item)
	}
}

func TestDeserialize(t *testing.T) {
	for _, check := range checks {
		bs, err := base64.StdEncoding.DecodeString(check)
		if err != nil {
			t.Fatal(err)
		}
		item, err := Deserialize(bs)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(item)
	}
}

func TestSerialize(t *testing.T) {
	for _, check := range checks {
		var result = check
		var history = make([]string, 10)
		var item Item
		var err error
		for i := 0; i < 10; i++ {
			bs, err := base64.StdEncoding.DecodeString(result)
			if err != nil {
				t.Fatal(err)
			}
			seed, err := GetSeedFromSerial(bs)
			if err != nil {
				t.Fatal(err)
			}
			item, err = Deserialize(bs)
			if err != nil {
				t.Fatal(err)
			}
			bs, err = Serialize(item, seed)
			if err != nil {
				t.Fatal(err)
			}
			result = base64.StdEncoding.EncodeToString(bs)
			history[i] = result
			i2, err := Deserialize(bs)
			if err != nil {
				t.Fatal(err)
			}
			if item.Level != i2.Level || item.Version != i2.Version {
				t.Fatal("component mismatch in re-serialized item")
			}
		}
		if result != check {
			log.Println(err)
			log.Println(check)
			log.Println(result)
			bs1, _ := base64.StdEncoding.DecodeString(check)
			bs2, _ := base64.StdEncoding.DecodeString(result)
			log.Println(hex.EncodeToString(bs1))
			log.Println(hex.EncodeToString(bs2))
			dec1, _ := DecryptSerial(bs1)
			dec2, _ := DecryptSerial(bs2)
			bs1, _ = base64.StdEncoding.DecodeString(check)
			bs2, _ = base64.StdEncoding.DecodeString(result)
			log.Println(hex.EncodeToString(dec1))
			log.Println(hex.EncodeToString(dec2))
			i1, _ := Deserialize(bs1)
			i2, _ := Deserialize(bs2)
			log.Println(i1.Version)
			log.Println(i2.Version)
			t.Fatal("invalid serial")
		}
	}
}

func TestAddPart(t *testing.T) {
	debug = true
	code := "AwAAAADuCYA3RhkBkWMalJ8AEtSYWC1gJmYIAQAAAAAAyhgA"
	part := "/Game/Gear/Weapons/Pistols/Vladof/_Shared/_Design/Parts/Barrels/Barrel_01/Part_PS_VLA_Barrel_01_B.Part_PS_VLA_Barrel_01_B"
	bs, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		t.Fatal(err)
	}
	seed, err := GetSeedFromSerial(bs)
	if err != nil {
		panic(err)
	}
	item, err := Deserialize(bs)
	if err != nil {
		t.Fatal(err)
	}
	item.Parts = append(item.Parts, part)
	bs, err = Serialize(item, seed)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(base64.StdEncoding.EncodeToString(bs))
	i2, err := Deserialize(bs)
	if err != nil {
		t.Fatal(err)
	}
	if len(i2.Parts) != 13 {
		t.Fatalf("invalid part length %v", len(i2.Parts))
	}

}

func TestAddAnointment(t *testing.T) {
	debug = true
	code := "AwAAAADuCYA3RhkBkWMalJ8AEtSYWC1gJmYIAQAAAAAAyhgA"
	anointment := "/Game/Gear/Weapons/_Shared/_Design/EndGameParts/Character/Operative/CloneSwapDamage/GPart_CloneSwap_WeaponDamage.GPart_CloneSwap_WeaponDamage"
	bs, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		t.Fatal(err)
	}
	seed, err := GetSeedFromSerial(bs)
	if err != nil {
		t.Fatal(err)
	}
	item, err := Deserialize(bs)
	if err != nil {
		t.Fatal(err)
	}
	item.Generics = append(item.Generics, anointment)
	bs, err = Serialize(item, seed)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(base64.StdEncoding.EncodeToString(bs))
	i2, err := Deserialize(bs)
	if err != nil {
		t.Fatal(err)
	}
	if len(i2.Generics) != 2 {
		t.Fatal("invalid anointment length")
	}

}
