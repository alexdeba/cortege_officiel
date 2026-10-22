package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type StyleBouton int

const (
	BoutonNormal    StyleBouton = iota
	BoutonPrincipal             // l'action attendue
	BoutonDiscret               // navigation, fermeture
)

type Bouton struct {
	X, Y, W, H float64
	Texte      string
	Sous       string // seconde ligne, plus petite
	Actif      bool
	Action     func()
	Style      StyleBouton
	Taille     float64
	Choisi     bool
	De         int // > 0 : un dé dessiné au lieu du texte
	Pion       int // AvecPion : le pion dessiné au lieu du texte
	AvecPion   bool
	Fond       *color.RGBA
}

func (b *Bouton) Contient(x, y int) bool {
	fx, fy := float64(x), float64(y)
	return fx >= b.X && fy >= b.Y && fx < b.X+b.W && fy < b.Y+b.H
}

func dessinerBouton(dst *ebiten.Image, b *Bouton, survol bool) {
	var fond, bord, encre color.RGBA
	switch b.Style {
	case BoutonPrincipal:
		fond, bord, encre = coulSignal, coulSignal2, coulCraie
		if survol {
			fond = coulSignal2
		}
	case BoutonDiscret:
		fond, bord, encre = coulAcier, coulAcier3, coulCraie
		if survol {
			fond = coulAcier2
		}
	default:
		fond, bord, encre = coulAcier2, coulAcier3, coulCraie
		if survol {
			fond, bord = coulAcier3, coulAmbre
		}
	}
	if b.Fond != nil {
		fond = *b.Fond
		if 299*int(fond.R)+587*int(fond.G)+114*int(fond.B) < 150000 {
			encre = coulCraie
		} else {
			encre = coulEncre
		}
		bord = melange(fond, coulNoir, 0.4)
		if survol {
			fond, bord = melange(fond, coulBlanc, 0.2), coulCraie
		}
	}
	if b.Choisi {
		fond, bord, encre = coulCraie, coulCraie, coulEncre
	}
	if !b.Actif {
		fond, bord, encre = melange(fond, coulAcier, 0.6), melange(bord, coulAcier, 0.6), melange(encre, coulAcier, 0.6)
	}
	if b.Actif {
		rect(dst, b.X+3, b.Y+3, b.W, b.H, alpha(coulNoir, 0.5))
	}
	arrondi(dst, b.X, b.Y, b.W, b.H, 3, fond, bord, 1.5)
	// le liseré rouge de l'option choisie
	if b.Choisi {
		rect(dst, b.X, b.Y+b.H-5, b.W, 5, coulSignal)
	}
	if b.De > 0 {
		s := b.H - 12
		de(dst, b.X+(b.W-s)/2, b.Y+6, s, b.De, coulBlanc, coulEncre)
		return
	}
	if b.Pion >= 0 && b.AvecPion {
		pion(dst, b.Pion, b.X+b.W/2, b.Y+b.H/2, b.H-8, 1)
		if b.Choisi {
			cadre(dst, b.X+3, b.Y+3, b.W-6, b.H-6, 3, coulSignal)
		}
		return
	}
	t, police := tailleTexte(b)
	f := face(t, police)
	m := f.Metrics()
	hl := m.HAscent + m.HDescent // hauteur de la ligne, pour centrer
	if b.Sous == "" {
		ecrireCentre(dst, b.Texte, b.X+b.W/2, b.Y+(b.H-hl)/2, f, encre)
		return
	}
	ts := 15.0
	for ts > 13 && largeur(b.Sous, face(ts, policeTexte)) > b.W-18 {
		ts--
	}
	ecrireCentre(dst, b.Texte, b.X+b.W/2, b.Y+b.H*0.3-hl/2, f, encre)
	ecrireCentre(dst, b.Sous, b.X+b.W/2, b.Y+b.H-ts*1.35-3, face(ts, policeTexte), encre)
}

// tailleTexte : la taille voulue, réduite seulement si le texte déborde du
// bouton.
func tailleTexte(b *Bouton) (float64, *Police) {
	t := b.Taille
	if t == 0 {
		t = 21
	}
	police := policeTitre
	if b.Style == BoutonDiscret {
		police = policeGras
	}
	hauteur := b.H - 6
	if b.Sous != "" {
		hauteur = b.H * 0.5
	}
	for t > 14 && (largeur(b.Texte, face(t, police)) > b.W-18 || t*1.2 > hauteur) {
		t--
	}
	return t, police
}
