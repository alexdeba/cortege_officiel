package jeu

import (
	"fmt"
	"sort"
)

// Le déroulement de la partie : préparation (rôles secrets, parcours, pose
// des pions), puis les tours de jeu (actions a à h de la règle).

// Jouer : la partie entière ; les contrôleurs sont appelés au fil des
// décisions.
func (e *Etat) Jouer() {
	e.preparer()
	for !e.Fini {
		e.Tour++
		e.tour()
		e.verifierFin()
	}
	e.terminer()
}

// candidatsSecrets : pions où le rôle choisit ses deux (ou quatre) pions.
func (e *Etat) candidatsSecrets(r Role) (cands []int, n int) {
	impair := r == RTerro
	for _, p := range e.Pions {
		ok := p.Genre == Manifestant || (e.NbJ == 2 && p.Genre == Policier)
		if ok && (p.Num%2 == 1) == impair {
			cands = append(cands, p.ID)
		}
	}
	n = 2
	if e.NbJ == 2 {
		n = 4
	}
	return
}

// secretsValides : deux manifestants (et, à deux joueurs, deux policiers)
// de la bonne parité, tous différents.
func (e *Etat) secretsValides(r Role, choix []int) bool {
	cands, n := e.candidatsSecrets(r)
	if len(choix) != n {
		return false
	}
	vu := map[int]bool{}
	nPol := 0
	for _, id := range choix {
		if vu[id] || !contient(cands, id) {
			return false
		}
		vu[id] = true
		if e.Pions[id].Genre == Policier {
			nPol++
		}
	}
	return e.NbJ != 2 || nPol == 2
}

func (e *Etat) tirerSecrets(r Role) []int {
	cands, n := e.candidatsSecrets(r)
	var pol, man []int
	for _, id := range cands {
		if e.Pions[id].Genre == Policier {
			pol = append(pol, id)
		} else {
			man = append(man, id)
		}
	}
	e.Rng.Shuffle(len(pol), func(i, j int) { pol[i], pol[j] = pol[j], pol[i] })
	e.Rng.Shuffle(len(man), func(i, j int) { man[i], man[j] = man[j], man[i] })
	if n == 4 {
		return []int{man[0], man[1], pol[0], pol[1]}
	}
	return []int{man[0], man[1]}
}

func (e *Etat) preparer() {
	e.Phase = PhPreparation
	for _, r := range []Role{RAgent, RTerro} {
		j := e.ParRole[r]
		cands, n := e.candidatsSecrets(r)
		choix := j.Ctrl.Secrets(e, j, r, cands, n)
		if !e.secretsValides(r, choix) {
			choix = e.tirerSecrets(r)
		}
		if r == RAgent {
			e.Agents = choix
		} else {
			e.Terroristes = choix
		}
	}
	jp := e.ParRole[RPolice]
	e.CouleurAller = jp.Ctrl.Couleur(e, jp, false)
	if e.CouleurAller != Jaune && e.CouleurAller != Rouge && e.CouleurAller != Bleu {
		e.CouleurAller = Rouge
	}
	e.Voiture.Pos = PtHotelDeVille
	e.Helico.Pos = PtHelicoAeroport
	for _, r := range []Role{RPolice, RManif} {
		j := e.ParRole[r]
		prop := e.ProposerPlacement(r)
		pl := j.Ctrl.Placer(e, j, r, prop)
		if e.ValiderPlacement(r, pl) != "" {
			pl = prop
		}
		e.appliquerPlacement(pl)
	}
	e.dire("Les pions sont en place. Le cortège va quitter l'Hôtel de Ville.")
}

// ---------------------------------------------------------------- pose

// posable : un point au sol, non colorié, hors tribune et hôtel de ville.
func posable(p int) bool {
	pc := C.Points[p]
	return !pc.Toit && pc.Couleurs == 0 && p != PtTribune
}

// ValiderPlacement : "" si la pose est légale, sinon la raison.
func (e *Etat) ValiderPlacement(r Role, pl []Placement) string {
	par := map[int]int{}
	vehic := map[int]int{}
	places := map[int]bool{}
	for _, x := range pl {
		if x.Pion < 0 || x.Pion >= len(e.Pions) || places[x.Pion] {
			return "pion inconnu ou posé deux fois"
		}
		places[x.Pion] = true
		p := e.Pions[x.Pion]
		if x.Point < 0 || x.Point >= len(C.Points) {
			return "point inconnu"
		}
		switch {
		case r == RPolice && (p.Genre == Policier || p.Genre == VoiturePolice || p.Genre == Moto):
		case r == RManif && p.Genre == Manifestant:
		default:
			return fmt.Sprintf("le %s n'est pas à vous", p.Nom())
		}
		if !posable(x.Point) {
			return "les pions se posent sur des points de rue non coloriés"
		}
		if r == RManif && C.Points[x.Point].GP {
			return "les manifestants ne peuvent pas être dans la Grand-Place"
		}
		if p.Genre.Vehicule() {
			vehic[x.Point]++
			if vehic[x.Point] > 1 {
				return "un seul véhicule par point"
			}
		} else if x.Dans < 0 {
			par[x.Point]++
			if par[x.Point] > 2 {
				return "deux pions au plus par point"
			}
		}
		if r == RManif && (e.Compte(x.Point, Policier) > 0 || e.VehiculeSur(x.Point) != nil) {
			return "jamais sur un point où se trouve la police"
		}
	}
	n := 0
	for _, p := range e.Pions {
		if r == RPolice && (p.Genre == Policier || p.Genre == VoiturePolice || p.Genre == Moto) || r == RManif && p.Genre == Manifestant {
			n++
		}
	}
	if len(places) != n {
		return "tous les pions doivent être posés"
	}
	// passagers : dans un véhicule posé au même point, capacité respectée
	cap := map[int]int{}
	for _, x := range pl {
		if x.Dans >= 0 {
			v := e.Pions[x.Dans]
			ok := false
			for _, y := range pl {
				if y.Pion == x.Dans && y.Point == x.Point {
					ok = true
				}
			}
			cap[x.Dans]++
			if !ok || !v.Genre.Vehicule() || cap[x.Dans] > v.Genre.Capacite() || e.Pions[x.Pion].Genre != Policier {
				return "passager mal placé"
			}
		}
	}
	return ""
}

func (e *Etat) appliquerPlacement(pl []Placement) {
	for _, x := range pl {
		p := e.Pions[x.Pion]
		p.Arrivee = 0
		if x.Dans >= 0 {
			p.Pos, p.Dans = -1, x.Dans
		} else {
			p.Pos = x.Point
		}
	}
}

// ProposerPlacement : une pose raisonnable, celle des joueurs de
// l'ordinateur (et la proposition faite au joueur humain).
func (e *Etat) ProposerPlacement(r Role) []Placement {
	if r == RPolice {
		return e.poserPolice()
	}
	return e.poserManifs()
}

// poserPolice : les policiers jalonnent le parcours choisi (points non
// coloriés voisins du chemin, dans l'ordre), puis la Grand-Place et les
// départs du retour ; voitures et motos près de l'Hôtel de Ville, avec des
// policiers à bord.
func (e *Etat) poserPolice() []Placement {
	par := map[int]int{}
	var pl []Placement
	var libres []*Pion
	for _, p := range e.Pions {
		if p.Genre == Policier {
			libres = append(libres, p)
		}
	}
	if e.NbJ == 2 {
		// des terroristes se cachent parmi les policiers impairs : les pairs
		// (dont les agents) bordent le parcours
		sort.SliceStable(libres, func(a, b int) bool { return libres[a].Num%2 == 0 && libres[b].Num%2 == 1 })
	}
	chemin := append([]int{PtHotelDeVille}, C.Chemin(e.CouleurAller, false, PtHotelDeVille)...)
	var cibles []int
	vu := map[int]bool{}
	ajoute := func(q int) {
		// les points 1-2-3 de la Grand-Place restent libres pour la voiture
		if !vu[q] && posable(q) && !C.Points[q].Toit && q != PtGP1 && q != PtGP2 && q != PtGP3 {
			vu[q] = true
			cibles = append(cibles, q)
		}
	}
	for _, c := range chemin {
		for _, q := range C.Voisins[c] {
			if C.Trait[c][q] {
				ajoute(q)
			}
		}
	}
	// le parcours d'abord : deux policiers par point voisin, puis un cran
	// plus loin, puis les départs du retour
	premiers := len(cibles)
	for _, c := range chemin {
		for _, q := range C.Voisins[c] {
			for _, r := range C.Voisins[q] {
				if C.Trait[c][q] && C.Trait[q][r] {
					ajoute(r)
				}
			}
		}
	}
	for _, col := range CouleursParcours {
		ch := C.Chemin(col, true, PtGP3)
		for _, c := range ch[:min(4, len(ch))] {
			for _, q := range C.Voisins[c] {
				ajoute(q)
			}
		}
	}
	i := 0
	for passe := range 3 {
		for k, q := range cibles {
			if passe == 0 && k >= premiers {
				break
			}
			if i < len(libres) && par[q] < 2 && (passe == 0 || par[q] == 0 || passe == 2) {
				pl = append(pl, Placement{Pion: libres[i].ID, Point: q, Dans: -1})
				par[q]++
				i++
			}
		}
	}
	// le reste : les points libres les plus proches du parcours
	for _, q := range e.pointsParProximite(chemin) {
		if q == PtGP1 || q == PtGP2 || q == PtGP3 {
			continue
		}
		for i < len(libres) && par[q] < 2 && posable(q) {
			pl = append(pl, Placement{Pion: libres[i].ID, Point: q, Dans: -1})
			par[q]++
			i++
		}
	}
	// véhicules : sur les points de policiers les plus proches de l'Hôtel de
	// Ville, les policiers montent à bord
	var vehic []*Pion
	for _, p := range e.Pions {
		if p.Genre == VoiturePolice || p.Genre == Moto {
			vehic = append(vehic, p)
		}
	}
	pts := make([]int, 0, len(par))
	for q := range par {
		pts = append(pts, q)
	}
	sort.Slice(pts, func(a, b int) bool {
		da, db := C.Dist[PtHotelDeVille][pts[a]], C.Dist[PtHotelDeVille][pts[b]]
		if da != db {
			return da < db
		}
		return pts[a] < pts[b]
	})
	k := 0
	for _, v := range vehic {
		for k < len(pts) && par[pts[k]] == 0 {
			k++
		}
		if k >= len(pts) {
			break
		}
		q := pts[k]
		pl = append(pl, Placement{Pion: v.ID, Point: q, Dans: -1})
		// un policier du point monte à bord
		for n := range pl {
			if pl[n].Point == q && pl[n].Dans < 0 && e.Pions[pl[n].Pion].Genre == Policier {
				pl[n].Dans = v.ID
				break
			}
		}
		k++
	}
	return pl
}

// pointsParProximite : les points posables triés par distance au chemin.
func (e *Etat) pointsParProximite(chemin []int) []int {
	var pts []int
	d := map[int]int{}
	for i := range C.Points {
		if !posable(i) {
			continue
		}
		m := 99
		for _, c := range chemin {
			if x := C.Dist[c][i]; x >= 0 && x < m {
				m = x
			}
		}
		d[i] = m
		pts = append(pts, i)
	}
	sort.SliceStable(pts, func(a, b int) bool { return d[pts[a]] < d[pts[b]] })
	return pts
}

// poserManifs : les manifestants ne connaissent pas encore le parcours ; ils
// se massent près des trois, et surtout de leur tronc commun.
func (e *Etat) poserManifs() []Placement {
	poids := map[int]float64{}
	for _, col := range CouleursParcours {
		for k, c := range C.Chemin(col, false, PtHotelDeVille) {
			for i := range C.Points {
				if d := C.Dist[c][i]; d >= 1 && d <= 2 {
					poids[i] += 1 / float64(d) / (1 + 0.05*float64(k))
				}
			}
		}
	}
	var pts []int
	for i := range C.Points {
		if posable(i) && !C.Points[i].GP && !e.PolicePresente(i) && e.VehiculeSur(i) == nil && e.Compte(i, Policier) == 0 {
			pts = append(pts, i)
		}
	}
	sort.SliceStable(pts, func(a, b int) bool {
		if poids[pts[a]] != poids[pts[b]] {
			return poids[pts[a]] > poids[pts[b]]
		}
		return pts[a] < pts[b]
	})
	var pl []Placement
	i := 0
	for _, q := range pts {
		for range 2 {
			if i < NbManifs {
				pl = append(pl, Placement{Pion: NbPoliciers + i, Point: q, Dans: -1})
				i++
			}
		}
	}
	return pl
}

// ---------------------------------------------------------------- tour

func (e *Etat) phase(ph Phase) {
	e.Phase = ph
	e.signaler(Ev{Type: EvPhase, Phase: ph})
}

func (e *Etat) chasse() bool { return e.TourFin > 0 }

func (e *Etat) tour() {
	e.signaler(Ev{Type: EvTour})
	clear(e.derniersMvtsT)
	if e.Tour == 1 {
		e.AllerPublic = true
		e.dire("La police dévoile sa feuille : le cortège suivra le parcours %s.", e.CouleurAller.Nom())
	}
	if !e.chasse() {
		e.phase(PhVoiture)
		e.avancerVoiture()
		if e.cortegeFini() {
			return
		}
		e.phase(PhAttentat)
		e.attentats(PhAttentat)
	}
	e.phase(PhArrestation)
	e.tirs(RAgent, PhArrestation)
	e.tirs(RPolice, PhArrestation)

	e.phase(PhPolice)
	e.deplacements(Quota{Role: RPolice, Policiers: true, Vehicules: true})
	e.tirs(RPolice, PhPolice)

	e.phase(PhManifs)
	e.deplacements(Quota{Role: RManif, Manifs: true, Vehicules: true})
	if e.NbJ == 3 { // la provocation : deux policiers déplacés par les manifestants
		e.deplacements(Quota{Role: RManif, Policiers: true, MaxPol: 2})
	}

	e.phase(PhAgents)
	e.deplacements(Quota{Role: RAgent, Manifs: true, Parite: 2, MaxManif: 4})
	e.tirs(RAgent, PhAgents)

	e.phase(PhTerros)
	switch e.NbJ {
	case 3:
		e.deplacements(Quota{Role: RTerro, Manifs: true, Policiers: true, MaxManif: 4, MaxPol: 4, Fuite: true})
	case 2:
		e.deplacements(Quota{Role: RTerro, Manifs: true, Policiers: true, Parite: 1, MaxTotal: 4, Fuite: true})
	default:
		e.deplacements(Quota{Role: RTerro, Manifs: true, Parite: 1, MaxManif: 4, Fuite: true})
	}
	e.majSoupcons()
	if !e.chasse() {
		e.attentats(PhTerros)
	}
}

// deplacements : le joueur du rôle déplace ses pions un à un.
func (e *Etat) deplacements(q Quota) {
	j := e.ParRole[q.Role]
	bouges := map[int]bool{}
	e.NbMvts = 0
	for range 400 {
		opts := e.Mouvements(q, bouges)
		if len(opts) == 0 {
			return
		}
		k := j.Ctrl.Deplacer(e, j, q, opts)
		if k < 0 || k >= len(opts) {
			return
		}
		e.NbMvts++
		for _, id := range e.Appliquer(opts[k]) {
			bouges[id] = true
			if q.Role == RTerro {
				e.derniersMvtsT[id]++
			}
		}
		if e.Fini {
			return
		}
	}
}

// tirs : le rôle tire autant qu'il veut, chaque tireur une fois.
func (e *Etat) tirs(r Role, ph Phase) {
	j := e.ParRole[r]
	deja := map[int]bool{}
	for range 40 {
		opts := e.TirsPossibles(r, ph, deja)
		if len(opts) == 0 {
			return
		}
		k := j.Ctrl.Tirer(e, j, ph, opts)
		if k < 0 || k >= len(opts) {
			return
		}
		deja[opts[k].Tireur] = true
		e.Tirer(opts[k], r)
	}
}

// attentats : les terroristes à portée tirent ; riposte aussitôt.
func (e *Etat) attentats(ph Phase) {
	avant := len(e.Lignes)
	tire := false
	j := e.ParRole[RTerro]
	deja := map[int]bool{}
	for range 4 {
		opts := e.TirsPossibles(RTerro, ph, deja)
		if len(opts) == 0 {
			break
		}
		k := j.Ctrl.Tirer(e, j, ph, opts)
		if k < 0 || k >= len(opts) {
			break
		}
		deja[opts[k].Tireur] = true
		e.Tirer(opts[k], RTerro)
		tire = true
		if e.PresTue {
			break
		}
	}
	_ = avant
	if tire {
		e.phase(PhRiposte)
		e.tirs(RAgent, PhRiposte)
		e.tirs(RPolice, PhRiposte)
		e.Phase = ph
	}
}

// ---------------------------------------------------------------- voiture

// ProchainsPoints : les points où la voiture peut avancer (libres ou non).
func (e *Etat) ProchainsPoints() []int {
	v := e.Voiture
	if v.Pos < 0 || e.Arrivee {
		return nil
	}
	switch e.Discours {
	case 1:
		return []int{PtGP2}
	case 2:
		return []int{PtGP3}
	}
	retour := e.Discours >= 3
	col := e.CouleurAller
	if retour {
		col = e.CouleurRetour
	}
	d := C.DistParcours(col, retour, v.Pos)
	var l []int
	for _, q := range C.Voisins[v.Pos] {
		if C.Trait[v.Pos][q] && d > 0 && C.DistParcours(col, retour, q) == d-1 {
			l = append(l, q)
		}
	}
	return l
}

// librePourVoiture : la voiture peut-elle entrer en q ? Entrer dans la
// Grand-Place, c'est se ranger au point 1.
func (e *Etat) librePourVoiture(q int) bool {
	if e.Discours == 0 && C.Points[q].GP {
		return !e.Occupe(PtGP1) && !e.Occupe(q) || q == PtGP1 && !e.Occupe(q)
	}
	return !e.Occupe(q)
}

func (e *Etat) avancerVoiture() {
	v := e.Voiture
	pas := 1
	if e.DoubleAvance && e.Discours != 1 && e.Discours != 2 {
		pas = 2
	}
	e.DoubleAvance = false
	for i := range pas {
		opts := e.ProchainsPoints()
		var libres []int
		for _, q := range opts {
			if e.librePourVoiture(q) {
				libres = append(libres, q)
			}
		}
		if len(libres) == 0 {
			if i == 0 {
				e.Blocages++
				e.signaler(Ev{Type: EvBlocage, Texte: "Le cortège est bloqué !"})
				e.marquer("Blocage du cortège officiel", ScBlocage)
			}
			return
		}
		q := libres[0]
		if len(libres) > 1 {
			jp := e.ParRole[RPolice]
			k := jp.Ctrl.Avancer(e, jp, libres)
			if k >= 0 && k < len(libres) {
				q = libres[k]
			}
		}
		de := v.Pos
		if e.Discours == 3 {
			e.Discours = 4
			e.RetourPublic = true
			e.dire("Le cortège quitte la Grand-Place par le parcours %s.", e.CouleurRetour.Nom())
		}
		v.Pos = q
		e.signaler(Ev{Type: EvMouvement, Pion: v.ID, Chemin: []int{de, q}})
		switch {
		case e.Discours == 0 && C.Points[q].GP:
			if q != PtGP1 {
				v.Pos = PtGP1
				e.signaler(Ev{Type: EvMouvement, Pion: v.ID, Chemin: []int{q, PtGP1}})
			}
			e.Discours = 1
			e.Pres.Pos = PtTribune
			e.dire("Le cortège arrive sur la Grand-Place : le président monte à la tribune pour son discours.")
			return
		case e.Discours == 1:
			e.Discours = 2
			jp := e.ParRole[RPolice]
			e.CouleurRetour = jp.Ctrl.Couleur(e, jp, true)
			if e.CouleurRetour != Jaune && e.CouleurRetour != Rouge && e.CouleurRetour != Bleu {
				e.CouleurRetour = Rouge
			}
			e.dire("Le discours continue ; la police a choisi en secret le parcours du retour.")
			return
		case e.Discours == 2:
			e.Discours = 3
			e.Pres.Pos = -1
			e.dire("Fin du discours : le président regagne sa voiture.")
			return
		case q == PtAeroport:
			e.Arrivee = true
			if !e.PresTue && !e.PresBlesse {
				e.dire("Le président arrive indemne à l'aéroport !")
				e.marquer("Président indemne à l'aéroport", ScPresAeroport)
			} else {
				e.dire("Le président, blessé, arrive à l'aéroport.")
			}
			return
		}
	}
}

// ---------------------------------------------------------------- fin

func (e *Etat) cortegeFini() bool {
	return e.Arrivee || e.PresTue || e.Tour >= e.TourMax
}

// terroristesEnFuite : terroristes démasqués encore libres dans la ville.
func (e *Etat) terroristesEnFuite() int {
	n := 0
	for _, id := range e.Terroristes {
		p := e.Pions[id]
		if p.Actif() && e.Revele[id] == RTerro {
			n++
		}
	}
	return n
}

func (e *Etat) verifierFin() {
	if !e.cortegeFini() {
		return
	}
	if e.TourFin == 0 {
		e.TourFin = e.Tour
		if e.Attentat && e.terroristesEnFuite() > 0 {
			e.dire("Le cortège est passé : la chasse aux terroristes commence (%d tours).", ToursChasse)
			return
		}
	}
	if e.Tour-e.TourFin < ToursChasse && e.Attentat && e.terroristesEnFuite() > 0 {
		return
	}
	e.Fini = true
}

func (e *Etat) terminer() {
	e.Fini = true
	if e.Blocages == 0 && e.Arrivee {
		e.marquer("Le cortège n'a jamais été bloqué", ScJamaisBloque)
	}
	if e.Attentat {
		n := 0
		for _, id := range e.Terroristes {
			p := e.Pions[id]
			// enfui, ou resté anonyme dans la foule
			if p.Enfui || p.Actif() && e.Revele[id] != RTerro {
				n++
			}
		}
		switch n {
		case 2:
			e.marquer("Les deux terroristes s'échappent", ScDeuxEchappent)
		case 1:
			e.marquer("Un seul terroriste s'échappe", ScUnSeulEchappe)
		}
	}
	e.signaler(Ev{Type: EvFin})
}

// TirageSecret : des pions secrets tirés au hasard (bouton « au hasard »).
func (e *Etat) TirageSecret(r Role) []int { return e.tirerSecrets(r) }

// SecretsValides : le choix respecte-t-il la règle des numéros ?
func (e *Etat) SecretsValides(r Role, choix []int) bool { return e.secretsValides(r, choix) }
