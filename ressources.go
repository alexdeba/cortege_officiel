package main

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"cortegeofficiel/jeu"
)

//go:embed assets/logo.png assets/gen/plateau.jpg assets/gen/couverture.jpg assets/gen/pions.png assets/polices regle/regle.txt
var fichiers embed.FS

var (
	imgPlateau *ebiten.Image // le plan redressé
	imgCouv    *ebiten.Image // l'illustration de l'écran d'accueil
	imgPions   *ebiten.Image // la planche de pions découpée sur l'encart
	imgLogo    *ebiten.Image // le titre de la page 64, en blanc : teinté au dessin
	blanc      *ebiten.Image // 3 × 3 blanc, pour les triangles

	policeTitre, policeTexte, policeGras, policeItalique, policeMachine *Police
)

// Palette du film d'espionnage des années 70, prise sur la couverture :
// graphite et acier pour les cadres, papier bulle des dossiers pour ce qui
// se lit, un seul rouge signal pour l'action, l'ambre des projecteurs pour
// ce qui se joue, le kaki des services secrets.
var (
	coulFond     = color.RGBA{0x12, 0x14, 0x17, 0xff} // derrière la carte
	coulAcier    = color.RGBA{0x1d, 0x20, 0x24, 0xff} // le panneau
	coulAcier2   = color.RGBA{0x2c, 0x31, 0x37, 0xff} // cartes et boutons
	coulAcier3   = color.RGBA{0x4a, 0x51, 0x59, 0xff} // filets, survols
	coulCendre   = color.RGBA{0x95, 0x9b, 0xa1, 0xff} // texte secondaire, sur sombre
	coulCraie    = color.RGBA{0xec, 0xe7, 0xdc, 0xff} // texte, sur sombre
	coulDossier  = color.RGBA{0xe4, 0xda, 0xc2, 0xff} // papier des dossiers
	coulDossier2 = color.RGBA{0xd4, 0xc7, 0xa8, 0xff}
	coulDossier3 = color.RGBA{0xa8, 0x9a, 0x7a, 0xff}
	coulEncre    = color.RGBA{0x1a, 0x19, 0x17, 0xff} // encre, sur papier
	coulEncre2   = color.RGBA{0x55, 0x4f, 0x45, 0xff}
	coulSignal   = color.RGBA{0xc8, 0x2f, 0x24, 0xff} // le rouge du titre
	coulSignal2  = color.RGBA{0xe0, 0x4b, 0x3c, 0xff}
	coulAmbre    = color.RGBA{0xe0, 0xa1, 0x30, 0xff} // ce qui peut jouer
	coulAmbre2   = color.RGBA{0xf2, 0xc8, 0x6a, 0xff}
	coulKaki     = color.RGBA{0x6c, 0x78, 0x45, 0xff}
	coulRouge    = color.RGBA{0xa8, 0x1e, 0x1a, 0xff} // le sang, les blessures
	coulBleu     = color.RGBA{0x2a, 0x4b, 0x82, 0xff} // la police
	coulBlanc    = color.RGBA{0xff, 0xff, 0xff, 0xff}
	coulNoir     = color.RGBA{0x00, 0x00, 0x00, 0xff}
)

// couleurs des rôles
var coulRoles = [jeu.NbRoles]color.RGBA{
	jeu.RPolice: coulBleu,
	jeu.RManif:  {0xb0, 0x3a, 0x2e, 0xff},
	jeu.RAgent:  coulKaki,
	jeu.RTerro:  {0x0c, 0x0c, 0x0d, 0xff},
}

// couleurs des parcours, éteintes comme celles du plan étalonné
func coulParcours(c jeu.Couleur) color.RGBA {
	switch c {
	case jeu.Jaune:
		return color.RGBA{0xd6, 0xb0, 0x3c, 0xff}
	case jeu.Rouge:
		return color.RGBA{0xb4, 0x34, 0x28, 0xff}
	}
	return color.RGBA{0x3c, 0x86, 0xb0, 0xff}
}

func lireImage(nom string) image.Image {
	b, err := fichiers.ReadFile(nom)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(nom, err)
	}
	return img
}

func chargerImage(nom string) *ebiten.Image { return ebiten.NewImageFromImage(lireImage(nom)) }

// blanchir : l'encre du logo en blanc, l'alpha gardé, pour pouvoir la
// teinter de n'importe quelle couleur.
func blanchir(src image.Image) *image.NRGBA {
	b := src.Bounds()
	dst := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			_, _, _, a := src.At(b.Min.X+x, b.Min.Y+y).RGBA()
			o := dst.PixOffset(x, y)
			dst.Pix[o], dst.Pix[o+1], dst.Pix[o+2], dst.Pix[o+3] = 0xff, 0xff, 0xff, uint8(a>>8)
		}
	}
	return dst
}

// logo : le titre du jeu, de hauteur h, coin haut gauche en (x, y).
func logo(dst *ebiten.Image, x, y, h float64, c color.Color) {
	s := h / float64(imgLogo.Bounds().Dy())
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(s, s)
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	dst.DrawImage(imgLogo, op)
}

// etalonner : le plan à l'ambiance du film, calculé une fois au
// chargement. Couleurs à moitié éteintes, ombres bleutées, un voile
// sombre ; les points coloriés restent assez saturés pour se distinguer.
func etalonner(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	var lut [256]float64 // un peu de contraste, en courbe douce
	for i := range lut {
		v := float64(i) / 255
		lut[i] = v + 0.35*v*(1-v)*(2*v-1)
	}
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r16, g16, b16, _ := src.At(b.Min.X+x, b.Min.Y+y).RGBA()
			r, g, bl := lut[r16>>8], lut[g16>>8], lut[b16>>8]
			l := 0.299*r + 0.587*g + 0.114*bl
			const sat = 0.5
			r, g, bl = l+(r-l)*sat, l+(g-l)*sat, l+(bl-l)*sat
			// ombres froides, lumières à peine chaudes
			r, g, bl = r*0.86+0.02*l, g*0.88+0.03*(1-l)*l, bl*0.9+0.08*(1-l)
			o := dst.PixOffset(x, y)
			dst.Pix[o] = uint8(min(max(r, 0), 1) * 255)
			dst.Pix[o+1] = uint8(min(max(g, 0), 1) * 255)
			dst.Pix[o+2] = uint8(min(max(bl, 0), 1) * 255)
			dst.Pix[o+3] = 0xff
		}
	}
	return dst
}

type Police struct{ src *text.GoTextFaceSource }

func chargerPolice(nom string) *Police {
	b, err := fichiers.ReadFile("assets/polices/" + nom)
	if err != nil {
		log.Fatal(err)
	}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	return &Police{s}
}

func chargerRessources() {
	imgPlateau = ebiten.NewImageFromImage(etalonner(lireImage("assets/gen/plateau.jpg")))
	imgCouv = chargerImage("assets/gen/couverture.jpg")
	imgPions = chargerImage("assets/gen/pions.png")
	imgLogo = ebiten.NewImageFromImage(blanchir(lireImage("assets/logo.png")))
	blanc = ebiten.NewImage(3, 3)
	blanc.Fill(color.White)
	policeTitre = chargerPolice("Anton-Regular.ttf")
	policeTexte = chargerPolice("Lato-Regular.ttf")
	policeGras = chargerPolice("Lato-Bold.ttf")
	policeItalique = chargerPolice("Lato-Italic.ttf")
	policeMachine = chargerPolice("SpecialElite-Regular.ttf")
}

// ---------------------------------------------------------------- texte

type cleFace struct {
	p *Police
	t float64
}

var faces = map[cleFace]*text.GoTextFace{}

func face(taille float64, p *Police) *text.GoTextFace {
	k := cleFace{p, taille}
	if f, ok := faces[k]; ok {
		return f
	}
	f := &text.GoTextFace{Source: p.src, Size: taille}
	faces[k] = f
	return f
}

func ecrire(dst *ebiten.Image, s string, x, y float64, f text.Face, c color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, s, f, op)
}

func ecrireCentre(dst *ebiten.Image, s string, cx, y float64, f text.Face, c color.Color) {
	ecrire(dst, s, cx-text.Advance(s, f)/2, y, f, c)
}

func ecrireDroite(dst *ebiten.Image, s string, x, y float64, f text.Face, c color.Color) {
	ecrire(dst, s, x-text.Advance(s, f), y, f, c)
}

// ecrireHalo : texte cerné.
func ecrireHalo(dst *ebiten.Image, s string, x, y float64, f text.Face, c, halo color.Color, ep float64) {
	for a := 0.0; a < 2*math.Pi; a += math.Pi / 6 {
		ecrire(dst, s, x+math.Cos(a)*ep, y+math.Sin(a)*ep, f, halo)
	}
	ecrire(dst, s, x, y, f, c)
}

func largeur(s string, f text.Face) float64 { return text.Advance(s, f) }

// mots : la ponctuation haute reste avec le mot qui précède.
func mots(s string) []string {
	var l []string
	for _, m := range strings.Fields(s) {
		if n := len(l); n > 0 && (m == ";" || m == ":" || m == "!" || m == "?" || m == "»" || m == "$" || m == "%" || l[n-1] == "«") {
			l[n-1] += " " + m
			continue
		}
		l = append(l, m)
	}
	return l
}

// couper : un texte en lignes de largeur maximale w.
func couper(s string, f text.Face, w float64) []string {
	var lignes []string
	for _, para := range strings.Split(s, "\n") {
		cur := ""
		for _, mot := range mots(para) {
			essai := mot
			if cur != "" {
				essai = cur + " " + mot
			}
			if text.Advance(essai, f) > w && cur != "" {
				lignes = append(lignes, cur)
				cur = mot
			} else {
				cur = essai
			}
		}
		lignes = append(lignes, cur)
	}
	return lignes
}

// paragraphe : écrit un texte coupé, renvoie la hauteur utilisée.
func paragraphe(dst *ebiten.Image, s string, x, y, w float64, f text.Face, c color.Color, inter float64) float64 {
	l := couper(s, f, w)
	for i, t := range l {
		ecrire(dst, t, x, y+float64(i)*inter, f, c)
	}
	return float64(len(l)) * inter
}

// ---------------------------------------------------------------- formes

func rgba(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func rect(dst *ebiten.Image, x, y, w, h float64, c color.Color) {
	k := rgba(c)
	vs := []ebiten.Vertex{sommet(x, y, k), sommet(x+w, y, k), sommet(x, y+h, k), sommet(x+w, y+h, k)}
	dst.DrawTriangles(vs, []uint16{0, 1, 2, 1, 3, 2}, blanc, &ebiten.DrawTrianglesOptions{})
}

func cadre(dst *ebiten.Image, x, y, w, h, ep float64, c color.Color) {
	k := rgba(c)
	rect(dst, x-ep/2, y-ep/2, w+ep, ep, k)
	rect(dst, x-ep/2, y+h-ep/2, w+ep, ep, k)
	rect(dst, x-ep/2, y+ep/2, ep, h-ep, k)
	rect(dst, x+w-ep/2, y+ep/2, ep, h-ep, k)
}

// disque : plein, bord adouci d'un pixel.
func disque(dst *ebiten.Image, x, y, r float64, c color.Color) {
	k := rgba(c)
	pastille(dst, x, y, max(r-0.5, 0.5), k)
	anneau(dst, x, y, r-0.25, 0.5, k)
}

func cercle(dst *ebiten.Image, x, y, r, ep float64, c color.Color) {
	anneau(dst, x, y, r, ep, rgba(c))
}

func trait(dst *ebiten.Image, x0, y0, x1, y1, ep float64, c color.Color) {
	polyligne(dst, [][2]float64{{x0, y0}, {x1, y1}}, ep, rgba(c))
}

// arrondi : rectangle aux coins arrondis.
func cheminArrondi(x, y, w, h, r float64) *vector.Path {
	var p vector.Path
	r = min(r, w/2, h/2)
	p.MoveTo(float32(x+r), float32(y))
	p.LineTo(float32(x+w-r), float32(y))
	p.ArcTo(float32(x+w), float32(y), float32(x+w), float32(y+r), float32(r))
	p.LineTo(float32(x+w), float32(y+h-r))
	p.ArcTo(float32(x+w), float32(y+h), float32(x+w-r), float32(y+h), float32(r))
	p.LineTo(float32(x+r), float32(y+h))
	p.ArcTo(float32(x), float32(y+h), float32(x), float32(y+h-r), float32(r))
	p.LineTo(float32(x), float32(y+r))
	p.ArcTo(float32(x), float32(y), float32(x+r), float32(y), float32(r))
	p.Close()
	return &p
}

func remplir(dst *ebiten.Image, p *vector.Path, c color.Color) {
	op := &vector.FillOptions{}
	dop := &vector.DrawPathOptions{AntiAlias: true}
	dop.ColorScale.ScaleWithColor(c)
	vector.FillPath(dst, p, op, dop)
}

func tracer(dst *ebiten.Image, p *vector.Path, ep float64, c color.Color) {
	op := &vector.StrokeOptions{Width: float32(ep), LineJoin: vector.LineJoinRound, LineCap: vector.LineCapRound}
	dop := &vector.DrawPathOptions{AntiAlias: true}
	dop.ColorScale.ScaleWithColor(c)
	vector.StrokePath(dst, p, op, dop)
}

// arrondi : rectangle aux coins arrondis, rendu une fois (anticrénelé)
// dans une image gardée en cache, puis simplement recopié.
func arrondi(dst *ebiten.Image, x, y, w, h, r float64, fond color.Color, bord color.Color, ep float64) {
	k := cleArrondi{w: int(math.Round(w)), h: int(math.Round(h)), r: int(math.Round(r)), ep: int(math.Round(ep * 2))}
	if fond != nil {
		k.fond, k.aFond = rgba(fond), true
	}
	if bord != nil && ep > 0 {
		k.bord, k.aBord = rgba(bord), true
	}
	img, ok := arrondis[k]
	if !ok {
		if len(arrondis) > 2000 { // tailles très variées : on repart à zéro
			for _, i := range arrondis {
				i.Deallocate()
			}
			arrondis = map[cleArrondi]*ebiten.Image{}
		}
		m := int(math.Ceil(ep)) + 2
		img = ebiten.NewImage(k.w+2*m, k.h+2*m)
		p := cheminArrondi(float64(m), float64(m), float64(k.w), float64(k.h), r)
		if fond != nil {
			remplir(img, p, fond)
		}
		if bord != nil && ep > 0 {
			tracer(img, p, ep, bord)
		}
		arrondis[k] = img
	}
	m := float64(int(math.Ceil(ep)) + 2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(math.Round(x)-m, math.Round(y)-m)
	dst.DrawImage(img, op)
}

type cleArrondi struct {
	w, h, r, ep  int
	fond, bord   color.RGBA
	aFond, aBord bool
}

var arrondis = map[cleArrondi]*ebiten.Image{}

func alpha(c color.RGBA, a float64) color.RGBA {
	return color.RGBA{uint8(float64(c.R) * a), uint8(float64(c.G) * a), uint8(float64(c.B) * a), uint8(float64(c.A) * a)}
}

func melange(a, b color.RGBA, t float64) color.RGBA {
	m := func(x, y uint8) uint8 { return uint8(float64(x)*(1-t) + float64(y)*t) }
	return color.RGBA{m(a.R, b.R), m(a.G, b.G), m(a.B, b.B), m(a.A, b.A)}
}

// ---------------------------------------------------------------- décor

// filet : le double trait des chemises de dossier, un gras rouge et un
// maigre, de hauteur h.
func filet(dst *ebiten.Image, x, y, w, h float64, maigre color.Color) {
	rect(dst, x, y, w, h*0.6, coulSignal)
	rect(dst, x, y+h*0.8, w, max(h*0.2, 1), maigre)
}

// feuilleDossier : une feuille de papier bulle sur la table, ombre portée,
// filet en tête et trous de classeur.
func feuilleDossier(dst *ebiten.Image, x, y, w, h float64) {
	rect(dst, x+8, y+8, w, h, alpha(coulNoir, 0.6))
	rect(dst, x, y, w, h, coulDossier)
	cadre(dst, x, y, w, h, 1, coulDossier3)
	filet(dst, x+20, y+18, w-40, 8, coulEncre)
	for _, f := range []float64{0.25, 0.75} {
		disque(dst, x+10, y+h*f, 4, coulDossier3)
	}
}

// tampon : un coup de tampon rouge, encadré, un peu de travers.
func tampon(dst *ebiten.Image, s string, cx, cy, taille, angle float64) {
	f := face(taille, policeTitre)
	w, h := largeur(s, f)+taille*0.9, taille*1.75
	k := cleFace{nil, taille}
	img, ok := tampons[s][k]
	if !ok {
		img = ebiten.NewImage(int(w)+8, int(h)+8)
		cadre(img, 4+1.5, 4+1.5, w-3, h-3, 3, coulSignal)
		cadre(img, 4+6, 4+6, w-12, h-12, 1.2, coulSignal)
		ecrireCentre(img, s, 4+w/2, 4+(h-taille*1.5)/2, f, coulSignal)
		if tampons[s] == nil {
			tampons[s] = map[cleFace]*ebiten.Image{}
		}
		tampons[s][k] = img
	}
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Translate(-(w+8)/2, -(h+8)/2)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(cx, cy)
	op.ColorScale.ScaleAlpha(0.85)
	dst.DrawImage(img, op)
}

var tampons = map[string]map[cleFace]*ebiten.Image{}

// de : un dé à jouer, valeur v, de côté s.
func de(dst *ebiten.Image, x, y, s float64, v int, fond, points color.Color) {
	arrondi(dst, x, y, s, s, s*0.18, fond, coulEncre, 1.6)
	r := s * 0.09
	pos := map[int][][2]float64{
		1: {{.5, .5}},
		2: {{.27, .27}, {.73, .73}},
		3: {{.25, .25}, {.5, .5}, {.75, .75}},
		4: {{.27, .27}, {.73, .27}, {.27, .73}, {.73, .73}},
		5: {{.25, .25}, {.75, .25}, {.5, .5}, {.25, .75}, {.75, .75}},
		6: {{.27, .23}, {.73, .23}, {.27, .5}, {.73, .5}, {.27, .77}, {.73, .77}},
	}[v]
	for _, q := range pos {
		disque(dst, x+q[0]*s, y+q[1]*s, r, points)
	}
}

// spritePion : la case du pion dans la planche (même ordre que les ID).
func spritePion(id int) *ebiten.Image {
	const c = 76
	x, y := (id%10)*c, (id/10)*c
	return imgPions.SubImage(image.Rect(x+2, y+2, x+c-2, y+c-2)).(*ebiten.Image)
}

// pion : le pion de l'encart, de côté t, centré en (x, y).
func pion(dst *ebiten.Image, id int, x, y, t float64, a float32) {
	img := spritePion(id)
	s := t / float64(img.Bounds().Dx())
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(s, s)
	op.GeoM.Translate(x-t/2, y-t/2)
	op.ColorScale.ScaleAlpha(a)
	dst.DrawImage(img, op)
}

// ---------------------------------------------------------------- triangles

// Formes dessinées directement en triangles, sans le rasteriseur
// vectoriel : pour tout ce qui se redessine à chaque image sur la carte
// (anneaux qui pulsent, trajets), bien plus rapide.

func sommet(x, y float64, c color.RGBA) ebiten.Vertex {
	return ebiten.Vertex{DstX: float32(x), DstY: float32(y), SrcX: 1, SrcY: 1,
		ColorR: float32(c.R) / 255, ColorG: float32(c.G) / 255, ColorB: float32(c.B) / 255, ColorA: float32(c.A) / 255}
}

// anneau : cercle d'épaisseur ep, bords adoucis par un fondu d'un pixel.
func anneau(dst *ebiten.Image, x, y, r, ep float64, c color.RGBA) {
	n := max(int(r*0.7), 24)
	var vs []ebiten.Vertex
	var is []uint16
	transp := color.RGBA{}
	rayons := []float64{r - ep/2 - 1, r - ep/2, r + ep/2, r + ep/2 + 1}
	couls := []color.RGBA{transp, c, c, transp}
	for i := 0; i <= n; i++ {
		a := 2 * math.Pi * float64(i) / float64(n)
		ca, sa := math.Cos(a), math.Sin(a)
		for k, rr := range rayons {
			vs = append(vs, sommet(x+rr*ca, y+rr*sa, couls[k]))
		}
		if i > 0 {
			b := uint16(4 * (i - 1))
			for k := uint16(0); k < 3; k++ {
				is = append(is, b+k, b+k+1, b+4+k, b+k+1, b+5+k, b+4+k)
			}
		}
	}
	dst.DrawTriangles(vs, is, blanc, &ebiten.DrawTrianglesOptions{})
}

// polyligne : trait épais le long de points (écran), joints arrondis.
func polyligne(dst *ebiten.Image, pts [][2]float64, ep float64, c color.RGBA) {
	var vs []ebiten.Vertex
	var is []uint16
	for i := 1; i < len(pts); i++ {
		x0, y0, x1, y1 := pts[i-1][0], pts[i-1][1], pts[i][0], pts[i][1]
		l := math.Hypot(x1-x0, y1-y0)
		if l == 0 {
			continue
		}
		nx, ny := -(y1-y0)/l*ep/2, (x1-x0)/l*ep/2
		b := uint16(len(vs))
		vs = append(vs, sommet(x0+nx, y0+ny, c), sommet(x0-nx, y0-ny, c), sommet(x1+nx, y1+ny, c), sommet(x1-nx, y1-ny, c))
		is = append(is, b, b+1, b+2, b+1, b+3, b+2)
		if len(vs) > 60000 {
			break
		}
	}
	dst.DrawTriangles(vs, is, blanc, &ebiten.DrawTrianglesOptions{})
	for _, p := range pts[1 : len(pts)-1] {
		pastille(dst, p[0], p[1], ep/2, c)
	}
}

// pastille : disque plein en triangles.
func pastille(dst *ebiten.Image, x, y, r float64, c color.RGBA) {
	n := max(int(r), 10)
	vs := []ebiten.Vertex{sommet(x, y, c)}
	var is []uint16
	for i := 0; i <= n; i++ {
		a := 2 * math.Pi * float64(i) / float64(n)
		vs = append(vs, sommet(x+r*math.Cos(a), y+r*math.Sin(a), c))
		if i > 0 {
			is = append(is, 0, uint16(i), uint16(i+1))
		}
	}
	dst.DrawTriangles(vs, is, blanc, &ebiten.DrawTrianglesOptions{})
}
