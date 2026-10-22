package jeu

import "fmt"

// Le plan : points relevés sur la photo (outils/releve), traits rouges qui
// les relient, flèches vers les toits, zones du quadrillage.

type Point struct{ X, Y float64 }

// Couleur d'un parcours du cortège (un point peut en porter plusieurs).
type Couleur uint8

const (
	Jaune Couleur = 1 << iota
	Rouge
	Bleu
)

var CouleursParcours = []Couleur{Jaune, Rouge, Bleu}

func (c Couleur) Nom() string {
	switch c {
	case Jaune:
		return "jaune"
	case Rouge:
		return "rouge"
	case Bleu:
		return "bleu"
	}
	return "?"
}

func (c Couleur) indice() int {
	switch c {
	case Jaune:
		return 0
	case Rouge:
		return 1
	}
	return 2
}

type PointCarte struct {
	X, Y, R        float64
	Couleurs       Couleur
	Toit, H        bool
	Bati, Col, Lig int
	GP, Sortie     bool
}

// Zone : indice col*7+lig, nom « A1 » à « G7 ».
func (p PointCarte) Zone() int { return p.Col*NbLig + p.Lig }

const (
	NbCol = 7
	NbLig = 7
)

func NomZone(z int) string { return fmt.Sprintf("%c%d", 'A'+z/NbLig, z%NbLig+1) }

// DistZones : distance entre zones en nombre de zones (diagonales comprises).
func DistZones(a, b int) int {
	dc, dl := a/NbLig-b/NbLig, a%NbLig-b%NbLig
	return max(dc, -dc, dl, -dl)
}

type Carte struct {
	Points  []PointCarte
	Voisins [][]int  // voisins de déplacement (traits, flèches, toits d'un même immeuble)
	Trait   [][]bool // relié par un trait rouge
	Fleche  [][]bool // point de rue → toit désigné par une de ses flèches
	Dist    [][]int  // distance en coups (tous voisins)
	// parcours de la voiture présidentielle : distance au but en restant sur
	// les points de la couleur, -1 hors parcours
	Aller, Retour [3][]int
	ToitsZone     [][]int // points de toit par zone
}

var C = construire()

func construire() *Carte {
	n := len(pointsCarte)
	c := &Carte{Points: pointsCarte}
	mat := func() [][]bool {
		m := make([][]bool, n)
		for i := range m {
			m[i] = make([]bool, n)
		}
		return m
	}
	c.Trait, c.Fleche = mat(), mat()
	vois := mat()
	for _, t := range traitsCarte {
		c.Trait[t[0]][t[1]], c.Trait[t[1]][t[0]] = true, true
		vois[t[0]][t[1]], vois[t[1]][t[0]] = true, true
	}
	for _, a := range accesCarte {
		c.Fleche[a[0]][a[1]] = true
		vois[a[0]][a[1]], vois[a[1]][a[0]] = true, true
	}
	for i, p := range c.Points {
		for j, q := range c.Points {
			if i != j && p.Toit && q.Toit && p.Bati != 0 && p.Bati == q.Bati {
				vois[i][j] = true
			}
		}
	}
	c.Voisins = make([][]int, n)
	for i := range n {
		for j := range n {
			if vois[i][j] {
				c.Voisins[i] = append(c.Voisins[i], j)
			}
		}
	}
	c.Dist = make([][]int, n)
	for i := range n {
		c.Dist[i] = c.bfs([]int{i}, func(int) bool { return true })
	}
	for _, col := range CouleursParcours {
		k := col.indice()
		sur := func(p int) bool { return c.Points[p].Couleurs&col != 0 || p == PtGP1 }
		var buts []int
		for i, p := range c.Points {
			if p.GP && sur(i) && i != PtGP3 {
				buts = append(buts, i)
			}
		}
		c.Aller[k] = c.bfsTrait(buts, sur)
		c.Retour[k] = c.bfsTrait([]int{PtAeroport}, func(p int) bool { return c.Points[p].Couleurs&col != 0 })
	}
	c.ToitsZone = make([][]int, NbCol*NbLig)
	for i, p := range c.Points {
		if p.Toit {
			c.ToitsZone[p.Zone()] = append(c.ToitsZone[p.Zone()], i)
		}
	}
	return c
}

func (c *Carte) bfs(depart []int, ok func(int) bool) []int {
	d := make([]int, len(c.Points))
	for i := range d {
		d[i] = -1
	}
	file := append([]int(nil), depart...)
	for _, p := range depart {
		d[p] = 0
	}
	for len(file) > 0 {
		p := file[0]
		file = file[1:]
		for _, q := range c.Voisins[p] {
			if d[q] < 0 && ok(q) {
				d[q] = d[p] + 1
				file = append(file, q)
			}
		}
	}
	return d
}

// bfsTrait : comme bfs, mais seulement le long des traits rouges (la voiture
// présidentielle reste dans les rues).
func (c *Carte) bfsTrait(depart []int, ok func(int) bool) []int {
	d := make([]int, len(c.Points))
	for i := range d {
		d[i] = -1
	}
	file := append([]int(nil), depart...)
	for _, p := range depart {
		d[p] = 0
	}
	for len(file) > 0 {
		p := file[0]
		file = file[1:]
		for _, q := range c.Voisins[p] {
			if c.Trait[p][q] && d[q] < 0 && ok(q) {
				d[q] = d[p] + 1
				file = append(file, q)
			}
		}
	}
	return d
}

// DistParcours : distance au but du parcours de couleur col (aller ou
// retour) depuis le point p, -1 si p n'est pas sur ce parcours.
func (c *Carte) DistParcours(col Couleur, retour bool, p int) int {
	if retour {
		return c.Retour[col.indice()][p]
	}
	return c.Aller[col.indice()][p]
}

// Chemin : les points que suivra la voiture depuis p (non compris), en
// prenant à chaque pas le premier voisin qui rapproche du but.
func (c *Carte) Chemin(col Couleur, retour bool, p int) []int {
	var ch []int
	for d := c.DistParcours(col, retour, p); d > 0; {
		suiv := -1
		for _, q := range c.Voisins[p] {
			if c.Trait[p][q] && c.DistParcours(col, retour, q) == d-1 {
				suiv = q
				break
			}
		}
		if suiv < 0 {
			break
		}
		ch = append(ch, suiv)
		p, d = suiv, d-1
	}
	return ch
}

// BatimentAccessible : les points de toit atteints par les flèches de p.
func (c *Carte) ToitsDepuis(p int) []int {
	var t []int
	for _, q := range c.Voisins[p] {
		if c.Fleche[p][q] {
			t = append(t, q)
		}
	}
	return t
}
