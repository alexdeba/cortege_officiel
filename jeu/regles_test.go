package jeu

import (
	"math"
	"testing"
	"time"
)

// etatVide : une partie à 4 joueurs de l'ordinateur, sans aucun pion posé.
func etatVide() *Etat {
	e := NouvelEtat(4, []string{"A", "B", "C", "D"}, 1)
	for _, j := range e.Joueurs {
		j.Ctrl = IA{}
	}
	e.Tour = 1
	e.Agents, e.Terroristes = []int{NbPoliciers + 1, NbPoliciers + 3}, []int{NbPoliciers + 0, NbPoliciers + 2}
	return e
}

// unPointLibre : un point de rue non colorié, hors Grand-Place, avec au
// moins n voisins par trait de même nature.
func unPointLibre(t *testing.T, n int) int {
	for i, p := range C.Points {
		if p.Toit || p.Couleurs != 0 || p.GP || i == PtTribune {
			continue
		}
		k := 0
		for _, v := range C.Voisins[i] {
			if C.Trait[i][v] && !C.Points[v].GP && !C.Points[v].Toit {
				k++
			}
		}
		if k >= n {
			return i
		}
	}
	t.Fatal("aucun point")
	return -1
}

func TestTableTirs(t *testing.T) {
	for k := range TableTirs {
		c := Chances(k)
		if s := c[0] + c[1] + c[2] + c[3]; math.Abs(s-1) > 1e-9 {
			t.Errorf("ligne %d : %v", k, s)
		}
	}
	// la règle : le terroriste ne tue que sur un 7 depuis un point au sol
	if TableTirs[TTerroNonColorie][5] != Tue || TableTirs[TPoliceHelico][4] != Tue {
		t.Error("table des tirs mal recopiée")
	}
}

// La police au contact : le manifestant recule vers un point sans police,
// ou il est interpellé s'il est cerné.
func TestContact(t *testing.T) {
	e := etatVide()
	a := unPointLibre(t, 2)
	m := e.Pions[NbPoliciers+5]
	m.Pos = a
	var b int
	for _, v := range C.Voisins[a] {
		if C.Trait[a][v] && !C.Points[v].Toit {
			b = v
			break
		}
	}
	pol := e.Pions[0]
	pol.Pos = b
	e.Appliquer(Mouvement{Type: MvtPas, Pion: pol.ID, Vers: a, Chemin: []int{b, a}})
	if m.Pos == a || m.Retire {
		t.Fatalf("le manifestant n'a pas reculé : pos %d, retiré %v", m.Pos, m.Retire)
	}
	if e.PolicePresente(m.Pos) {
		t.Fatal("recul vers la police")
	}
	// cerné : un policier sur chaque voisin
	e2 := etatVide()
	m2 := e2.Pions[NbPoliciers+5]
	m2.Pos = a
	k := 1
	for _, v := range C.Voisins[a] {
		if v == b {
			continue
		}
		e2.Pions[k].Pos = v
		k++
	}
	e2.Pions[0].Pos = b
	e2.Pions[k].Pos = b // le point que quitte le policier reste gardé
	e2.Appliquer(Mouvement{Type: MvtPas, Pion: 0, Vers: a, Chemin: []int{b, a}})
	if !m2.Retire || !m2.Arrete {
		t.Fatal("le manifestant cerné n'est pas interpellé")
	}
}

// Les manifestants ne franchissent pas les barrières de la Grand-Place et
// ne marchent jamais vers la police.
func TestBarrieres(t *testing.T) {
	e := etatVide()
	for i, p := range C.Points {
		if p.GP {
			continue
		}
		for _, v := range C.Voisins[i] {
			if C.Points[v].GP {
				m := e.Pions[NbPoliciers]
				m.Pos = i
				if e.PeutAller(m, i, v, true) {
					t.Fatalf("un manifestant entre dans la Grand-Place par %d → %d", i, v)
				}
				return
			}
		}
	}
	t.Fatal("aucun point au bord de la Grand-Place")
}

// Un manifestant ne monte sur un toit qu'après un tour sur le point fléché.
func TestToit(t *testing.T) {
	e := etatVide()
	for _, a := range accesCarte {
		r, toit := a[0], a[1]
		m := e.Pions[NbPoliciers]
		m.Pos, m.Arrivee = r, e.Tour
		if e.PeutAller(m, r, toit, true) {
			t.Fatal("monté sur le toit sans attendre")
		}
		m.Arrivee = e.Tour - 1
		if !e.PeutAller(m, r, toit, true) {
			t.Fatal("accès au toit refusé après un tour")
		}
		return
	}
}

// Point occupé devant la voiture : blocage, et la table des scores.
func TestBlocage(t *testing.T) {
	e := etatVide()
	e.CouleurAller, e.AllerPublic = Rouge, true
	e.Voiture.Pos = PtHotelDeVille
	suiv := e.ProchainsPoints()
	if len(suiv) == 0 {
		t.Fatal("pas de point suivant")
	}
	for _, q := range suiv {
		e.Pions[NbPoliciers+7].Pos = q
	}
	e.avancerVoiture()
	if e.Blocages != 1 || e.Score[RManif] != 80 || e.Score[RPolice] != -20 {
		t.Fatalf("blocage %d, scores %v", e.Blocages, e.Score)
	}
	e.Pions[NbPoliciers+7].Pos = -1
	e.avancerVoiture()
	if e.Voiture.Pos == PtHotelDeVille {
		t.Fatal("la voiture n'avance pas sur un point libre")
	}
}

// Le discours : point 1 → président à la tribune, 2 → choix du retour, 3 →
// fin du discours, puis le retour vers l'aéroport.
func TestDiscours(t *testing.T) {
	e := etatVide()
	e.CouleurAller, e.AllerPublic = Bleu, true
	e.Voiture.Pos = PtHotelDeVille
	for range 40 {
		e.avancerVoiture()
		if e.Discours == 1 {
			break
		}
	}
	if e.Voiture.Pos != PtGP1 || e.Pres.Pos != PtTribune {
		t.Fatalf("arrivée sur la Grand-Place : voiture %d, président %d", e.Voiture.Pos, e.Pres.Pos)
	}
	e.avancerVoiture()
	e.avancerVoiture()
	if e.Voiture.Pos != PtGP3 || e.Pres.Pos != -1 || e.CouleurRetour == 0 {
		t.Fatalf("fin du discours : voiture %d, président %d", e.Voiture.Pos, e.Pres.Pos)
	}
	for range 40 {
		e.avancerVoiture()
		if e.Arrivee {
			break
		}
	}
	if !e.Arrivee || e.Voiture.Pos != PtAeroport || e.Score[RPolice] != 200 {
		t.Fatalf("retour : arrivée %v, scores %v", e.Arrivee, e.Score)
	}
}

// Un véhicule vide et sans garde est pris par le manifestant qui y entre.
func TestVehiculePris(t *testing.T) {
	e := etatVide()
	a := unPointLibre(t, 2)
	var b int
	for _, v := range C.Voisins[a] {
		if C.Trait[a][v] && !C.Points[v].Toit && !C.Points[v].GP {
			b = v
			break
		}
	}
	moto := e.Pions[NbPoliciers+NbManifs+2]
	moto.Pos = a
	m := e.Pions[NbPoliciers+4]
	m.Pos = b
	e.Appliquer(Mouvement{Type: MvtPas, Pion: m.ID, Vers: a, Chemin: []int{b, a}})
	if m.Dans != moto.ID {
		t.Fatal("la moto n'est pas prise")
	}
	// au même tour, elle ne roule pas encore
	for _, mv := range e.Mouvements(Quota{Role: RManif, Manifs: true, Vehicules: true}, map[int]bool{}) {
		if mv.Pion == moto.ID {
			t.Fatal("la moto prise roule le tour même")
		}
	}
	e.Tour++
	n := 0
	for _, mv := range e.Mouvements(Quota{Role: RManif, Manifs: true, Vehicules: true}, map[int]bool{}) {
		if mv.Pion == moto.ID {
			n++
		}
	}
	if n == 0 {
		t.Fatal("la moto prise ne roule pas au tour suivant")
	}
}

// Un terroriste blessé reste en jeu, immobile ; la seconde blessure le tue.
func TestBlessures(t *testing.T) {
	e := etatVide()
	ter := e.Pions[e.Terroristes[0]]
	ter.Pos = unPointLibre(t, 1)
	e.toucher(ter, 1, false)
	if ter.Retire || ter.Blessures != 1 || e.Revele[ter.ID] != RTerro {
		t.Fatalf("terroriste blessé : %+v", ter)
	}
	for _, mv := range e.Mouvements(Quota{Role: RTerro, Manifs: true, Parite: 1}, map[int]bool{}) {
		if mv.Pion == ter.ID {
			t.Fatal("un terroriste blessé se déplace")
		}
	}
	e.toucher(ter, 1, false)
	if !ter.Retire || e.Score[RAgent] != 200 {
		t.Fatalf("seconde blessure : retiré %v, scores %v", ter.Retire, e.Score)
	}
}

// Les décisions de l'ordinateur restent rapides (il tient le verrou de
// l'interface pendant qu'il réfléchit).
func TestTempsIA(t *testing.T) {
	e, _ := partieIA(4, 7)
	_ = e
	pire := time.Duration(0)
	e2 := NouvelEtat(4, []string{"A", "B", "C", "D"}, 3)
	for _, j := range e2.Joueurs {
		j.Ctrl = chrono{IA{}, &pire}
	}
	e2.Jouer()
	if pire > 80*time.Millisecond {
		t.Errorf("décision la plus lente : %v", pire)
	}
	t.Logf("décision la plus lente : %v", pire)
}

type chrono struct {
	IA
	pire *time.Duration
}

func (c chrono) Deplacer(e *Etat, j *Joueur, q Quota, opts []Mouvement) int {
	d := time.Now()
	k := c.IA.Deplacer(e, j, q, opts)
	*c.pire = max(*c.pire, time.Since(d))
	return k
}
