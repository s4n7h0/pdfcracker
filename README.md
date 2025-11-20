# pdfcracker

A fast and customizable **PDF password brute-force tool** written in Go.  
Supports advanced cracking options such as:

- Custom character sets (`a-z`, `0-9`, `abcxyz`, `a-f0-9`, etc.)
- Minimum and maximum password length
- Prefix and suffix support for pattern-based cracking  
  (e.g., brute forcing *DDMM* + fixed year *2025*)
- Verbose mode for real-time password attempts
- Streaming, memory-efficient brute-force generation
- Built on the `pdfcpu` library for accurate PDF password validation

Example formats:

- `a-z`
- `0-9`
- `a-x,0-9`
- `abcxyz`


### Build

`go build -o pdf_cracker pdf_cracker.go`

### Run 

`./pdf_cracker -pdf=secret.pdf -charset=0-9 -min=4 -max=4`


### Examples 

Crack alphabetic passwords

`./pdf_cracker -pdf=secret.pdf -charset=a-z -min=1 -max=5`

Verbose mode

`./pdf_cracker -v -pdf=secret.pdf -charset=a-z`

Crack all 4-digit numeric passwords

`./pdf_cracker -pdf=secret.pdf -charset=0-9 -min=4 -max=4`


Crack DDMMYYYY date-based passwords with predictable year 2025

`./pdf_cracker -pdf=secret.pdf -charset=0-9 -min=4 -max=4 -suffix=2025`