package main

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"cortegeofficiel/jeu"
)

// Partie : le moteur tourne dans sa propre goroutine et garde le verrou
// tant qu'il calcule ; il le rend pendant qu'il attend le joueur humain ou
// qu'une animation se joue. L'interface ne lit l'état que sous ce verrou.

type TypeReq int

const (
	ReqAucune TypeReq = iota
	ReqSecrets
	ReqCouleur
	ReqPlacer
	ReqAvancer
	ReqDeplacer
	ReqTirer
	ReqReculer
)

var nomsReq = map[string]TypeReq{"secrets": ReqSecrets, "couleur": ReqCouleur, "placer": ReqPlacer,
	"avancer": ReqAvancer, "deplacer": ReqDeplacer, "tirer": ReqTirer, "reculer": ReqReculer}

// Demande : ce que le moteur attend du joueur humain.
type Demande struct {
	Type   TypeReq
	Num    int // change à chaque demande
	Role   jeu.Role
	Retour bool
	Cands  []int
	N      int
	Prop   []jeu.Placement
	Points []int
	Quota  jeu.Quota
	Mvts   []jeu.Mouvement
	Phase  jeu.Phase
	Tirs   []jeu.Tir
}

// Anim : un pion qui glisse le long d'un trajet (coordonnées du plan).
type Anim struct {
	Pion  int
	Trace []jeu.Point
	Debut time.Time
	Duree time.Duration
}

type LigneJournal struct {
	Tour  int
	Texte string
	Coul  color.RGBA
}

type Partie struct {
	mu  sync.Mutex
	E   *jeu.Etat
	Moi *jeu.Joueur

	Req     Demande
	nReq    int
	reponse chan any

	Journal    []LigneJournal
	DernierTir *jeu.Ev
	Anim       *Anim
	Rapide     bool
	Termine    bool
	sansPause  bool // tests : aucune attente

	// captures : l'ordinateur joue pour l'humain jusqu'à cette demande
	auto   bool
	arret  TypeReq
	arretN int
	ia     jeu.IA
}

// nouvellePartie : nbJ joueurs, l'humain prend la place moi (-1 : aucun).
func nouvellePartie(nbJ, moi int, graine uint64) *Partie {
	noms := nomsJoueurs(nbJ, moi)
	p := &Partie{E: jeu.NouvelEtat(nbJ, noms, graine), reponse: make(chan any, 1)}
	for i, j := range p.E.Joueurs {
		if i == moi {
			j.Humain = true
			j.Ctrl = &Humain{p}
			p.Moi = j
		} else {
			j.Ctrl = jeu.IA{}
		}
	}
	p.E.Obs = p
	return p
}

// nomsJoueurs : « Vous », et des noms de l'époque pour l'ordinateur.
func nomsJoueurs(nbJ, moi int) []string {
	autres := []string{"Commissaire Vallat", "Le Comité", "Colonel Marceau", "Le Chacal"}
	var noms []string
	for i := range nbJ {
		if i == moi {
			noms = append(noms, "Vous")
			continue
		}
		noms = append(noms, autres[min(i, 3)])
	}
	if nbJ == 3 && moi != 2 {
		noms[2] = "Le Chacal"
	}
	if nbJ == 2 && moi != 1 {
		noms[1] = "Le Comité"
	}
	return noms
}

func (p *Partie) lancer() {
	go func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.E.Jouer()
		p.Termine = true
	}()
}

// attendre : rend le verrou pendant d (animation, lecture).
func (p *Partie) attendre(d time.Duration) {
	if p.auto || p.sansPause || d <= 0 {
		return
	}
	if p.Rapide {
		d /= 4
	}
	p.mu.Unlock()
	time.Sleep(d)
	p.mu.Lock()
}

// demander : pose une demande et attend la réponse de l'interface.
func (p *Partie) demander(d Demande) any {
	p.nReq++
	d.Num = p.nReq
	p.Req = d
	p.mu.Unlock()
	r := <-p.reponse
	p.mu.Lock()
	p.Req = Demande{Num: p.nReq}
	return r
}

// repondre : appelé par l'interface, verrou tenu.
func (p *Partie) repondre(r any) {
	if p.Req.Type == ReqAucune {
		return
	}
	p.Req.Type = ReqAucune
	select {
	case p.reponse <- r:
	default:
	}
}

// pilote : l'ordinateur joue-t-il encore pour l'humain (captures) ?
func (p *Partie) pilote(t TypeReq) bool {
	if !p.auto {
		return false
	}
	if t == p.arret {
		p.arretN--
		if p.arretN <= 0 {
			p.auto = false
			return false
		}
	}
	return true
}

func (p *Partie) noter(texte string, c color.RGBA) {
	p.Journal = append(p.Journal, LigneJournal{p.E.Tour, texte, c})
	if len(p.Journal) > 400 {
		p.Journal = p.Journal[100:]
	}
}

// joueurActif : le joueur qui agit pendant la phase en cours.
func (p *Partie) humainJoue() bool {
	return p.Req.Type != ReqAucune
}

// Evenement : l'observateur du moteur.
func (p *Partie) Evenement(e *jeu.Etat, ev jeu.Ev) {
	switch ev.Type {
	case jeu.EvTour:
		p.noter(fmt.Sprintf("— Tour %d —", e.Tour), coulSignal)
	case jeu.EvPhase:
		p.attendre(120 * time.Millisecond)
	case jeu.EvTexte:
		p.noter(ev.Texte, coulEncre)
		p.attendre(250 * time.Millisecond)
	case jeu.EvBlocage:
		p.noter(ev.Texte, coulRouge)
		p.attendre(500 * time.Millisecond)
	case jeu.EvRevele:
		p.noter(ev.Texte, coulRouge)
		p.attendre(700 * time.Millisecond)
	case jeu.EvScore:
		p.noter("Score : "+ev.Texte, coulEncre2)
	case jeu.EvTir:
		c := ev
		p.DernierTir = &c
		p.noter(ev.Texte, coulRouge)
		p.attendre(1100 * time.Millisecond)
	case jeu.EvMouvement:
		tr := tracePion(e, ev)
		if len(tr) < 2 {
			return
		}
		duree := 160 * time.Millisecond
		if pi := e.Pions[ev.Pion]; pi.Genre == jeu.VoitureOfficielle {
			duree = 600 * time.Millisecond
		} else if !pi.Genre.APied() {
			duree = 400 * time.Millisecond
		}
		if p.Req.Type != ReqAucune || p.Moi != nil && e.Phase != jeu.PhContact && p.Moi == p.joueurPhase(e) {
			duree = 110 * time.Millisecond
		}
		p.Anim = &Anim{Pion: ev.Pion, Trace: tr, Debut: time.Now(), Duree: duree}
		p.attendre(duree)
		p.Anim = nil
	case jeu.EvFin:
		p.noter("Fin de la partie.", coulSignal)
	}
}

// joueurPhase : qui joue pendant la phase.
func (p *Partie) joueurPhase(e *jeu.Etat) *jeu.Joueur {
	switch e.Phase {
	case jeu.PhVoiture, jeu.PhPolice:
		return e.ParRole[jeu.RPolice]
	case jeu.PhManifs, jeu.PhContact:
		return e.ParRole[jeu.RManif]
	case jeu.PhAgents:
		return e.ParRole[jeu.RAgent]
	case jeu.PhTerros, jeu.PhAttentat:
		return e.ParRole[jeu.RTerro]
	}
	return nil
}

// tracePion : le dessin d'un mouvement (zones survolées : leur centre).
func tracePion(e *jeu.Etat, ev jeu.Ev) []jeu.Point {
	var tr []jeu.Point
	for _, q := range ev.Chemin {
		if q >= 0 {
			tr = append(tr, jeu.Point{X: jeu.C.Points[q].X, Y: jeu.C.Points[q].Y})
		}
	}
	return tr
}

// ---------------------------------------------------------------------------
// Humain : le contrôleur du joueur humain, qui passe par l'interface.

type Humain struct{ p *Partie }

func (h *Humain) Secrets(e *jeu.Etat, j *jeu.Joueur, r jeu.Role, cands []int, n int) []int {
	if h.p.pilote(ReqSecrets) {
		return h.p.ia.Secrets(e, j, r, cands, n)
	}
	return h.p.demander(Demande{Type: ReqSecrets, Role: r, Cands: cands, N: n}).([]int)
}

func (h *Humain) Couleur(e *jeu.Etat, j *jeu.Joueur, retour bool) jeu.Couleur {
	if h.p.pilote(ReqCouleur) {
		return h.p.ia.Couleur(e, j, retour)
	}
	return h.p.demander(Demande{Type: ReqCouleur, Retour: retour}).(jeu.Couleur)
}

func (h *Humain) Placer(e *jeu.Etat, j *jeu.Joueur, r jeu.Role, prop []jeu.Placement) []jeu.Placement {
	if h.p.pilote(ReqPlacer) {
		return prop
	}
	return h.p.demander(Demande{Type: ReqPlacer, Role: r, Prop: prop}).([]jeu.Placement)
}

func (h *Humain) Avancer(e *jeu.Etat, j *jeu.Joueur, opts []int) int {
	if h.p.pilote(ReqAvancer) {
		return h.p.ia.Avancer(e, j, opts)
	}
	return h.p.demander(Demande{Type: ReqAvancer, Points: opts}).(int)
}

func (h *Humain) Deplacer(e *jeu.Etat, j *jeu.Joueur, q jeu.Quota, opts []jeu.Mouvement) int {
	if h.p.pilote(ReqDeplacer) {
		return h.p.ia.Deplacer(e, j, q, opts)
	}
	return h.p.demander(Demande{Type: ReqDeplacer, Quota: q, Mvts: opts, Role: q.Role}).(int)
}

func (h *Humain) Tirer(e *jeu.Etat, j *jeu.Joueur, ph jeu.Phase, opts []jeu.Tir) int {
	if h.p.pilote(ReqTirer) {
		return h.p.ia.Tirer(e, j, ph, opts)
	}
	return h.p.demander(Demande{Type: ReqTirer, Phase: ph, Tirs: opts}).(int)
}

func (h *Humain) Reculer(e *jeu.Etat, j *jeu.Joueur, opts []jeu.Mouvement) int {
	if h.p.pilote(ReqReculer) {
		return h.p.ia.Reculer(e, j, opts)
	}
	return h.p.demander(Demande{Type: ReqReculer, Mvts: opts}).(int)
}
