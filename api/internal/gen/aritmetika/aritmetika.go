// Package aritmetika generates social-arithmetic TIU questions: profit and
// loss, discount, tax and simple interest (pattern TIU.numerik.aritmetika-sosial).
package aritmetika

import (
	"fmt"
	"math/rand"
	"strings"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/optkit"
)

const (
	PatternID   = "TIU.numerik.aritmetika-sosial"
	maxAttempts = 50
)

// Misconception ids attached to wrong options.
const (
	MisArahTerbalik = "arah-terbalik"
	MisHanyaNominal = "hanya-nominal"
	MisHanyaBunga   = "hanya-bunga"
	MisLupaTahun    = "lupa-kali-tahun"
	MisBungaMajemuk = "bunga-majemuk"
	MisPersenJumlah = "persen-dijumlah"
	MisLupaPajak    = "lupa-pajak"
	MisBasisSalah   = "dasar-persen-salah"
	MisSalahHitung  = "salah-hitung-dekat"
)

var (
	items    = []string{"sepatu", "tas", "jam tangan", "buku", "kipas angin", "sepeda", "jaket"}
	percents = []int{5, 10, 15, 20, 25, 30, 40, 50}
)

// Family implements gen.Family for social-arithmetic questions.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string    { return PatternID }
func (Family) Label() string { return "Aritmetika sosial" }

func (Family) Misconceptions() map[string]string {
	return map[string]string{
		MisArahTerbalik: "Membalik arah: menambah persen padahal dikurangi, atau sebaliknya",
		MisHanyaNominal: "Menjawab besar potongan atau laba saja, bukan harga akhirnya",
		MisHanyaBunga:   "Menjawab besar bunga saja, tanpa menambahkan modal",
		MisLupaTahun:    "Lupa mengalikan bunga dengan lama menabung",
		MisBungaMajemuk: "Memakai bunga majemuk padahal soal meminta bunga tunggal",
		MisPersenJumlah: "Menjumlahkan atau mengurangkan persen dari dua tahap berbeda",
		MisLupaPajak:    "Lupa menambahkan pajak pada harga setelah diskon",
		MisBasisSalah:   "Menghitung persen dari dasar yang salah",
		MisSalahHitung:  "Salah hitung kecil: hasil mendekati jawaban benar",
	}
}

func rp(v int) string {
	s := fmt.Sprint(v)
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	return "Rp " + strings.Join(parts, ".")
}

// Generate returns the question for seed and band: band 1 one-step profit,
// loss or discount; band 2 simple interest or percentage of profit or loss;
// band 3 two-step discount, discount plus tax, or a reversed discount.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return optkit.Retry("aritmetika", seed, band, maxAttempts, func(rng *rand.Rand) (gen.Draft, bool) {
		item := items[rng.Intn(len(items))]
		price := (rng.Intn(29) + 2) * 10000
		pc := func() int { return percents[rng.Intn(len(percents))] }
		pct := func(v, p int) int { return v * p / 100 }
		var stem, key string
		var steps []string
		var cands []optkit.Cand
		money := func(v int, mis string) {
			if v > 0 {
				cands = append(cands, optkit.Cand{Text: rp(v), Mis: mis})
			}
		}
		var keyVal int
		isMoney := true
		switch band {
		case gen.Band1:
			l := pc()
			switch rng.Intn(3) {
			case 0:
				keyVal = pct(price, 100+l)
				stem = fmt.Sprintf("Seorang pedagang membeli %s seharga %s, lalu menjualnya dengan laba %d%%. Harga jual %s tersebut adalah ...", item, rp(price), l, item)
				steps = []string{fmt.Sprintf("Laba = %d%% x %s = %s.", l, rp(price), rp(pct(price, l))), fmt.Sprintf("Harga jual = %s + %s = %s.", rp(price), rp(pct(price, l)), rp(keyVal))}
				money(pct(price, 100-l), MisArahTerbalik)
				money(pct(price, l), MisHanyaNominal)
			case 1:
				keyVal = pct(price, 100-l)
				stem = fmt.Sprintf("Seorang pedagang membeli %s seharga %s, lalu menjualnya dengan rugi %d%%. Harga jual %s tersebut adalah ...", item, rp(price), l, item)
				steps = []string{fmt.Sprintf("Rugi = %d%% x %s = %s.", l, rp(price), rp(pct(price, l))), fmt.Sprintf("Harga jual = %s - %s = %s.", rp(price), rp(pct(price, l)), rp(keyVal))}
				money(pct(price, 100+l), MisArahTerbalik)
				money(pct(price, l), MisHanyaNominal)
			default:
				keyVal = pct(price, 100-l)
				stem = fmt.Sprintf("Harga %s sebelum diskon adalah %s. Selama promo, %s itu didiskon %d%%. Harga yang harus dibayar adalah ...", item, rp(price), item, l)
				steps = []string{fmt.Sprintf("Diskon = %d%% x %s = %s.", l, rp(price), rp(pct(price, l))), fmt.Sprintf("Harga dibayar = %s - %s = %s.", rp(price), rp(pct(price, l)), rp(keyVal))}
				money(pct(price, 100+l), MisArahTerbalik)
				money(pct(price, l), MisHanyaNominal)
			}
		case gen.Band2:
			if rng.Intn(2) == 0 {
				r, y := rng.Intn(11)+2, rng.Intn(5)+1
				principal := (rng.Intn(20) + 2) * 100000
				interest := principal * r * y / 100
				keyVal = principal + interest
				stem = fmt.Sprintf("Dina menabung %s di bank dengan bunga tunggal %d%% per tahun. Setelah %d tahun, jumlah tabungan Dina adalah ...", rp(principal), r, y)
				steps = []string{fmt.Sprintf("Bunga = %s x %d%% x %d = %s.", rp(principal), r, y, rp(interest)), fmt.Sprintf("Tabungan = %s + %s = %s.", rp(principal), rp(interest), rp(keyVal))}
				money(interest, MisHanyaBunga)
				money(principal+principal*r/100, MisLupaTahun)
				if y > 1 {
					comp := principal
					for i := 0; i < y; i++ {
						if comp*(100+r)%100 != 0 {
							comp = 0
							break
						}
						comp = comp * (100 + r) / 100
					}
					money(comp, MisBungaMajemuk)
				}
			} else {
				isMoney = false
				l := []int{10, 20, 25, 30, 40, 50}[rng.Intn(6)]
				loss := rng.Intn(2) == 0
				sell := pct(price, 100+l)
				word := "laba"
				if loss {
					sell, word = pct(price, 100-l), "rugi"
				}
				keyVal = l
				key = fmt.Sprintf("%d%%", l)
				stem = fmt.Sprintf("Seorang pedagang membeli %s seharga %s dan menjualnya seharga %s. Persentase %s pedagang itu adalah ...", item, rp(price), rp(sell), word)
				diff := sell - price
				if loss {
					diff = price - sell
				}
				steps = []string{fmt.Sprintf("%s = %s - %s = %s.", strings.ToUpper(word[:1])+word[1:], rp(max(sell, price)), rp(min(sell, price)), rp(diff)), fmt.Sprintf("Persentase %s = %s : %s x 100%% = %d%%.", word, rp(diff), rp(price), l)}
				if diff*100%sell == 0 {
					cands = append(cands, optkit.Cand{Text: fmt.Sprintf("%d%%", diff*100/sell), Mis: MisBasisSalah})
				}
				cands = append(cands, optkit.Cand{Text: fmt.Sprintf("%d%%", 100-l), Mis: MisArahTerbalik})
			}
		default:
			d1, d2 := pc(), pc()
			switch rng.Intn(3) {
			case 0:
				first := pct(price, 100-d1)
				keyVal = pct(first, 100-d2)
				stem = fmt.Sprintf("Harga %s adalah %s. Toko memberi diskon %d%%, kemudian diskon tambahan %d%% dari harga setelah diskon pertama. Harga yang harus dibayar adalah ...", item, rp(price), d1, d2)
				steps = []string{fmt.Sprintf("Setelah diskon pertama = %d%% x %s = %s.", 100-d1, rp(price), rp(first)), fmt.Sprintf("Setelah diskon kedua = %d%% x %s = %s.", 100-d2, rp(first), rp(keyVal))}
				if d1+d2 < 100 {
					money(pct(price, 100-d1-d2), MisPersenJumlah)
				}
				money(pct(price, 100-d2), MisBasisSalah)
			case 1:
				tax := []int{10, 11}[rng.Intn(2)]
				after := pct(price, 100-d1)
				keyVal = pct(after, 100+tax)
				stem = fmt.Sprintf("Harga %s adalah %s. Setelah didiskon %d%%, pembeli dikenai pajak %d%% dari harga setelah diskon. Total yang dibayar adalah ...", item, rp(price), d1, tax)
				steps = []string{fmt.Sprintf("Harga setelah diskon = %d%% x %s = %s.", 100-d1, rp(price), rp(after)), fmt.Sprintf("Total = %d%% x %s = %s.", 100+tax, rp(after), rp(keyVal))}
				money(after, MisLupaPajak)
				money(pct(price, 100-d1+tax), MisPersenJumlah)
			default:
				after := pct(price, 100-d1)
				keyVal = price
				stem = fmt.Sprintf("Setelah mendapat diskon %d%%, harga %s menjadi %s. Harga %s sebelum diskon adalah ...", d1, item, rp(after), item)
				steps = []string{fmt.Sprintf("%d%% dari harga semula = %s.", 100-d1, rp(after)), fmt.Sprintf("Harga semula = %s : %d%% = %s.", rp(after), 100-d1, rp(price))}
				money(pct(after, 100+d1), MisBasisSalah)
				money(after+pct(after, d1), MisBasisSalah)
			}
		}
		if isMoney {
			key = rp(keyVal)
			for d := 1; d <= 6; d++ {
				money(keyVal+d*1000, MisSalahHitung)
				money(keyVal-d*1000, MisSalahHitung)
			}
		} else {
			for d := 1; d <= 6; d++ {
				cands = append(cands, optkit.Cand{Text: fmt.Sprintf("%d%%", keyVal+d), Mis: MisSalahHitung}, optkit.Cand{Text: fmt.Sprintf("%d%%", keyVal-d), Mis: MisSalahHitung})
			}
		}
		draft, err := optkit.Build(PatternID, stem, seed, key, cands, steps)
		return draft, err == nil
	})
}
