package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const baseURL = "https://advokat-sank-peterburg.ru"

type SiteConfig struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	ThemeColor  string `json:"themeColor"`
}

type NavItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Page string `json:"page"`
}

type Contact struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Text2 string `json:"text2"`
	Link  string `json:"link"`
	Icon  string `json:"icon"`
}

type TextItem struct {
	ID              int    `json:"id"`
	Paragraph       string `json:"paragraph"`
	ParagraphBefore string `json:"paragraphBefore"`
	ParagraphAfter  string `json:"paragraphAfter"`
	LinkName        string `json:"linkName"`
	LinkHref        string `json:"linkHref"`
	Target          string `json:"target"`
	Download        bool   `json:"download"`
	Align           string `json:"align"`
	Variant         string `json:"variant"`
	TextDecoration  string `json:"textDecoration"`
}

type Button struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
	Text string `json:"text"`
	Text2 string `json:"text2"`
	Link string `json:"link"`
	Href string `json:"href"`
	Icon string `json:"icon"`
}

type TableItem struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Price  string `json:"price"`
	Price2 string `json:"price2"`
}

type ListItem struct {
	ID         int    `json:"id"`
	Text       string `json:"text"`
	Subheader  string `json:"subheader"`
	NotDivider bool   `json:"notDivider"`
}

type Tab struct {
	ID         int         `json:"id"`
	TabName    string      `json:"tabName"`
	ButtonName string      `json:"buttonName"`
	Text       []TextItem  `json:"text"`
	Table      []TableItem `json:"table"`
}

type AboutData struct {
	Header    string     `json:"header"`
	SubHeader string     `json:"subHeader"`
	Text      []TextItem `json:"text"`
}

type AssistanceData struct {
	Header  string     `json:"header"`
	Text    []TextItem `json:"text"`
	Buttons []Button   `json:"buttons"`
}

type CriminalData struct {
	Header string     `json:"header"`
	Text   []TextItem `json:"text"`
	Tabs   []Tab      `json:"tabs"`
}

type TablePageData struct {
	Header string      `json:"header"`
	Text   []TextItem  `json:"text"`
	Table  []TableItem `json:"table"`
}

type ListPageData struct {
	Header string     `json:"header"`
	List   []ListItem `json:"list"`
}

type ConsultingData struct {
	Text      []TextItem  `json:"text"`
	AlertText []TextItem  `json:"alertText"`
	Table     []TableItem `json:"table"`
}

type ContactsPageData struct {
	Header  string     `json:"header"`
	Text    []TextItem `json:"text"`
	Text2   []TextItem `json:"text2"`
	Buttons []Button   `json:"buttons"`
}

type NotFoundData struct {
	Header string     `json:"header"`
	Text   []TextItem `json:"text"`
}

type UsefulLink struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Breadcrumb struct {
	Name string
	URL  string
}

type AppData struct {
	Site         SiteConfig
	Nav          []NavItem
	Contacts     []Contact
	HeaderIcons  []string
	UsefulLinks  []UsefulLink
	About        AboutData
	Assistance   AssistanceData
	Criminal     CriminalData
	Arbitraj     TablePageData
	Civil        ListPageData
	Family       ListPageData
	House        ListPageData
	Legacy       ListPageData
	Tax          TablePageData
	Consulting   ConsultingData
	ContactsPage ContactsPageData
	NotFound     NotFoundData
}

type PageContext struct {
	AppData
	PageTitle       string
	PageHeader      string
	PageSubHeader   string
	PageDescription string
	ActiveNav       string
	Breadcrumbs     []Breadcrumb
	UseTabs         bool
	Year            int
}

type PageDef struct {
	Page string
	Path string
	Ctx  PageContext
}

var (
	appData   AppData
	templates map[string]*template.Template
	pageDefs  []PageDef
)

func loadJSON(filename string, v interface{}) {
	data, err := os.ReadFile("_data/" + filename)
	if err != nil {
		log.Fatalf("Failed to load %s: %v", filename, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		log.Fatalf("Failed to parse %s: %v", filename, err)
	}
}

func loadData() {
	loadJSON("site.json", &appData.Site)
	loadJSON("nav.json", &appData.Nav)
	loadJSON("contacts.json", &appData.Contacts)
	loadJSON("headerIcons.json", &appData.HeaderIcons)
	loadJSON("usefulLinks.json", &appData.UsefulLinks)
	loadJSON("about.json", &appData.About)
	loadJSON("assistance.json", &appData.Assistance)
	loadJSON("criminal.json", &appData.Criminal)
	loadJSON("arbitraj.json", &appData.Arbitraj)
	loadJSON("civil.json", &appData.Civil)
	loadJSON("family.json", &appData.Family)
	loadJSON("house.json", &appData.House)
	loadJSON("legacy.json", &appData.Legacy)
	loadJSON("tax.json", &appData.Tax)
	loadJSON("consulting.json", &appData.Consulting)
	loadJSON("contactsPage.json", &appData.ContactsPage)
	loadJSON("notFound.json", &appData.NotFound)

	// Generate contact page buttons from shared contacts data
	for _, c := range appData.Contacts {
		if c.Icon == "location" {
			continue
		}
		appData.ContactsPage.Buttons = append(appData.ContactsPage.Buttons, Button{
			Name: c.Name,
			Text: c.Text,
			Href: c.Link,
			Icon: c.Icon,
		})
	}
}

func newFuncMap() template.FuncMap {
	return template.FuncMap{
		"year": func() int { return time.Now().Year() },
		"hasPrefix": strings.HasPrefix,
		"isActiveNav": func(activeNav, page string) bool {
			if activeNav == page {
				return true
			}
			if page == "/assistance/" && strings.HasPrefix(activeNav, "/assistance/") {
				return true
			}
			return false
		},
		"defaultAlign": func(align string) string {
			if align != "" {
				return align
			}
			return "left"
		},
		"defaultAlignCenter": func(align string) string {
			if align != "" {
				return align
			}
			return "center"
		},
		"buttonLink": func(b Button) string {
			if b.Link != "" {
				return b.Link
			}
			return b.Href
		},
		"isExternal": func(b Button) bool {
			return b.Href != "" && b.Link == ""
		},
		"contactByIcon": func(icon string, contacts []Contact) Contact {
			for _, c := range contacts {
				if c.Icon == icon {
					return c
				}
			}
			return Contact{}
		},
		"isFirst": func(i int) bool { return i == 0 },
		"safeURL": func(s string) template.URL { return template.URL(s) },
	}
}

func parseTemplates() {
	templates = make(map[string]*template.Template)

	base := template.Must(
		template.New("").Funcs(newFuncMap()).ParseFiles(
			"templates/base.html",
			"templates/partials/components.html",
		),
	)

	pages := []string{
		"index", "about", "contacts", "404", "assistance",
		"criminal", "arbitraj", "civil", "family", "house",
		"legacy", "tax", "consulting",
	}

	for _, page := range pages {
		t := template.Must(template.Must(base.Clone()).ParseFiles(
			fmt.Sprintf("templates/pages/%s.html", page),
		))
		templates[page] = t
	}
}

func render(w http.ResponseWriter, page string, ctx PageContext) {
	ctx.Year = time.Now().Year()
	ctx.AppData = appData
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates[page].ExecuteTemplate(w, "base", ctx); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func initPageDefs() {
	pageDefs = []PageDef{
		{"index", "index.html", PageContext{
			PageTitle:       "Адвокат в Санкт-Петербурге — Кынтэрец Людмила Николаевна",
			PageDescription: "Адвокат в Санкт-Петербурге. Защита по уголовным делам, арбитраж, гражданские и семейные споры. Опыт более 15 лет.",
			ActiveNav:       "/",
		}},
		{"about", "about/index.html", PageContext{
			PageHeader:      "Кынтэрец Людмила Николаевна",
			PageSubHeader:   "Адвокат",
			PageDescription: "Адвокат Кынтэрец Людмила Николаевна. Более 15 лет практики защиты прав и интересов в Санкт-Петербурге.",
			ActiveNav:       "/about/",
			Breadcrumbs:     []Breadcrumb{{Name: "Кынтэрец Людмила Николаевна"}},
		}},
		{"contacts", "contacts/index.html", PageContext{
			PageHeader:      "Контактная информация",
			PageDescription: "Контакты адвоката Кынтэрец Л.Н. в Санкт-Петербурге: телефон, Telegram, WhatsApp, email. Адрес: 5-я Советская, д. 44.",
			ActiveNav:       "/contacts/",
			Breadcrumbs:     []Breadcrumb{{Name: "Контактная информация"}},
		}},
		{"assistance", "assistance/index.html", PageContext{
			PageHeader:      "Юридическая помощь",
			PageDescription: "Виды юридической помощи адвоката в Санкт-Петербурге: уголовные дела, арбитраж, гражданские и семейные споры, наследство, жилищные и налоговые вопросы.",
			ActiveNav:       "/assistance/",
			Breadcrumbs:     []Breadcrumb{{Name: "Юридическая помощь"}},
		}},
		{"criminal", "assistance/criminal/index.html", PageContext{
			PageHeader:      "Оказание юридической помощи по уголовным делам",
			PageDescription: "Защита по уголовным делам в Санкт-Петербурге. Адвокат Кынтэрец Л.Н. — опыт ведения дел любой сложности и категории.",
			ActiveNav:       "/assistance/",
			UseTabs:         true,
			Breadcrumbs: []Breadcrumb{
				{Name: "Юридическая помощь", URL: "/assistance/"},
				{Name: "Уголовные дела"},
			},
		}},
		{"arbitraj", "assistance/arbitraj/index.html", PageContext{
			PageHeader:      "Оказание юридической помощи по арбитражным делам",
			PageDescription: "Арбитражные споры в Санкт-Петербурге. Представительство в арбитражных судах, защита интересов юридических лиц.",
			ActiveNav:       "/assistance/",
			Breadcrumbs: []Breadcrumb{
				{Name: "Юридическая помощь", URL: "/assistance/"},
				{Name: "Арбитражные дела"},
			},
		}},
		{"civil", "assistance/civil/index.html", PageContext{
			PageHeader:      "Оказание юридической помощи по гражданским делам",
			PageDescription: "Гражданские дела в Санкт-Петербурге. Представительство в судах общей юрисдикции, защита прав и интересов.",
			ActiveNav:       "/assistance/",
			Breadcrumbs: []Breadcrumb{
				{Name: "Юридическая помощь", URL: "/assistance/"},
				{Name: "Гражданские дела"},
			},
		}},
		{"family", "assistance/family/index.html", PageContext{
			PageHeader:      "Оказание юридической помощи по семейным делам",
			PageDescription: "Семейные споры в Санкт-Петербурге. Разводы, раздел имущества, определение места жительства детей.",
			ActiveNav:       "/assistance/",
			Breadcrumbs: []Breadcrumb{
				{Name: "Юридическая помощь", URL: "/assistance/"},
				{Name: "Семейные дела"},
			},
		}},
		{"house", "assistance/house/index.html", PageContext{
			PageHeader:      "Оказание юридической помощи по жилищным делам",
			PageDescription: "Жилищные споры в Санкт-Петербурге. Споры о праве собственности, выселение, приватизация.",
			ActiveNav:       "/assistance/",
			Breadcrumbs: []Breadcrumb{
				{Name: "Юридическая помощь", URL: "/assistance/"},
				{Name: "Жилищные дела"},
			},
		}},
		{"legacy", "assistance/legacy/index.html", PageContext{
			PageHeader:      "Оказание юридической помощи по наследственным делам",
			PageDescription: "Наследственные споры в Санкт-Петербурге. Оспаривание завещаний, восстановление сроков, раздел наследства.",
			ActiveNav:       "/assistance/",
			Breadcrumbs: []Breadcrumb{
				{Name: "Юридическая помощь", URL: "/assistance/"},
				{Name: "Наследственные дела"},
			},
		}},
		{"tax", "assistance/tax/index.html", PageContext{
			PageHeader:      "Оказание юридической помощи по налоговым делам",
			PageDescription: "Налоговые споры в Санкт-Петербурге. Защита интересов в налоговых органах и судах.",
			ActiveNav:       "/assistance/",
			Breadcrumbs: []Breadcrumb{
				{Name: "Юридическая помощь", URL: "/assistance/"},
				{Name: "Налоговые дела"},
			},
		}},
		{"consulting", "assistance/consulting/index.html", PageContext{
			PageHeader:      "Консультирование",
			PageDescription: "Юридические консультации адвоката в Санкт-Петербурге. Устные и письменные консультации по всем отраслям права.",
			ActiveNav:       "/assistance/",
			Breadcrumbs: []Breadcrumb{
				{Name: "Юридическая помощь", URL: "/assistance/"},
				{Name: "Консультирование"},
			},
		}},
		{"404", "404.html", PageContext{
			PageHeader: "Ошибка 404",
		}},
	}
}

func getPageDef(page string) PageContext {
	for _, p := range pageDefs {
		if p.Page == page {
			return p.Ctx
		}
	}
	return PageContext{}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		handle404(w, r)
		return
	}
	render(w, "index", getPageDef("index"))
}

func pageHandler(page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, page, getPageDef(page))
	}
}

func handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	render(w, "404", getPageDef("404"))
}

func renderToFile(page string, ctx PageContext, outPath string) error {
	ctx.Year = time.Now().Year()
	ctx.AppData = appData

	var buf bytes.Buffer
	if err := templates[page].ExecuteTemplate(&buf, "base", ctx); err != nil {
		return fmt.Errorf("template %s: %w", page, err)
	}

	dir := filepath.Dir(outPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	return os.WriteFile(outPath, buf.Bytes(), 0644)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}

func sitemapXML() []byte {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	buf.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	for _, p := range pageDefs {
		if p.Page == "404" {
			continue
		}
		loc := baseURL + "/"
		if p.Path != "index.html" {
			dir := strings.TrimSuffix(p.Path, "index.html")
			loc = baseURL + "/" + dir
		}
		buf.WriteString(fmt.Sprintf("  <url><loc>%s</loc></url>\n", loc))
	}
	buf.WriteString("</urlset>\n")
	return buf.Bytes()
}

func generateSitemap(outDir string) {
	outPath := filepath.Join(outDir, "sitemap.xml")
	if err := os.WriteFile(outPath, sitemapXML(), 0644); err != nil {
		log.Fatalf("Sitemap: %v", err)
	}
	log.Printf("  sitemap.xml")
}

func buildStatic(outDir string) {
	// Remove old output
	os.RemoveAll(outDir)

	// Render all pages from single source of truth
	for _, p := range pageDefs {
		outPath := filepath.Join(outDir, p.Path)
		if err := renderToFile(p.Page, p.Ctx, outPath); err != nil {
			log.Fatalf("Build error: %v", err)
		}
		log.Printf("  %s", p.Path)
	}

	// Generate sitemap.xml
	generateSitemap(outDir)

	// Copy static assets
	staticDirs := []string{"css", "js", "images"}
	for _, dir := range staticDirs {
		if _, err := os.Stat(dir); err == nil {
			if err := copyDir(dir, filepath.Join(outDir, dir)); err != nil {
				log.Fatalf("Copy %s: %v", dir, err)
			}
			log.Printf("  %s/", dir)
		}
	}

	staticFiles := []string{"favicon.ico", "robots.txt", "agreement.docx", "manifest.json"}
	for _, f := range staticFiles {
		if _, err := os.Stat(f); err == nil {
			if err := copyFile(f, filepath.Join(outDir, f)); err != nil {
				log.Fatalf("Copy %s: %v", f, err)
			}
			log.Printf("  %s", f)
		}
	}

	log.Printf("Build complete → %s/", outDir)
}

func main() {
	loadData()
	parseTemplates()
	initPageDefs()

	// Static build mode
	if len(os.Args) > 1 && os.Args[1] == "--build" {
		outDir := "dist"
		if len(os.Args) > 2 {
			outDir = os.Args[2]
		}
		log.Printf("Building static site...")
		buildStatic(outDir)
		return
	}

	// Dev server mode
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("GET /about/", pageHandler("about"))
	mux.HandleFunc("GET /contacts/", pageHandler("contacts"))
	mux.HandleFunc("GET /assistance/", pageHandler("assistance"))
	mux.HandleFunc("GET /assistance/criminal/", pageHandler("criminal"))
	mux.HandleFunc("GET /assistance/arbitraj/", pageHandler("arbitraj"))
	mux.HandleFunc("GET /assistance/civil/", pageHandler("civil"))
	mux.HandleFunc("GET /assistance/family/", pageHandler("family"))
	mux.HandleFunc("GET /assistance/house/", pageHandler("house"))
	mux.HandleFunc("GET /assistance/legacy/", pageHandler("legacy"))
	mux.HandleFunc("GET /assistance/tax/", pageHandler("tax"))
	mux.HandleFunc("GET /assistance/consulting/", pageHandler("consulting"))

	mux.Handle("GET /css/", http.FileServer(http.Dir(".")))
	mux.Handle("GET /js/", http.FileServer(http.Dir(".")))
	mux.Handle("GET /images/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})
	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "robots.txt")
	})
	mux.HandleFunc("GET /agreement.docx", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "agreement.docx")
	})
	mux.HandleFunc("GET /manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "manifest.json")
	})
	mux.HandleFunc("GET /sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Write(sitemapXML())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
