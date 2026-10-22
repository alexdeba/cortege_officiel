package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"cortegeofficiel/jeu"
)

// Options de test : capture d'écran et mise en place directe d'une partie.
var (
	optCapture = flag.String("shot", "", "enregistre l'écran dans ce fichier PNG puis quitte")
	optImages  = flag.Int("frames", 30, "nombre d'images avant la capture")
	optEcran   = flag.String("ecran", "", "jeu : saute l'accueil")
	optJoueurs = flag.Int("joueurs", 4, "nombre de joueurs (2 à 4)")
	optCamp    = flag.Int("camp", 1, "votre place (1 : police…)")
	optGraine  = flag.Uint64("graine", 0, "graine du hasard (0 = aléatoire)")
	optRegle   = flag.Int("regle", 0, "afficher cette page de la règle")
	optArret   = flag.String("arret", "", "captures : l'ordinateur joue pour vous jusqu'à cette demande (secrets, couleur, placer, avancer, deplacer, tirer, reculer, fin ; deplacer3 : la troisième)")
	optSel     = flag.Int("sel", -1, "captures : pion choisi (ID) ; -2 : le premier pion mobile")
	optSurvol  = flag.Int("survol", -1, "captures : souris sur ce point")
	optZoom    = flag.Float64("zoom", 0, "captures : zoom de la carte")
	optCentre  = flag.Int("centre", -1, "captures : vue centrée sur ce point (-2 : la voiture)")
	optScores  = flag.Bool("scores", false, "captures : la table des scores")
	optProfil  = flag.String("cpuprofile", "", "profil CPU dans ce fichier")
	optPerf    = flag.Bool("perf", false, "mesure le temps de dessin (moyenne sur les images)")
)

func (g *Game) preparer() {
	if *optProfil != "" {
		f, err := os.Create(*optProfil)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
	g.nbJ = min(max(*optJoueurs, 2), 4)
	g.place = min(max(*optCamp-1, 0), g.nbJ-1)
	if *optRegle > 0 {
		g.regles = true
		preparerRegle()
		g.pageRgl = min(*optRegle, len(pagesRegle)) - 1
	}
	if *optEcran != "jeu" {
		return
	}
	g.commencer(*optGraine)
	if *optArret != "" {
		nom, n := strings.TrimRight(*optArret, "0123456789"), 1
		fmt.Sscan(strings.TrimPrefix(*optArret, nom), &n)
		g.p.auto = true
		g.p.arret, g.p.arretN = nomsReq[nom], n
		if nom == "fin" {
			g.p.arret = -1
		}
	}
	g.scores = *optScores
	if *optZoom > 0 {
		g.cam.Z = *optZoom
		g.cam.borner()
	}
	g.p.lancer()
}

// preparerCapture : à chaque nouvelle demande, les options de sélection.
func (g *Game) preparerCapture() {
	if *optCapture == "" || g.p.auto {
		return
	}
	if *optSel == -2 {
		for _, m := range g.p.Req.Mvts {
			g.sel = m.Pion
			break
		}
		for _, t := range g.p.Req.Tirs {
			g.sel = t.Tireur
			break
		}
		if g.p.Req.Type == ReqPlacer && len(g.pl) > 0 {
			g.sel = g.pl[0].Pion
		}
	} else if *optSel >= 0 {
		g.sel = *optSel
	}
	if *optCentre >= 0 {
		pt := jeu.C.Points[*optCentre]
		g.cam.centrer(pt.X, pt.Y)
	} else if *optCentre == -2 && g.p.E.Voiture.Pos >= 0 {
		pt := jeu.C.Points[g.p.E.Voiture.Pos]
		g.cam.centrer(pt.X, pt.Y)
	}
}

// survolCapture : option -survol, la souris posée sur un point.
func (g *Game) survolCapture() {
	if *optSurvol >= 0 && *optCapture != "" {
		g.survol = *optSurvol
	}
}

var tempsDessin time.Duration

// capturer : après -frames images et une fois la demande attendue posée
// (ou la partie finie), enregistre l'écran.
func (g *Game) capturer(dst *ebiten.Image) {
	if *optCapture == "" {
		return
	}
	g.images++
	if g.p != nil && g.p.auto && !g.p.E.Fini {
		if g.images > 60000 {
			log.Fatal("capture : demande jamais atteinte")
		}
		g.images = min(g.images, 2)
		return
	}
	if g.p != nil && *optEcran == "jeu" && g.p.Req.Type == ReqAucune && !g.p.E.Fini && g.images < 4000 {
		g.images = min(g.images, *optImages-1)
		return
	}
	if g.images < *optImages {
		return
	}
	b := dst.Bounds()
	img := image.NewRGBA(b)
	dst.ReadPixels(img.Pix)
	f, err := os.Create(*optCapture)
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(f, img)
	f.Close()
	pprof.StopCPUProfile()
	if *optPerf {
		fmt.Printf("dessin : %.2f ms par image en moyenne (CPU) ; %.1f images/s mesurées\n",
			float64(tempsDessin.Microseconds())/1000/float64(g.images), ebiten.ActualFPS())
	}
	os.Exit(0)
}
