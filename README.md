# Cortège officiel

Adaptation en Go / [Ebitengine](https://ebitengine.org) du jeu de
Jean-Charles et Sylvie Rodriguez paru en encart dans *Jeux & Stratégie*
(1983), plan dessiné par Philippe Fassier et Alain Meyer. Police,
manifestants, agents secrets et terroristes autour du cortège d'un président
étranger, de l'Hôtel de Ville à la Grand-Place puis à l'aéroport. Un joueur
humain contre l'ordinateur, à 2, 3 ou 4 joueurs, dans le camp de son choix.

## Compiler

Il faut [Go](https://go.dev/dl/) 1.27.1 ou plus récent. Le plan, les pions,
les polices et la règle sont intégrés au binaire : il se suffit à lui-même.

```
git clone https://github.com/alexdeba/cortege_officiel.git
cd cortege_officiel
go build -o cortegeofficiel .
./cortegeofficiel                        # F11 : plein écran
```

`go run .` lance le jeu sans garder de binaire.

Selon le système :

- **Linux** : Ebitengine a besoin d'un compilateur C et des en-têtes OpenGL
  et X11. Sous Debian ou Ubuntu :
  ```
  sudo apt install gcc pkg-config libgl1-mesa-dev libxcursor-dev libxi-dev \
      libxinerama-dev libxrandr-dev libxxf86vm-dev
  ```
- **macOS** : les outils en ligne de commande de Xcode
  (`xcode-select --install`).
- **Windows** : rien de plus. `go build -o cortegeofficiel.exe .` donne le
  binaire. On peut aussi le produire depuis Linux ou macOS :
  `GOOS=windows go build -o cortegeofficiel.exe .`

## Comment jouer

- Une consigne à chaque étape, dans le panneau de droite, avec ses boutons.
- Coups possibles surlignés : les pions jouables sont encadrés d'ambre,
  le pion choisi de rouge ; ses destinations (ou ses cibles) sont entourées
  sur la carte ; le prochain point de la voiture est cerclé de blanc.
  Survoler un point montre sa zone, sa nature et ses occupants.
- Barre d'espace : zoom sur la carte (là où est la souris) ; une mini-carte
  apparaît en haut à droite (clic ou glisser pour s'y déplacer).
- R : la règle complète, les deux tables, les choix d'adaptation et les
  commandes. S : la table des scores. A : accélérer l'ordinateur.

## Choix d'adaptation

La règle laisse quelques points dans l'ombre ; le jeu les tranche ainsi
(ces choix sont aussi dans la règle consultable en jeu).

- **Voisins** : deux points de rue sont voisins s'ils sont reliés par un
  trait rouge. Un point muni d'une flèche est voisin des points de toit de
  l'immeuble qu'elle désigne ; les toits d'un même immeuble sont voisins.
- **Grand-Place** : quand la voiture présidentielle entre dans les barrières
  par son parcours, elle est rangée au point 1 (le parcours rouge passe par
  la tribune, qu'aucun trait ne relie au point 1). Les chiffres des autres
  ronds blancs (1 à 6), que la règle n'explique pas, ne servent pas.
- La **tribune** est réservée au président ; personne ne se pose sur le
  point de la voiture présidentielle ou du président.
- Les **véhicules** ne comptent pas dans la limite de deux pions par point
  (un seul véhicule par point), leurs passagers non plus ; un véhicule ne
  traverse que des points libres. Un manifestant au volant d'un véhicule
  pris, rejoint par la police, l'abandonne pour reculer.
- **Pose** : les manifestants sont posés hors de la Grand-Place.
- **Toits** : un manifestant y monte s'il est arrivé sur le point fléché à
  un tour précédent ; un policier, sans attendre. D'un toit (ou de
  l'hélicoptère) on tire sur toute la zone ; « toujours de toit à toit » est
  lu comme : vers les toits des zones voisines aussi.
- **Hélicoptère** : il ne se pose que sur un toit à H ; il vole de 4 zones
  par tour, 3 s'il décolle ou se pose ; en l'air au-dessus d'une zone,
  personne ne peut le viser.
- **Bavures** : la bavure (BV) d'un terroriste touche un manifestant voisin
  de la cible, sans pénalité (la règle réserve les bavures à la police).
  Pour la police : un tir réussi sur un pion qui n'est pas terroriste ; un
  tir raté ou une bavure touche le manifestant qui partage le point de la
  cible. Un agent qui touche un innocent ne commet pas de bavure, mais le
  pion est retiré.
- Le président blessé une seconde fois est tué. Un terroriste **interpellé**
  (cerné par la police) compte comme abattu.
- **Tirs** : chacun tire au plus une fois par phase. Les tirs « à tout
  moment » ont lieu aux phases d (arrestation), e (la police, sur les
  terroristes démasqués) et g (les agents) ; la riposte suit aussitôt chaque
  tentative d'attentat. Après une tentative, la voiture avance de deux
  points au tour suivant (sauf pendant le discours).
- **Fuite** : un terroriste démasqué quitte la ville depuis un point de rue
  du bord du plan ; un terroriste resté anonyme se fond dans la foule et
  compte comme échappé. Quand le cortège arrive (ou si le président est
  tué), trois tours de chasse à l'homme suivent si des terroristes démasqués
  sont en fuite.
- La partie s'arrête au **tour 30** si le cortège n'est pas arrivé. « Le
  cortège n'a jamais été bloqué » ne compte que si le président atteint
  l'aéroport.
- **À 3 joueurs**, les policiers déplacés par les manifestants (provocation)
  ou les terroristes (faux ordres) ne vont ni sur un point colorié ni sur
  les points 1, 2, 3 de la Grand-Place : sinon le cortège serait bloqué sans
  fin. Le terroriste déplace quatre manifestants de n'importe quel numéro.
- **À 2 joueurs**, chaque joueur additionne les colonnes de ses deux rôles ;
  à la phase h, le terroriste déplace quatre pions impairs, manifestants ou
  policiers.
- **Ordinateur** : il ne lit que ce que son rôle sait (ses pions secrets, ce
  qui a été démasqué, et les pions que les terroristes déplacent, que tout
  le monde voit). La foule de l'ordinateur ne déplace que 12 manifestants
  par tour.

## Crédits et licence

*Cortège officiel* est un jeu de Jean-Charles et Sylvie Rodriguez, paru en
encart dans *Jeux & Stratégie* en 1983, avec un plan dessiné par Philippe
Fassier et Alain Meyer. Le jeu original, sa règle et ses illustrations
appartiennent à leurs auteurs. Cette adaptation est un hommage non
commercial.

Les scans du magazine ne sont pas distribués dans ce dépôt. Seules les
images dont le jeu a besoin le sont (`assets/gen`, `assets/logo.png`).

Le code est sous licence MIT (voir `LICENSE`). Les polices de
`assets/polices` gardent leurs propres licences : SIL OFL pour Anton et Lato,
Apache 2.0 pour Special Elite.
