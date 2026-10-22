package main

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"cortegeofficiel/jeu"
)

// La règle, consultable en jeu (R) : le texte de la règle (regle/regle.txt,
// transcrit des pages 64, 73, 74 et 75 du magazine), ses deux tables, et
// les choix faits pour l'adaptation.

type pageRegle struct {
	titre   string
	section string // titre du sommaire
	lignes  []string
	special string // "tirs", "scores"
}

const rX, rY, rW, rH = 70.0, 30.0, 1460.0, 940.0

var pagesRegle []pageRegle

// sections : le texte découpé à chaque « ## ».
func sections() [][2]string {
	b, _ := fichiers.ReadFile("regle/regle.txt")
	var res [][2]string
	titre := ""
	var corps []string
	for _, l := range strings.Split(string(b), "\n") {
		switch {
		case strings.HasPrefix(l, "# "):
		case strings.HasPrefix(l, "## "):
			if titre != "" || len(corps) > 0 {
				res = append(res, [2]string{titre, strings.Join(corps, "\n")})
			}
			titre, corps = strings.TrimPrefix(l, "## "), nil
		default:
			corps = append(corps, l)
		}
	}
	res = append(res, [2]string{titre, strings.Join(corps, "\n")})
	if res[0][0] == "" {
		res[0][0] = "Présentation"
	}
	return res
}

// preparerRegle : les pages, coupées à la hauteur de la fenêtre.
func preparerRegle() {
	w := rW - 330 - 50
	f := face(17, policeTexte)
	hr := rH - 200.0
	maxH := hr
	ajoute := func(titre, texte string) {
		var lignes []string
		for _, para := range strings.Split(strings.TrimSpace(texte), "\n") {
			if strings.TrimSpace(para) == "" {
				continue
			}
			lignes = append(lignes, couper(para, f, w)...)
			lignes = append(lignes, "") // un demi-interligne entre paragraphes
		}
		n := 0
		for len(lignes) > 0 {
			k, h := 0, 0.0
			for k < len(lignes) {
				dh := 24.0
				if lignes[k] == "" {
					dh = 9
				}
				if h+dh > maxH {
					break
				}
				h += dh
				k++
			}
			t := titre
			if n > 0 {
				t += " (suite)"
			}
			pagesRegle = append(pagesRegle, pageRegle{titre: t, section: titre, lignes: lignes[:k]})
			lignes = lignes[k:]
			n++
		}
	}
	for _, s := range sections() {
		switch s[0] {
		case "Table des tirs":
			pagesRegle = append(pagesRegle, pageRegle{titre: "Table des tirs", section: s[0], special: "tirs"})
		case "Table des scores":
			pagesRegle = append(pagesRegle, pageRegle{titre: "Table des scores", section: s[0], special: "scores"})
		default:
			ajoute(s[0], s[1])
		}
	}
	ajoute("Choix de l'adaptation", choixAdaptation)
	ajoute("Commandes", commandes)
}

const commandes = `Clic gauche : choisir un pion surligné, puis sa destination (ou sa cible). Un second clic sur le même point passe au pion suivant.
Entrée : terminer la phase en cours (ou ne pas tirer).
Barre d'espace : zoom sur la carte (là où est la souris) ou retour à la vue d'ensemble ; une mini-carte apparaît en haut à droite, cliquez ou glissez dedans pour vous déplacer.
Molette : zoom ; clic droit glissé ou flèches : déplacer la vue ; 0 : vue d'ensemble.
R : la règle ; S : la table des scores ; A : accélérer l'ordinateur ; F11 : plein écran.
Survolez un point pour voir sa zone, sa nature et les pions qui s'y trouvent. Le prochain point de la voiture présidentielle est cerclé de blanc.
Sur les pions : « T » et « A » marquent les terroristes et agents secrets que vous connaissez (pastille pleine : démasqués devant tous) ; un trait rouge, un pion blessé ; un chiffre, les passagers d'un véhicule.`

func (g *Game) touchesRegle() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.regles = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		g.pageRgl = min(g.pageRgl+1, len(pagesRegle)-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		g.pageRgl = max(g.pageRgl-1, 0)
	}
}

func (g *Game) boutonsRegles() {
	if pagesRegle == nil {
		preparerRegle()
	}
	y := rY + 110
	vu := map[string]bool{}
	for i, p := range pagesRegle {
		if vu[p.section] {
			continue
		}
		vu[p.section] = true
		i := i
		b := g.bouton(rX+24, y, 290, 46, p.section, true, func() { g.pageRgl = i })
		b.Style, b.Taille = BoutonDiscret, 19
		b.Choisi = pagesRegle[g.pageRgl].section == p.section
		y += 52
	}
	b := g.bouton(rX+24, rY+rH-66, 290, 48, "Fermer (R)", true, func() { g.regles = false })
	b.Style = BoutonPrincipal
	g.bouton(rX+rW-200, rY+rH-66, 80, 48, "←", g.pageRgl > 0, func() { g.pageRgl-- }).Style = BoutonDiscret
	g.bouton(rX+rW-110, rY+rH-66, 80, 48, "→", g.pageRgl < len(pagesRegle)-1, func() { g.pageRgl++ }).Style = BoutonDiscret
}

func (g *Game) dessinerRegles(dst *ebiten.Image) {
	// le clic sur « Règle » ouvre la règle après la construction des boutons :
	// Draw peut passer avant boutonsRegles.
	if pagesRegle == nil {
		preparerRegle()
	}
	rect(dst, 0, 0, LargeurEcran, HauteurEcran, alpha(coulNoir, 0.7))
	feuilleDossier(dst, rX, rY, rW, rH)
	ecrire(dst, "RÈGLE DU JEU", rX+28, rY+36, face(38, policeTitre), coulEncre)
	ecrire(dst, "PIÈCE N° 1 · À DÉTRUIRE APRÈS LECTURE", rX+30, rY+92, face(13, policeMachine), coulEncre2)
	p := pagesRegle[g.pageRgl]
	x := rX + 340
	w := rW - 340 - 40
	ecrire(dst, strings.ToUpper(p.titre), x, rY+44, face(28, policeTitre), coulSignal)
	ecrireDroite(dst, fmt.Sprintf("page %d / %d · ← →", g.pageRgl+1, len(pagesRegle)), rX+rW-36, rY+60, face(15, policeItalique), coulEncre2)
	trait(dst, x-18, rY+110, x-18, rY+rH-30, 1.5, coulDossier3)
	y := rY + 110
	switch p.special {
	case "tirs":
		g.tableTirs(dst, x, y, w)
		return
	case "scores":
		g.tableScores(dst, x, y, w)
		return
	}
	f := face(17, policeTexte)
	for _, l := range p.lignes {
		if l == "" {
			y += 9
			continue
		}
		ecrire(dst, l, x, y, f, coulEncre)
		y += 24
	}
}

func (g *Game) tableTirs(dst *ebiten.Image, x, y, w float64) {
	ecrire(dst, "Résultat de deux dés. — : raté · B : blessé · T : tué · BV : bavure", x, y, face(16, policeGras), coulEncre2)
	y += 40
	cx := x + 420
	pas := (w - 420) / 11
	for s := 2; s <= 12; s++ {
		ecrireCentre(dst, fmt.Sprint(s), cx+float64(s-2)*pas+pas/2, y, face(20, policeTitre), coulSignal)
	}
	y += 34
	groupes := []string{"Le terroriste tire", "L'agent secret tire", "La police tire"}
	lib := []string{"d'un point non colorié", "d'un point colorié", "à partir d'un véhicule",
		"vers un point quelconque", "d'un véhicule", "de l'hélicoptère",
		"d'un point quelconque", "d'un véhicule", "de l'hélicoptère"}
	sym := map[jeu.Resultat]string{jeu.Rate: "—", jeu.Blesse: "B", jeu.Tue: "T", jeu.Bavure: "BV"}
	for k, ligne := range jeu.TableTirs {
		if k%3 == 0 {
			y += 8
			ecrire(dst, strings.ToUpper(groupes[k/3]), x, y, face(19, policeTitre), coulEncre)
			y += 32
		}
		if k%2 == 0 {
			rect(dst, x-6, y-3, w, 30, coulDossier2)
		}
		ecrire(dst, lib[k], x+16, y, face(17, policeTexte), coulEncre)
		for i, r := range ligne {
			c := coulEncre
			if r == jeu.Tue {
				c = coulRouge
			}
			ecrireCentre(dst, sym[r], cx+float64(i)*pas+pas/2, y, face(17, policeGras), c)
		}
		y += 30
	}
}

func (g *Game) tableScores(dst *ebiten.Image, x, y, w float64) {
	cols := []jeu.Role{jeu.RPolice, jeu.RAgent, jeu.RManif, jeu.RTerro}
	cx := x + w - 4*130
	for k, r := range cols {
		ecrireCentre(dst, jeu.NomsRoles[r], cx+float64(k)*130+65, y, face(17, policeGras), coulRoles[r])
	}
	y += 40
	lignes := []struct {
		t string
		p [jeu.NbRoles]int
	}{
		{"Pour chaque blocage du cortège officiel", jeu.ScBlocage},
		{"Le cortège n'a jamais été bloqué", jeu.ScJamaisBloque},
		{"Pour chaque bavure", jeu.ScBavure},
		{"Pour un terroriste abattu avant un attentat", jeu.ScTerroAvant},
		{"Président blessé", jeu.ScPresBlesse},
		{"Président tué", jeu.ScPresTue},
		{"Pour un terroriste abattu après tentative d'attentat", jeu.ScTerroApres},
		{"Président indemne à l'aéroport", jeu.ScPresAeroport},
		{"Après une tentative d'attentat les 2 terroristes s'échappent", jeu.ScDeuxEchappent},
		{"Un seul s'échappe", jeu.ScUnSeulEchappe},
	}
	for i, l := range lignes {
		if i%2 == 0 {
			rect(dst, x-6, y-4, w, 40, coulDossier2)
		}
		paragraphe(dst, l.t, x+10, y, cx-x-30, face(17, policeTexte), coulEncre, 20)
		for k, r := range cols {
			s := "—"
			if v := l.p[r]; v != 0 {
				s = fmt.Sprintf("%+d", v)
			}
			ecrireCentre(dst, s, cx+float64(k)*130+65, y+4, face(18, policeGras), coulEncre)
		}
		y += 46
	}
}
