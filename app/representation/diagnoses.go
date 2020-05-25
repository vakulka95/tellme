package representation

var Diagnoses = map[string]string{
	"panic":                         "Панічний розлад",
	"health_concerns":               "Тривога за здоров'я",
	"unknown_fear":                  "Генералізований тривожний розлад(страх невідомості і невизначеності)",
	"eating_disorders":              "Розлади харчової поведінки(булімія, анорексія, переїдання)",
	"experience_of_loss":            "Важке переживання втрати",
	"ptsr":                          "ПТСР(травма)",
	"depression":                    "Депресія",
	"phobias":                       "Фобії",
	"obsessive_compulsive_disorder": "Обсесивно-компульсивний розлад(нав'язливі думки та дії)",
	"other":                         "Інше",
}

var DiagnosesStripped = map[string]string{
	"panic":                         "Панічний розлад",
	"health_concerns":               "Тривога за здоров'я",
	"unknown_fear":                  "Генералізований тривожний розлад",
	"eating_disorders":              "Розлади харчової поведінки",
	"experience_of_loss":            "Важке переживання втрати",
	"ptsr":                          "ПТСР",
	"depression":                    "Депресія",
	"phobias":                       "Фобії",
	"obsessive_compulsive_disorder": "Обсесивно-компульсивний розлад",
	"other":                         "Інше",
}

type DiagnosesOptions struct {
	Translated string
	Selected   bool
}

func GenerateDiagnosesOptions(sel []string) map[string]*DiagnosesOptions {
	res := map[string]*DiagnosesOptions{
		"panic":                         {Translated: "Панічний розлад"},
		"health_concerns":               {Translated: "Тривога за здоров'я"},
		"unknown_fear":                  {Translated: "Генералізований тривожний розлад"},
		"eating_disorders":              {Translated: "Розлади харчової поведінки"},
		"experience_of_loss":            {Translated: "Важке переживання втрати"},
		"ptsr":                          {Translated: "ПТСР"},
		"depression":                    {Translated: "Депресія"},
		"phobias":                       {Translated: "Фобії"},
		"obsessive_compulsive_disorder": {Translated: "Обсесивно-компульсивний розлад"},
		"other":                         {Translated: "Інше"},
	}
	for _, v := range sel {
		if opt, ok := res[v]; ok {
			opt.Selected = true
		}
	}
	return res
}
