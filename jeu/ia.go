package jeu

import (
	"math"
	"math/rand/v2"
)

// IA : un joueur de l'ordinateur, simple mais honnête : il ne lit que ce
// que son rôle lui permet de savoir (Connu, soupçons publics, parcours
// dévoilés).
type IA struct{}

func (IA) Secrets(e *Etat, j *Joueur, r Role, cands []int, n int) []int {
	return e.tirerSecrets(r)
}

func (IA) Couleur(e *Etat, j *Joueur, retour bool) Couleur {
	if !retour {
		return CouleursParcours[e.Rng.IntN(3)]
	}
	// le retour : le parcours dont le départ est le moins encombré
	best, bm := Rouge, math.Inf(1)
	for _, col := range CouleursParcours {
		ch := C.Chemin(col, true, PtGP3)
		m := 0.0
		for k, c := range ch[:min(5, len(ch))] {
			for _, p := range e.Pions {
				if p.Genre == Manifestant && p.Actif() && e.Lieu(p) >= 0 && C.Dist[c][e.Lieu(p)] <= 2 {
					m += 1 / (1 + float64(k))
				}
			}
		}
		m += 0.05 * float64(len(ch))
		if m < bm {
			best, bm = col, m
		}
	}
	return best
}

func (IA) Placer(e *Etat, j *Joueur, r Role, prop []Placement) []Placement { return prop }

func (IA) Avancer(e *Etat, j *Joueur, opts []int) int {
	best, bm := 0, math.Inf(1)
	for i, q := range opts {
		m := 0.0
		for _, n := range C.Voisins[q] {
			m += float64(len(e.ManifsSur(n)))
		}
		if m < bm {
			best, bm = i, m
		}
	}
	return best
}

// ---------------------------------------------------------------- parcours vu

// routeConnue : les prochains points de la voiture tels que j peut les
// prévoir ; pendant le discours, si le retour est secret, les débuts des
// trois parcours.
func (e *Etat) routeConnue(j *Joueur, n int) [][]int {
	v := e.Voiture
	if v.Pos < 0 || e.Arrivee {
		return nil
	}
	if e.Discours == 0 {
		var cols []Couleur
		if e.AllerPublic || j.A(RPolice) {
			cols = []Couleur{e.CouleurAller}
		} else {
			cols = CouleursParcours
		}
		var l [][]int
		for _, c := range cols {
			ch := C.Chemin(c, false, v.Pos)
			// entrer dans la Grand-Place, c'est se ranger au point 1 : il
			// doit être libre en même temps
			for k, q := range ch {
				if C.Points[q].GP {
					ch = append(ch[:k+1:k+1], PtGP2, PtGP3)
					if q != PtGP1 {
						l = append(l, []int{PtGP1})
					}
					break
				}
			}
			l = append(l, tronquer(ch, n))
		}
		return l
	}
	var tete []int
	switch e.Discours {
	case 1:
		tete = []int{PtGP2, PtGP3}
	case 2:
		tete = []int{PtGP3}
	}
	if e.RetourPublic || (j.A(RPolice) && e.Discours >= 2) {
		return [][]int{tronquer(append(tete, C.Chemin(e.CouleurRetour, true, v.Pos)...), n)}
	}
	var l [][]int
	for _, c := range CouleursParcours {
		l = append(l, tronquer(append(append([]int(nil), tete...), C.Chemin(c, true, PtGP3)...), n))
	}
	return l
}

func tronquer(l []int, n int) []int {
	if len(l) > n {
		return l[:n]
	}
	return l
}

// ---------------------------------------------------------------- soupçons

// majSoupcons : après la phase des terroristes, tout le monde a vu quels
// pions ils ont déplacés.
func (e *Etat) majSoupcons() {
	cible := e.CiblePresident()
	for id, n := range e.derniersMvtsT {
		p := e.Pions[id]
		if !p.Actif() || !p.Genre.APied() {
			continue
		}
		s := 0.6 * float64(n)
		if l := e.Lieu(p); l >= 0 && cible != nil && e.cibleTir(cible) >= 0 {
			if d := C.Dist[l][e.cibleTir(cible)]; d >= 0 && d <= 3 {
				s += 0.5 * float64(4-d) / 3
			}
		}
		e.SoupcPublic[id] += s
	}
}

// ProbTerro : la probabilité, pour j, que le pion id soit un terroriste.
func (e *Etat) ProbTerro(j *Joueur, id int) float64 {
	p := e.Pions[id]
	switch e.Connu(j, id) {
	case RTerro:
		return 1
	case RAgent:
		return 0
	}
	if p.Num%2 == 0 || !p.Genre.APied() || p.Genre == Policier && e.NbJ != 2 {
		return 0
	}
	if j.A(RPolice) && p.Genre == Policier && e.NbJ != 2 {
		return 0
	}
	tot := 0.0
	nCand := 0
	for _, q := range e.Pions {
		if q.Actif() && q.Genre == p.Genre && q.Num%2 == 1 && e.Connu(j, q.ID) < 0 {
			tot += e.SoupcPublic[q.ID] + 0.3
			nCand++
		}
	}
	restants := 0
	for _, t := range e.Terroristes {
		// le nombre de terroristes encore anonymes de ce genre est public :
		// la règle en fixe le nombre, les démasqués sont connus
		if e.Pions[t].Genre == p.Genre && e.Revele[t] != RTerro {
			restants++
		}
	}
	if tot == 0 || restants == 0 {
		return 0
	}
	return min(float64(restants)*(e.SoupcPublic[id]+0.3)/tot, 0.97)
}

// ---------------------------------------------------------------- tirs

func (IA) Tirer(e *Etat, j *Joueur, ph Phase, opts []Tir) int {
	best, bv := -1, 0.0
	for i, t := range opts {
		ch := Chances(t.Table)
		touche := ch[Blesse] + ch[Tue]
		cible := e.Pions[t.Cible]
		var v float64
		switch {
		case e.EstTerroriste(t.Tireur) && j.A(RTerro) && (cible.ID == e.Pres.ID || cible.ID == e.Voiture.ID):
			v = 1 + touche
		case e.Revele[cible.ID] == RTerro:
			v = 1 + touche
			if !e.EstAgent(t.Tireur) {
				// la police craint la bavure sur un voisin
				if len(e.ManifsSur(e.cibleTir(cible))) > 1 {
					v -= 0.4
				}
			}
		default:
			pr := e.ProbTerro(j, cible.ID)
			if e.EstAgent(t.Tireur) {
				if pr < 0.55 {
					continue
				}
				v = pr * touche
			} else {
				if pr < 0.6 || len(e.ManifsSur(e.cibleTir(cible))) > 1 {
					continue
				}
				v = pr*touche - (1-pr)*touche
			}
		}
		if v > bv {
			best, bv = i, v
		}
	}
	return best
}

// ---------------------------------------------------------------- déplacements

func (ia IA) Deplacer(e *Etat, j *Joueur, q Quota, opts []Mouvement) int {
	switch {
	case q.Role == RPolice:
		return ia.police(e, j, opts)
	case q.Role == RManif && q.Policiers:
		return ia.provocation(e, j, opts)
	case q.Role == RManif:
		return ia.manifs(e, j, opts)
	case q.Role == RAgent:
		return ia.agents(e, j, opts)
	}
	return ia.terros(e, j, q, opts)
}

func (ia IA) Reculer(e *Etat, j *Joueur, opts []Mouvement) int {
	rt := e.routeConnue(j, 6)
	best, bv := 0, math.Inf(-1)
	for i, m := range opts {
		v := valeurManif(e, rt, m.Vers) - 0.5*float64(len(e.ManifsSur(m.Vers)))
		if v > bv {
			best, bv = i, v
		}
	}
	return best
}

// valeurPolice : l'intérêt d'avoir un policier en x : ne pas gêner la
// voiture, tenir le parcours en avant et le border (cordon).
func valeurPolice(e *Etat, rt [][]int, x int) float64 {
	if x < 0 {
		return 0
	}
	v := 0.0
	for _, r := range rt {
		for k, c := range r {
			if x == c {
				if k == 0 {
					return -100
				}
				v += float64(7 - k)
			} else if C.Trait[x][c] && C.Points[x].Couleurs == 0 && k < 5 {
				v += float64(5-k) * 1.5
			}
		}
	}
	if c := e.CiblePresident(); c != nil {
		if pt := e.cibleTir(c); pt >= 0 && C.Trait[x][pt] {
			v += 2
		}
	}
	return v
}

// menace : manifestants sur x ou à côté qui pourraient se mettre en travers.
func valeurContact(e *Etat, rt [][]int, x int) float64 {
	n := float64(len(e.ManifsSur(x)))
	if n == 0 {
		return 0
	}
	v := 0.0
	for _, r := range rt {
		for k, c := range r {
			if k >= 5 {
				break
			}
			if x == c {
				v += float64(6-k) * 3
			} else if C.Trait[x][c] {
				v += float64(5-k) * 1.5
			}
		}
	}
	return v * n
}

// police : chaque mouvement est essayé sur une copie de la partie (contacts
// compris) et jugé par evalPolice ; le meilleur est joué tant qu'il
// améliore la situation.
func (IA) police(e *Etat, j *Joueur, opts []Mouvement) int {
	rt := e.routeConnue(j, 7)
	base := evalPolice(e, j, rt)
	best, bv := -1, 0.1
	for i, m := range opts {
		if m.Type != MvtPas && m.Type != MvtConduite && m.Type != MvtDescente {
			continue
		}
		s := e.simuler()
		s.Appliquer(m)
		if g := evalPolice(s, j, rt) - base; g > bv {
			best, bv = i, g
		}
	}
	return best
}

// menaceSur : manifestants capables d'entrer en x au prochain pas si x
// était libre.
func (e *Etat) menaceSur(x int) int {
	if e.interditAuxPions(x) {
		return 0
	}
	n := 0
	for _, y := range C.Voisins[x] {
		for _, m := range e.ManifsSur(y) {
			if m.Dans < 0 && m.Blessures == 0 && C.Points[y].GP == C.Points[x].GP && e.accesToit(m, y, x) {
				n++
			}
		}
	}
	return n
}

// evalPolice : la situation vue par la police à la fin de ses
// déplacements. Le prochain point de la voiture doit rester libre et hors
// de portée des manifestants ; les suivants aussi, autant que possible.
func evalPolice(e *Etat, j *Joueur, rt [][]int) float64 {
	v := 0.0
	surRoute := map[int]bool{}
	for _, r := range rt {
		for k, c := range r {
			surRoute[c] = true
			poids := 1 / float64(1+k)
			nm := len(e.ManifsSur(c))
			switch {
			case nm > 0:
				v -= 60 * poids
			case e.Occupe(c):
				// chaque policier compte : sinon ils ne partent jamais un par un
				n := float64(e.Compte(c, Policier))
				if e.VehiculeSur(c) != nil {
					n++
				}
				v -= (25 + 20*n) * poids
				v -= 4 * float64(e.menaceSur(c)) * poids
			default:
				// libre mais à la merci d'un manifestant : autant dire bloqué
				t := e.menaceSur(c)
				if t > 0 {
					v -= (42 + 4*float64(t)) * poids
				}
			}
			if k > 2 {
				continue
			}
			// le cordon : policiers sur les points voisins du parcours
			for _, y := range C.Voisins[c] {
				if e.PolicePresente(y) && !surRoute[y] {
					v += 4 * poids
				}
			}
		}
	}
	// le danger : des manifestants (des terroristes, peut-être) à portée de
	// tir de la voiture ou du point où elle va
	if c := e.CiblePresident(); c != nil {
		vises := []int{e.cibleTir(c)}
		if len(rt) > 0 && len(rt[0]) > 0 {
			vises = append(vises, rt[0][0])
		}
		for _, m := range e.Pions {
			if m.Genre != Manifestant || !m.Actif() || e.Lieu(m) < 0 {
				continue
			}
			l := e.Lieu(m)
			for _, x := range vises {
				if x >= 0 && (C.Trait[l][x] || C.Points[l].Toit && C.Points[l].Zone() == C.Points[x].Zone()) {
					v -= 7 * (0.3 + e.ProbTerro(j, m.ID))
				}
			}
		}
	}
	// repousser la foule loin des prochains points (et l'interpeller)
	for _, r := range rt {
		tete := r[:min(3, len(r))]
		for _, m := range e.Pions {
			if m.Genre != Manifestant || !m.Actif() || e.Lieu(m) < 0 {
				continue
			}
			d := 9
			for _, c := range tete {
				if dd := C.Dist[e.Lieu(m)][c]; dd >= 0 && dd < d {
					d = dd
				}
			}
			if d <= 3 {
				v -= float64(4-d) * 2
			}
		}
	}
	// les policiers inoccupés se rapprochent du cortège
	if len(rt) > 0 && len(rt[0]) > 1 {
		cible := rt[0][min(2, len(rt[0])-1)]
		for _, p := range e.Pions {
			if p.Genre == Policier && p.Actif() && e.Lieu(p) >= 0 {
				if d := C.Dist[e.Lieu(p)][cible]; d > 1 {
					v -= 0.6 * float64(d)
				}
			}
		}
	}
	// les terroristes démasqués
	for _, t := range e.Terroristes {
		if e.Revele[t] == RTerro && e.Pions[t].Actif() {
			if l := e.cibleTir(e.Pions[t]); l >= 0 {
				for _, y := range C.Voisins[l] {
					if C.Trait[y][l] && e.PolicePresente(y) {
						v += 5
					}
				}
			}
		}
	}
	return v
}

// simuler : une copie de la partie pour essayer un coup, jouée par
// l'ordinateur et sans observateur.
func (e *Etat) simuler() *Etat {
	c := *e
	c.Pions = make([]*Pion, len(e.Pions))
	for i, p := range e.Pions {
		q := *p
		c.Pions[i] = &q
	}
	c.Revele = make(map[int]Role, len(e.Revele))
	for k, v := range e.Revele {
		c.Revele[k] = v
	}
	c.Joueurs = nil
	for _, j := range e.Joueurs {
		jj := *j
		jj.Ctrl = IA{}
		c.Joueurs = append(c.Joueurs, &jj)
		for _, r := range jj.Roles {
			c.ParRole[r] = &jj
		}
	}
	c.Obs, c.Lignes, c.derniersMvtsT = nil, nil, map[int]int{}
	c.Voiture, c.Pres, c.Helico = c.Pions[e.Voiture.ID], c.Pions[e.Pres.ID], c.Pions[e.Helico.ID]
	c.Rng = rand.New(rand.NewPCG(1, 2))
	return &c
}

// valeurManif : l'intérêt d'avoir un manifestant en x : sur le prochain
// point de la voiture, sur le parcours, ou tout près.
func valeurManif(e *Etat, rt [][]int, x int) float64 {
	v := 0.0
	for _, r := range rt {
		for k, c := range r {
			switch {
			case x == c && k == 0:
				v += 40
			case x == c:
				v += float64(12 - k)
			case C.Trait[x][c] && k < 4:
				v += float64(6-k) * 1.2
			}
		}
	}
	// rapprochement général
	d := 99
	for _, r := range rt {
		for _, c := range r {
			if dd := C.Dist[x][c]; dd >= 0 && dd < d {
				d = dd
			}
		}
	}
	if d < 99 {
		v -= 0.8 * float64(d)
	}
	return v
}

// MaxMvtsFoule : la foule de l'ordinateur ne bouge pas d'un bloc.
const MaxMvtsFoule = 12

func (IA) manifs(e *Etat, j *Joueur, opts []Mouvement) int {
	if e.NbMvts >= MaxMvtsFoule {
		return -1
	}
	rt := e.routeConnue(j, 8)
	// un seul manifestant suffit sur le prochain point
	tenus := map[int]bool{}
	for _, r := range rt {
		if len(r) > 0 && len(e.ManifsSur(r[0])) > 0 {
			tenus[r[0]] = true
		}
	}
	best, bv := -1, 0.3
	for i, m := range opts {
		if m.Type != MvtPas {
			continue
		}
		p := e.Pions[m.Pion]
		if tenus[p.Pos] && len(e.ManifsSur(p.Pos)) == 1 {
			continue // il tient le point : il ne bouge pas
		}
		vv := valeurManif(e, rt, m.Vers)
		if tenus[m.Vers] {
			vv -= 30
		}
		gain := vv - valeurManif(e, rt, p.Pos)
		// ne pas se jeter contre la police
		for _, n := range C.Voisins[m.Vers] {
			if e.PolicePresente(n) {
				gain -= 0.4
			}
		}
		if gain > bv {
			best, bv = i, gain
		}
	}
	return best
}

// provocation (3 joueurs) : deux policiers envoyés là où ils gênent.
func (IA) provocation(e *Etat, j *Joueur, opts []Mouvement) int {
	rt := e.routeConnue(j, 6)
	best, bv := -1, 1.0
	for i, m := range opts {
		p := e.Pions[m.Pion]
		gain := valeurPolice(e, rt, p.Pos) - valeurPolice(e, rt, m.Vers)
		if gain > bv {
			best, bv = i, gain
		}
	}
	return best
}

// ciblesTir : les points d'où l'on peut tirer sur le président maintenant
// (pres) ou après le prochain pas de la voiture (proche).
func (e *Etat) postesDeTir(j *Joueur) (maintenant, ensuite map[int]bool) {
	maintenant, ensuite = map[int]bool{}, map[int]bool{}
	ajoute := func(m map[int]bool, c int) {
		for x := range C.Points {
			if x == c {
				continue
			}
			if C.Trait[x][c] || C.Points[x].Toit && C.Points[x].Zone() == C.Points[c].Zone() {
				m[x] = true
			}
		}
	}
	if c := e.CiblePresident(); c != nil {
		if pt := e.cibleTir(c); pt >= 0 {
			ajoute(maintenant, pt)
		}
	}
	for _, r := range e.routeConnue(j, 8) {
		for k, c := range r {
			// tout près : les points reliés ; plus loin, les toits (qui
			// portent sur toute la zone) valent d'y monter à l'avance
			if k < 3 {
				ajoute(ensuite, c)
			} else {
				for _, t := range C.ToitsZone[C.Points[c].Zone()] {
					ensuite[t] = true
				}
			}
		}
	}
	if e.Discours == 0 {
		ajoute(ensuite, PtTribune)
	}
	return
}

func distA(x int, m map[int]bool) int {
	d := 99
	for y := range m {
		if dd := C.Dist[x][y]; dd >= 0 && dd < d {
			d = dd
		}
	}
	return d
}

func (IA) terros(e *Etat, j *Joueur, q Quota, opts []Mouvement) int {
	maint, ens := e.postesDeTir(j)
	val := func(id, x int) float64 {
		if x < 0 {
			return 0
		}
		v := 0.0
		if maint[x] {
			v += 30
		}
		if ens[x] {
			v += 15
		}
		return v - float64(min(distA(x, ens), 20))
	}
	sorties := map[int]bool{}
	for i, p := range C.Points {
		if p.Sortie {
			sorties[i] = true
		}
	}
	best, bv := -1, 0.2
	for i, m := range opts {
		p := e.Pions[m.Pion]
		var gain float64
		switch {
		case m.Type == MvtFuite:
			gain = 100
		case m.Type != MvtPas:
			continue
		case e.EstTerroriste(p.ID) && e.Revele[p.ID] == RTerro:
			// démasqué : il file vers la sortie la plus proche
			gain = float64(distA(p.Pos, sorties)-distA(m.Vers, sorties)) * 10
		case e.EstTerroriste(p.ID):
			gain = (val(p.ID, m.Vers) - val(p.ID, p.Pos)) * 2
		case p.Genre == Policier && !e.EstTerroriste(p.ID):
			// 3 joueurs : écarter les policiers qui serrent un terroriste (les
			// autres escortent le cortège jusqu'aux tireurs : on les laisse)
			for _, t := range e.Terroristes {
				if l := e.Lieu(e.Pions[t]); l >= 0 && e.Pions[t].Actif() && C.Dist[l][p.Pos] <= 1 {
					gain += float64(C.Dist[l][m.Vers]-C.Dist[l][p.Pos]) * 2
				}
			}
		default:
			// les leurres imitent les terroristes
			gain = (val(p.ID, m.Vers) - val(p.ID, p.Pos)) * 0.4
			gain += e.Rng.Float64() * 0.3
		}
		if gain > bv {
			best, bv = i, gain
		}
	}
	return best
}

func (IA) agents(e *Etat, j *Joueur, opts []Mouvement) int {
	// la cible : un terroriste démasqué, sinon le plus suspect, sinon la
	// voiture
	cible := -1
	bp := 0.25
	for _, p := range e.Pions {
		if !p.Actif() || !p.Genre.APied() {
			continue
		}
		pr := e.ProbTerro(j, p.ID)
		if e.Revele[p.ID] == RTerro {
			pr = 2
		}
		if pr > bp {
			bp, cible = pr, e.Lieu(p)
		}
	}
	if cible < 0 {
		if c := e.CiblePresident(); c != nil {
			cible = e.cibleTir(c)
		}
	}
	best, bv := -1, 0.2
	for i, m := range opts {
		if m.Type != MvtPas || cible < 0 {
			continue
		}
		p := e.Pions[m.Pion]
		var gain float64
		if e.EstAgent(p.ID) {
			d0, d1 := C.Dist[p.Pos][cible], C.Dist[m.Vers][cible]
			if d0 == 1 && C.Trait[p.Pos][cible] {
				continue // déjà à portée
			}
			gain = float64(d0-d1) * 3
			if d1 == 1 && C.Trait[m.Vers][cible] {
				gain += 3
			}
		} else {
			gain = 0.25 + e.Rng.Float64()*0.2 // camouflage
		}
		if gain > bv {
			best, bv = i, gain
		}
	}
	return best
}
