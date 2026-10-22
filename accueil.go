package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"cortegeofficiel/jeu"
)

// L'écran d'accueil : l'illustration entière (assets/cover.jpg), et sous
// elle le nombre de joueurs, le camp du joueur humain, Jouer et la règle.

const hBandeau = 107.0

// camps : les places proposées selon le nombre de joueurs.
func camps(nbJ int) []string {
	switch nbJ {
	case 2:
		return []string{"Police + agents", "Manifs + terroristes"}
	case 3:
		return []string{"Police + agents", "Manifestants", "Terroristes"}
	}
	return []string{"Police", "Manifestants", "Agents secrets", "Terroristes"}
}

func (g *Game) boutonsAccueil() {
	y := float64(HauteurEcran) - hBandeau + 40
	x := 20.0
	for _, n := range []int{2, 3, 4} {
		n := n
		b := g.bouton(x, y, 118, 50, fmt.Sprintf("%d joueurs", n), true, func() {
			g.nbJ = n
			g.place = min(g.place, n-1)
		})
		b.Style, b.Choisi = BoutonDiscret, g.nbJ == n
		x += 126
	}
	x += 20
	cs := camps(g.nbJ)
	lw := (1600 - x - 20 - 2*(170+10) - 20 - float64(len(cs)-1)*8) / float64(len(cs))
	for i, c := range cs {
		i := i
		b := g.bouton(x, y, lw, 50, c, true, func() { g.place = i })
		b.Choisi = g.place == i
		x += lw + 8
	}
	x += 12
	g.bouton(x, y, 170, 50, "Règle (R)", true, func() { g.regles = true }).Style = BoutonDiscret
	b := g.bouton(x+180, y, 170, 50, "Jouer", true, func() {
		g.commencer(0)
		g.p.lancer()
	})
	b.Style = BoutonPrincipal
}

func (g *Game) dessinerAccueil(dst *ebiten.Image) {
	s := float64(LargeurEcran) / float64(imgCouv.Bounds().Dx())
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(s, s)
	dst.DrawImage(imgCouv, op)
	h := float64(imgCouv.Bounds().Dy()) * s
	y0 := float64(HauteurEcran) - hBandeau
	rect(dst, 0, h, LargeurEcran, float64(HauteurEcran)-h, coulAcier)
	filet(dst, 0, h, LargeurEcran, 7, coulNoir)
	ecrire(dst, "NOMBRE DE JOUEURS", 22, y0+12, face(15, policeMachine), coulCendre)
	ecrire(dst, "VOTRE CAMP", 20+3*126+22, y0+12, face(15, policeMachine), coulCendre)
	ecrireDroite(dst, "Jean-Charles et Sylvie Rodriguez · Jeux & Stratégie, 1983", LargeurEcran-22, y0+12, face(15, policeItalique), coulCendre)
	_ = jeu.NbRoles
}
