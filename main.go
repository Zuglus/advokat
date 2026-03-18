package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

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

var (
	appData   AppData
	templates map[string]*template.Template
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
}

func newFuncMap() template.FuncMap {
	return template.FuncMap{
		"year": func() int { return time.Now().Year() },
		"hasPrefix": strings.HasPrefix,
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

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		handle404(w, r)
		return
	}
	render(w, "index", PageContext{
		PageTitle:       "Адвокат в Санкт-Петербурге — Кынтэрец Людмила Николаевна",
		PageDescription: "Адвокат в Санкт-Петербурге. Защита по уголовным делам, арбитраж, гражданские и семейные споры. Опыт более 15 лет.",
		ActiveNav:       "/",
	})
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	render(w, "about", PageContext{
		PageHeader:      "Кынтэрец Людмила Николаевна",
		PageSubHeader:   "Адвокат",
		PageDescription: "Адвокат Кынтэрец Людмила Николаевна. Более 15 лет практики защиты прав и интересов в Санкт-Петербурге.",
		ActiveNav:       "/about/",
		Breadcrumbs:     []Breadcrumb{{Name: "Кынтэрец Людмила Николаевна"}},
	})
}

func handleContacts(w http.ResponseWriter, r *http.Request) {
	render(w, "contacts", PageContext{
		PageHeader:      "Контактная информация",
		PageDescription: "Контакты адвоката Кынтэрец Л.Н. в Санкт-Петербурге: телефон, Telegram, WhatsApp, email. Адрес: 5-я Советская, д. 44.",
		ActiveNav:       "/contacts/",
		Breadcrumbs:     []Breadcrumb{{Name: "Контактная информация"}},
	})
}

func handleAssistance(w http.ResponseWriter, r *http.Request) {
	render(w, "assistance", PageContext{
		PageHeader:      "Юридическая помощь",
		PageDescription: "Виды юридической помощи адвоката в Санкт-Петербурге: уголовные дела, арбитраж, гражданские и семейные споры, наследство, жилищные и налоговые вопросы.",
		ActiveNav:       "/assistance/",
		Breadcrumbs:     []Breadcrumb{{Name: "Юридическая помощь"}},
	})
}

func handleCriminal(w http.ResponseWriter, r *http.Request) {
	render(w, "criminal", PageContext{
		PageHeader:      "Оказание юридической помощи по уголовным делам",
		PageDescription: "Защита по уголовным делам в Санкт-Петербурге. Адвокат Кынтэрец Л.Н. — опыт ведения дел любой сложности и категории.",
		ActiveNav:       "/assistance/",
		UseTabs:         true,
		Breadcrumbs: []Breadcrumb{
			{Name: "Юридическая помощь", URL: "/assistance/"},
			{Name: "Уголовные дела"},
		},
	})
}

func handleArbitraj(w http.ResponseWriter, r *http.Request) {
	render(w, "arbitraj", PageContext{
		PageHeader:      "Оказание юридической помощи по арбитражным делам",
		PageDescription: "Арбитражные споры в Санкт-Петербурге. Представительство в арбитражных судах, защита интересов юридических лиц.",
		ActiveNav:       "/assistance/",
		Breadcrumbs: []Breadcrumb{
			{Name: "Юридическая помощь", URL: "/assistance/"},
			{Name: "Арбитражные дела"},
		},
	})
}

func handleCivil(w http.ResponseWriter, r *http.Request) {
	render(w, "civil", PageContext{
		PageHeader:      "Оказание юридической помощи по гражданским делам",
		PageDescription: "Гражданские дела в Санкт-Петербурге. Представительство в судах общей юрисдикции, защита прав и интересов.",
		ActiveNav:       "/assistance/",
		Breadcrumbs: []Breadcrumb{
			{Name: "Юридическая помощь", URL: "/assistance/"},
			{Name: "Гражданские дела"},
		},
	})
}

func handleFamily(w http.ResponseWriter, r *http.Request) {
	render(w, "family", PageContext{
		PageHeader:      "Оказание юридической помощи по семейным делам",
		PageDescription: "Семейные споры в Санкт-Петербурге. Разводы, раздел имущества, определение места жительства детей.",
		ActiveNav:       "/assistance/",
		Breadcrumbs: []Breadcrumb{
			{Name: "Юридическая помощь", URL: "/assistance/"},
			{Name: "Семейные дела"},
		},
	})
}

func handleHouse(w http.ResponseWriter, r *http.Request) {
	render(w, "house", PageContext{
		PageHeader:      "Оказание юридической помощи по жилищным делам",
		PageDescription: "Жилищные споры в Санкт-Петербурге. Споры о праве собственности, выселение, приватизация.",
		ActiveNav:       "/assistance/",
		Breadcrumbs: []Breadcrumb{
			{Name: "Юридическая помощь", URL: "/assistance/"},
			{Name: "Жилищные дела"},
		},
	})
}

func handleLegacy(w http.ResponseWriter, r *http.Request) {
	render(w, "legacy", PageContext{
		PageHeader:      "Оказание юридической помощи по наследственным делам",
		PageDescription: "Наследственные споры в Санкт-Петербурге. Оспаривание завещаний, восстановление сроков, раздел наследства.",
		ActiveNav:       "/assistance/",
		Breadcrumbs: []Breadcrumb{
			{Name: "Юридическая помощь", URL: "/assistance/"},
			{Name: "Наследственные дела"},
		},
	})
}

func handleTax(w http.ResponseWriter, r *http.Request) {
	render(w, "tax", PageContext{
		PageHeader:      "Оказание юридической помощи по налоговым делам",
		PageDescription: "Налоговые споры в Санкт-Петербурге. Защита интересов в налоговых органах и судах.",
		ActiveNav:       "/assistance/",
		Breadcrumbs: []Breadcrumb{
			{Name: "Юридическая помощь", URL: "/assistance/"},
			{Name: "Налоговые дела"},
		},
	})
}

func handleConsulting(w http.ResponseWriter, r *http.Request) {
	render(w, "consulting", PageContext{
		PageDescription: "Юридические консультации адвоката в Санкт-Петербурге. Устные и письменные консультации по всем отраслям права.",
		ActiveNav:       "/assistance/",
		Breadcrumbs: []Breadcrumb{
			{Name: "Юридическая помощь", URL: "/assistance/"},
			{Name: "Консультирование"},
		},
	})
}

func handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	render(w, "404", PageContext{
		PageHeader: "Ошибка 404",
	})
}

func main() {
	loadData()
	parseTemplates()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("GET /about/", handleAbout)
	mux.HandleFunc("GET /contacts/", handleContacts)
	mux.HandleFunc("GET /assistance/", handleAssistance)
	mux.HandleFunc("GET /assistance/criminal/", handleCriminal)
	mux.HandleFunc("GET /assistance/arbitraj/", handleArbitraj)
	mux.HandleFunc("GET /assistance/civil/", handleCivil)
	mux.HandleFunc("GET /assistance/family/", handleFamily)
	mux.HandleFunc("GET /assistance/house/", handleHouse)
	mux.HandleFunc("GET /assistance/legacy/", handleLegacy)
	mux.HandleFunc("GET /assistance/tax/", handleTax)
	mux.HandleFunc("GET /assistance/consulting/", handleConsulting)

	mux.Handle("GET /css/", http.FileServer(http.Dir(".")))
	mux.Handle("GET /js/", http.FileServer(http.Dir(".")))
	mux.Handle("GET /images/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})
	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "robots.txt")
	})
	mux.HandleFunc("GET /agreement.doc", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "agreement.doc")
	})
	mux.HandleFunc("GET /manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "manifest.json")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
