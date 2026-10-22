package main

// choixAdaptation : les points ambigus de la règle et la façon dont le jeu
// les tranche (repris dans le README).
const choixAdaptation = `• Voisins : deux points de rue sont voisins s'ils sont reliés par un trait rouge. Un point muni d'une flèche est voisin des points de toit de l'immeuble qu'elle désigne ; les toits d'un même immeuble sont voisins entre eux.
• La Grand-Place : dès que la voiture présidentielle entre dans les barrières par son parcours, elle est rangée au point 1 (le parcours rouge passe par la tribune, qu'aucun trait ne relie au point 1).
• La tribune est réservée au président ; personne ne se pose sur le point de la voiture présidentielle ou du président.
• Les véhicules ne comptent pas dans la limite de deux pions par point (un seul véhicule par point) ; leurs passagers non plus. Un véhicule ne traverse que des points libres.
• Les manifestants sont posés hors de la Grand-Place (ils ne peuvent pas franchir ses barrières).
• Un manifestant monte sur un toit s'il est arrivé sur le point fléché à un tour précédent ; un policier, sans attendre.
• D'un toit (ou de l'hélicoptère), on tire sur toute la zone ; « toujours de toit à toit » est lu comme : vers les toits des zones voisines aussi.
• La bavure d'un terroriste (BV) touche un manifestant voisin de la cible, sans pénalité : la règle réserve les bavures à la police.
• Bavures de la police : un tir réussi sur un pion qui n'est pas terroriste ; un tir raté ou une bavure touche le manifestant qui partage le point de la cible, s'il y en a un. Un agent secret qui touche un innocent ne commet pas de bavure, mais le pion est retiré.
• Le président blessé une seconde fois est tué. Un terroriste interpellé (cerné par la police) compte comme abattu.
• Chacun tire au plus une fois par phase. Les tirs « à tout moment » ont lieu aux phases d (arrestation), e (la police, sur les terroristes démasqués) et g (les agents) ; la riposte suit aussitôt chaque tentative d'attentat.
• Après une tentative, la voiture avance de deux points au tour suivant (sauf pendant le discours).
• Fuite : un terroriste démasqué quitte la ville depuis un point de rue du bord du plan ; un terroriste resté anonyme se fond dans la foule et compte comme échappé. Quand le cortège arrive (ou si le président est tué), trois tours de chasse à l'homme suivent si des terroristes démasqués sont en fuite.
• La partie s'arrête au tour 30 si le cortège n'est pas arrivé. « Le cortège n'a jamais été bloqué » ne compte que si le président atteint l'aéroport.
• À 3 joueurs, les policiers déplacés par les manifestants (provocation) ou les terroristes (faux ordres) ne vont pas sur un point colorié ni sur les points 1, 2, 3 de la Grand-Place : sinon le cortège serait bloqué sans fin. Le terroriste déplace quatre manifestants de n'importe quel numéro.
• À 2 joueurs, chaque joueur additionne les colonnes de ses deux rôles ; à la phase h, quatre pions impairs (manifestants ou policiers).
• L'hélicoptère ne se pose que sur un toit à H. Il vole de 4 zones par tour, 3 s'il décolle ou se pose ; il peut rester en l'air au-dessus d'une zone, où personne ne peut le viser.
• Un manifestant au volant d'un véhicule pris, rejoint par la police, l'abandonne sur place pour reculer.
• Les chiffres des ronds blancs (1 à 6), que la règle n'explique pas, ne servent pas ; seuls comptent les points 1, 2 et 3 de la Grand-Place.
• L'ordinateur ne lit que ce que son rôle sait : ses propres pions secrets, ce qui a été démasqué, et les pions que les terroristes déplacent (tout le monde les voit).`
