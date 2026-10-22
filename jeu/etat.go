package jeu

import (
	"fmt"
	"math/rand/v2"
)

// Genre d'un pion.
type Genre int

const (
	Policier Genre = iota
	Manifestant
	VoiturePolice
	Moto
	Helico
	VoitureOfficielle
	President
)

var NomsGenres = [...]string{"policier", "manifestant", "voiture de police", "moto", "hélicoptère", "voiture présidentielle", "président"}

// Capacite : pions police transportés ; Portee : points (ou zones pour
// l'hélicoptère) parcourus par tour.
func (g Genre) Capacite() int {
	switch g {
	case VoiturePolice:
		return 3
	case Moto:
		return 2
	case Helico:
		return 4
	}
	return 0
}

func (g Genre) Portee() int {
	switch g {
	case VoiturePolice:
		return 2
	case Moto:
		return 3
	case Helico:
		return 4
	}
	return 1
}

func (g Genre) Vehicule() bool { return g == VoiturePolice || g == Moto || g == Helico }

// AuSol : police et manifestants, les pions « à pied ».
func (g Genre) APied() bool { return g == Policier || g == Manifestant }

// Role : les quatre forces en présence.
type Role int

const (
	RPolice Role = iota
	RManif
	RAgent
	RTerro
	NbRoles
)

var NomsRoles = [...]string{"Police", "Manifestants", "Agents secrets", "Terroristes"}

type Pion struct {
	ID    int
	Genre Genre
	Num   int // numéro imprimé (policiers, manifestants)

	Pos       int // point ; -1 : hors du plan (dans un véhicule, en vol, retiré)
	Dans      int // véhicule qui le porte, -1
	Vol       int // hélicoptère en vol : zone survolée ; -1 posé
	Retire    bool
	Blessures int
	Arrivee   int // tour où le pion est arrivé sur son point
	PrisTour  int // véhicule pris par un manifestant à ce tour (-1)
	Enfui     bool
	Arrete    bool
}

func (p *Pion) Actif() bool { return !p.Retire && !p.Enfui }

// Nom : « policier 12 », « manifestant 7 », « moto »…
func (p *Pion) Nom() string {
	if p.Genre.APied() {
		return fmt.Sprintf("%s %d", NomsGenres[p.Genre], p.Num)
	}
	return NomsGenres[p.Genre]
}

type Joueur struct {
	Num    int
	Nom    string
	Roles  []Role
	Humain bool
	Ctrl   Controleur
}

func (j *Joueur) A(r Role) bool {
	for _, x := range j.Roles {
		if x == r {
			return true
		}
	}
	return false
}

// Etat : toute la partie, vérité comprise (les identités secrètes) ; les
// joueurs de l'ordinateur ne lisent que ce qu'ils ont le droit de savoir
// (Connu, Public).
type Etat struct {
	Pions   []*Pion
	Joueurs []*Joueur
	ParRole [NbRoles]*Joueur
	NbJ     int

	Terroristes []int // identités secrètes (ID de pions)
	Agents      []int
	Revele      map[int]Role // identités dévoilées : pion → RTerro ou RAgent

	CouleurAller  Couleur
	CouleurRetour Couleur
	AllerPublic   bool
	RetourPublic  bool
	Retour        bool // la voiture est repartie de la Grand-Place
	Discours      int  // 0 : pas encore ; 1–3 : voiture sur le point 1, 2, 3 ; 4 : fini
	DoubleAvance  bool // une tentative d'attentat : la voiture avance deux fois
	Tour          int
	Phase         Phase
	Blocages      int
	Bavures       int
	Attentat      bool // une tentative a eu lieu
	PresBlesse    bool
	PresTue       bool
	Arrivee       bool // président à l'aéroport
	TourFin       int  // fin du cortège : tour où commence la chasse à l'homme
	Fini          bool
	Score         [NbRoles]int
	Lignes        []LigneScore
	TourMax       int
	NbMvts        int // mouvements joués dans la phase en cours
	Rng           *rand.Rand
	Obs           Observateur
	derniersMvtsT map[int]int // pions déplacés à la phase des terroristes, ce tour
	SoupcPublic   []float64   // soupçon public sur chaque pion (ID)
	Voiture, Pres *Pion
	Helico        *Pion
}

// LigneScore : une ligne de la table des scores appliquée pendant la partie.
type LigneScore struct {
	Tour  int
	Texte string
	Pts   [NbRoles]int
}

const (
	NbPoliciers   = 30
	NbManifs      = 40
	TourMaxDefaut = 30
	ToursChasse   = 3
)

// NouvelEtat : la partie à nbJ joueurs (2, 3 ou 4) ; humain : le numéro du
// joueur humain (-1 : aucun).
func NouvelEtat(nbJ int, noms []string, graine uint64) *Etat {
	e := &Etat{NbJ: nbJ, Revele: map[int]Role{}, TourMax: TourMaxDefaut,
		Rng: rand.New(rand.NewPCG(graine, graine^0x5eed)), derniersMvtsT: map[int]int{}}
	ajoute := func(g Genre, num int) *Pion {
		p := &Pion{ID: len(e.Pions), Genre: g, Num: num, Pos: -1, Dans: -1, Vol: -1, PrisTour: -1}
		e.Pions = append(e.Pions, p)
		return p
	}
	for i := 1; i <= NbPoliciers; i++ {
		ajoute(Policier, i)
	}
	for i := 1; i <= NbManifs; i++ {
		ajoute(Manifestant, i)
	}
	ajoute(VoiturePolice, 1)
	ajoute(VoiturePolice, 2)
	for i := 1; i <= 4; i++ {
		ajoute(Moto, i)
	}
	e.Helico = ajoute(Helico, 1)
	e.Voiture = ajoute(VoitureOfficielle, 0)
	e.Pres = ajoute(President, 0)
	e.SoupcPublic = make([]float64, len(e.Pions))

	var roles [][]Role
	switch nbJ {
	case 2:
		roles = [][]Role{{RPolice, RAgent}, {RManif, RTerro}}
	case 3:
		roles = [][]Role{{RPolice, RAgent}, {RManif}, {RTerro}}
	default:
		roles = [][]Role{{RPolice}, {RManif}, {RAgent}, {RTerro}}
	}
	for i, rs := range roles {
		j := &Joueur{Num: i, Nom: noms[i], Roles: rs}
		e.Joueurs = append(e.Joueurs, j)
		for _, r := range rs {
			e.ParRole[r] = j
		}
	}
	return e
}

// ---------------------------------------------------------------- occupation

// APied : les pions à pied (police et manifestants) posés sur le point p.
func (e *Etat) APied(p int) []*Pion {
	var l []*Pion
	for _, q := range e.Pions {
		if q.Pos == p && q.Actif() && q.Dans < 0 && q.Genre.APied() {
			l = append(l, q)
		}
	}
	return l
}

func (e *Etat) Compte(p int, g Genre) int {
	n := 0
	for _, q := range e.Pions {
		if q.Pos == p && q.Actif() && q.Dans < 0 && q.Genre == g {
			n++
		}
	}
	return n
}

// VehiculeSur : le véhicule posé sur p (nil si aucun).
func (e *Etat) VehiculeSur(p int) *Pion {
	for _, q := range e.Pions {
		if q.Pos == p && q.Actif() && q.Genre.Vehicule() {
			return q
		}
	}
	return nil
}

// Passagers : les pions dans le véhicule v.
func (e *Etat) Passagers(v *Pion) []*Pion {
	var l []*Pion
	for _, q := range e.Pions {
		if q.Dans == v.ID && q.Actif() {
			l = append(l, q)
		}
	}
	return l
}

// Occupe : un pion, un véhicule, la voiture ou le président sur p ?
func (e *Etat) Occupe(p int) bool {
	for _, q := range e.Pions {
		if q.Pos == p && q.Actif() && q.Dans < 0 {
			return true
		}
	}
	return false
}

// Lieu : le point où se trouve un pion (celui de son véhicule s'il est
// dedans) ; -1 en vol ou hors jeu.
func (e *Etat) Lieu(p *Pion) int {
	if p.Dans >= 0 {
		return e.Pions[p.Dans].Pos
	}
	return p.Pos
}

// Zone d'un pion (hélicoptère en vol compris), -1 hors du plan.
func (e *Etat) ZonePion(p *Pion) int {
	v := p
	if p.Dans >= 0 {
		v = e.Pions[p.Dans]
	}
	if v.Vol >= 0 {
		return v.Vol
	}
	if v.Pos < 0 {
		return -1
	}
	return C.Points[v.Pos].Zone()
}

// PolicePresente : un policier à pied ou un véhicule tenu par la police sur p.
func (e *Etat) PolicePresente(p int) bool {
	if e.Compte(p, Policier) > 0 {
		return true
	}
	if v := e.VehiculeSur(p); v != nil {
		for _, q := range e.Passagers(v) {
			if q.Genre == Policier {
				return true
			}
		}
	}
	return false
}

// ManifsSur : manifestants sur p, à pied ou au volant d'un véhicule pris.
func (e *Etat) ManifsSur(p int) []*Pion {
	var l []*Pion
	for _, q := range e.Pions {
		if q.Genre == Manifestant && q.Actif() && e.Lieu(q) == p && (q.Dans < 0 || e.Pions[q.Dans].Vol < 0) {
			l = append(l, q)
		}
	}
	return l
}

// ---------------------------------------------------------------- identités

func (e *Etat) EstTerroriste(id int) bool { return contient(e.Terroristes, id) }
func (e *Etat) EstAgent(id int) bool      { return contient(e.Agents, id) }

func contient(l []int, x int) bool {
	for _, y := range l {
		if y == x {
			return true
		}
	}
	return false
}

// Connu : ce que le joueur j sait de l'identité du pion id (RTerro, RAgent,
// ou -1 : rien de plus que son uniforme).
func (e *Etat) Connu(j *Joueur, id int) Role {
	if r, ok := e.Revele[id]; ok {
		return r
	}
	if j.A(RTerro) && e.EstTerroriste(id) {
		return RTerro
	}
	if j.A(RAgent) && e.EstAgent(id) {
		return RAgent
	}
	return -1
}

// Pair : la règle des numéros (agents pairs, terroristes impairs).
func Pair(p *Pion) bool { return p.Num%2 == 0 }

// ---------------------------------------------------------------- score

func (e *Etat) marquer(texte string, pts [NbRoles]int) {
	for r := range NbRoles {
		e.Score[r] += pts[r]
	}
	e.Lignes = append(e.Lignes, LigneScore{e.Tour, texte, pts})
	e.signaler(Ev{Type: EvScore, Texte: texte, Pts: pts})
}

// Lignes de la table des scores (police, agent secret, manifestant,
// terroriste ; réordonnées selon Role).
func pts(police, agent, manif, terro int) [NbRoles]int {
	var p [NbRoles]int
	p[RPolice], p[RAgent], p[RManif], p[RTerro] = police, agent, manif, terro
	return p
}

var (
	ScBlocage       = pts(-20, 0, 80, 0)
	ScJamaisBloque  = pts(100, 0, -30, 0)
	ScBavure        = pts(-30, 0, -20, 0)
	ScTerroAvant    = pts(200, 200, 0, -50)
	ScPresBlesse    = pts(-100, 0, -50, 200)
	ScPresTue       = pts(-150, -50, -100, 200)
	ScTerroApres    = pts(70, 100, -50, -30)
	ScPresAeroport  = pts(200, 0, 100, -50)
	ScDeuxEchappent = pts(-100, -150, 0, 100)
	ScUnSeulEchappe = pts(-50, -100, 0, 200)
)

// ScoreJoueur : la somme des rôles du joueur.
func (e *Etat) ScoreJoueur(j *Joueur) int {
	s := 0
	for _, r := range j.Roles {
		s += e.Score[r]
	}
	return s
}

// ---------------------------------------------------------------- journal

func (e *Etat) signaler(ev Ev) {
	ev.Tour = e.Tour
	if e.Obs != nil {
		e.Obs.Evenement(e, ev)
	}
}

func (e *Etat) dire(format string, a ...any) {
	e.signaler(Ev{Type: EvTexte, Texte: fmt.Sprintf(format, a...)})
}
