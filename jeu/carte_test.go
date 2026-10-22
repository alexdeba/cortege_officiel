package jeu

import "testing"

// Les trois parcours relient l'Hôtel de Ville à la Grand-Place, puis le
// point 3 de la Grand-Place à l'aéroport.
func TestParcours(t *testing.T) {
	for _, col := range CouleursParcours {
		a := C.DistParcours(col, false, PtHotelDeVille)
		r := C.DistParcours(col, true, PtGP3)
		if a <= 0 || r <= 0 {
			t.Errorf("parcours %s : aller %d, retour %d", col.Nom(), a, r)
		}
		ch := C.Chemin(col, true, PtGP3)
		if len(ch) == 0 || ch[len(ch)-1] != PtAeroport {
			t.Errorf("retour %s n'aboutit pas à l'aéroport : %v", col.Nom(), ch)
		}
		t.Logf("%s : aller %d pas, retour %d pas", col.Nom(), a, r)
	}
}

func TestZones(t *testing.T) {
	cas := map[int]string{PtHotelDeVille: "G1", PtAeroport: "A7", PtGP1: "E3", PtTribune: "E3", PtHelicoAeroport: "A7"}
	for p, z := range cas {
		if n := NomZone(C.Points[p].Zone()); n != z {
			t.Errorf("point %d en %s, attendu %s", p, n, z)
		}
	}
	if !C.Points[PtHelicoAeroport].H {
		t.Error("le toit de l'aéroport doit porter un H")
	}
}
