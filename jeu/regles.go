package jeu

import "fmt"

// Règles de déplacement, de contact et de tir.

// interditAuxPions : la voiture présidentielle, le président et la tribune
// officielle (réservée au président) ne laissent la place à personne.
func (e *Etat) interditAuxPions(b int) bool {
	return b == PtTribune || b == e.Voiture.Pos || b == e.Pres.Pos
}

// accesToit : passage entre la rue et un toit (par une flèche) ou entre deux
// toits d'un même immeuble.
func (e *Etat) accesToit(p *Pion, a, b int) bool {
	pa, pb := C.Points[a], C.Points[b]
	switch {
	case !pa.Toit && pb.Toit:
		if !C.Fleche[a][b] {
			return false
		}
		// un manifestant doit avoir séjourné un tour sur le point fléché
		return p.Genre != Manifestant || p.Arrivee < e.Tour
	case pa.Toit && !pb.Toit:
		return C.Fleche[b][a]
	}
	return true
}

// PeutAller : un pas du pion à pied p de a vers b (b voisin de a).
func (e *Etat) PeutAller(p *Pion, a, b int, limite bool) bool {
	if e.interditAuxPions(b) || !e.accesToit(p, a, b) {
		return false
	}
	v := e.VehiculeSur(b)
	switch p.Genre {
	case Policier:
		if e.Compte(b, Policier) >= 2 {
			return false
		}
		// un véhicule pris par un manifestant : contact
		return true
	case Manifestant:
		if e.PolicePresente(b) || C.Points[a].GP != C.Points[b].GP {
			return false
		}
		if v != nil && len(e.Passagers(v)) > 0 && e.Passagers(v)[0].Genre != Manifestant {
			return false
		}
		if v != nil && v.Vol >= 0 {
			return false
		}
		return !limite || len(e.ManifsSur(b)) < 2
	}
	return false
}

// deplacable : le pion peut-il bouger (vivant, pas blessé, sur le plan) ?
func deplacable(p *Pion) bool { return p.Actif() && p.Blessures == 0 }

func pariteOK(q Quota, p *Pion) bool {
	return q.Parite == 0 || (q.Parite == 1) == (p.Num%2 == 1)
}

// unite : le pion qui compte pour le quota (le conducteur d'un véhicule
// pris, le véhicule de police lui-même).
func (e *Etat) conducteur(v *Pion) *Pion {
	ps := e.Passagers(v)
	if len(ps) == 0 {
		return nil
	}
	return ps[0]
}

// Mouvements : les déplacements possibles pour le quota, les pions de
// bouges ayant déjà joué.
func (e *Etat) Mouvements(q Quota, bouges map[int]bool) []Mouvement {
	nPol, nMan := 0, 0
	for id := range bouges {
		switch e.Pions[id].Genre {
		case Policier:
			nPol++
		case Manifestant:
			nMan++
		}
	}
	tot := nPol + nMan
	polOK := q.Policiers && (q.MaxPol == 0 || nPol < q.MaxPol) && (q.MaxTotal == 0 || tot < q.MaxTotal)
	manOK := q.Manifs && (q.MaxManif == 0 || nMan < q.MaxManif) && (q.MaxTotal == 0 || tot < q.MaxTotal)
	var l []Mouvement
	for _, p := range e.Pions {
		if !deplacable(p) || bouges[p.ID] {
			continue
		}
		switch {
		case p.Genre == Policier && polOK && pariteOK(q, p):
			if p.Dans >= 0 {
				if q.Vehicules && q.Role == RPolice {
					l = append(l, e.descentes(p)...)
				}
				continue
			}
			if p.Pos < 0 {
				continue
			}
			for _, b := range C.Voisins[p.Pos] {
				// un policier manipulé par l'adversaire (provocation, faux
				// ordres) n'est pas envoyé sur le parcours du cortège
				if q.Role != RPolice && (C.Points[b].Couleurs != 0 || b == PtGP1 || b == PtGP2 || b == PtGP3) {
					continue
				}
				if e.PeutAller(p, p.Pos, b, true) {
					l = append(l, Mouvement{Type: MvtPas, Pion: p.ID, Vers: b, Chemin: []int{p.Pos, b}})
					if q.Role == RPolice && e.peutMonter(p, b) {
						l = append(l, Mouvement{Type: MvtMonter, Pion: p.ID, Vers: b, Chemin: []int{p.Pos, b}})
					}
				}
			}
			if q.Role == RPolice && e.peutMonter(p, p.Pos) {
				l = append(l, Mouvement{Type: MvtMonter, Pion: p.ID, Vers: p.Pos, Chemin: []int{p.Pos}})
			}
		case p.Genre == Manifestant && manOK && pariteOK(q, p):
			if p.Dans >= 0 {
				v := e.Pions[p.Dans]
				if q.Vehicules && v.PrisTour < e.Tour && !bouges[v.ID] {
					l = append(l, e.conduites(v, p)...)
				}
				continue
			}
			if p.Pos < 0 {
				continue
			}
			for _, b := range C.Voisins[p.Pos] {
				if e.PeutAller(p, p.Pos, b, true) {
					l = append(l, Mouvement{Type: MvtPas, Pion: p.ID, Vers: b, Chemin: []int{p.Pos, b}})
				}
			}
			if q.Fuite && e.Attentat && C.Points[p.Pos].Sortie && e.EstTerroriste(p.ID) {
				l = append(l, Mouvement{Type: MvtFuite, Pion: p.ID, Vers: p.Pos, Chemin: []int{p.Pos}})
			}
		case p.Genre.Vehicule() && q.Vehicules && q.Role == RPolice && polOK:
			c := e.conducteur(p)
			if c != nil && c.Genre == Policier && !bouges[c.ID] {
				l = append(l, e.conduites(p, c)...)
			}
		}
	}
	return l
}

// peutMonter : le policier p peut-il monter dans le véhicule du point b ?
func (e *Etat) peutMonter(p *Pion, b int) bool {
	v := e.VehiculeSur(b)
	if v == nil || v.Vol >= 0 {
		return false
	}
	ps := e.Passagers(v)
	if len(ps) > 0 && ps[0].Genre != Policier {
		return false
	}
	return len(ps) < v.Genre.Capacite()
}

// descentes : où peut aller un policier qui sort de son véhicule.
func (e *Etat) descentes(p *Pion) []Mouvement {
	v := e.Pions[p.Dans]
	if v.Vol >= 0 || v.Pos < 0 {
		return nil
	}
	cibles := []int{v.Pos}
	if v.Genre == Helico {
		// le toit et les points voisins de l'immeuble
		for _, q := range C.Voisins[v.Pos] {
			cibles = append(cibles, q)
		}
	}
	var l []Mouvement
	for _, b := range cibles {
		if b != v.Pos && e.interditAuxPions(b) {
			continue
		}
		if e.Compte(b, Policier) < 2 && len(e.ManifsSur(b)) == 0 {
			l = append(l, Mouvement{Type: MvtDescente, Pion: p.ID, Vers: b, Chemin: []int{v.Pos, b}})
		}
	}
	return l
}

// conduites : les trajets du véhicule v (conducteur c).
func (e *Etat) conduites(v, c *Pion) []Mouvement {
	if v.Genre == Helico {
		return e.vols(v, c)
	}
	police := c.Genre == Policier
	type etape struct{ p, d int }
	prec := map[int]int{v.Pos: -1}
	file := []etape{{v.Pos, 0}}
	var l []Mouvement
	for len(file) > 0 {
		x := file[0]
		file = file[1:]
		if x.d == v.Genre.Portee() {
			continue
		}
		for _, b := range C.Voisins[x.p] {
			if _, vu := prec[b]; vu || C.Points[b].Toit || e.interditAuxPions(b) {
				continue
			}
			if !police && C.Points[x.p].GP != C.Points[b].GP {
				continue
			}
			prec[b] = x.p
			chemin := func() []int {
				var ch []int
				for y := b; y >= 0; y = prec[y] {
					ch = append([]int{y}, ch...)
				}
				return ch
			}
			arrivee := e.VehiculeSur(b) == nil
			if police {
				arrivee = arrivee && e.Compte(b, Policier) <= 2
			} else {
				arrivee = arrivee && !e.PolicePresente(b) && len(e.ManifsSur(b)) < 2
			}
			if arrivee {
				l = append(l, Mouvement{Type: MvtConduite, Pion: v.ID, Vers: b, Chemin: chemin()})
			}
			// on ne traverse que des points libres
			if !e.Occupe(b) {
				file = append(file, etape{b, x.d + 1})
			}
		}
	}
	return l
}

// vols : l'hélicoptère se pose sur un toit à H ou survole une zone.
func (e *Etat) vols(v, c *Pion) []Mouvement {
	police := c.Genre == Policier
	z := e.ZonePion(v)
	enVol := v.Vol >= 0
	var l []Mouvement
	for i, p := range C.Points {
		if !p.H || i == v.Pos || DistZones(z, p.Zone()) > 3 {
			continue
		}
		if e.VehiculeSur(i) != nil || e.interditAuxPions(i) {
			continue
		}
		if police && len(e.ManifsSur(i)) > 0 || !police && e.PolicePresente(i) {
			continue
		}
		l = append(l, Mouvement{Type: MvtVol, Pion: v.ID, Vers: i, Zone: p.Zone()})
	}
	portee := 3
	if enVol {
		portee = 4
	}
	for zz := range NbCol * NbLig {
		if DistZones(z, zz) <= portee && (zz != z || !enVol) {
			l = append(l, Mouvement{Type: MvtVol, Pion: v.ID, Vers: -1, Zone: zz})
		}
	}
	return l
}

// ---------------------------------------------------------------- application

// Appliquer : exécute un mouvement légal, contacts compris ; renvoie les
// pions qui ont bougé.
func (e *Etat) Appliquer(m Mouvement) []int {
	p := e.Pions[m.Pion]
	bouges := []int{p.ID}
	switch m.Type {
	case MvtPas:
		v := e.VehiculeSur(m.Vers)
		p.Pos, p.Arrivee = m.Vers, e.Tour
		if p.Genre == Manifestant && v != nil && len(e.Passagers(v)) == 0 && !e.PolicePresente(m.Vers) && len(e.ManifsSur(m.Vers)) == 1 {
			// le véhicule vide et sans garde est pris
			p.Dans, p.Pos = v.ID, -1
			v.PrisTour = e.Tour
			e.dire("Le %s s'empare de la %s !", p.Nom(), v.Nom())
		}
		e.signaler(Ev{Type: EvMouvement, Pion: p.ID, Chemin: m.Chemin})
		if p.Genre == Policier {
			e.contact(m.Vers)
		}
	case MvtMonter:
		v := e.VehiculeSur(m.Vers)
		from := p.Pos
		p.Pos, p.Dans, p.Arrivee = -1, v.ID, e.Tour
		e.signaler(Ev{Type: EvMouvement, Pion: p.ID, Chemin: []int{from, m.Vers}})
		e.dire("Le %s monte dans la %s.", p.Nom(), v.Nom())
		e.contact(m.Vers)
	case MvtDescente:
		v := e.Pions[p.Dans]
		p.Dans, p.Pos, p.Arrivee = -1, m.Vers, e.Tour
		e.signaler(Ev{Type: EvMouvement, Pion: p.ID, Chemin: []int{v.Pos, m.Vers}})
	case MvtConduite:
		p.Pos, p.Arrivee = m.Vers, e.Tour
		for _, q := range e.Passagers(p) {
			bouges = append(bouges, q.ID)
			q.Arrivee = e.Tour
		}
		e.signaler(Ev{Type: EvMouvement, Pion: p.ID, Chemin: m.Chemin})
		if c := e.conducteur(p); c != nil && c.Genre == Policier {
			e.contact(m.Vers)
		}
	case MvtVol:
		depart := e.ZonePion(p)
		p.Pos, p.Vol, p.Arrivee = m.Vers, -1, e.Tour
		if m.Vers < 0 {
			p.Vol = m.Zone
		}
		for _, q := range e.Passagers(p) {
			bouges = append(bouges, q.ID)
		}
		e.signaler(Ev{Type: EvMouvement, Pion: p.ID, Chemin: m.Chemin, Texte: fmt.Sprintf("%s→%s", NomZone(depart), NomZone(m.Zone))})
		if m.Vers >= 0 {
			// posé : les policiers du toit montent à bord
			c := e.conducteur(p)
			for _, q := range e.APied(m.Vers) {
				if c != nil && c.Genre == Policier && q.Genre == Policier && len(e.Passagers(p)) < p.Genre.Capacite() {
					q.Pos, q.Dans = -1, p.ID
				}
			}
		}
	case MvtFuite:
		p.Pos, p.Enfui = -1, true
		e.Revele[p.ID] = RTerro
		e.signaler(Ev{Type: EvRevele, Pion: p.ID, Texte: fmt.Sprintf("Le %s était un terroriste : il quitte la ville !", p.Nom())})
	}
	return bouges
}

// contact : la police arrive sur le point b ; les manifestants qui s'y
// trouvent doivent reculer, ou sont interpellés.
func (e *Etat) contact(b int) {
	for range 8 {
		ms := e.ManifsSur(b)
		if len(ms) == 0 {
			return
		}
		m := ms[0]
		if m.Dans >= 0 { // il abandonne le véhicule pris
			m.Dans, m.Pos = -1, b
		}
		opts := e.reculs(m, b)
		if len(opts) == 0 {
			e.interpeller(m)
			continue
		}
		ph := e.Phase
		e.Phase = PhContact
		j := e.ParRole[RManif]
		k := j.Ctrl.Reculer(e, j, opts)
		e.Phase = ph
		if k < 0 || k >= len(opts) {
			k = 0
		}
		e.dire("Contact ! Le %s recule.", m.Nom())
		m.Pos, m.Arrivee = opts[k].Vers, e.Tour
		e.signaler(Ev{Type: EvMouvement, Pion: m.ID, Chemin: opts[k].Chemin})
		e.desengorger(opts[k].Vers)
	}
}

// reculs : les points voisins sans police où le manifestant m peut se
// replier (il peut momentanément y être le troisième).
func (e *Etat) reculs(m *Pion, b int) []Mouvement {
	if m.Blessures > 0 {
		return nil
	}
	var l []Mouvement
	for _, c := range C.Voisins[b] {
		if e.PeutAller(m, b, c, false) && len(e.ManifsSur(c)) < 3 && e.VehiculeSur(c) == nil {
			l = append(l, Mouvement{Type: MvtPas, Pion: m.ID, Vers: c, Chemin: []int{b, c}})
		}
	}
	return l
}

// desengorger : jamais plus de deux manifestants par point.
func (e *Etat) desengorger(c int) {
	for range 4 {
		ms := e.ManifsSur(c)
		if len(ms) <= 2 {
			return
		}
		var opts []Mouvement
		for _, m := range ms {
			if m.Dans >= 0 || m.Blessures > 0 {
				continue
			}
			for _, d := range C.Voisins[c] {
				if e.PeutAller(m, c, d, true) && e.VehiculeSur(d) == nil {
					opts = append(opts, Mouvement{Type: MvtPas, Pion: m.ID, Vers: d, Chemin: []int{c, d}})
				}
			}
		}
		if len(opts) == 0 {
			e.interpeller(ms[len(ms)-1])
			continue
		}
		j := e.ParRole[RManif]
		k := j.Ctrl.Reculer(e, j, opts)
		if k < 0 || k >= len(opts) {
			k = 0
		}
		m := e.Pions[opts[k].Pion]
		m.Pos, m.Arrivee = opts[k].Vers, e.Tour
		e.dire("Trop de monde : le %s s'écarte.", m.Nom())
		e.signaler(Ev{Type: EvMouvement, Pion: m.ID, Chemin: opts[k].Chemin})
	}
}

// interpeller : le manifestant cerné est retiré du jeu.
func (e *Etat) interpeller(m *Pion) {
	m.Retire, m.Arrete, m.Pos, m.Dans = true, true, -1, -1
	txt := fmt.Sprintf("Le %s, cerné par la police, est interpellé.", m.Nom())
	switch {
	case e.EstTerroriste(m.ID):
		e.Revele[m.ID] = RTerro
		txt += " C'était un terroriste !"
		e.signaler(Ev{Type: EvRevele, Pion: m.ID, Texte: txt})
		e.terroristeNeutralise(m)
		return
	case e.EstAgent(m.ID):
		e.Revele[m.ID] = RAgent
		txt += " C'était un agent secret."
	}
	e.signaler(Ev{Type: EvRevele, Pion: m.ID, Texte: txt})
}

func (e *Etat) terroristeNeutralise(m *Pion) {
	if e.Attentat {
		e.marquer(fmt.Sprintf("Terroriste (%s) neutralisé après tentative d'attentat", m.Nom()), ScTerroApres)
	} else {
		e.marquer(fmt.Sprintf("Terroriste (%s) neutralisé avant un attentat", m.Nom()), ScTerroAvant)
	}
}

// ---------------------------------------------------------------- tirs

// enHauteur : le tireur est sur un toit ou dans l'hélicoptère.
func (e *Etat) enHauteur(p *Pion) bool {
	if p.Dans >= 0 && e.Pions[p.Dans].Genre == Helico {
		return true
	}
	l := e.Lieu(p)
	return l >= 0 && C.Points[l].Toit
}

// APortee : le pion t peut-il tirer sur le point c ?
func (e *Etat) APortee(t *Pion, c int) bool {
	if c < 0 {
		return false
	}
	if e.enHauteur(t) {
		z := e.ZonePion(t)
		if C.Points[c].Zone() == z {
			return true
		}
		return C.Points[c].Toit && DistZones(C.Points[c].Zone(), z) <= 1
	}
	l := e.Lieu(t)
	return l >= 0 && l != c && C.Trait[l][c]
}

// cibleTir : le point visé (le lieu du pion, même dans un véhicule) ; -1 si
// le pion n'est pas atteignable (en vol).
func (e *Etat) cibleTir(p *Pion) int {
	if p.Dans >= 0 && e.Pions[p.Dans].Vol >= 0 {
		return -1
	}
	return e.Lieu(p)
}

func (e *Etat) ligneTable(t *Pion, r Role) int {
	dans := -1
	if t.Dans >= 0 {
		dans = int(e.Pions[t.Dans].Genre)
	}
	switch r {
	case RTerro:
		if dans >= 0 {
			return TTerroVehicule
		}
		if C.Points[t.Pos].Couleurs != 0 {
			return TTerroColorie
		}
		return TTerroNonColorie
	case RAgent:
		switch {
		case dans == int(Helico):
			return TAgentHelico
		case dans >= 0:
			return TAgentVehicule
		}
		return TAgentQuelconque
	}
	switch {
	case dans == int(Helico):
		return TPoliceHelico
	case dans >= 0:
		return TPoliceVehicule
	}
	return TPoliceQuelconque
}

// CiblePresident : le pion à abattre (le président à la tribune, ou la
// voiture qui le transporte), nil si hors d'atteinte.
func (e *Etat) CiblePresident() *Pion {
	if e.PresTue || e.Arrivee {
		return nil
	}
	if e.Pres.Pos >= 0 {
		return e.Pres
	}
	if e.Voiture.Pos >= 0 {
		return e.Voiture
	}
	return nil
}

// tireurs : les pions qui tirent pour le rôle r.
func (e *Etat) tireurs(j *Joueur, r Role) []*Pion {
	var l []*Pion
	for _, p := range e.Pions {
		if !p.Actif() || (p.Pos < 0 && p.Dans < 0) {
			continue
		}
		switch r {
		case RTerro:
			if e.EstTerroriste(p.ID) {
				l = append(l, p)
			}
		case RAgent:
			if e.EstAgent(p.ID) {
				l = append(l, p)
			}
		case RPolice:
			if p.Genre == Policier && e.Connu(j, p.ID) != RAgent {
				l = append(l, p)
			}
		}
	}
	return l
}

// TirsPossibles : les tirs du rôle r à la phase ph ; deja : tireurs qui ont
// déjà tiré pendant la phase.
func (e *Etat) TirsPossibles(r Role, ph Phase, deja map[int]bool) []Tir {
	j := e.ParRole[r]
	var l []Tir
	for _, t := range e.tireurs(j, r) {
		if deja[t.ID] {
			continue
		}
		vise := func(c *Pion) {
			if c.ID == t.ID || !c.Actif() {
				return
			}
			if pt := e.cibleTir(c); pt >= 0 && e.APortee(t, pt) {
				l = append(l, Tir{Tireur: t.ID, Cible: c.ID, Table: e.ligneTable(t, r)})
			}
		}
		switch {
		case r == RTerro:
			if c := e.CiblePresident(); c != nil {
				vise(c)
			}
		case ph == PhRiposte || ph == PhPolice:
			for _, c := range e.Pions {
				if e.Revele[c.ID] == RTerro && e.EstTerroriste(c.ID) && c.Actif() {
					vise(c)
				}
			}
		default: // arrestation, agents : tout pion à pied suspect
			for _, c := range e.Pions {
				if !c.Genre.APied() || e.Connu(j, c.ID) == RAgent {
					continue
				}
				if r == RPolice && c.Genre == Policier && e.Revele[c.ID] != RTerro {
					continue
				}
				vise(c)
			}
		}
	}
	return l
}

func (e *Etat) des() [2]int { return [2]int{e.Rng.IntN(6) + 1, e.Rng.IntN(6) + 1} }

// Tirer : le tir, ses conséquences et la révélation du tireur.
func (e *Etat) Tirer(t Tir, r Role) Resultat {
	tireur, cible := e.Pions[t.Tireur], e.Pions[t.Cible]
	d := e.des()
	res := TableTirs[t.Table][d[0]+d[1]-2]
	if r == RTerro || r == RAgent {
		e.Revele[tireur.ID] = r
	}
	e.signaler(Ev{Type: EvTir, Tir: t, Des: d, Result: res, Pion: tireur.ID,
		Texte: fmt.Sprintf("%s (%s) tire sur %s : %d + %d, %s.", nomTireur(tireur, r), NomsTables[t.Table], cible.Nom(), d[0], d[1], NomsResultats[res])})
	switch r {
	case RTerro:
		e.Attentat, e.DoubleAvance = true, true
		switch res {
		case Blesse:
			if e.PresBlesse {
				e.presidentTue()
			} else {
				e.PresBlesse = true
				e.dire("Le président est blessé !")
				e.marquer("Président blessé", ScPresBlesse)
			}
		case Tue:
			e.presidentTue()
		case Bavure:
			// une balle perdue : un manifestant proche de la cible
			pt := e.cibleTir(cible)
			for _, q := range append([]int{pt}, C.Voisins[pt]...) {
				for _, m := range e.ManifsSur(q) {
					if m.ID != tireur.ID {
						e.dire("Balle perdue : le %s est touché.", m.Nom())
						e.toucher(m, 1, false)
						return res
					}
				}
			}
		}
	case RAgent:
		switch res {
		case Tue:
			e.toucher(cible, 2, false)
		case Blesse:
			e.toucher(cible, 1, false)
		}
	case RPolice:
		switch res {
		case Tue, Blesse:
			g := 1
			if res == Tue {
				g = 2
			}
			e.toucher(cible, g, !e.EstTerroriste(cible.ID))
		default:
			// raté ou bavure : le manifestant voisin du visé est touché
			for _, m := range e.ManifsSur(e.cibleTir(cible)) {
				if m.ID != cible.ID {
					e.dire("Bavure : c'est le %s, à côté, qui est touché !", m.Nom())
					e.toucher(m, 1, !e.EstTerroriste(m.ID))
					break
				}
			}
		}
	}
	return res
}

func nomTireur(p *Pion, r Role) string {
	switch r {
	case RTerro:
		return "Le terroriste (" + p.Nom() + ")"
	case RAgent:
		return "L'agent secret (" + p.Nom() + ")"
	}
	return "Le " + p.Nom()
}

func (e *Etat) presidentTue() {
	e.PresTue, e.PresBlesse = true, true
	e.dire("Le président est tué !")
	e.marquer("Président tué", ScPresTue)
}

// toucher : g = 1 blessé, 2 tué. Un terroriste ou un agent blessé reste en
// jeu sans pouvoir bouger ; la seconde blessure est fatale ; tout autre
// pion touché est retiré. bavure : la police a touché un innocent.
func (e *Etat) toucher(p *Pion, g int, bavure bool) {
	if bavure {
		e.Bavures++
		e.marquer(fmt.Sprintf("Bavure (%s)", p.Nom()), ScBavure)
	}
	secret := e.EstTerroriste(p.ID) || e.EstAgent(p.ID)
	if secret {
		r := RAgent
		if e.EstTerroriste(p.ID) {
			r = RTerro
		}
		e.Revele[p.ID] = r
		p.Blessures += g
		etat := "blessé"
		if p.Blessures >= 2 {
			etat = "tué"
		}
		e.signaler(Ev{Type: EvRevele, Pion: p.ID, Texte: fmt.Sprintf("Le %s était %s : %s.", p.Nom(), map[Role]string{RAgent: "un agent secret", RTerro: "un terroriste"}[r], etat)})
		if p.Blessures >= 2 {
			e.retirer(p)
			if r == RTerro {
				e.terroristeNeutralise(p)
			}
		}
		return
	}
	e.dire("Le %s est %s et retiré du jeu.", p.Nom(), map[int]string{1: "blessé", 2: "tué"}[g])
	e.retirer(p)
}

func (e *Etat) retirer(p *Pion) {
	if p.Dans >= 0 {
		p.Dans = -1
	}
	p.Retire, p.Pos = true, -1
}
