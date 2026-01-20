package vo

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
)

type DocumentNumber struct {
	document string
}

func NewDocumentNumber(document string) (documentNumber DocumentNumber, err error) {
	normalized := normalizeCPF(document)
	if !isValidCPF(normalized) {
		return documentNumber, domain.InvalidInputError("document_number", "invalid document number")
	}
	return DocumentNumber{document: normalized}, nil
}

func (e DocumentNumber) String() string {
	return e.document
}

var repeatedNumbersRE = regexp.MustCompile(`^0+$|^1+$|^2+$|^3+$|^4+$|^5+$|^6+$|^7+$|^8+$|^9+$`)

func generateCPFChecks(nums int64) string {
	cpf := make([]int16, 9)
	scpf := strings.Split(strconv.FormatInt(nums, 10), "")
	scpfI := len(scpf) - 1

	for i := 0; i < 9 && scpfI > -1; i++ {
		if n, err := strconv.ParseInt(scpf[scpfI], 10, 8); err == nil {
			cpf[i] = int16(n)
		}
		scpfI--
	}

	var v1, v2 int16
	for i := 0; i < 9; i++ {
		v1 = v1 + cpf[i]*int16(9-(i%10))
		v2 = v2 + cpf[i]*int16(9-((i+1)%10))
	}

	v1 = (v1 % 11) % 10
	v2 = ((v2 + v1*9) % 11) % 10
	return strconv.Itoa(int(v1)) + strconv.Itoa(int(v2))
}

func isValidCPF(text string) bool {
	cpfLen := len(text)

	if cpfLen != 11 && cpfLen != 14 {
		return false
	}

	if repeatedNumbersRE.MatchString(text) {
		return false
	}

	cpfNums, err := strconv.ParseInt(text[:9], 10, 64)
	if err != nil {
		return false
	}

	checkDigits := generateCPFChecks(cpfNums)
	return (text[9:] == checkDigits)
}

func normalizeCPF(cpf string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(cpf, "")
}
