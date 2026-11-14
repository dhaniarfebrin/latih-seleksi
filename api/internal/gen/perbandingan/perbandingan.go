// Package perbandingan generates ratio and proportion TIU questions
// (pattern TIU.numerik.perbandingan).
package perbandingan

import (
	"fmt"
	"math/rand"
	"strings"

	"latihseleksi/api/internal/gen"
)

const (
	PatternID     = "TIU.numerik.perbandingan"
	maxAttempts   = 50
	shuffleSalt   = 1 << 20
	optionCount   = 5
	maxValue      = 1000 // upper bound of any answer
	maxChainSum   = 60   // upper bound of the summed weights of a chained ratio
	maxShiftRatio = 20   // upper bound of each number of the ratio after the change
)

// Misconception labels attached to wrong options.
const (
	MisBagiTanpaJumlahRasio = "bagi-tanpa-jumlah-rasio"
	MisBagianTertukar       = "bagian-tertukar"
	MisSenilaiBukanBerbalik = "senilai-bukan-berbalik"
	MisSelisihBukanRasio    = "selisih-bukan-rasio"
	MisTengahTakDisamakan   = "tengah-tak-disamakan"
	MisSalahHitungDekat     = "salah-hitung-dekat"
)

// kind selects the question variant.
type kind int

const (
	kindSplit   kind = iota // total shared in a ratio of two (band 1)
	kindThree               // three-part ratio with a known difference or part (band 2)
	kindInverse             // direct or inverse proportion (band 2)
	kindShift               // ratio after the same change on both sides (band 3)
	kindChain               // A:B and B:C combined into A:B:C (band 3)
)

// given tells what params.gv is for kindSplit and kindThree.
type given int

const (
	givenTotal given = iota // the sum of all parts
	givenDiff               // the difference of parts gi and gj
	givenPart               // part gi
)

const askTotal = -1

// params are the numbers behind one question.
//
//	kindSplit, kindThree: w holds the ratio weights, gv the given value (see
//	  given), ask the part asked for or askTotal.
//	kindChain: w holds p, q, r, s of A:B = p:q and B:C = r:s; gv is the sum of
//	  A, B and C; ask is the index (0..2) of the part asked for.
//	kindInverse: the stem reads x1, y1, x2 in that order and asks y2.
//	kindShift: w holds a, b (ratio now) and cd holds c, d (ratio after the
//	  change); n is the change, sign is +1 for added or later and -1 for
//	  removed or earlier; ask is 0 or 1.
type params struct {
	kind       kind
	context    int
	w          []int
	given      given
	gi, gj     int
	gv         int
	ask        int
	x1, y1, x2 int
	cd         [2]int
	n, sign    int
}

type ratioContext struct {
	measure string    // what is compared, one word
	items   [3]string // the things compared
	unit    string    // suffix of every value, including the leading space
}

var ratioContexts = []ratioContext{
	{"banyak", [3]string{"kelereng merah", "kelereng biru", "kelereng hijau"}, " butir"},
	{"banyak", [3]string{"siswa kelas A", "siswa kelas B", "siswa kelas C"}, " orang"},
	{"berat", [3]string{"beras", "jagung", "kacang tanah"}, " kg"},
	{"panjang", [3]string{"tali merah", "tali biru", "tali kuning"}, " meter"},
	{"banyak", [3]string{"buku cerita", "buku pelajaran", "majalah"}, " buah"},
	{"banyak", [3]string{"pegawai tetap", "pegawai kontrak", "pegawai magang"}, " orang"},
	{"banyak", [3]string{"kursi plastik", "kursi kayu", "kursi besi"}, " buah"},
	{"volume", [3]string{"air", "sirup", "susu"}, " liter"},
	{"banyak", [3]string{"apel", "jeruk", "mangga"}, " buah"},
	{"berat", [3]string{"tepung terigu", "gula pasir", "mentega"}, " kg"},
}

type inverseContext struct {
	inverse bool   // true: inverse proportion, false: direct
	format  string // stem with %d for x1, y1, x2
	unit    string // suffix of the answer, including the leading space
	scale   int    // factor applied to the displayed x values
}

var inverseContexts = []inverseContext{
	{true, "%d mesin pengepak dapat menyelesaikan satu pesanan dalam %d hari. Jika dipakai %d mesin yang sama, pesanan itu selesai dalam ... hari.", " hari", 1},
	{true, "Persediaan beras sebuah asrama cukup untuk %d orang selama %d hari. Jika penghuninya menjadi %d orang, persediaan itu cukup untuk ... hari.", " hari", 1},
	{true, "Dengan kecepatan %d km/jam, sebuah bus menempuh jarak antarkota dalam %d jam. Jika kecepatannya menjadi %d km/jam, waktu tempuhnya ... jam.", " jam", 5},
	{true, "%d pompa air dapat mengosongkan sebuah kolam dalam %d jam. Jika dipakai %d pompa yang sama, kolam itu kosong dalam ... jam.", " jam", 1},
	{false, "Untuk membuat %d kue dibutuhkan %d butir telur. Untuk membuat %d kue dibutuhkan ... butir telur.", " butir", 1},
	{false, "Untuk menempuh %d km, sebuah mobil menghabiskan %d liter bensin. Untuk menempuh %d km, bensin yang dihabiskan adalah ... liter.", " liter", 10},
	{false, "Dalam %d menit, sebuah mesin cetak menghasilkan %d lembar. Dalam %d menit, mesin itu menghasilkan ... lembar.", " lembar", 1},
	{false, "Untuk mengecat dinding seluas %d meter persegi diperlukan %d kaleng cat. Untuk dinding seluas %d meter persegi diperlukan ... kaleng cat.", " kaleng", 5},
}

type shiftContext struct {
	age    bool
	n1, n2 string
	unit   string // suffix of every value, including the leading space
	max    int    // upper bound of any value shown or asked
	minOld int    // lowest value of the first side now; the first side is also the larger one when age is true
}

var shiftContexts = []shiftContext{
	{true, "ayah", "anak", " tahun", 70, 28},
	{true, "ibu", "putrinya", " tahun", 70, 28},
	{true, "kakak", "adik", " tahun", 50, 8},
	{true, "paman", "keponakannya", " tahun", 70, 28},
	{false, "buku di rak kiri", "buku di rak kanan", " buku", 150, 0},
	{false, "siswa kelas A", "siswa kelas B", " siswa", 100, 0},
	{false, "kelereng merah", "kelereng biru", " butir", 150, 0},
}

// Family implements gen.Family for ratio questions.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string { return PatternID }

func (Family) Label() string { return "Perbandingan" }

// Misconceptions returns the labels of the ids in the Mis* constants.
func (Family) Misconceptions() map[string]string {
	return map[string]string{
		MisBagiTanpaJumlahRasio: "Menentukan satu bagian dengan pembagi yang salah (satu bobot saja atau tanpa jumlah rasio), bukan dengan jumlah rasio atau selisih bobot yang tepat",
		MisBagianTertukar:       "Mengambil nilai bagian yang lain, bukan bagian yang ditanyakan",
		MisSenilaiBukanBerbalik: "Memakai perbandingan senilai padahal soal berbalik nilai, atau sebaliknya",
		MisSelisihBukanRasio:    "Salah menerapkan penambahan atau pengurangan yang sama: menjawab nilai sesudah perubahan atau membalik arahnya",
		MisTengahTakDisamakan:   "Menggabungkan dua perbandingan tanpa menyamakan dulu bagian yang dipakai bersama",
		MisSalahHitungDekat:     "Salah hitung kecil: hasil mendekati jawaban benar",
	}
}

// Generate returns the question for seed and band. Parameters whose result is
// not a positive whole number are rejected; attempt k uses a seed derived from
// (seed, k), and after maxAttempts it returns an error.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return generate(seed, band, maxAttempts, valid)
}

func valid(p params) bool {
	res, ok := solve(p)
	return ok && res.key <= maxValue && res.max <= maxValue
}

func generate(seed int64, band gen.Band, attempts int, accept func(params) bool) (gen.Draft, error) {
	if band < gen.Band1 || band > gen.Band3 {
		return gen.Draft{}, fmt.Errorf("perbandingan: band %d tidak dikenal", band)
	}
	for k := 0; k < attempts; k++ {
		rng := rand.New(rand.NewSource(derive(seed, k)))
		p := pickParams(rng, band)
		if accept(p) {
			return build(p, seed), nil
		}
	}
	return gen.Draft{}, fmt.Errorf("perbandingan: tidak ada parameter valid untuk seed %d band %d dalam %d percobaan", seed, band, attempts)
}

// derive mixes (seed, k) into a new seed (splitmix64).
func derive(seed int64, k int) int64 {
	z := uint64(seed) + uint64(k+1)*0x9E3779B97F4A7C15
	z = (z ^ z>>30) * 0xBF58476D1CE4E5B9
	z = (z ^ z>>27) * 0x94D049BB133111EB
	return int64(z ^ z>>31)
}

func pickParams(rng *rand.Rand, band gen.Band) params {
	switch band {
	case gen.Band1:
		return pickSplit(rng)
	case gen.Band2:
		if rng.Intn(2) == 0 {
			return pickThree(rng)
		}
		return pickInverse(rng)
	default:
		if rng.Intn(2) == 0 {
			return pickShift(rng)
		}
		return pickChain(rng)
	}
}

func pickSplit(rng *rand.Rand) params {
	ctx := rng.Intn(len(ratioContexts))
	a, b := rng.Intn(9)+1, rng.Intn(9)+1
	m := rng.Intn(19) + 2
	ask := rng.Intn(2)
	if a == b || gcd(a, b) != 1 {
		return params{kind: kindSplit, context: ctx} // rejected by solve
	}
	return params{kind: kindSplit, context: ctx, w: []int{a, b}, given: givenTotal, gv: m * (a + b), ask: ask}
}

func pickThree(rng *rand.Rand) params {
	ctx := rng.Intn(len(ratioContexts))
	w := []int{rng.Intn(9) + 1, rng.Intn(9) + 1, rng.Intn(9) + 1}
	m := rng.Intn(12) + 2
	ask := rng.Intn(4) - 1 // askTotal or a part
	i := rng.Intn(3)
	j := (i + 1 + rng.Intn(2)) % 3
	if gcd(gcd(w[0], w[1]), w[2]) != 1 {
		return params{kind: kindThree, context: ctx} // rejected by solve
	}
	if rng.Intn(2) == 0 {
		return params{kind: kindThree, context: ctx, w: w, given: givenDiff, gi: i, gj: j, gv: m * abs(w[i]-w[j]), ask: ask}
	}
	return params{kind: kindThree, context: ctx, w: w, given: givenPart, gi: i, gv: m * w[i], ask: ask}
}

func pickInverse(rng *rand.Rand) params {
	ctx := rng.Intn(len(inverseContexts))
	c := inverseContexts[ctx]
	x1, x2 := rng.Intn(19)+2, rng.Intn(19)+2
	if x1 == x2 {
		return params{kind: kindInverse, context: ctx} // rejected by solve
	}
	// inverse: x1*y1 = x2*y2 = K, a multiple of lcm(x1,x2), so y1 = t*x2/g and
	// y2 = t*x1/g; direct: y1/x1 = y2/x2 = t/g, so y1 = t*x1/g and y2 = t*x2/g.
	g := gcd(x1, x2)
	tMax := 60 / (max(x1, x2) / g)
	if tMax < 1 {
		return params{kind: kindInverse, context: ctx}
	}
	t := rng.Intn(tMax) + 1
	y1 := t * x1 / g
	if c.inverse {
		y1 = t * x2 / g
	} else if t == g {
		return params{kind: kindInverse, context: ctx} // y = x makes the question trivial
	}
	return params{kind: kindInverse, context: ctx, x1: x1 * c.scale, y1: y1, x2: x2 * c.scale}
}

func pickShift(rng *rand.Rand) params {
	ctx := rng.Intn(len(shiftContexts))
	c := shiftContexts[ctx]
	a, b := rng.Intn(9)+1, rng.Intn(9)+1
	if a == b || gcd(a, b) != 1 {
		return params{kind: kindShift, context: ctx} // rejected by solve
	}
	if c.age && a < b {
		a, b = b, a // the first side is the older one
	}
	// k keeps both values now inside [minOld, two thirds of max]
	lo, hi := max((c.minOld+a-1)/a, 2), c.max*2/3/a
	if hi < lo {
		return params{kind: kindShift, context: ctx}
	}
	k := lo + rng.Intn(hi-lo+1)
	sign := 1
	if rng.Intn(2) == 0 {
		sign = -1
	}
	nMax := 12
	if sign < 0 {
		nMax = min(nMax, min(a, b)*k-1)
	}
	if nMax < 1 {
		return params{kind: kindShift, context: ctx}
	}
	n := rng.Intn(nMax) + 1
	x, y := a*k+sign*n, b*k+sign*n
	g := gcd(x, y)
	cd := [2]int{x / g, y / g}
	if cd[0] > maxShiftRatio || cd[1] > maxShiftRatio {
		return params{kind: kindShift, context: ctx}
	}
	return params{kind: kindShift, context: ctx, w: []int{a, b}, cd: cd, n: n, sign: sign, ask: rng.Intn(2)}
}

func pickChain(rng *rand.Rand) params {
	ctx := rng.Intn(len(ratioContexts))
	p, q, r, s := rng.Intn(8)+1, rng.Intn(8)+1, rng.Intn(8)+1, rng.Intn(8)+1
	m := rng.Intn(9) + 2
	ask := rng.Intn(3)
	if p == q || r == s || q == r || gcd(p, q) != 1 || gcd(r, s) != 1 {
		return params{kind: kindChain, context: ctx} // rejected by solve
	}
	return params{kind: kindChain, context: ctx, w: []int{p, q, r, s}, gv: m * sumInts(chainWeights(p, q, r, s)), ask: ask}
}

// chainWeights returns A:B:C in lowest terms for A:B = p:q and B:C = r:s.
func chainWeights(p, q, r, s int) []int {
	w := []int{p * r, q * r, q * s}
	g := gcd(gcd(w[0], w[1]), w[2])
	return []int{w[0] / g, w[1] / g, w[2] / g}
}

type result struct {
	key   int
	base  int // value of one ratio part (kindSplit, kindThree, kindChain) or k (kindShift)
	max   int // largest quantity shown in the stem or the answer
	steps []string
}

// solve returns the key and solution steps; ok is false when params are
// malformed or the key is not a positive whole number.
func solve(p params) (result, bool) {
	switch p.kind {
	case kindSplit, kindThree:
		return solveRatio(p)
	case kindInverse:
		return solveInverse(p)
	case kindShift:
		return solveShift(p)
	case kindChain:
		return solveChain(p)
	}
	return result{}, false
}

func solveRatio(p params) (result, bool) {
	ctx := ratioContexts[p.context%len(ratioContexts)]
	if len(p.w) < 2 || len(p.w) > 3 || p.gv <= 0 || p.ask < askTotal || p.ask >= len(p.w) {
		return result{}, false
	}
	for _, x := range p.w {
		if x <= 0 {
			return result{}, false
		}
	}
	sum := sumInts(p.w)
	var den int
	var steps []string
	switch p.given {
	case givenTotal:
		if p.ask == askTotal {
			return result{}, false
		}
		den = sum
		steps = append(steps, fmt.Sprintf("Jumlah perbandingan = %s = %d.", joinInts(p.w, " + "), sum))
	case givenDiff:
		if p.gi < 0 || p.gi >= len(p.w) || p.gj < 0 || p.gj >= len(p.w) || p.gi == p.gj {
			return result{}, false
		}
		den = abs(p.w[p.gi] - p.w[p.gj])
		if den == 0 {
			return result{}, false
		}
		steps = append(steps, fmt.Sprintf("Selisih perbandingan %s dan %s = |%d - %d| = %d.", ctx.items[p.gi], ctx.items[p.gj], p.w[p.gi], p.w[p.gj], den))
	case givenPart:
		if p.gi < 0 || p.gi >= len(p.w) || p.ask == p.gi {
			return result{}, false
		}
		den = p.w[p.gi]
		steps = append(steps, fmt.Sprintf("%s mewakili %d bagian perbandingan.", capitalize(ctx.items[p.gi]), den))
	default:
		return result{}, false
	}
	if p.gv%den != 0 {
		return result{}, false
	}
	u := p.gv / den
	steps = append(steps, fmt.Sprintf("Satu bagian = %d : %d = %d.", p.gv, den, u))
	var key int
	if p.ask == askTotal {
		key = u * sum
		steps = append(steps, fmt.Sprintf("Jumlah ketiganya = %d x %d = %d.", sum, u, key))
	} else {
		key = u * p.w[p.ask]
		steps = append(steps, fmt.Sprintf("%s = %d x %d = %d.", capitalize(ctx.items[p.ask]), p.w[p.ask], u, key))
	}
	return result{key: key, base: u, max: max(key, p.gv), steps: steps}, true
}

func solveInverse(p params) (result, bool) {
	c := inverseContexts[p.context%len(inverseContexts)]
	if p.x1 <= 0 || p.y1 <= 0 || p.x2 <= 0 || p.x1 == p.x2 {
		return result{}, false
	}
	var num, den int
	var steps []string
	if c.inverse {
		num, den = p.x1*p.y1, p.x2
		steps = append(steps,
			"Berbalik nilai: x1 x y1 = x2 x y2, sehingga y2 = (x1 x y1) : x2.",
			fmt.Sprintf("y2 = (%d x %d) : %d = %d : %d.", p.x1, p.y1, p.x2, num, den))
	} else {
		num, den = p.y1*p.x2, p.x1
		steps = append(steps,
			"Senilai: y1 : x1 = y2 : x2, sehingga y2 = (y1 x x2) : x1.",
			fmt.Sprintf("y2 = (%d x %d) : %d = %d : %d.", p.y1, p.x2, p.x1, num, den))
	}
	if num%den != 0 {
		return result{}, false
	}
	key := num / den
	steps = append(steps, fmt.Sprintf("y2 = %d.", key))
	return result{key: key, max: max(key, p.y1), steps: steps}, key > 0
}

func solveShift(p params) (result, bool) {
	ctx := shiftContexts[p.context%len(shiftContexts)]
	if len(p.w) != 2 || p.w[0] <= 0 || p.w[1] <= 0 || p.cd[0] <= 0 || p.cd[1] <= 0 ||
		p.n <= 0 || (p.sign != 1 && p.sign != -1) || (p.ask != 0 && p.ask != 1) {
		return result{}, false
	}
	a, b, c, d := p.w[0], p.w[1], p.cd[0], p.cd[1]
	// d(ak + sn) = c(bk + sn) gives k(ad - bc) = sn(c - d)
	den := a*d - b*c
	num := p.sign * p.n * (c - d)
	if den == 0 || num%den != 0 {
		return result{}, false
	}
	k := num / den
	x, y := a*k+p.sign*p.n, b*k+p.sign*p.n
	if k <= 0 || x <= 0 || y <= 0 || a*k < ctx.minOld {
		return result{}, false
	}
	op := "+"
	if p.sign < 0 {
		op = "-"
	}
	kDen, kNum := den, num
	if kDen < 0 {
		kDen, kNum = -kDen, -kNum
	}
	solved := fmt.Sprintf("jadi k x %d = %d dan k = %d.", kDen, kNum, k)
	if den < 0 {
		solved = fmt.Sprintf("jadi k x %d = %d; kedua ruas dikalikan -1 menjadi k x %d = %d, sehingga k = %d.", den, num, kDen, kNum, k)
	}
	steps := []string{
		fmt.Sprintf("Misalkan nilai sekarang = %dk dan %dk.", a, b),
		fmt.Sprintf("Setelah perubahan: (%dk %s %d) : (%dk %s %d) = %d : %d, sehingga %d(%dk %s %d) = %d(%dk %s %d).", a, op, p.n, b, op, p.n, c, d, d, a, op, p.n, c, b, op, p.n),
		fmt.Sprintf("k x (%d x %d - %d x %d) = %s%d x (%d - %d), %s", a, d, b, c, signPrefix(p.sign), p.n, c, d, solved),
	}
	key := p.w[p.ask] * k
	steps = append(steps, fmt.Sprintf("%s = %d x %d = %d.", capitalize(shiftItem(ctx, p.ask)), p.w[p.ask], k, key))
	return result{key: key, base: k, max: max(max(a*k, b*k), max(x, y)), steps: steps}, key <= ctx.max && max(max(a*k, b*k), max(x, y)) <= ctx.max
}

func solveChain(p params) (result, bool) {
	ctx := ratioContexts[p.context%len(ratioContexts)]
	if len(p.w) != 4 || p.gv <= 0 || p.ask < 0 || p.ask > 2 {
		return result{}, false
	}
	for _, x := range p.w {
		if x <= 0 {
			return result{}, false
		}
	}
	pp, q, r, s := p.w[0], p.w[1], p.w[2], p.w[3]
	w := chainWeights(pp, q, r, s)
	sum := sumInts(w)
	if sum > maxChainSum || p.gv%sum != 0 {
		return result{}, false
	}
	u := p.gv / sum
	key := u * w[p.ask]
	merged := fmt.Sprintf("%d : %d : %d", pp*r, q*r, q*s)
	if reduced := joinInts(w, " : "); reduced != merged {
		merged += " = " + reduced
	}
	steps := []string{
		fmt.Sprintf("Samakan bagian %s: kalikan %d : %d dengan %d dan %d : %d dengan %d.", ctx.items[1], pp, q, r, r, s, q),
		fmt.Sprintf("%s : %s : %s = %s.", ctx.items[0], ctx.items[1], ctx.items[2], merged),
		fmt.Sprintf("Jumlah perbandingan = %s = %d, satu bagian = %d : %d = %d.", joinInts(w, " + "), sum, p.gv, sum, u),
		fmt.Sprintf("%s = %d x %d = %d.", capitalize(ctx.items[p.ask]), w[p.ask], u, key),
	}
	return result{key: key, base: u, max: max(key, p.gv), steps: steps}, true
}

type cand struct {
	v   int
	mis string
}

// distractors lists wrong values in priority order; build keeps the first
// ones that differ from the key and from each other.
func distractors(p params, res result) []cand {
	var cands []cand
	switch p.kind {
	case kindSplit, kindThree:
		sum := sumInts(p.w)
		weight := sum
		if p.ask != askTotal {
			weight = p.w[p.ask]
		}
		switch p.given {
		case givenTotal:
			if p.gv%p.w[p.ask] == 0 {
				cands = append(cands, cand{p.gv / p.w[p.ask], MisBagiTanpaJumlahRasio})
			}
			for j, wj := range p.w {
				if j != p.ask && p.gv%wj == 0 {
					cands = append(cands, cand{p.gv / wj, MisBagiTanpaJumlahRasio})
				}
			}
		case givenDiff:
			for _, j := range []int{p.gi, p.gj} {
				if p.gv%p.w[j] == 0 {
					cands = append(cands, cand{p.gv / p.w[j] * weight, MisBagiTanpaJumlahRasio})
				}
			}
		case givenPart:
			if p.gv%sum == 0 {
				cands = append(cands, cand{p.gv / sum * weight, MisBagiTanpaJumlahRasio})
			}
		}
		if p.ask != askTotal {
			for j, wj := range p.w {
				if j != p.ask {
					cands = append(cands, cand{res.base * wj, MisBagianTertukar})
				}
			}
		}
	case kindInverse:
		c := inverseContexts[p.context%len(inverseContexts)]
		// the other proportion type applied to the same numbers
		if c.inverse && p.y1*p.x2%p.x1 == 0 {
			cands = append(cands, cand{p.y1 * p.x2 / p.x1, MisSenilaiBukanBerbalik})
		}
		if !c.inverse && p.x1*p.y1%p.x2 == 0 {
			cands = append(cands, cand{p.x1 * p.y1 / p.x2, MisSenilaiBukanBerbalik})
		}
	case kindShift:
		other := 1 - p.ask
		cands = append(cands,
			cand{res.key + p.sign*p.n, MisSelisihBukanRasio},
			cand{res.key - p.sign*p.n, MisSelisihBukanRasio},
			cand{p.w[other] * res.base, MisBagianTertukar})
	case kindChain:
		w := chainWeights(p.w[0], p.w[1], p.w[2], p.w[3])
		for j, wj := range w {
			if j != p.ask {
				cands = append(cands, cand{res.base * wj, MisBagianTertukar})
			}
		}
		// B is not matched: A:B:C read as p : q : s
		naive := []int{p.w[0], p.w[1], p.w[3]}
		if ns := sumInts(naive); p.gv%ns == 0 {
			cands = append(cands, cand{p.gv / ns * naive[p.ask], MisTengahTakDisamakan})
		}
	}
	return cands
}

func unitOf(p params) string {
	switch p.kind {
	case kindInverse:
		return inverseContexts[p.context%len(inverseContexts)].unit
	case kindShift:
		return shiftContexts[p.context%len(shiftContexts)].unit
	}
	return ratioContexts[p.context%len(ratioContexts)].unit
}

func build(p params, seed int64) gen.Draft {
	res, _ := solve(p)
	key := res.key
	unit := unitOf(p)

	cands := distractors(p, res)
	for d := 1; d <= key+optionCount; d++ {
		cands = append(cands, cand{key + d, MisSalahHitungDekat}, cand{key - d, MisSalahHitungDekat})
	}

	type opt struct{ text, mis string }
	text := func(v int) string { return fmt.Sprintf("%d%s", v, unit) }
	opts := []opt{{text(key), ""}}
	seen := map[string]bool{text(key): true}
	for _, c := range cands {
		if len(opts) == optionCount {
			break
		}
		if c.v <= 0 {
			continue
		}
		t := text(c.v)
		if seen[t] {
			continue
		}
		seen[t] = true
		opts = append(opts, opt{t, c.mis})
	}

	rng := rand.New(rand.NewSource(derive(seed, shuffleSalt)))
	rng.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })
	options := make([]gen.Option, len(opts))
	answerKey := ""
	for i, o := range opts {
		k := string(rune('A' + i))
		options[i] = gen.Option{Key: k, Text: o.text, MisconceptionID: o.mis}
		if o.mis == "" {
			answerKey = k
		}
	}
	return gen.Draft{
		PatternID: PatternID,
		Stem:      stem(p),
		Options:   options,
		AnswerKey: answerKey,
		Steps:     res.steps,
		Seed:      seed,
	}
}

func stem(p params) string {
	switch p.kind {
	case kindSplit, kindThree:
		return stemRatio(p)
	case kindInverse:
		c := inverseContexts[p.context%len(inverseContexts)]
		return fmt.Sprintf(c.format, p.x1, p.y1, p.x2)
	case kindShift:
		return stemShift(p)
	default:
		ctx := ratioContexts[p.context%len(ratioContexts)]
		return fmt.Sprintf("Perbandingan %s %s dan %s adalah %d : %d, sedangkan perbandingan %s %s dan %s adalah %d : %d. Jika jumlah ketiganya %d%s, %s %s adalah ...",
			ctx.measure, ctx.items[0], ctx.items[1], p.w[0], p.w[1],
			ctx.measure, ctx.items[1], ctx.items[2], p.w[2], p.w[3],
			p.gv, ctx.unit, ctx.measure, ctx.items[p.ask])
	}
}

func stemRatio(p params) string {
	ctx := ratioContexts[p.context%len(ratioContexts)]
	var b strings.Builder
	if len(p.w) == 2 {
		fmt.Fprintf(&b, "Perbandingan %s %s dan %s adalah %s. ", ctx.measure, ctx.items[0], ctx.items[1], joinInts(p.w, " : "))
	} else {
		fmt.Fprintf(&b, "Perbandingan %s %s, %s, dan %s adalah %s. ", ctx.measure, ctx.items[0], ctx.items[1], ctx.items[2], joinInts(p.w, " : "))
	}
	all := "keduanya"
	if len(p.w) == 3 {
		all = "ketiganya"
	}
	switch p.given {
	case givenTotal:
		fmt.Fprintf(&b, "Jika jumlah %s %d%s, ", all, p.gv, ctx.unit)
	case givenDiff:
		fmt.Fprintf(&b, "Jika selisih %s %s dan %s adalah %d%s, ", ctx.measure, ctx.items[p.gi], ctx.items[p.gj], p.gv, ctx.unit)
	case givenPart:
		fmt.Fprintf(&b, "Jika %s %s adalah %d%s, ", ctx.measure, ctx.items[p.gi], p.gv, ctx.unit)
	}
	if p.ask == askTotal {
		fmt.Fprintf(&b, "jumlah %s adalah ...", all)
	} else {
		fmt.Fprintf(&b, "%s %s adalah ...", ctx.measure, ctx.items[p.ask])
	}
	return b.String()
}

func stemShift(p params) string {
	ctx := shiftContexts[p.context%len(shiftContexts)]
	if ctx.age {
		when, verb := "lagi", "menjadi"
		if p.sign < 0 {
			when, verb = "yang lalu", "adalah"
		}
		return fmt.Sprintf("Perbandingan umur %s dan %s sekarang adalah %d : %d. %d tahun %s, perbandingan umur mereka %s %d : %d. Umur %s sekarang adalah ...%s.",
			ctx.n1, ctx.n2, p.w[0], p.w[1], p.n, when, verb, p.cd[0], p.cd[1], shiftItem(ctx, p.ask), ctx.unit)
	}
	change := "ditambah"
	if p.sign < 0 {
		change = "dikurangi"
	}
	return fmt.Sprintf("Perbandingan banyak %s dan %s adalah %d : %d. Setelah masing-masing %s %d%s, perbandingannya menjadi %d : %d. Banyak %s semula adalah ...%s.",
		ctx.n1, ctx.n2, p.w[0], p.w[1], change, p.n, ctx.unit, p.cd[0], p.cd[1], shiftItem(ctx, p.ask), ctx.unit)
}

func shiftItem(ctx shiftContext, i int) string {
	if i == 0 {
		return ctx.n1
	}
	return ctx.n2
}

func signPrefix(sign int) string {
	if sign < 0 {
		return "-"
	}
	return ""
}

func joinInts(v []int, sep string) string {
	parts := make([]string, len(v))
	for i, x := range v {
		parts[i] = fmt.Sprint(x)
	}
	return strings.Join(parts, sep)
}

func sumInts(v []int) int {
	s := 0
	for _, x := range v {
		s += x
	}
	return s
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
