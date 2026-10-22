package jeu

// Table des tirs de la règle : selon le tireur et sa position, le résultat
// de deux dés (2 à 12).

type Resultat int

const (
	Rate Resultat = iota
	Blesse
	Tue
	Bavure
)

var NomsResultats = [...]string{"raté", "blessé", "tué", "bavure"}

// Lignes de la table des tirs.
const (
	TTerroNonColorie = iota
	TTerroColorie
	TTerroVehicule
	TAgentQuelconque
	TAgentVehicule
	TAgentHelico
	TPoliceQuelconque
	TPoliceVehicule
	TPoliceHelico
)

var NomsTables = [...]string{
	"le terroriste tire d'un point non colorié",
	"le terroriste tire d'un point colorié",
	"le terroriste tire à partir d'un véhicule",
	"l'agent secret tire",
	"l'agent secret tire d'un véhicule",
	"l'agent secret tire de l'hélicoptère",
	"la police tire",
	"la police tire d'un véhicule",
	"la police tire de l'hélicoptère",
}

const (
	r_ = Rate
	rB = Blesse
	rT = Tue
	rV = Bavure
)

// TableTirs[ligne][dés − 2]
var TableTirs = [...][11]Resultat{
	{r_, r_, rV, r_, rB, rT, rB, r_, rV, r_, r_},
	{r_, r_, rV, rB, rB, rT, rB, rB, rV, r_, r_},
	{r_, rT, r_, rB, rV, rB, rV, rB, r_, rT, r_},
	{rT, rT, rT, rB, r_, r_, r_, rB, rT, rT, rT},
	{rT, rT, rB, r_, r_, r_, r_, r_, rB, rT, rT},
	{rT, rB, r_, r_, r_, r_, r_, r_, r_, rB, rT},
	{rT, rT, rB, r_, rV, rV, rV, r_, rB, rT, rT},
	{rT, rT, rB, rV, rV, r_, rV, rV, rB, rT, rT},
	{rV, rB, rV, rV, rT, rV, rB, rV, rV, rB, rV},
}

// Chances : probabilité de chaque résultat sur une ligne.
func Chances(ligne int) [4]float64 {
	var c [4]float64
	for i, r := range TableTirs[ligne] {
		s := i + 2
		n := min(s-1, 13-s)
		c[r] += float64(n) / 36
	}
	return c
}
