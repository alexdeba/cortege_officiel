package main

import (
	"fmt"
	"image/color"
	"math"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"cortegeofficiel/jeu"
)

// La carte : la photo redressée du plan, zoomable, avec les pions de
// l'encart, les coups possibles surlignés et l'étiquette du point survolé.

type Camera struct{ X, Y, Z float64 }

func (c *Camera) zMin() float64 {
	return min(LargVue/float64(jeu.LargeurCarte), float64(HauteurEcran)/float64(jeu.HauteurCarte))
}

func (c *Camera) toutVoir() {
	c.Z = c.zMin()
	c.X = -(LargVue/c.Z - jeu.LargeurCarte) / 2
	c.Y = -(HauteurEcran/c.Z - jeu.HauteurCarte) / 2
}

func (c *Camera) borner() {
	c.Z = min(max(c.Z, c.zMin()), 1.6)
	w, h := LargVue/c.Z, HauteurEcran/c.Z
	if w >= jeu.LargeurCarte {
		c.X = -(w - jeu.LargeurCarte) / 2
	} else {
		c.X = min(max(c.X, 0), jeu.LargeurCarte-w)
	}
	if h >= jeu.HauteurCarte {
		c.Y = -(h - jeu.HauteurCarte) / 2
	} else {
		c.Y = min(max(c.Y, 0), jeu.HauteurCarte-h)
	}
}

func (c *Camera) zoomer(f, sx, sy float64) {
	mx, my := c.X+sx/c.Z, c.Y+sy/c.Z
	c.Z *= f
	c.borner()
	c.X, c.Y = mx-sx/c.Z, my-sy/c.Z
	c.borner()
}

func (c *Camera) centrer(x, y float64) {
	c.X, c.Y = x-LargVue/2/c.Z, y-HauteurEcran/2/c.Z
	c.borner()
}

func (c *Camera) ecran(x, y float64) (float64, float64)   { return (x - c.X) * c.Z, (y - c.Y) * c.Z }
func (c *Camera) carte(sx, sy float64) (float64, float64) { return c.X + sx/c.Z, c.Y + sy/c.Z }

// tailleJeton : côté d'un pion sur le plan (pixels de la photo).
const tailleJeton = 70.0

// pointSous : le point du plan sous la souris, -1 si aucun.
func (g *Game) pointSous(sx, sy float64) int {
	x, y := g.cam.carte(sx, sy)
	best, bd := -1, math.Inf(1)
	for i, p := range jeu.C.Points {
		d := math.Hypot(x-p.X, y-p.Y)
		if d < max(p.R+26, 14/g.cam.Z) && d < bd {
			best, bd = i, d
		}
	}
	return best
}

// zoneSous : la zone du quadrillage sous la souris.
func (g *Game) zoneSous(sx, sy float64) int {
	x, y := g.cam.carte(sx, sy)
	c, l := -1, -1
	for i := 0; i < len(jeu.LignesX)-1; i++ {
		if x >= jeu.LignesX[i] && x < jeu.LignesX[i+1] {
			c = i
		}
	}
	for i := 0; i < len(jeu.LignesY)-1; i++ {
		if y >= jeu.LignesY[i] && y < jeu.LignesY[i+1] {
			l = i
		}
	}
	if c < 0 || l < 0 {
		return -1
	}
	return c*jeu.NbLig + l
}

func centreZone(z int) (float64, float64) {
	c, l := z/jeu.NbLig, z%jeu.NbLig
	return (jeu.LignesX[c] + jeu.LignesX[c+1]) / 2, (jeu.LignesY[l] + jeu.LignesY[l+1]) / 2
}

// ---------------------------------------------------------------- pions

// occupants : ce qu'il faut dessiner sur un point, dans l'ordre (véhicule
// d'abord, puis les pions à pied, puis voiture et président).
func occupants(e *jeu.Etat, pt int, pl []jeu.Placement) []*jeu.Pion {
	var l []*jeu.Pion
	if pl != nil {
		for _, x := range pl {
			if x.Point == pt && x.Dans < 0 {
				l = append(l, e.Pions[x.Pion])
			}
		}
	} else {
		for _, p := range e.Pions {
			if p.Pos == pt && p.Actif() && p.Dans < 0 {
				l = append(l, p)
			}
		}
	}
	// les véhicules dessous
	var v, a []*jeu.Pion
	for _, p := range l {
		if p.Genre.Vehicule() {
			v = append(v, p)
		} else {
			a = append(a, p)
		}
	}
	return append(v, a...)
}

// placeJeton : la position (plan) du k-ième des n occupants d'un point.
func placeJeton(pt, k, n int) (float64, float64) {
	p := jeu.C.Points[pt]
	if n <= 1 {
		return p.X, p.Y
	}
	d := tailleJeton * 0.52
	switch n {
	case 2:
		return p.X + (float64(k)-0.5)*tailleJeton*1.04, p.Y
	default:
		a := -math.Pi/2 + 2*math.Pi*float64(k)/float64(n)
		return p.X + math.Cos(a)*d*1.05, p.Y + math.Sin(a)*d*1.05
	}
}

// positions : où est dessiné chaque pion (plan) ; pendant la pose, les
// pions du joueur sont à la place qu'il leur donne.
func (g *Game) positions() map[int][2]float64 {
	e := g.p.E
	occ := map[int][]*jeu.Pion{}
	enPose := map[int]bool{}
	if g.p.Req.Type == ReqPlacer {
		for _, x := range g.pl {
			enPose[x.Pion] = true
			if x.Dans < 0 {
				occ[x.Point] = append(occ[x.Point], e.Pions[x.Pion])
			}
		}
	}
	for _, p := range e.Pions {
		if !enPose[p.ID] && p.Actif() && p.Pos >= 0 && p.Dans < 0 {
			occ[p.Pos] = append(occ[p.Pos], p)
		}
	}
	res := map[int][2]float64{}
	for pt, l := range occ {
		var v, a []*jeu.Pion
		for _, p := range l {
			if p.Genre.Vehicule() {
				v = append(v, p)
			} else {
				a = append(a, p)
			}
		}
		l = append(v, a...)
		for k, p := range l {
			x, y := placeJeton(pt, k, len(l))
			res[p.ID] = [2]float64{x, y}
		}
	}
	for _, p := range e.Pions {
		if p.Genre == jeu.Helico && p.Vol >= 0 && p.Actif() {
			x, y := centreZone(p.Vol)
			res[p.ID] = [2]float64{x, y}
		}
	}
	if a := g.p.Anim; a != nil {
		f := min(float64(time.Since(a.Debut))/float64(a.Duree), 1)
		x, y := pointSur(a.Trace, f)
		res[a.Pion] = [2]float64{x, y}
	}
	return res
}

// passagers : combien de pions dans le véhicule v (pose en cours comprise),
// et de quel camp.
func (g *Game) passagers(v *jeu.Pion) (int, jeu.Genre) {
	e := g.p.E
	if g.p.Req.Type == ReqPlacer {
		n := 0
		for _, x := range g.pl {
			if x.Dans == v.ID {
				n++
			}
		}
		if n > 0 {
			return n, jeu.Policier
		}
	}
	ps := e.Passagers(v)
	if len(ps) == 0 {
		return 0, 0
	}
	return len(ps), ps[0].Genre
}

func pointSur(tr []jeu.Point, f float64) (float64, float64) {
	long := 0.0
	for i := 1; i < len(tr); i++ {
		long += math.Hypot(tr[i].X-tr[i-1].X, tr[i].Y-tr[i-1].Y)
	}
	cible := f * long
	for i := 1; i < len(tr); i++ {
		d := math.Hypot(tr[i].X-tr[i-1].X, tr[i].Y-tr[i-1].Y)
		if cible <= d && d > 0 {
			t := cible / d
			return tr[i-1].X + (tr[i].X-tr[i-1].X)*t, tr[i-1].Y + (tr[i].Y-tr[i-1].Y)*t
		}
		cible -= d
	}
	q := tr[len(tr)-1]
	return q.X, q.Y
}

// ---------------------------------------------------------------- surlignage

// mouvementsDe : les mouvements de la demande pour le pion id.
func (g *Game) mouvementsDe(id int) []int {
	var l []int
	for i, m := range g.p.Req.Mvts {
		if m.Pion == id {
			l = append(l, i)
		}
	}
	return l
}

// pionsActifs : les pions que la demande permet de choisir.
func (g *Game) pionsActifs() map[int]bool {
	res := map[int]bool{}
	switch g.p.Req.Type {
	case ReqDeplacer, ReqReculer:
		for _, m := range g.p.Req.Mvts {
			res[m.Pion] = true
		}
	case ReqTirer:
		for _, t := range g.p.Req.Tirs {
			res[t.Tireur] = true
		}
	}
	return res
}

func (g *Game) dessinerCarte(dst *ebiten.Image) {
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Translate(-g.cam.X, -g.cam.Y)
	op.GeoM.Scale(g.cam.Z, g.cam.Z)
	dst.DrawImage(imgPlateau, op)
	e := g.p.E
	t := g.temps()
	pulse := 0.5 + 0.5*math.Sin(t*5)
	z := g.cam.Z

	// le parcours survolé lors du choix de la couleur
	if g.vueCouleur != 0 {
		g.dessinerParcours(dst, g.vueCouleur, g.p.Req.Retour)
	}
	// le prochain pas de la voiture
	if !e.Fini && e.Phase != jeu.PhPreparation {
		for _, q := range g.prochainsVus() {
			x, y := g.cam.ecran(jeu.C.Points[q].X, jeu.C.Points[q].Y)
			r := (jeu.C.Points[q].R + 10) * z
			anneau(dst, x, y, r+2, max(3, 6*z), alpha(coulBlanc, 0.8))
			anneau(dst, x, y, r+2, max(1.5, 3*z), coulEncre)
		}
	}

	g.dessinerPions(dst, pulse)
	switch g.p.Req.Type {
	case ReqAvancer:
		for _, q := range g.p.Req.Points {
			g.anneauCible(dst, q, pulse, coulSignal)
		}
	case ReqDeplacer, ReqReculer:
		if g.sel >= 0 {
			for _, i := range g.mouvementsDe(g.sel) {
				m := g.p.Req.Mvts[i]
				if m.Vers >= 0 {
					c := coulSignal
					if m.Type == jeu.MvtFuite {
						c = coulRouge
					}
					g.anneauCible(dst, m.Vers, pulse, c)
				} else {
					g.cadreZone(dst, m.Zone, pulse)
				}
			}
		}
	case ReqTirer:
		if g.sel >= 0 {
			for _, tr := range g.p.Req.Tirs {
				if tr.Tireur == g.sel {
					g.anneauCible(dst, e.Lieu(e.Pions[tr.Cible]), pulse, coulRouge)
				}
			}
		}
	case ReqPlacer:
		if g.sel >= 0 {
			for i := range jeu.C.Points {
				if g.posePossible(g.sel, i) {
					x, y := g.cam.ecran(jeu.C.Points[i].X, jeu.C.Points[i].Y)
					anneau(dst, x, y, (jeu.C.Points[i].R+6)*z, max(2, 4*z), alpha(coulSignal, 0.8))
				}
			}
		}
	}
	if d := g.p.DernierTir; d != nil && time.Since(g.tirVu) < 2500*time.Millisecond {
		g.dessinerTir(dst, d)
	}
	if g.survol >= 0 {
		g.etiquette(dst, g.survol)
	}
}

// prochainsVus : le prochain point de la voiture, tel que le joueur le
// connaît (le parcours du retour reste secret un moment).
func (g *Game) prochainsVus() []int {
	e := g.p.E
	if e.Discours >= 3 && !e.RetourPublic && !g.p.Moi.A(jeu.RPolice) {
		return nil
	}
	return e.ProchainsPoints()
}

func (g *Game) anneauCible(dst *ebiten.Image, q int, pulse float64, c color.RGBA) {
	if q < 0 {
		return
	}
	x, y := g.cam.ecran(jeu.C.Points[q].X, jeu.C.Points[q].Y)
	r := (jeu.C.Points[q].R + 8) * g.cam.Z
	if g.p.E.Occupe(q) {
		r = tailleJeton * 0.75 * g.cam.Z // autour des pions qui l'occupent
	}
	r = max(r, 9)
	anneau(dst, x, y, r+3+2*pulse, 6, alpha(coulAmbre2, 0.85))
	anneau(dst, x, y, r+3+2*pulse, 2.5, c)
}

func (g *Game) cadreZone(dst *ebiten.Image, zn int, pulse float64) {
	c, l := zn/jeu.NbLig, zn%jeu.NbLig
	x0, y0 := g.cam.ecran(jeu.LignesX[c], jeu.LignesY[l])
	x1, y1 := g.cam.ecran(jeu.LignesX[c+1], jeu.LignesY[l+1])
	rect(dst, x0, y0, x1-x0, y1-y0, alpha(coulAmbre2, 0.12+0.1*pulse))
	cadre(dst, x0+3, y0+3, x1-x0-6, y1-y0-6, 2, alpha(coulSignal, 0.9))
}

// dessinerParcours : un parcours coloré, de l'Hôtel de Ville à la
// Grand-Place ou de la Grand-Place à l'aéroport.
func (g *Game) dessinerParcours(dst *ebiten.Image, col jeu.Couleur, retour bool) {
	dep := jeu.PtHotelDeVille
	if retour {
		dep = jeu.PtGP3
	}
	ch := append([]int{dep}, jeu.C.Chemin(col, retour, dep)...)
	var pts [][2]float64
	for _, q := range ch {
		x, y := g.cam.ecran(jeu.C.Points[q].X, jeu.C.Points[q].Y)
		pts = append(pts, [2]float64{x, y})
	}
	if len(pts) < 2 {
		return
	}
	polyligne(dst, pts, max(8, 22*g.cam.Z), alpha(coulEncre, 0.6))
	polyligne(dst, pts, max(5, 14*g.cam.Z), alpha(coulParcours(col), 0.9))
}

func (g *Game) dessinerPions(dst *ebiten.Image, pulse float64) {
	e := g.p.E
	pos := g.positions()
	actifs := g.pionsActifs()
	t := tailleJeton * g.cam.Z
	moi := g.p.Moi
	ordre := make([]*jeu.Pion, 0, len(e.Pions))
	for _, p := range e.Pions {
		if p.Genre.Vehicule() {
			ordre = append(ordre, p)
		}
	}
	for _, p := range e.Pions {
		if !p.Genre.Vehicule() && p.Genre != jeu.VoitureOfficielle && p.Genre != jeu.President {
			ordre = append(ordre, p)
		}
	}
	ordre = append(ordre, e.Voiture, e.Pres)
	for _, p := range ordre {
		xy, ok := pos[p.ID]
		if !ok {
			continue
		}
		x, y := g.cam.ecran(xy[0], xy[1])
		if x < -t || y < -t || x > LargVue+t || y > HauteurEcran+t {
			continue
		}
		tt := t
		if p.Genre == jeu.VoitureOfficielle || p.Genre == jeu.President {
			tt = t * 1.25
			disque(dst, x, y, tt*0.78, alpha(coulAmbre2, 0.55+0.3*pulse))
		}
		if p.Genre == jeu.Helico && p.Vol >= 0 {
			tt = t * 1.35
			rect(dst, x-tt/2+8, y-tt/2+10, tt, tt, alpha(coulNoir, 0.3))
		}
		// ombre et liseré
		rect(dst, x-tt/2+2, y-tt/2+2, tt, tt, alpha(coulNoir, 0.35))
		if actifs[p.ID] {
			c := coulAmbre
			if p.ID == g.sel {
				c = coulSignal
			}
			rect(dst, x-tt/2-3-pulse*2, y-tt/2-3-pulse*2, tt+6+pulse*4, tt+6+pulse*4, alpha(c, 0.9))
		} else if p.ID == g.sel {
			rect(dst, x-tt/2-4, y-tt/2-4, tt+8, tt+8, coulSignal)
		}
		pion(dst, p.ID, x, y, tt, 1)
		// ce que le joueur sait de son identité
		if moi != nil && p.Genre.APied() {
			switch e.Connu(moi, p.ID) {
			case jeu.RTerro:
				badge(dst, x+tt/2, y-tt/2, max(tt*0.32, 9), "T", coulRoles[jeu.RTerro], e.Revele[p.ID] == jeu.RTerro)
			case jeu.RAgent:
				badge(dst, x+tt/2, y-tt/2, max(tt*0.32, 9), "A", coulRoles[jeu.RAgent], e.Revele[p.ID] == jeu.RAgent)
			}
		}
		if p.Blessures > 0 {
			trait(dst, x-tt/2, y+tt/2, x+tt/2, y-tt/2, max(2, tt*0.08), coulRouge)
		}
		if p.Genre.Vehicule() {
			if n, gen := g.passagers(p); n > 0 {
				c := coulBleu
				if gen == jeu.Manifestant {
					c = coulRouge
				}
				badge(dst, x-tt/2, y-tt/2, max(tt*0.3, 9), fmt.Sprint(n), c, true)
			}
		}
	}
}

// badge : une pastille avec une lettre, au coin d'un pion.
func badge(dst *ebiten.Image, x, y, r float64, s string, c color.RGBA, plein bool) {
	if plein {
		disque(dst, x, y, r+1.5, coulBlanc)
		disque(dst, x, y, r, c)
		ecrireCentre(dst, s, x, y-r*0.95, face(r*1.45, policeGras), coulBlanc)
		return
	}
	disque(dst, x, y, r+1.5, c)
	disque(dst, x, y, r-0.5, coulBlanc)
	ecrireCentre(dst, s, x, y-r*0.95, face(r*1.45, policeGras), c)
}

// dessinerTir : le trait du dernier tir, et ses dés.
func (g *Game) dessinerTir(dst *ebiten.Image, ev *jeu.Ev) {
	e := g.p.E
	pos := g.positions()
	a, ok1 := pos[ev.Tir.Tireur]
	b, ok2 := pos[ev.Tir.Cible]
	if !ok2 {
		if l := e.Lieu(e.Pions[ev.Tir.Cible]); l >= 0 {
			b, ok2 = [2]float64{jeu.C.Points[l].X, jeu.C.Points[l].Y}, true
		}
	}
	if !ok1 || !ok2 {
		return
	}
	x0, y0 := g.cam.ecran(a[0], a[1])
	x1, y1 := g.cam.ecran(b[0], b[1])
	trait(dst, x0, y0, x1, y1, 5, alpha(coulNoir, 0.5))
	trait(dst, x0, y0, x1, y1, 2.5, coulRouge)
	disque(dst, x1, y1, 9, coulRouge)
	s := 34.0
	bx, by := (x0+x1)/2-s-6, (y0+y1)/2-s-14
	rect(dst, bx-4, by-4, 2*s+28, s+16+30, alpha(coulNoir, 0.5))
	arrondi(dst, bx-8, by-8, 2*s+28, s+16+30, 2, alpha(coulAcier, 0.96), coulAcier3, 1)
	de(dst, bx, by, s, ev.Des[0], coulCraie, coulEncre)
	de(dst, bx+s+12, by, s, ev.Des[1], coulCraie, coulEncre)
	ecrireCentre(dst, strings.ToUpper(jeu.NomsResultats[ev.Result]), bx+s+6, by+s+2, face(19, policeTitre), coulSignal2)
}

// etiquette : le point survolé, sa zone et ce qui s'y trouve.
func (g *Game) etiquette(dst *ebiten.Image, q int) {
	e := g.p.E
	pc := jeu.C.Points[q]
	titre := "Point de rue"
	switch {
	case q == jeu.PtTribune:
		titre = "La tribune officielle"
	case q == jeu.PtHotelDeVille:
		titre = "Devant l'Hôtel de Ville"
	case q == jeu.PtAeroport:
		titre = "L'aéroport"
	case q == jeu.PtGP1 || q == jeu.PtGP2 || q == jeu.PtGP3:
		titre = fmt.Sprintf("Grand-Place, point %d", map[int]int{jeu.PtGP1: 1, jeu.PtGP2: 2, jeu.PtGP3: 3}[q])
	case pc.H:
		titre = "Toit avec héliport (H)"
	case pc.Toit:
		titre = "Toit d'immeuble"
	case pc.Couleurs != 0:
		var c []string
		for _, col := range jeu.CouleursParcours {
			if pc.Couleurs&col != 0 {
				c = append(c, col.Nom())
			}
		}
		titre = "Parcours " + strings.Join(c, ", ")
	}
	sous := "Zone " + jeu.NomZone(pc.Zone())
	if pc.GP {
		sous += " · dans les barrières de la Grand-Place"
	}
	if len(jeu.C.ToitsDepuis(q)) > 0 {
		sous += " · flèche vers un toit"
	}
	var lignes []string
	for _, p := range occupants(e, q, nil) {
		l := p.Nom()
		if p.Genre.Vehicule() {
			if ps := e.Passagers(p); len(ps) > 0 {
				var n []string
				for _, x := range ps {
					n = append(n, x.Nom())
				}
				l += " (" + strings.Join(n, ", ") + ")"
			} else {
				l += " (vide)"
			}
		}
		if p.Blessures > 0 {
			l += ", blessé"
		}
		lignes = append(lignes, l)
	}
	if q == e.Pres.Pos {
		lignes = append(lignes, "le président, à la tribune")
	}
	x, y := g.cam.ecran(pc.X, pc.Y)
	ft, fs := face(19, policeTitre), face(15, policeTexte)
	w := max(largeur(titre, ft), largeur(sous, fs)) + 26
	for _, l := range lignes {
		w = max(w, largeur(l, fs)+26)
	}
	h := 56 + 20*float64(len(lignes))
	bx := min(max(x-w/2, 6), LargVue-w-6)
	by := y - tailleJeton*g.cam.Z - h - 10
	if by < 6 {
		by = y + tailleJeton*g.cam.Z + 10
	}
	rect(dst, bx+4, by+4, w, h, alpha(coulNoir, 0.5))
	arrondi(dst, bx, by, w, h, 2, alpha(coulAcier, 0.96), coulAcier3, 1)
	rect(dst, bx, by, 4, h, coulSignal)
	ecrire(dst, titre, bx+13, by+3, ft, coulCraie)
	ecrire(dst, sous, bx+13, by+32, fs, coulCendre)
	for i, l := range lignes {
		ecrire(dst, "• "+l, bx+13, by+52+20*float64(i), fs, coulCraie)
	}
}

// ---------------------------------------------------------------- clics

func (g *Game) cliquerCarte(q int) {
	r := g.p.Req
	switch r.Type {
	case ReqAvancer:
		for i, x := range r.Points {
			if x == q {
				g.p.repondre(i)
			}
		}
	case ReqPlacer:
		g.cliquerPose(q)
	case ReqDeplacer, ReqReculer:
		g.cliquerMouvement(q)
	case ReqTirer:
		g.cliquerTir(q)
	}
}

// pionsSur : les pions dessinés sur le point q (et ceux qu'il transporte).
func (g *Game) pionsSur(q int) []int {
	e := g.p.E
	var l []int
	if q < 0 {
		return l
	}
	for _, p := range occupants(e, q, nil) {
		l = append(l, p.ID)
		if p.Genre.Vehicule() {
			for _, x := range e.Passagers(p) {
				l = append(l, x.ID)
			}
		}
	}
	return l
}

func (g *Game) cliquerMouvement(q int) {
	r := g.p.Req
	// une destination du pion choisi ? D'abord le point cliqué ; pour
	// l'hélicoptère, sinon la zone sous la souris.
	if g.sel >= 0 {
		var ok []int
		for _, i := range g.mouvementsDe(g.sel) {
			if m := r.Mvts[i]; q >= 0 && m.Vers == q {
				ok = append(ok, i)
			}
		}
		if len(ok) == 0 {
			z := g.zoneSous(float64(g.mx), float64(g.my))
			for _, i := range g.mouvementsDe(g.sel) {
				if m := r.Mvts[i]; m.Vers < 0 && m.Type == jeu.MvtVol && m.Zone == z {
					ok = append(ok, i)
				}
			}
		}
		switch len(ok) {
		case 0:
		case 1:
			g.jouerMouvement(ok[0])
			return
		default:
			g.choixMvt = ok
			return
		}
	}
	// sinon : le pion mobile dessiné sous la souris, ou ceux du point (un
	// clic de plus passe au suivant)
	actifs := g.pionsActifs()
	if id := g.pionSous(actifs); id >= 0 && id != g.sel {
		g.sel, g.choixMvt = id, nil
		return
	}
	var cands []int
	for _, id := range g.pionsSur(q) {
		if actifs[id] {
			cands = append(cands, id)
		}
	}
	if len(cands) == 0 {
		if q < 0 && g.sel >= 0 {
			return
		}
		g.sel, g.choixMvt = -1, nil
		return
	}
	k := 0
	for i, id := range cands {
		if id == g.sel {
			k = (i + 1) % len(cands)
		}
	}
	g.sel, g.choixMvt = cands[k], nil
}

// pionSous : le pion choisissable dessiné sous la souris, -1 si aucun.
func (g *Game) pionSous(actifs map[int]bool) int {
	t := tailleJeton * g.cam.Z / 2
	sx, sy := float64(g.mx), float64(g.my)
	for id, xy := range g.positions() {
		if !actifs[id] {
			continue
		}
		x, y := g.cam.ecran(xy[0], xy[1])
		if math.Abs(sx-x) <= t && math.Abs(sy-y) <= t {
			return id
		}
	}
	return -1
}

func (g *Game) jouerMouvement(i int) {
	g.p.repondre(i)
	g.sel, g.choixMvt = -1, nil
}

func (g *Game) cliquerTir(q int) {
	r := g.p.Req
	e := g.p.E
	if g.sel >= 0 {
		for i, t := range r.Tirs {
			if t.Tireur == g.sel && e.Lieu(e.Pions[t.Cible]) == q {
				g.p.repondre(i)
				g.sel = -1
				return
			}
		}
	}
	actifs := g.pionsActifs()
	var cands []int
	for _, id := range g.pionsSur(q) {
		if actifs[id] {
			cands = append(cands, id)
		}
	}
	if len(cands) == 0 {
		g.sel = -1
		return
	}
	k := 0
	for i, id := range cands {
		if id == g.sel {
			k = (i + 1) % len(cands)
		}
	}
	g.sel = cands[k]
}

// ---------------------------------------------------------------- pose

// posePossible : le pion id peut-il être posé en q (pose en cours) ?
func (g *Game) posePossible(id, q int) bool {
	e := g.p.E
	pc := jeu.C.Points[q]
	if pc.Toit || pc.Couleurs != 0 || q == jeu.PtTribune {
		return false
	}
	p := e.Pions[id]
	if p.Genre == jeu.Manifestant && (pc.GP || e.PolicePresente(q) || e.VehiculeSur(q) != nil) {
		return false
	}
	n, v := 0, 0
	for _, x := range g.pl {
		if x.Point == q && x.Pion != id {
			if e.Pions[x.Pion].Genre.Vehicule() {
				v++
			} else if x.Dans < 0 {
				n++
			}
		}
	}
	if p.Genre.Vehicule() {
		return v == 0
	}
	return n < 2
}

func (g *Game) cliquerPose(q int) {
	if q < 0 {
		return
	}
	e := g.p.E
	if g.sel >= 0 && g.posePossible(g.sel, q) {
		for i := range g.pl {
			x := &g.pl[i]
			if x.Pion == g.sel {
				x.Point, x.Dans = q, -1
			}
			// les passagers suivent leur véhicule
			if x.Dans == g.sel {
				x.Point = q
			}
		}
		g.sel, g.plMsg = -1, ""
		return
	}
	// choisir un pion posé sur ce point (y compris dans un véhicule)
	var cands []int
	for _, x := range g.pl {
		if x.Point == q {
			cands = append(cands, x.Pion)
		}
	}
	if len(cands) == 0 {
		g.sel = -1
		return
	}
	k := 0
	for i, id := range cands {
		if id == g.sel {
			k = (i + 1) % len(cands)
		}
	}
	g.sel = cands[k]
	_ = e
}

func (g *Game) validerPose() {
	if msg := g.p.E.ValiderPlacement(g.p.Req.Role, g.pl); msg != "" {
		g.plMsg = msg
		return
	}
	g.p.repondre(append([]jeu.Placement(nil), g.pl...))
}
