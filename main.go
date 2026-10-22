// Cortège officiel — d'après le jeu de Jean-Charles et Sylvie Rodriguez
// paru en encart dans Jeux & Stratégie (1983), plan dessiné par Philippe
// Fassier et Alain Meyer (assets/Regles1.jpg, Regles2.jpg, Encart.jpg).
// Police, manifestants, agents secrets et terroristes autour du cortège
// d'un président étranger, de l'Hôtel de Ville à l'aéroport ; un joueur
// humain contre l'ordinateur.
package main

import (
	"flag"
	"log"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"cortegeofficiel/jeu"
)

const (
	LargeurEcran = 1600
	HauteurEcran = 1000
	LargVue      = 1140.0 // la carte
	PanneauX     = LargVue
	LargPanneau  = LargeurEcran - PanneauX
)

type Ecran int

const (
	EcranAccueil Ecran = iota
	EcranJeu
)

type Game struct {
	ecran            Ecran
	nbJ              int
	place            int // le joueur humain (indice dans la partie)
	p                *Partie
	numReq           int
	cam              Camera
	glisse           bool
	glisseMini       bool
	gx, gy, gcx, gcy float64

	survol     int   // point sous la souris (-1)
	sel        int   // pion choisi (-1)
	choixMvt   []int // mouvements possibles vers le point cliqué (plusieurs types)
	pl         []jeu.Placement
	plMsg      string
	selSec     []int
	vueCouleur jeu.Couleur // parcours survolé (choix de la couleur)
	scores     bool
	regles     bool
	pageRgl    int
	defil      int

	tirVu      time.Time
	dernierTir *jeu.Ev

	boutons []*Bouton
	mx, my  int
	t0      time.Time
	images  int
}

func nouveauJeu() *Game {
	g := &Game{nbJ: 4, place: 0, t0: time.Now(), survol: -1, sel: -1}
	g.cam.toutVoir()
	return g
}

func (g *Game) Layout(int, int) (int, int) { return LargeurEcran, HauteurEcran }

func (g *Game) temps() float64 { return time.Since(g.t0).Seconds() }

func (g *Game) commencer(graine uint64) {
	if graine == 0 {
		graine = rand.Uint64()
	}
	g.p = nouvellePartie(g.nbJ, g.place, graine)
	g.ecran = EcranJeu
	g.numReq = -1
	g.cam.toutVoir()
}

func (g *Game) Update() error {
	g.mx, g.my = ebiten.CursorPosition()
	clic := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if g.p != nil {
		g.p.mu.Lock()
		defer g.p.mu.Unlock()
		if g.p.DernierTir != g.dernierTir {
			g.dernierTir, g.tirVu = g.p.DernierTir, time.Now()
		}
		if g.p.Req.Num != g.numReq {
			g.numReq = g.p.Req.Num
			g.nouvelleDemande()
		}
	}
	g.construireBoutons()
	if clic {
		for _, b := range g.boutons {
			if b.Actif && b.Contient(g.mx, g.my) {
				b.Action()
				return nil
			}
		}
	}
	if g.regles {
		g.touchesRegle()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.regles = true
		return nil
	}
	if g.ecran == EcranAccueil {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.commencer(0)
			g.p.lancer()
		}
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.p.Rapide = !g.p.Rapide
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.scores = !g.scores
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		switch {
		case g.scores:
			g.scores = false
		default:
			g.sel, g.choixMvt = -1, nil
		}
	}
	surMini := g.majMiniCarte()
	if !surMini {
		g.cameraClavierSouris()
	}
	g.survol = -1
	if !surMini && !g.scores && float64(g.mx) < LargVue {
		g.survol = g.pointSous(float64(g.mx), float64(g.my))
	}
	g.survolCapture()
	if clic && !surMini && !g.scores && float64(g.mx) < LargVue {
		g.cliquerCarte(g.survol)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.entree()
	}
	return nil
}

// nouvelleDemande : remise à zéro de la saisie à chaque demande du moteur.
func (g *Game) nouvelleDemande() {
	g.sel, g.choixMvt = -1, nil
	r := g.p.Req
	switch r.Type {
	case ReqPlacer:
		g.pl = append([]jeu.Placement(nil), r.Prop...)
		g.plMsg = ""
	case ReqSecrets:
		g.selSec = nil
	case ReqDeplacer, ReqReculer:
		// un seul pion concerné : il est déjà choisi
		ids := map[int]bool{}
		for _, m := range r.Mvts {
			ids[m.Pion] = true
		}
		if len(ids) == 1 || r.Type == ReqReculer {
			g.sel = r.Mvts[0].Pion
		}
		if r.Type == ReqReculer && g.zoomee() {
			// le contact peut avoir lieu hors de la vue
			pt := jeu.C.Points[r.Mvts[0].Chemin[0]]
			if sx, sy := g.cam.ecran(pt.X, pt.Y); sx < 80 || sy < 80 || sx > LargVue-80 || sy > HauteurEcran-80 {
				g.cam.centrer(pt.X, pt.Y)
			}
		}
	}
	g.vueCouleur = 0
	g.preparerCapture()
}

// entree : la touche Entrée valide ou termine l'étape en cours.
func (g *Game) entree() {
	switch g.p.Req.Type {
	case ReqDeplacer:
		g.p.repondre(-1)
	case ReqTirer:
		g.p.repondre(-1)
	case ReqPlacer:
		g.validerPose()
	case ReqSecrets:
		if len(g.selSec) == g.p.Req.N {
			g.p.repondre(append([]int(nil), g.selSec...))
		}
	}
}

func (g *Game) cameraClavierSouris() {
	sx, sy := float64(g.mx), float64(g.my)
	if _, dy := ebiten.Wheel(); dy != 0 {
		if sx < LargVue && !g.scores {
			g.cam.zoomer(1+0.12*dy, sx, sy)
		} else if sx >= LargVue {
			g.defil = max(g.defil+int(dy), 0)
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		if sx < LargVue {
			g.glisse, g.gx, g.gy, g.gcx, g.gcy = true, sx, sy, g.cam.X, g.cam.Y
		}
	}
	if g.glisse {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
			g.cam.X = g.gcx - (sx-g.gx)/g.cam.Z
			g.cam.Y = g.gcy - (sy-g.gy)/g.cam.Z
			g.cam.borner()
		} else {
			g.glisse = false
		}
	}
	v := 14 / g.cam.Z
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.cam.X -= v
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.cam.X += v
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.cam.Y -= v
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.cam.Y += v
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKPAdd) || inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		g.cam.zoomer(1.25, LargVue/2, HauteurEcran/2)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKPSubtract) || inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		g.cam.zoomer(1/1.25, LargVue/2, HauteurEcran/2)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit0) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad0) || inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		g.cam.toutVoir()
	}
	g.cam.borner()
}

func (g *Game) Draw(dst *ebiten.Image) {
	debut := time.Now()
	defer func() { tempsDessin += time.Since(debut) }()
	if g.p != nil {
		g.p.mu.Lock()
		defer g.p.mu.Unlock()
	}
	dst.Fill(coulFond)
	switch g.ecran {
	case EcranAccueil:
		g.dessinerAccueil(dst)
	case EcranJeu:
		g.dessinerCarte(dst)
		g.dessinerMiniCarte(dst)
		if g.p.E.Fini || g.scores {
			g.dessinerScores(dst)
		}
		g.dessinerPanneau(dst)
	}
	if g.regles {
		g.dessinerRegles(dst)
	}
	for _, b := range g.boutons {
		dessinerBouton(dst, b, b.Actif && b.Contient(g.mx, g.my))
	}
	g.capturer(dst)
}

func (g *Game) bouton(x, y, w, h float64, t string, actif bool, f func()) *Bouton {
	b := &Bouton{X: x, Y: y, W: w, H: h, Texte: t, Actif: actif, Action: f}
	g.boutons = append(g.boutons, b)
	return b
}

func (g *Game) construireBoutons() {
	g.boutons = g.boutons[:0]
	if g.regles {
		g.boutonsRegles()
		return
	}
	switch g.ecran {
	case EcranAccueil:
		g.boutonsAccueil()
	case EcranJeu:
		g.boutonsJeu()
	}
}

var optPleinEcran = flag.Bool("fullscreen", false, "démarre en plein écran (F11 bascule)")

func main() {
	flag.Parse()
	chargerRessources()
	g := nouveauJeu()
	g.preparer()
	ebiten.SetWindowSize(LargeurEcran, HauteurEcran)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Cortège officiel")
	ebiten.SetFullscreen(*optPleinEcran)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
