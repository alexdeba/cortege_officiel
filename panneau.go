package main

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"cortegeofficiel/jeu"
)

// Le panneau de droite : le camp du joueur, la consigne de l'étape en cours
// et ses boutons, les scores, et le télex des événements (journal).

const margeP = 18.0

var (
	yConsigne, hConsigne, yScores, yJournal float64
)

// nomCamp : les rôles d'un joueur, en clair.
func nomCamp(j *jeu.Joueur) string {
	var n []string
	for _, r := range j.Roles {
		n = append(n, jeu.NomsRoles[r])
	}
	return strings.Join(n, " + ")
}

func roleDeTable(t int) jeu.Role {
	switch {
	case t <= jeu.TTerroVehicule:
		return jeu.RTerro
	case t <= jeu.TAgentHelico:
		return jeu.RAgent
	}
	return jeu.RPolice
}

// consigne : titre et texte de l'étape en cours.
func (g *Game) consigne() (string, string) {
	e := g.p.E
	r := g.p.Req
	if e.Fini {
		return "Fin de la partie", "Le cortège est passé. La table des scores est affichée sur la carte : le vainqueur est le joueur qui a le plus de points."
	}
	switch r.Type {
	case ReqSecrets:
		if r.Role == jeu.RAgent {
			t := "Choisissez en secret vos deux agents secrets parmi les manifestants de numéro pair. Personne d'autre ne les connaît : ils se mêlent à la foule, prêts à abattre les terroristes."
			if e.NbJ == 2 {
				t = "Choisissez en secret quatre agents secrets de numéro pair : deux manifestants et deux policiers (des agents déguisés)."
			}
			return "Vos agents secrets", t
		}
		t := "Choisissez en secret vos deux terroristes parmi les manifestants de numéro impair. Ils se cachent dans la foule ; leur but : abattre le président, puis quitter la ville."
		if e.NbJ == 2 {
			t = "Choisissez en secret quatre terroristes de numéro impair : deux manifestants et deux policiers (des terroristes déguisés)."
		}
		return "Vos terroristes", t
	case ReqCouleur:
		if r.Retour {
			return "Le parcours du retour", "Le président parle à la tribune. Choisissez en secret le parcours qui mènera le cortège de la Grand-Place à l'aéroport : il ne sera dévoilé qu'au départ du point 3. Survolez un bouton pour voir le trajet."
		}
		return "Le parcours du cortège", "Choisissez la couleur des points que suivra la voiture présidentielle de l'Hôtel de Ville (G1) à la Grand-Place. Votre feuille sera montrée au premier tour. Survolez un bouton pour voir le trajet."
	case ReqPlacer:
		if r.Role == jeu.RPolice {
			return "La pose de la police", "Une pose vous est proposée : cliquez un pion puis un point entouré pour le déplacer. Policiers sur des points de rue non coloriés, deux au plus par point ; voitures et motos à côté des policiers ou avec eux à bord. L'hélicoptère attend sur le toit de l'aéroport. Validez quand c'est prêt."
		}
		return "La pose des manifestants", "Une pose vous est proposée : cliquez un manifestant puis un point entouré pour le déplacer. Points de rue non coloriés, hors de la Grand-Place, jamais avec la police, deux au plus par point. Le parcours du cortège est encore secret !"
	case ReqAvancer:
		return "a. La voiture présidentielle", "Plusieurs points de votre parcours mènent vers le but : cliquez celui où avance la voiture (entouré)."
	case ReqDeplacer:
		return g.consigneDeplacement()
	case ReqTirer:
		return g.consigneTir()
	case ReqReculer:
		return "Contact avec la police !", "Des policiers arrivent sur vos manifestants : choisissez où recule chacun (points entourés ; d'abord le pion s'il y en a plusieurs). Jamais plus de deux manifestants par point. Sans issue, un manifestant est interpellé."
	}
	if e.Phase == jeu.PhPreparation {
		return "Préparation", "Les autres joueurs choisissent leurs pions secrets et posent leurs pions…"
	}
	j := g.p.joueurPhase(e)
	qui := "L'ordinateur"
	if j != nil && !j.Humain {
		qui = j.Nom + " (" + nomCamp(j) + ")"
	}
	return "Tour " + fmt.Sprint(e.Tour) + " · " + jeu.LettresPhases[e.Phase] + ". " + jeu.NomsPhases[e.Phase], qui + " joue. Suivez la carte et les dépêches ci-dessous (touche A : accélérer)."
}

func (g *Game) consigneDeplacement() (string, string) {
	e := g.p.E
	q := g.p.Req.Quota
	n := e.NbMvts
	switch {
	case q.Role == jeu.RPolice:
		return "e. Déplacement de la police", fmt.Sprintf("Chaque policier avance d'un point (cliquez un pion surligné, puis sa destination). Arriver sur des manifestants les fait reculer, ou les interpelle s'ils sont cernés. Voitures (2 points) et motos (3) roulent avec leurs policiers ; l'hélicoptère vole de zone en zone. Gardez libre le prochain point de la voiture (cerclé de blanc). %d mouvement(s) joué(s) ; Entrée pour finir.", n)
	case q.Role == jeu.RManif && q.Policiers:
		return "Provocation !", fmt.Sprintf("Jets de pierres : vous pouvez déplacer deux policiers d'un point (pas sur le parcours du cortège). %d/2. Entrée pour finir.", n)
	case q.Role == jeu.RManif:
		return "f. Les manifestants", fmt.Sprintf("Déplacez vos manifestants d'un point chacun, vers un point libre ou tenu par un autre manifestant, jamais vers la police. Un manifestant sur le prochain point de la voiture (cerclé de blanc) la bloque ! Les barrières de la Grand-Place ne se franchissent pas ; un toit se gagne après un tour sur un point fléché. %d joué(s) ; Entrée pour finir.", n)
	case q.Role == jeu.RAgent:
		return "g. Les agents secrets", fmt.Sprintf("Déplacez jusqu'à quatre manifestants de numéro pair : vos agents (marqués A) ou d'autres, pour brouiller les pistes. Approchez-vous des suspects : ensuite, vous pourrez tirer. %d/4 ; Entrée pour finir.", n)
	}
	t := "Déplacez jusqu'à quatre manifestants de numéro impair : vos terroristes (marqués T) ou d'autres, pour brouiller les pistes. Gagnez un point relié par un trait à la voiture, ou un toit de la même zone : vous pourrez tirer."
	switch e.NbJ {
	case 3:
		t = "Déplacez jusqu'à quatre manifestants et quatre policiers (faux ordres par talkie-walkie). Placez vos terroristes (T) à portée de la voiture : un point relié par un trait, ou un toit de la même zone."
	case 2:
		t = "Déplacez jusqu'à quatre pions de numéro impair, manifestants ou policiers : vos terroristes (T) ou d'autres. Placez-les à portée de la voiture."
	}
	if e.Attentat {
		t += " Un terroriste démasqué sur un point du bord du plan peut quitter la ville (point rouge)."
	}
	return "h. Les terroristes", fmt.Sprintf("%s %d joué(s) ; Entrée pour finir.", t, n)
}

func (g *Game) consigneTir() (string, string) {
	r := g.p.Req
	e := g.p.E
	titre := "Tirs"
	t := ""
	switch r.Phase {
	case jeu.PhAttentat, jeu.PhTerros:
		titre = "« Je tire ! »"
		t = "Un de vos terroristes est à portée du président. Cliquez le tireur surligné, puis la cible : le tireur sera démasqué, et la riposte suivra aussitôt. Vous pouvez aussi renoncer."
	case jeu.PhRiposte:
		titre = "c. Riposte !"
		t = "Un terroriste vient de tirer : il est démasqué. Cliquez un de vos tireurs surlignés, puis le terroriste."
	case jeu.PhArrestation:
		titre = "d. Tentative d'arrestation"
		t = "Si vous croyez avoir localisé un terroriste, tirez : cliquez un tireur surligné, puis sa cible. Attention : toucher un innocent est une bavure (police), et un agent qui tire est démasqué."
	case jeu.PhPolice:
		titre = "Tirs de la police"
		t = "Un terroriste démasqué est à portée de vos policiers : cliquez un tireur, puis le terroriste."
	case jeu.PhAgents:
		titre = "Tirs des agents secrets"
		t = "Vos agents peuvent tirer sur un pion voisin relié par un trait (d'un toit : toute la zone). L'agent qui tire est démasqué."
	}
	if g.sel >= 0 {
		for _, tr := range r.Tirs {
			if tr.Tireur == g.sel {
				ch := jeu.Chances(tr.Table)
				t += fmt.Sprintf("\nLe %s tire (%s) : touché %d %%, tué %d %%, bavure %d %%.", e.Pions[g.sel].Nom(), jeu.NomsTables[tr.Table],
					int(100*(ch[jeu.Blesse]+ch[jeu.Tue])+.5), int(100*ch[jeu.Tue]+.5), int(100*ch[jeu.Bavure]+.5))
				break
			}
		}
	}
	return titre, t + " Entrée : ne pas tirer."
}

// boutonsJeu : les boutons du panneau selon la demande.
func (g *Game) boutonsJeu() {
	x0, w := PanneauX+margeP, LargPanneau-2*margeP
	yb := float64(HauteurEcran) - 54
	lb := (w - 24) / 4
	g.bouton(x0, yb, lb, 42, "Règle", true, func() { g.regles = true }).Style = BoutonDiscret
	sb := g.bouton(x0+lb+8, yb, lb, 42, "Scores", true, func() { g.scores = !g.scores })
	sb.Style, sb.Choisi = BoutonDiscret, g.scores
	vt := "Lent"
	if g.p.Rapide {
		vt = "Rapide"
	}
	g.bouton(x0+2*(lb+8), yb, lb, 42, vt, true, func() { g.p.Rapide = !g.p.Rapide }).Style = BoutonDiscret
	g.bouton(x0+3*(lb+8), yb, lb, 42, "Plan", true, func() { g.cam.toutVoir() }).Style = BoutonDiscret

	if g.p.E.Fini {
		b := g.bouton(LargVue/2-160, HauteurEcran-96, 320, 54, "Nouvelle partie", true, func() {
			n, pl := g.nbJ, g.place
			*g = *nouveauJeu()
			g.nbJ, g.place = n, pl
		})
		b.Style = BoutonPrincipal
		return
	}

	// la consigne, puis les boutons de la demande
	_, texte := g.consigne()
	yConsigne = 168
	hConsigne = 46 + float64(len(couper(texte, face(17, policeTexte), w-28)))*23 + 14
	y := yConsigne + hConsigne + 12
	r := g.p.Req
	switch r.Type {
	case ReqSecrets:
		c := 6
		s := (w - float64(c-1)*6) / float64(c)
		for i, id := range r.Cands {
			id := id
			pris := false
			for _, x := range g.selSec {
				pris = pris || x == id
			}
			b := g.bouton(x0+float64(i%c)*(s+6), y+float64(i/c)*(s+6), s, s, "", true, func() {
				for k, x := range g.selSec {
					if x == id {
						g.selSec = append(g.selSec[:k], g.selSec[k+1:]...)
						return
					}
				}
				if len(g.selSec) < r.N {
					g.selSec = append(g.selSec, id)
				}
			})
			b.Pion, b.AvecPion, b.Choisi, b.Style = id, true, pris, BoutonDiscret
		}
		y += float64((len(r.Cands)+c-1)/c)*(s+6) + 8
		g.bouton(x0, y, w/2-6, 50, "Au hasard", true, func() {
			g.selSec = append([]int(nil), g.p.E.TirageSecret(r.Role)...)
		}).Style = BoutonDiscret
		b := g.bouton(x0+w/2+6, y, w/2-6, 50, fmt.Sprintf("Valider (%d/%d)", len(g.selSec), r.N), len(g.selSec) == r.N && g.p.E.SecretsValides(r.Role, g.selSec), func() {
			g.p.repondre(append([]int(nil), g.selSec...))
		})
		b.Style = BoutonPrincipal
	case ReqCouleur:
		g.vueCouleur = 0
		for i, col := range jeu.CouleursParcours {
			col := col
			c := coulParcours(col)
			b := g.bouton(x0, y+float64(i)*62, w, 54, "Parcours "+col.Nom(), true, func() { g.p.repondre(col) })
			dep := jeu.PtHotelDeVille
			if r.Retour {
				dep = jeu.PtGP3
			}
			b.Sous = fmt.Sprintf("%d points jusqu'%s", jeu.C.DistParcours(col, r.Retour, dep), map[bool]string{false: "à la Grand-Place", true: "à l'aéroport"}[r.Retour])
			b.Fond = &c
			if b.Contient(g.mx, g.my) {
				g.vueCouleur = col
			}
		}
	case ReqPlacer:
		g.bouton(x0, y, w/2-6, 50, "Pose proposée", true, func() {
			g.pl = append([]jeu.Placement(nil), r.Prop...)
			g.sel, g.plMsg = -1, ""
		}).Style = BoutonDiscret
		g.bouton(x0+w/2+6, y, w/2-6, 50, "Valider la pose", true, g.validerPose).Style = BoutonPrincipal
	case ReqDeplacer:
		if len(g.choixMvt) > 1 {
			for k, i := range g.choixMvt {
				i := i
				m := r.Mvts[i]
				t := map[jeu.TypeMvt]string{jeu.MvtPas: "Garder le véhicule (à côté)", jeu.MvtMonter: "Monter à bord", jeu.MvtConduite: "Rouler", jeu.MvtVol: "Voler", jeu.MvtDescente: "Descendre", jeu.MvtFuite: "Quitter la ville"}[m.Type]
				g.bouton(x0, y+float64(k)*58, w, 50, t, true, func() { g.jouerMouvement(i) })
			}
			y += float64(len(g.choixMvt)) * 58
		}
		b := g.bouton(x0, y, w, 50, "Terminer la phase (Entrée)", true, func() { g.p.repondre(-1) })
		b.Style = BoutonPrincipal
		// descendre d'un véhicule, fuir : pas de point à cliquer ailleurs
		if g.sel >= 0 {
			for _, i := range g.mouvementsDe(g.sel) {
				i := i
				m := r.Mvts[i]
				if m.Type == jeu.MvtFuite {
					g.bouton(x0, y+58, w, 50, "Quitter la ville", true, func() { g.jouerMouvement(i) }).Fond = &coulRouge
				}
			}
		}
	case ReqTirer:
		b := g.bouton(x0, y, w, 50, "Ne pas tirer (Entrée)", true, func() { g.p.repondre(-1) })
		b.Style = BoutonPrincipal
	}
}

func (g *Game) dessinerPanneau(dst *ebiten.Image) {
	x0 := PanneauX
	rect(dst, x0, 0, LargPanneau, HauteurEcran, coulAcier)
	rect(dst, x0, 0, 3, HauteurEcran, coulNoir)
	filet(dst, x0+3, 0, LargPanneau-3, 8, coulNoir)
	// le logo de la page 64, ombré de rouge comme le titre de la couverture
	logo(dst, x0+margeP+3, 15, 70, coulSignal)
	logo(dst, x0+margeP, 12, 70, coulCraie)
	xt := x0 + margeP + 70*float64(imgLogo.Bounds().Dx())/float64(imgLogo.Bounds().Dy()) + 22
	ecrire(dst, "DOSSIER 83-64", xt, 20, face(14, policeMachine), coulCendre)
	ecrire(dst, "OPÉRATION GRAND-PLACE", xt, 40, face(14, policeMachine), coulCendre)
	tampon(dst, "TRÈS SECRET", x0+LargPanneau-100, 70, 18, -0.08)
	e := g.p.E
	w := LargPanneau - 2*margeP

	// le camp du joueur
	moi := g.p.Moi
	y := 100.0
	arrondi(dst, x0+margeP, y, w, 56, 3, coulAcier2, coulAcier3, 1.5)
	for i, r := range moi.Roles { // une barre par rôle tenu
		hb := 56.0 / float64(len(moi.Roles))
		rect(dst, x0+margeP, y+float64(i)*hb, 8, hb, coulRoles[r])
	}
	ecrire(dst, "VOUS : "+strings.ToUpper(nomCamp(moi)), x0+margeP+22, y+2, face(22, policeTitre), coulCraie)
	etat := g.etatCortege()
	ecrire(dst, etat, x0+margeP+22, y+33, face(15, policeGras), coulAmbre2)

	// consigne
	if !e.Fini {
		titre, texte := g.consigne()
		y = yConsigne
		rect(dst, x0+margeP+4, y+4, w, hConsigne, alpha(coulNoir, 0.55))
		arrondi(dst, x0+margeP, y, w, hConsigne, 2, coulDossier, coulDossier3, 1)
		ecrire(dst, titre, x0+margeP+14, y+3, face(24, policeTitre), coulEncre)
		trait(dst, x0+margeP+14, y+39, x0+margeP+w-14, y+39, 1.5, coulSignal)
		paragraphe(dst, texte, x0+margeP+14, y+46, w-28, face(17, policeTexte), coulEncre, 23)
		if g.p.Req.Type == ReqPlacer && g.plMsg != "" {
			ecrire(dst, "⚠ "+g.plMsg, x0+margeP+4, y+hConsigne+70, face(16, policeGras), coulSignal2)
		}
	}

	// le verdict, en travers du panneau vide
	if e.Fini {
		v := "AFFAIRE CLASSÉE"
		if e.PresTue {
			v = "ATTENTAT"
		}
		tampon(dst, v, x0+LargPanneau/2, 330, 44, -0.14)
	}

	// scores et journal, en bas
	yJournal = g.hautJournal()
	yS := yJournal - 140
	if yS > g.basBoutons()+10 {
		g.dessinerMiniScores(dst, x0+margeP, yS, w)
	} else {
		yS = yJournal
	}
	g.dessinerJournal(dst, x0+margeP, max(yJournal, g.basBoutons()+10), w, float64(HauteurEcran)-66)
}

// basBoutons : le bas des boutons de la demande.
func (g *Game) basBoutons() float64 {
	b := yConsigne + hConsigne
	for _, x := range g.boutons {
		if x.X >= PanneauX && x.Y < HauteurEcran-60 {
			b = max(b, x.Y+x.H)
		}
	}
	return b
}

func (g *Game) hautJournal() float64 { return float64(HauteurEcran) - 66 - 250 }

// etatCortege : où en est le cortège.
func (g *Game) etatCortege() string {
	e := g.p.E
	switch {
	case e.Phase == jeu.PhPreparation:
		return "Préparation de la partie"
	case e.PresTue:
		return fmt.Sprintf("Tour %d · le président a été tué", e.Tour)
	case e.Arrivee:
		return fmt.Sprintf("Tour %d · le président est à l'aéroport", e.Tour)
	case e.Discours >= 1 && e.Discours <= 2:
		return fmt.Sprintf("Tour %d · discours sur la Grand-Place", e.Tour)
	}
	parc := "parcours secret"
	if e.Discours >= 3 {
		if e.RetourPublic || g.p.Moi.A(jeu.RPolice) {
			parc = "retour " + e.CouleurRetour.Nom()
		} else {
			parc = "retour secret"
		}
	} else if e.AllerPublic || g.p.Moi.A(jeu.RPolice) {
		parc = "parcours " + e.CouleurAller.Nom()
	}
	return fmt.Sprintf("Tour %d/%d · %s · %d blocage(s)", e.Tour, e.TourMax, parc, e.Blocages)
}

func (g *Game) dessinerMiniScores(dst *ebiten.Image, x, y, w float64) {
	e := g.p.E
	arrondi(dst, x, y, w, 128, 3, coulAcier2, coulAcier3, 1)
	ecrire(dst, "SCORES", x+12, y+3, face(18, policeTitre), coulCraie)
	for i, j := range e.Joueurs {
		yy := y + 34 + float64(i)*22
		rect(dst, x+11, yy+3, 14, 14, coulCendre)
		rect(dst, x+12, yy+4, 12, 12, coulRoles[j.Roles[0]])
		nom := j.Nom + " · " + nomCamp(j)
		ecrire(dst, nom, x+32, yy, face(15, policeTexte), coulCraie)
		ecrireDroite(dst, fmt.Sprint(e.ScoreJoueur(j)), x+w-12, yy, face(16, policeGras), coulCraie)
	}
}

// dessinerJournal : les dépêches, les plus récentes en bas.
func (g *Game) dessinerJournal(dst *ebiten.Image, x, y, w, bas float64) {
	if bas-y < 60 {
		return
	}
	arrondi(dst, x, y, w, bas-y, 3, coulFond, coulAcier3, 1)
	ecrire(dst, "TÉLEX", x+12, y+3, face(17, policeTitre), coulAmbre)
	ecrireDroite(dst, "· · · réception · · ·", x+w-12, y+6, face(13, policeMachine), coulCendre)
	f := face(14, policeMachine)
	var lignes []LigneJournal
	for _, l := range g.p.Journal {
		for _, t := range couper(l.Texte, f, w-24) {
			lignes = append(lignes, LigneJournal{l.Tour, t, l.Coul})
		}
	}
	n := int((bas - y - 34) / 18)
	fin := max(len(lignes)-g.defil, 0)
	deb := max(fin-n, 0)
	for i, l := range lignes[deb:fin] {
		c := coulCraie
		switch l.Coul {
		case coulRouge:
			c = coulSignal2
		case coulSignal:
			c = coulAmbre
		case coulEncre2:
			c = coulCendre
		}
		ecrire(dst, l.Texte, x+12, y+30+float64(i)*18, f, c)
	}
}

// ---------------------------------------------------------------- scores

func (g *Game) dessinerScores(dst *ebiten.Image) {
	e := g.p.E
	x, y, w, h := 70.0, 60.0, LargVue-140, float64(HauteurEcran)-190
	rect(dst, 0, 0, LargVue, HauteurEcran, alpha(coulNoir, 0.55))
	feuilleDossier(dst, x, y, w, h)
	titre := "TABLE DES SCORES"
	if e.Fini {
		titre = "FIN DE LA PARTIE"
	}
	ecrire(dst, titre, x+28, y+34, face(36, policeTitre), coulEncre)
	tampon(dst, "RAPPORT", x+w-120, y+62, 22, 0.08)
	// colonnes : police, agent, manifestant, terroriste (comme la règle)
	cols := []jeu.Role{jeu.RPolice, jeu.RAgent, jeu.RManif, jeu.RTerro}
	cx := x + w - 4*110 - 20
	yy := y + 100
	for k, r := range cols {
		ecrireCentre(dst, strings.Split(jeu.NomsRoles[r], " ")[0], cx+float64(k)*110+55, yy, face(17, policeGras), coulRoles[r])
	}
	yy += 30
	f := face(16, policeTexte)
	lignes := e.Lignes
	nmax := int((h - 290) / 24)
	if len(lignes) > nmax {
		lignes = lignes[len(lignes)-nmax:]
	}
	for i, l := range lignes {
		if i%2 == 0 {
			rect(dst, x+20, yy-2, w-40, 24, coulDossier2)
		}
		ecrire(dst, fmt.Sprintf("Tour %d · %s", l.Tour, l.Texte), x+30, yy, f, coulEncre)
		for k, r := range cols {
			if v := l.Pts[r]; v != 0 {
				ecrireCentre(dst, fmt.Sprintf("%+d", v), cx+float64(k)*110+55, yy, face(16, policeGras), coulEncre)
			}
		}
		yy += 24
	}
	trait(dst, x+20, yy+4, x+w-20, yy+4, 2, coulEncre)
	yy += 12
	ecrire(dst, "Total par rôle", x+30, yy, face(17, policeGras), coulEncre)
	for k, r := range cols {
		ecrireCentre(dst, fmt.Sprint(e.Score[r]), cx+float64(k)*110+55, yy, face(18, policeTitre), coulRoles[r])
	}
	yy += 46
	// classement des joueurs
	ordre := append([]*jeu.Joueur(nil), e.Joueurs...)
	for i := range ordre {
		for k := i + 1; k < len(ordre); k++ {
			if e.ScoreJoueur(ordre[k]) > e.ScoreJoueur(ordre[i]) {
				ordre[i], ordre[k] = ordre[k], ordre[i]
			}
		}
	}
	for i, j := range ordre {
		t := fmt.Sprintf("%d. %s (%s) : %d points", i+1, j.Nom, nomCamp(j), e.ScoreJoueur(j))
		fc := face(19, policeGras)
		if i == 0 && e.Fini {
			fc = face(25, policeTitre)
			t = "Vainqueur : " + t[3:]
		}
		ecrire(dst, t, x+30, yy, fc, coulEncre)
		yy += 32
		if i == 0 && e.Fini {
			yy += 10
		}
	}
	if !e.Fini {
		ecrireDroite(dst, "S ou Échap : fermer", x+w-24, y+h-34, face(15, policeItalique), coulEncre2)
	}
}
