// Code généré par outils/releve ; NE PAS MODIFIER.

package jeu

// Dimensions de assets/gen/plateau.jpg (cadre + marge de 48 px).
const (
	LargeurCarte = 2784
	HauteurCarte = 2362
)

// Quadrillage des zones.
var (
	LignesX = []float64{48, 428, 811, 1197, 1582, 1967, 2350, 2736}
	LignesY = []float64{48, 378, 699, 1027, 1348, 1672, 1994, 2314}
)

// Barrières de la Grand-Place.
var BarriereGP = []Point{{1440, 865}, {1645, 775}, {1745, 805}, {1860, 935}, {1835, 1150}, {1830, 1212}, {1720, 1180}, {1560, 1185}, {1520, 1225}, {1450, 1240}, {1430, 950}}

var pointsCarte = []PointCarte{
	{X: 739, Y: 110, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 1, Lig: 0, GP: false, Sortie: true},                      // 0
	{X: 1757, Y: 110, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 0, GP: false, Sortie: true},                     // 1
	{X: 2054, Y: 100, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 0, GP: false, Sortie: true},                     // 2
	{X: 2285, Y: 102, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 0, GP: false, Sortie: true},                     // 3
	{X: 1148, Y: 120, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 0, GP: false, Sortie: true},                     // 4
	{X: 1499, Y: 140, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 0, GP: false, Sortie: true},                     // 5
	{X: 1901, Y: 172, R: 20, Couleurs: 0, Toit: true, H: true, Bati: 55, Col: 4, Lig: 0, GP: false, Sortie: false},                     // 6
	{X: 220, Y: 218, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 0, Lig: 0, GP: false, Sortie: false},                     // 7
	{X: 1847, Y: 285, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 0, GP: false, Sortie: false},                    // 8
	{X: 2046, Y: 279, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 0, GP: false, Sortie: false},                    // 9
	{X: 2450, Y: 280, R: 16, Couleurs: Jaune | Rouge | Bleu, Toit: false, H: false, Bati: 0, Col: 6, Lig: 0, GP: false, Sortie: false}, // 10
	{X: 495, Y: 306, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 1, Lig: 0, GP: false, Sortie: false},                     // 11
	{X: 1041, Y: 300, R: 19, Couleurs: 0, Toit: true, H: true, Bati: 62, Col: 2, Lig: 0, GP: false, Sortie: false},                     // 12
	{X: 2116, Y: 310, R: 17, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 5, Lig: 0, GP: false, Sortie: false},                 // 13
	{X: 1748, Y: 374, R: 17, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 4, Lig: 0, GP: false, Sortie: false},                 // 14
	{X: 1968, Y: 377, R: 16, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 5, Lig: 0, GP: false, Sortie: false},                 // 15
	{X: 2350, Y: 377, R: 16, Couleurs: Jaune | Rouge | Bleu, Toit: false, H: false, Bati: 0, Col: 6, Lig: 0, GP: false, Sortie: false}, // 16
	{X: 2204, Y: 408, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 1, GP: false, Sortie: false},                    // 17
	{X: 736, Y: 430, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 1, Lig: 1, GP: false, Sortie: false},                     // 18
	{X: 1297, Y: 420, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 1, GP: false, Sortie: false},                    // 19
	{X: 1697, Y: 422, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 1, GP: false, Sortie: false},                    // 20
	{X: 2079, Y: 412, R: 20, Couleurs: 0, Toit: true, H: false, Bati: 119, Col: 5, Lig: 1, GP: false, Sortie: false},                   // 21
	{X: 2532, Y: 429, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 6, Lig: 1, GP: false, Sortie: false},                    // 22
	{X: 1582, Y: 460, R: 16, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 4, Lig: 1, GP: false, Sortie: false},                 // 23
	{X: 2668, Y: 476, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 6, Lig: 1, GP: false, Sortie: true},                     // 24
	{X: 1854, Y: 498, R: 19, Couleurs: 0, Toit: true, H: false, Bati: 95, Col: 4, Lig: 1, GP: false, Sortie: false},                    // 25
	{X: 2043, Y: 531, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 1, GP: false, Sortie: false},                    // 26
	{X: 2156, Y: 540, R: 16, Couleurs: Rouge, Toit: false, H: false, Bati: 0, Col: 5, Lig: 1, GP: false, Sortie: false},                // 27
	{X: 2286, Y: 533, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 1, GP: false, Sortie: false},                    // 28
	{X: 2351, Y: 533, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 6, Lig: 1, GP: false, Sortie: false},                // 29
	{X: 374, Y: 556, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 0, Lig: 1, GP: false, Sortie: false},                     // 30
	{X: 976, Y: 573, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 1, GP: false, Sortie: false},                     // 31
	{X: 1524, Y: 584, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 1, GP: false, Sortie: false},                    // 32
	{X: 1726, Y: 574, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 1, GP: false, Sortie: false},                    // 33
	{X: 1942, Y: 589, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 1, GP: false, Sortie: false},                    // 34
	{X: 2268, Y: 638, R: 20, Couleurs: 0, Toit: true, H: false, Bati: 200, Col: 5, Lig: 1, GP: false, Sortie: false},                   // 35
	{X: 2512, Y: 628, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 6, Lig: 1, GP: false, Sortie: false},                // 36
	{X: 2138, Y: 656, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 1, GP: false, Sortie: false},                    // 37
	{X: 1582, Y: 696, R: 17, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 4, Lig: 1, GP: false, Sortie: false},                 // 38
	{X: 1842, Y: 697, R: 16, Couleurs: Rouge, Toit: false, H: false, Bati: 0, Col: 4, Lig: 1, GP: false, Sortie: false},                // 39
	{X: 1968, Y: 698, R: 16, Couleurs: Rouge, Toit: false, H: false, Bati: 0, Col: 5, Lig: 1, GP: false, Sortie: false},                // 40
	{X: 1694, Y: 721, R: 20, Couleurs: 0, Toit: true, H: false, Bati: 210, Col: 4, Lig: 2, GP: false, Sortie: false},                   // 41
	{X: 2553, Y: 727, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 6, Lig: 2, GP: false, Sortie: false},                    // 42
	{X: 636, Y: 748, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 1, Lig: 2, GP: false, Sortie: false},                     // 43
	{X: 104, Y: 764, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 0, Lig: 2, GP: false, Sortie: true},                      // 44
	{X: 320, Y: 754, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 0, Lig: 2, GP: false, Sortie: false},                     // 45
	{X: 479, Y: 754, R: 20, Couleurs: 0, Toit: true, H: true, Bati: 214, Col: 1, Lig: 2, GP: false, Sortie: false},                     // 46
	{X: 885, Y: 754, R: 19, Couleurs: 0, Toit: true, H: true, Bati: 218, Col: 2, Lig: 2, GP: false, Sortie: false},                     // 47
	{X: 1330, Y: 762, R: 19, Couleurs: 0, Toit: true, H: true, Bati: 168, Col: 3, Lig: 2, GP: false, Sortie: false},                    // 48
	{X: 1905, Y: 753, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 2, GP: false, Sortie: false},                    // 49
	{X: 1146, Y: 774, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 2, GP: false, Sortie: false},                    // 50
	{X: 2312, Y: 764, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 2, GP: false, Sortie: false},                    // 51
	{X: 1486, Y: 799, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 2, GP: false, Sortie: false},                    // 52
	{X: 2108, Y: 795, R: 20, Couleurs: 0, Toit: true, H: true, Bati: 242, Col: 5, Lig: 2, GP: false, Sortie: false},                    // 53
	{X: 2462, Y: 788, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 6, Lig: 2, GP: false, Sortie: false},                // 54
	{X: 812, Y: 858, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 2, Lig: 2, GP: false, Sortie: false},                 // 55
	{X: 1584, Y: 858, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 2, GP: true, Sortie: false},                     // 56
	{X: 1698, Y: 859, R: 16, Couleurs: Jaune | Rouge | Bleu, Toit: false, H: false, Bati: 0, Col: 4, Lig: 2, GP: true, Sortie: false},  // 57
	{X: 2672, Y: 869, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 6, Lig: 2, GP: false, Sortie: true},                     // 58
	{X: 916, Y: 892, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 2, Lig: 2, GP: false, Sortie: false},                 // 59
	{X: 2542, Y: 904, R: 20, Couleurs: 0, Toit: true, H: false, Bati: 260, Col: 6, Lig: 2, GP: false, Sortie: false},                   // 60
	{X: 582, Y: 936, R: 17, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 1, Lig: 2, GP: false, Sortie: false},                 // 61
	{X: 1078, Y: 932, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 2, Lig: 2, GP: false, Sortie: false},                // 62
	{X: 1583, Y: 936, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 2, GP: true, Sortie: false},                     // 63
	{X: 2317, Y: 930, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 5, Lig: 2, GP: false, Sortie: false},                // 64
	{X: 1454, Y: 950, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 2, GP: true, Sortie: false},                     // 65
	{X: 2110, Y: 956, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 2, GP: false, Sortie: false},                    // 66
	{X: 1248, Y: 978, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 3, Lig: 2, GP: false, Sortie: false},                // 67
	{X: 1810, Y: 979, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 4, Lig: 2, GP: true, Sortie: false},                 // 68
	{X: 259, Y: 996, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 0, Lig: 2, GP: false, Sortie: false},                     // 69
	{X: 1156, Y: 1025, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 2, GP: false, Sortie: false},                   // 70
	{X: 1582, Y: 1025, R: 17, Couleurs: Jaune | Rouge | Bleu, Toit: false, H: false, Bati: 0, Col: 4, Lig: 2, GP: true, Sortie: false}, // 71
	{X: 536, Y: 1064, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 1, Lig: 3, GP: false, Sortie: false},                // 72
	{X: 1354, Y: 1056, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 3, Lig: 3, GP: false, Sortie: false},               // 73
	{X: 2038, Y: 1070, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 5, Lig: 3, GP: false, Sortie: false},               // 74
	{X: 2330, Y: 1082, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 5, Lig: 3, GP: false, Sortie: false},               // 75
	{X: 698, Y: 1101, R: 19, Couleurs: 0, Toit: true, H: true, Bati: 332, Col: 1, Lig: 3, GP: false, Sortie: false},                    // 76
	{X: 1490, Y: 1102, R: 16, Couleurs: Jaune | Rouge, Toit: false, H: false, Bati: 0, Col: 3, Lig: 3, GP: true, Sortie: false},        // 77
	{X: 1903, Y: 1098, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 3, GP: false, Sortie: false},                   // 78
	{X: 2630, Y: 1092, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 6, Lig: 3, GP: false, Sortie: true},                    // 79
	{X: 2254, Y: 1114, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 5, Lig: 3, GP: false, Sortie: false},               // 80
	{X: 424, Y: 1146, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 0, Lig: 3, GP: false, Sortie: false},                // 81
	{X: 884, Y: 1158, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 3, GP: false, Sortie: false},                    // 82
	{X: 1282, Y: 1170, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 3, GP: false, Sortie: false},                   // 83
	{X: 1076, Y: 1183, R: 20, Couleurs: 0, Toit: true, H: false, Bati: 383, Col: 2, Lig: 3, GP: false, Sortie: false},                  // 84
	{X: 1788, Y: 1172, R: 17, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 4, Lig: 3, GP: true, Sortie: false},                 // 85
	{X: 2162, Y: 1174, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 5, Lig: 3, GP: false, Sortie: false},               // 86
	{X: 1388, Y: 1188, R: 16, Couleurs: Rouge, Toit: false, H: false, Bati: 0, Col: 3, Lig: 3, GP: false, Sortie: false},               // 87
	{X: 1498, Y: 1192, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 3, GP: true, Sortie: false},                    // 88
	{X: 274, Y: 1258, R: 18, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 0, Lig: 3, GP: false, Sortie: false},                    // 89
	{X: 1660, Y: 1274, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 3, GP: false, Sortie: false},                   // 90
	{X: 1924, Y: 1262, R: 20, Couleurs: 0, Toit: true, H: true, Bati: 430, Col: 4, Lig: 3, GP: false, Sortie: false},                   // 91
	{X: 2052, Y: 1265, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 3, GP: false, Sortie: false},                   // 92
	{X: 738, Y: 1280, R: 18, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 1, Lig: 3, GP: false, Sortie: false},                    // 93
	{X: 1388, Y: 1287, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 3, GP: false, Sortie: false},                   // 94
	{X: 2281, Y: 1278, R: 19, Couleurs: 0, Toit: true, H: true, Bati: 427, Col: 5, Lig: 3, GP: false, Sortie: false},                   // 95
	{X: 2479, Y: 1281, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 6, Lig: 3, GP: false, Sortie: false},                   // 96
	{X: 100, Y: 1298, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 0, Lig: 3, GP: false, Sortie: true},                     // 97
	{X: 1522, Y: 1308, R: 19, Couleurs: 0, Toit: true, H: false, Bati: 471, Col: 3, Lig: 3, GP: false, Sortie: false},                  // 98
	{X: 537, Y: 1347, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 1, Lig: 3, GP: false, Sortie: false},                // 99
	{X: 1196, Y: 1348, R: 15, Couleurs: Rouge, Toit: false, H: false, Bati: 0, Col: 2, Lig: 3, GP: false, Sortie: false},               // 100
	{X: 875, Y: 1387, R: 20, Couleurs: 0, Toit: true, H: false, Bati: 462, Col: 2, Lig: 4, GP: false, Sortie: false},                   // 101
	{X: 1027, Y: 1382, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 4, GP: false, Sortie: false},                   // 102
	{X: 1846, Y: 1402, R: 16, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 4, Lig: 4, GP: false, Sortie: false},                // 103
	{X: 302, Y: 1434, R: 19, Couleurs: 0, Toit: true, H: true, Bati: 479, Col: 0, Lig: 4, GP: false, Sortie: false},                    // 104
	{X: 1926, Y: 1436, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 4, GP: false, Sortie: false},                   // 105
	{X: 1340, Y: 1473, R: 20, Couleurs: 0, Toit: true, H: false, Bati: 494, Col: 3, Lig: 4, GP: false, Sortie: false},                  // 106
	{X: 1535, Y: 1461, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 4, GP: false, Sortie: false},                   // 107
	{X: 500, Y: 1480, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 1, Lig: 4, GP: false, Sortie: false},                    // 108
	{X: 2230, Y: 1474, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 4, GP: false, Sortie: false},                   // 109
	{X: 674, Y: 1494, R: 19, Couleurs: 0, Toit: true, H: false, Bati: 493, Col: 1, Lig: 4, GP: false, Sortie: false},                   // 110
	{X: 427, Y: 1512, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 0, Lig: 4, GP: false, Sortie: false},                // 111
	{X: 870, Y: 1514, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 4, GP: false, Sortie: false},                    // 112
	{X: 1004, Y: 1510, R: 16, Couleurs: Rouge, Toit: false, H: false, Bati: 0, Col: 2, Lig: 4, GP: false, Sortie: false},               // 113
	{X: 1130, Y: 1506, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 4, GP: false, Sortie: false},                   // 114
	{X: 2600, Y: 1542, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 6, Lig: 4, GP: false, Sortie: true},                    // 115
	{X: 2131, Y: 1572, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 4, GP: false, Sortie: false},                   // 116
	{X: 100, Y: 1604, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 0, Lig: 4, GP: false, Sortie: true},                     // 117
	{X: 1970, Y: 1594, R: 17, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 5, Lig: 4, GP: false, Sortie: false},                // 118
	{X: 1522, Y: 1623, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 4, GP: false, Sortie: false},                   // 119
	{X: 978, Y: 1636, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 4, GP: false, Sortie: false},                    // 120
	{X: 1135, Y: 1640, R: 19, Couleurs: 0, Toit: true, H: false, Bati: 549, Col: 2, Lig: 4, GP: false, Sortie: false},                  // 121
	{X: 1344, Y: 1646, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 4, GP: false, Sortie: false},                   // 122
	{X: 332, Y: 1670, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 0, Lig: 4, GP: false, Sortie: false},                // 123
	{X: 811, Y: 1672, R: 15, Couleurs: Rouge, Toit: false, H: false, Bati: 0, Col: 2, Lig: 4, GP: false, Sortie: false},                // 124
	{X: 489, Y: 1701, R: 19, Couleurs: 0, Toit: true, H: false, Bati: 542, Col: 1, Lig: 5, GP: false, Sortie: false},                   // 125
	{X: 643, Y: 1704, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 1, Lig: 5, GP: false, Sortie: false},                    // 126
	{X: 1808, Y: 1710, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 5, GP: false, Sortie: false},                   // 127
	{X: 1968, Y: 1744, R: 17, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 5, Lig: 5, GP: false, Sortie: false},                // 128
	{X: 250, Y: 1766, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 0, Lig: 5, GP: false, Sortie: false},                    // 129
	{X: 946, Y: 1801, R: 19, Couleurs: 0, Toit: true, H: false, Bati: 577, Col: 2, Lig: 5, GP: false, Sortie: false},                   // 130
	{X: 1162, Y: 1794, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 5, GP: false, Sortie: false},                   // 131
	{X: 484, Y: 1842, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 1, Lig: 5, GP: false, Sortie: false},                    // 132
	{X: 618, Y: 1834, R: 15, Couleurs: Rouge, Toit: false, H: false, Bati: 0, Col: 1, Lig: 5, GP: false, Sortie: false},                // 133
	{X: 749, Y: 1828, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 1, Lig: 5, GP: false, Sortie: false},                    // 134
	{X: 1658, Y: 1830, R: 20, Couleurs: 0, Toit: true, H: true, Bati: 572, Col: 4, Lig: 5, GP: false, Sortie: false},                   // 135
	{X: 2286, Y: 1854, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 5, GP: false, Sortie: false},                   // 136
	{X: 2162, Y: 1888, R: 19, Couleurs: 0, Toit: true, H: true, Bati: 545, Col: 5, Lig: 5, GP: false, Sortie: false},                   // 137
	{X: 286, Y: 1902, R: 16, Couleurs: Jaune, Toit: false, H: false, Bati: 0, Col: 0, Lig: 5, GP: false, Sortie: false},                // 138
	{X: 106, Y: 1934, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 0, Lig: 5, GP: false, Sortie: true},                     // 139
	{X: 1452, Y: 1930, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 5, GP: false, Sortie: false},                   // 140
	{X: 1744, Y: 1924, R: 17, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 4, Lig: 5, GP: false, Sortie: false},                // 141
	{X: 1882, Y: 1920, R: 20, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 5, GP: false, Sortie: false},                   // 142
	{X: 586, Y: 1962, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 1, Lig: 5, GP: false, Sortie: false},                    // 143
	{X: 718, Y: 1956, R: 19, Couleurs: 0, Toit: true, H: false, Bati: 648, Col: 1, Lig: 5, GP: false, Sortie: false},                   // 144
	{X: 885, Y: 1960, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 5, GP: false, Sortie: false},                    // 145
	{X: 1580, Y: 1960, R: 17, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 3, Lig: 5, GP: false, Sortie: false},                // 146
	{X: 430, Y: 1994, R: 14, Couleurs: Rouge, Toit: false, H: false, Bati: 0, Col: 1, Lig: 5, GP: false, Sortie: false},                // 147
	{X: 1390, Y: 1998, R: 17, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 3, Lig: 6, GP: false, Sortie: false},                // 148
	{X: 1458, Y: 2036, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 6, GP: false, Sortie: false},                   // 149
	{X: 987, Y: 2064, R: 16, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 2, Lig: 6, GP: false, Sortie: false},                 // 150
	{X: 1169, Y: 2067, R: 16, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 2, Lig: 6, GP: false, Sortie: false},                // 151
	{X: 684, Y: 2080, R: 16, Couleurs: Bleu, Toit: false, H: false, Bati: 0, Col: 1, Lig: 6, GP: false, Sortie: false},                 // 152
	{X: 306, Y: 2096, R: 15, Couleurs: Jaune | Rouge | Bleu, Toit: false, H: false, Bati: 0, Col: 0, Lig: 6, GP: false, Sortie: false}, // 153
	{X: 2061, Y: 2088, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 5, Lig: 6, GP: false, Sortie: false},                   // 154
	{X: 2478, Y: 2110, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 6, Lig: 6, GP: false, Sortie: false},                   // 155
	{X: 88, Y: 2188, R: 18, Couleurs: 0, Toit: true, H: true, Bati: 711, Col: 0, Lig: 6, GP: false, Sortie: false},                     // 156
	{X: 1006, Y: 2226, R: 19, Couleurs: 0, Toit: true, H: true, Bati: 742, Col: 2, Lig: 6, GP: false, Sortie: false},                   // 157
	{X: 1670, Y: 2220, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 4, Lig: 6, GP: false, Sortie: true},                    // 158
	{X: 520, Y: 2259, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 1, Lig: 6, GP: false, Sortie: true},                     // 159
	{X: 876, Y: 2258, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 2, Lig: 6, GP: false, Sortie: true},                     // 160
	{X: 1324, Y: 2257, R: 19, Couleurs: 0, Toit: false, H: false, Bati: 0, Col: 3, Lig: 6, GP: false, Sortie: true},                    // 161
}

// Traits rouges entre points (voisins au sol).
var traitsCarte = [][2]int{
	{0, 4}, {0, 7}, {0, 11}, {0, 18}, {0, 31}, {1, 2}, {1, 5}, {1, 8}, {2, 3}, {2, 9},
	{3, 13}, {4, 5}, {4, 19}, {5, 14}, {5, 19}, {7, 11}, {7, 30}, {7, 44}, {8, 9}, {8, 14},
	{8, 26}, {9, 13}, {9, 15}, {10, 16}, {11, 18}, {11, 30}, {13, 15}, {13, 16}, {14, 15}, {14, 20},
	{14, 23}, {15, 16}, {15, 20}, {15, 26}, {15, 40}, {16, 17}, {16, 22}, {16, 27}, {16, 28}, {16, 29},
	{17, 26}, {17, 27}, {17, 28}, {17, 29}, {18, 31}, {18, 43}, {19, 23}, {19, 31}, {19, 32}, {19, 50},
	{19, 52}, {20, 23}, {20, 33}, {22, 24}, {22, 36}, {23, 32}, {23, 33}, {23, 38}, {24, 36}, {24, 58},
	{26, 27}, {26, 34}, {27, 28}, {27, 34}, {27, 37}, {27, 40}, {28, 29}, {28, 37}, {29, 36}, {29, 51},
	{30, 43}, {30, 45}, {31, 43}, {31, 50}, {32, 38}, {32, 52}, {33, 34}, {33, 38}, {33, 39}, {34, 37},
	{34, 39}, {34, 40}, {34, 49}, {36, 42}, {36, 54}, {37, 40}, {37, 51}, {38, 52}, {38, 57}, {39, 49},
	{39, 57}, {39, 68}, {40, 49}, {42, 54}, {43, 55}, {43, 61}, {44, 45}, {44, 69}, {44, 97}, {45, 61},
	{45, 69}, {49, 57}, {49, 66}, {50, 62}, {50, 65}, {51, 54}, {51, 64}, {51, 66}, {52, 56}, {52, 57},
	{52, 65}, {52, 71}, {54, 64}, {55, 59}, {55, 61}, {56, 63}, {56, 68}, {57, 65}, {57, 68}, {57, 71},
	{57, 85}, {58, 79}, {59, 62}, {59, 82}, {61, 72}, {62, 67}, {64, 66}, {64, 75}, {65, 71}, {65, 73},
	{65, 77}, {65, 87}, {66, 74}, {66, 80}, {67, 70}, {67, 73}, {68, 71}, {68, 74}, {68, 77}, {68, 78},
	{68, 85}, {68, 88}, {69, 81}, {69, 89}, {69, 97}, {70, 71}, {70, 82}, {70, 83}, {71, 77}, {71, 85},
	{72, 81}, {72, 93}, {73, 77}, {73, 83}, {73, 87}, {73, 88}, {74, 78}, {74, 86}, {75, 79}, {75, 80},
	{75, 96}, {77, 85}, {77, 87}, {77, 88}, {78, 85}, {78, 92}, {79, 96}, {79, 115}, {80, 86}, {82, 93},
	{82, 102}, {83, 87}, {83, 94}, {83, 100}, {83, 102}, {85, 90}, {85, 103}, {86, 92}, {87, 88}, {87, 94},
	{87, 100}, {88, 90}, {88, 94}, {89, 97}, {89, 99}, {90, 103}, {90, 107}, {92, 105}, {92, 109}, {93, 108},
	{93, 112}, {94, 100}, {94, 107}, {94, 114}, {96, 109}, {96, 115}, {97, 117}, {99, 108}, {99, 111}, {100, 102},
	{100, 113}, {100, 114}, {102, 112}, {102, 113}, {102, 114}, {103, 105}, {103, 107}, {103, 118}, {105, 116}, {105, 118},
	{105, 127}, {107, 119}, {107, 122}, {108, 111}, {108, 126}, {109, 116}, {109, 136}, {111, 123}, {111, 147}, {112, 113},
	{112, 120}, {112, 124}, {112, 126}, {113, 114}, {113, 120}, {113, 124}, {114, 120}, {114, 122}, {115, 136}, {115, 155},
	{116, 118}, {116, 128}, {117, 123}, {117, 129}, {117, 139}, {118, 128}, {119, 122}, {119, 127}, {119, 140}, {120, 124},
	{120, 131}, {120, 134}, {122, 131}, {122, 140}, {123, 129}, {123, 132}, {123, 138}, {124, 126}, {124, 133}, {124, 134},
	{126, 132}, {126, 133}, {126, 134}, {127, 128}, {127, 141}, {128, 141}, {128, 142}, {129, 138}, {129, 139}, {129, 147},
	{131, 145}, {131, 148}, {132, 133}, {132, 138}, {132, 143}, {132, 147}, {133, 134}, {133, 143}, {133, 147}, {134, 143},
	{134, 145}, {136, 154}, {136, 155}, {138, 139}, {138, 147}, {138, 153}, {140, 146}, {140, 148}, {140, 149}, {140, 151},
	{141, 142}, {141, 146}, {141, 149}, {142, 149}, {142, 154}, {143, 147}, {143, 152}, {145, 150}, {145, 152}, {146, 148},
	{146, 149}, {146, 158}, {147, 153}, {148, 149}, {148, 151}, {149, 151}, {149, 161}, {150, 151}, {150, 152}, {150, 160},
	{151, 161}, {152, 153}, {152, 159}, {152, 160}, {154, 155}, {154, 158}, {158, 161}, {159, 160}, {160, 161}, {38, 56},
	{63, 71}, {81, 99}, {39, 40}}

// Flèches : point de rue → point de toit de l'immeuble désigné.
var accesCarte = [][2]int{
	{0, 12}, {1, 6}, {2, 6}, {4, 12}, {8, 6}, {9, 6}, {17, 21}, {19, 48}, {19, 12}, {20, 25},
	{26, 25}, {26, 21}, {28, 35}, {30, 46}, {31, 47}, {31, 12}, {33, 41}, {33, 25}, {34, 25}, {42, 60},
	{43, 46}, {43, 47}, {45, 46}, {49, 53}, {50, 47}, {50, 48}, {51, 53}, {51, 35}, {52, 48}, {65, 48},
	{66, 53}, {70, 84}, {78, 91}, {79, 60}, {82, 101}, {82, 76}, {83, 84}, {88, 98}, {92, 91}, {92, 95},
	{93, 110}, {93, 76}, {93, 101}, {94, 98}, {94, 106}, {96, 95}, {102, 101}, {102, 84}, {105, 91}, {107, 106},
	{107, 98}, {108, 110}, {108, 125}, {109, 137}, {109, 95}, {112, 110}, {112, 101}, {114, 121}, {114, 106}, {116, 137},
	{117, 104}, {119, 135}, {120, 121}, {120, 130}, {122, 121}, {122, 106}, {126, 125}, {126, 110}, {127, 135}, {131, 130},
	{131, 121}, {132, 125}, {134, 144}, {134, 130}, {136, 137}, {140, 135}, {142, 137}, {143, 144}, {145, 144}, {145, 130},
	{154, 137}, {160, 157}, {161, 157}}

// Points remarquables.
const (
	PtHotelDeVille   = 10
	PtGP1            = 56
	PtGP2            = 63
	PtGP3            = 71
	PtTribune        = 57
	PtAeroport       = 153
	PtHelicoAeroport = 156
)
