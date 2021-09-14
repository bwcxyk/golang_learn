module ocr

go 1.17

require IdCard v0.0.0
require config v0.0.0

replace IdCard v0.0.0 => ./ocr/IdCard
replace config v0.0.0 => ./configs