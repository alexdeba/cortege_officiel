package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"cortegeofficiel/jeu"
)

// Zoom à la barre d'espace et mini-carte : quand la carte est zoomée, un
// médaillon en haut à droite de la vue montre toute la carte, les pions et
// le cadre de ce qu'on voit ; un clic ou un glisser dedans y déplace la vue.

const (
	zoomEspace = 0.9 // zoom de la barre d'espace
	miniW      = 270.0
	miniMarge  = 12.0
)

func miniH() float64 { return miniW * float64(jeu.HauteurCarte) / float64(jeu.LargeurCarte) }

func miniX() float64 { return LargVue - miniW - miniMarge }

const miniY = miniMarge

func (g *Game) zoomee() bool { return g.cam.Z > g.cam.zMin()*1.01 }

// surMini : la souris est-elle sur la mini-carte (affichée) ?
func (g *Game) surMini() bool {
	if !g.zoomee() || g.scores || g.p.E.Fini {
		return false
	}
	x, y := float64(g.mx), float64(g.my)
	return x >= miniX() && x < miniX()+miniW && y >= miniY && y < miniY+miniH()
}

// basculerZoom : barre d'espace, vue d'ensemble ↔ zoom sur la souris (ou
// sur votre pion si la souris n'est pas sur la carte).
func (g *Game) basculerZoom() {
	if g.zoomee() {
		g.cam.toutVoir()
		return
	}
	var x, y float64
	if sx, sy := float64(g.mx), float64(g.my); sx >= 0 && sx < LargVue && sy >= 0 && sy < HauteurEcran {
		x, y = g.cam.carte(sx, sy)
	} else if v := g.p.E.Voiture; v.Pos >= 0 {
		x, y = jeu.C.Points[v.Pos].X, jeu.C.Points[v.Pos].Y
	}
	g.cam.Z = zoomEspace
	g.cam.centrer(x, y)
}

// majMiniCarte : clic et glisser sur la mini-carte ; vrai si la souris y est.
func (g *Game) majMiniCarte() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.basculerZoom()
	}
	sur := g.surMini()
	if sur && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.glisseMini = true
	}
	if g.glisseMini {
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.glisseMini = false
		} else {
			f := float64(jeu.LargeurCarte) / miniW
			g.cam.centrer((float64(g.mx)-miniX())*f, (float64(g.my)-miniY)*f)
		}
	}
	return sur || g.glisseMini
}

func (g *Game) dessinerMiniCarte(dst *ebiten.Image) {
	if !g.zoomee() || g.scores || g.p.E.Fini {
		return
	}
	x0, y0, w, h := miniX(), miniY, miniW, miniH()
	rect(dst, x0-6+4, y0-6+4, w+12, h+12, alpha(coulNoir, 0.6))
	arrondi(dst, x0-6, y0-6, w+12, h+12, 2, coulAcier, coulAcier3, 1.5)
	s := w / float64(jeu.LargeurCarte)
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(s, s)
	op.GeoM.Translate(x0, y0)
	dst.DrawImage(imgPlateau, op)
	// la foule, la police, et la voiture présidentielle
	for id, xy := range g.positions() {
		p := g.p.E.Pions[id]
		c := coulRouge
		switch p.Genre {
		case jeu.Policier, jeu.VoiturePolice, jeu.Moto, jeu.Helico:
			c = coulBleu
		case jeu.VoitureOfficielle:
			continue
		}
		rect(dst, x0+xy[0]*s-2, y0+xy[1]*s-2, 4, 4, c)
	}
	if v := g.p.E.Voiture; v.Pos >= 0 {
		pt := jeu.C.Points[v.Pos]
		disque(dst, x0+pt.X*s, y0+pt.Y*s, 6, coulNoir)
		disque(dst, x0+pt.X*s, y0+pt.Y*s, 4.5, coulAmbre2)
	}
	// ce que montre la vue
	vx, vy := g.cam.X*s, g.cam.Y*s
	vw, vh := LargVue/g.cam.Z*s, HauteurEcran/g.cam.Z*s
	rx, ry := max(x0+vx, x0), max(y0+vy, y0)
	rw, rh := min(x0+vx+vw, x0+w)-rx, min(y0+vy+vh, y0+h)-ry
	rect(dst, rx, ry, rw, rh, alpha(coulAmbre2, 0.15))
	cadre(dst, rx, ry, rw, rh, 2.5, coulAmbre)
}
