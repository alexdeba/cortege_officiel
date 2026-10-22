package jeu

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type journal struct {
	lignes []string
	tirs   int
}

func (o *journal) Evenement(e *Etat, ev Ev) {
	switch ev.Type {
	case EvTexte, EvScore, EvRevele, EvBlocage, EvTir:
		if ev.Type == EvTir {
			o.tirs++
		}
		o.lignes = append(o.lignes, fmt.Sprintf("T%d %s %s", ev.Tour, NomsPhases[e.Phase], ev.Texte))
	}
}

func partieIA(nbJ int, graine uint64) (*Etat, *journal) {
	noms := []string{"A", "B", "C", "D"}
	e := NouvelEtat(nbJ, noms[:nbJ], graine)
	for _, j := range e.Joueurs {
		j.Ctrl = IA{}
	}
	o := &journal{}
	e.Obs = o
	e.Jouer()
	return e, o
}

// verifier : les invariants de la règle à la fin de la partie.
func verifier(t *testing.T, e *Etat) {
	t.Helper()
	for i := range C.Points {
		if n := e.Compte(i, Policier); n > 2 {
			t.Errorf("%d policiers sur le point %d", n, i)
		}
		if n := len(e.ManifsSur(i)); n > 2 {
			t.Errorf("%d manifestants sur le point %d", n, i)
		}
		if e.PolicePresente(i) && len(e.ManifsSur(i)) > 0 {
			t.Errorf("police et manifestants ensemble au point %d", i)
		}
	}
	for _, p := range e.Pions {
		if p.Genre == Manifestant && p.Actif() && p.Pos >= 0 && C.Points[p.Pos].GP {
			t.Errorf("le %s est entré dans la Grand-Place", p.Nom())
		}
	}
}

func TestPartiesIA(t *testing.T) {
	for _, nbJ := range []int{4, 3, 2} {
		for g := uint64(1); g <= 6; g++ {
			e, o := partieIA(nbJ, g)
			verifier(t, e)
			if os.Getenv("VERBEUX") != "" && g == 1 {
				t.Log(strings.Join(o.lignes, "\n"))
			}
			t.Logf("%d joueurs, graine %d : %d tours, blocages %d, bavures %d, tirs %d, attentat %v, blessé %v, tué %v, arrivé %v, scores %v",
				nbJ, g, e.Tour, e.Blocages, e.Bavures, o.tirs, e.Attentat, e.PresBlesse, e.PresTue, e.Arrivee, e.Score)
			if !e.Fini {
				t.Errorf("partie non finie")
			}
		}
	}
}
