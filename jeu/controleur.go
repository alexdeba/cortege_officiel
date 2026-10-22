package jeu

// Phase du tour de jeu (règle, « le tour de jeu », a à h).
type Phase int

const (
	PhPreparation Phase = iota
	PhVoiture           // a. la voiture présidentielle avance
	PhAttentat          // b. tentative d'assassinat
	PhRiposte           // c. riposte
	PhArrestation       // d. tentative d'arrestation
	PhPolice            // e. déplacement des policiers
	PhManifs            // f. déplacement des manifestants
	PhAgents            // g. les agents secrets déplacent quatre pions pairs
	PhTerros            // h. les terroristes déplacent quatre pions impairs
	PhContact           // un manifestant doit reculer
)

var NomsPhases = [...]string{"Préparation", "Voiture présidentielle", "Tentative d'assassinat", "Riposte",
	"Tentative d'arrestation", "Déplacement de la police", "Déplacement des manifestants",
	"Les agents secrets", "Les terroristes", "Contact"}

// Lettre de la phase dans la règle.
var LettresPhases = [...]string{"", "a", "b", "c", "d", "e", "f", "g", "h", ""}

type TypeMvt int

const (
	MvtPas      TypeMvt = iota // un pas vers un point voisin
	MvtConduite                // un véhicule et ses passagers
	MvtVol                     // l'hélicoptère vers une zone ou un toit à H
	MvtDescente                // un passager descend
	MvtMonter                  // un policier monte dans le véhicule de son point
	MvtFuite                   // un terroriste quitte la ville
)

// Mouvement : un déplacement possible. Vers : point d'arrivée (-1 : en vol,
// voir Zone).
type Mouvement struct {
	Type   TypeMvt
	Pion   int
	Vers   int
	Zone   int
	Chemin []int // points traversés, départ et arrivée compris
}

// Tir : un tireur et sa cible (pion ; la voiture ou le président).
type Tir struct {
	Tireur int
	Cible  int
	Table  int // ligne de la table des tirs (voir tables.go)
}

// Quota : ce qu'un joueur peut déplacer pendant une phase.
type Quota struct {
	Role      Role // le rôle qui joue
	Policiers bool
	Manifs    bool
	Vehicules bool
	Parite    int // 0 : tous ; 1 : impairs ; 2 : pairs
	// nombre de pions différents déplacés (0 : chacun une fois, sans limite)
	MaxManif, MaxPol, MaxTotal int
	Fuite                      bool
}

// Placement : la pose initiale d'un pion.
type Placement struct {
	Pion  int
	Point int
	Dans  int // véhicule posé sur ce point qui l'embarque (-1)
}

// Controleur : les décisions d'un joueur, demandées par le moteur au fil du
// tour. Le moteur vérifie chaque réponse.
type Controleur interface {
	// Secrets : n pions parmi les candidats (agents ou terroristes).
	Secrets(e *Etat, j *Joueur, r Role, candidats []int, n int) []int
	// Couleur : le parcours de la voiture présidentielle (aller ou retour).
	Couleur(e *Etat, j *Joueur, retour bool) Couleur
	// Placer : la pose des pions du rôle (police ou manifestants) ;
	// proposition : une pose légale déjà calculée.
	Placer(e *Etat, j *Joueur, r Role, proposition []Placement) []Placement
	// Avancer : le prochain point de la voiture quand il y a le choix.
	Avancer(e *Etat, j *Joueur, options []int) int
	// Deplacer : un mouvement parmi les options, -1 pour finir la phase.
	Deplacer(e *Etat, j *Joueur, q Quota, options []Mouvement) int
	// Tirer : un tir parmi les options, -1 pour ne pas (plus) tirer.
	Tirer(e *Etat, j *Joueur, ph Phase, options []Tir) int
	// Reculer : où envoyer un manifestant au contact de la police
	// (options : mouvements de recul possibles, un seul doit être choisi).
	Reculer(e *Etat, j *Joueur, options []Mouvement) int
}

// Observateur : suit la partie (affichage, animations, journal).
type Observateur interface {
	Evenement(e *Etat, ev Ev)
}

type TypeEv int

const (
	EvTexte TypeEv = iota
	EvTour
	EvPhase
	EvMouvement
	EvTir
	EvScore
	EvRevele
	EvBlocage
	EvFin
)

type Ev struct {
	Type   TypeEv
	Tour   int
	Phase  Phase
	Texte  string
	Pion   int
	Chemin []int
	Tir    Tir
	Des    [2]int
	Result Resultat
	Pts    [NbRoles]int
}
