package main

import (
	"testing"
	"time"

	"cortegeofficiel/jeu"

	"github.com/hajimehoshi/ebiten/v2"
)

// Le joueur humain répond par l'interface, à chaque place et pour chaque
// nombre de joueurs : une partie entière se joue en répondant aux demandes
// comme le feraient les clics.
func TestHumainParDemandes(t *testing.T) {
	for nbJ := 2; nbJ <= 4; nbJ++ {
		for place := range nbJ {
			p := nouvellePartie(nbJ, place, uint64(10*nbJ+place))
			p.sansPause = true
			p.lancer()
			ia := jeu.IA{}
			limite := time.Now().Add(60 * time.Second)
			demandes := map[TypeReq]int{}
			for {
				p.mu.Lock()
				if p.Termine {
					p.mu.Unlock()
					break
				}
				e, moi, r := p.E, p.Moi, p.Req
				switch r.Type {
				case ReqSecrets:
					p.repondre(ia.Secrets(e, moi, r.Role, r.Cands, r.N))
				case ReqCouleur:
					p.repondre(ia.Couleur(e, moi, r.Retour))
				case ReqPlacer:
					if msg := e.ValiderPlacement(r.Role, r.Prop); msg != "" {
						t.Errorf("pose proposée illégale : %s", msg)
					}
					p.repondre(r.Prop)
				case ReqAvancer:
					p.repondre(0)
				case ReqDeplacer:
					p.repondre(ia.Deplacer(e, moi, r.Quota, r.Mvts))
				case ReqTirer:
					p.repondre(ia.Tirer(e, moi, r.Phase, r.Tirs))
				case ReqReculer:
					p.repondre(0)
				}
				demandes[r.Type]++
				p.mu.Unlock()
				if time.Now().After(limite) {
					t.Fatal("partie trop longue")
				}
				time.Sleep(20 * time.Microsecond)
			}
			if demandes[ReqDeplacer] == 0 {
				t.Errorf("%d joueurs, place %d : jamais de déplacement demandé", nbJ, place)
			}
			t.Logf("%d joueurs, %s : %d tours, %d points", nbJ, nomCamp(p.Moi), p.E.Tour, p.E.ScoreJoueur(p.Moi))
		}
	}
}

// Barre d'espace : zoom sur la voiture (souris hors de la carte) puis retour.
func TestZoomEspace(t *testing.T) {
	g := nouveauJeu()
	g.commencer(1)
	g.p.E.Voiture.Pos = jeu.PtHotelDeVille
	g.mx, g.my = -1, -1
	g.basculerZoom()
	if !g.zoomee() || g.cam.Z != zoomEspace {
		t.Fatalf("pas de zoom : %+v", g.cam)
	}
	hv := jeu.C.Points[jeu.PtHotelDeVille]
	sx, sy := g.cam.ecran(hv.X, hv.Y)
	if sx < 0 || sy < 0 || sx > LargVue || sy > HauteurEcran {
		t.Errorf("l'Hôtel de Ville hors de la vue : %v %v", sx, sy)
	}
	g.basculerZoom()
	if g.zoomee() {
		t.Fatalf("retour à la vue d'ensemble raté : %+v", g.cam)
	}
}

// Les boutons : texte de 17 px au moins une fois ajusté, et jamais plus
// large que le bouton.
func TestTexteBoutons(t *testing.T) {
	chargerRessources()
	g := nouveauJeu()
	verifier := func(quoi string) {
		g.construireBoutons()
		for _, b := range g.boutons {
			if b.AvecPion || b.Texte == "" {
				continue
			}
			taille, police := tailleTexte(b)
			if taille < 17 {
				t.Errorf("%s : « %s » écrit en %.0f px", quoi, b.Texte, taille)
			}
			if largeur(b.Texte, face(taille, police)) > b.W-12 {
				t.Errorf("%s : « %s » déborde", quoi, b.Texte)
			}
		}
	}
	for _, n := range []int{2, 3, 4} {
		g.nbJ = n
		verifier("accueil")
	}
	g.regles = true
	verifier("règle")
	g.regles = false
	g.commencer(5)
	reqs := []Demande{
		{Type: ReqCouleur},
		{Type: ReqPlacer, Role: jeu.RPolice},
		{Type: ReqSecrets, Role: jeu.RAgent, Cands: []int{31, 33}, N: 2},
		{Type: ReqDeplacer, Quota: jeu.Quota{Role: jeu.RPolice}},
		{Type: ReqTirer, Phase: jeu.PhRiposte},
	}
	for _, r := range reqs {
		g.p.Req = r
		verifier("partie")
	}
}

// Le clic sur « Règle » ouvre la règle sans que boutonsRegles soit passé :
// Draw doit préparer les pages lui-même.
func TestRegleDessineeAvantBoutons(t *testing.T) {
	g := nouveauJeu()
	pagesRegle = nil
	g.regles = true
	g.Draw(ebiten.NewImage(LargeurEcran, HauteurEcran))
	if len(pagesRegle) == 0 {
		t.Fatal("pas de pages de règle")
	}
}
