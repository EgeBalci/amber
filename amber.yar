import "pe"

rule Amber_Reflective_PE_Packer {
	meta:
		description = "Amber Packer - POC Reflective PE Packer"
		author = "Ege Balcı"
		reference = "https://github.com/egebalci/Amber"
		date = "11.28.2017"
		sample_filetype = "exe"
	strings:
		$s = {3C 41 6D 62 65 72 3A 32 37 61 30 31 64 34 37 37 32 30 33 38 61 33 66 38 33 35 35 32  39 30 38 65  30 34 37 30  36 30 34 65  37 37 33 66 38 61 66}
		$s1 = "<Amber:27a01d4772038a3f83552908e0470604e773f8af>" ascii wide
		$s2 = "Mingw-w64 runtime failure:" ascii wide 

	condition:
		($s and $s1 and $s2) and (pe.number_of_sections < 17 and pe.number_of_sections > 6)

}
