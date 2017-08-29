package ethereal

import (
	"github.com/qor/i18n"
	"strings"
)

type Locale interface {
	Fill()
	Merge()
}

type graphQLI18n struct {
	structure map[string]map[string]string
}

func I18nGraphQL() (graphQL graphQLI18n) {
	graphQL = graphQLI18n{
		structure: map[string]map[string]string{
			"en-US": map[string]string{
				"graphQL.User.Description":  "List of users of your application.",
				"graphQL.UserType.id":       "ID your user.",
				"graphQL.UserType.email":    "Email your user.",
				"graphQL.UserType.name":     "Name your user.",
				"graphQL.UserType.password": "Hashed password your user.",
				"graphQL.UserType.role":     "Concrete role your user.",
			},
			"ru-RU": map[string]string{
				"graphQL.User.Description": "Список пользователей вашего приложения.",
				"graphQL.UserType.id":       "ID вашего пользователя.",
				"graphQL.UserType.email":    "Email вашего пользователя.",
				"graphQL.UserType.name":     "Имя вашего пользователя.",
				"graphQL.UserType.password": "Захэшированный пароль",
				"graphQL.UserType.role":     "Роль которой принадлежить пользователь",
			},
		},
	}
	return
}

/**
/ Function merge structure graph i!8n
*/
func (schema graphQLI18n) Merge(merge map[string]map[string]string) graphQLI18n {
	for locale, mapValues := range merge {
		for key, value := range mapValues {
			if _, exist := schema.structure[locale]; !exist {
				schema.structure[locale] = map[string]string{
					key: value,
				}
			} else {
				schema.structure[locale][key] = value
			}
		}
	}
	return schema
}

func (schema graphQLI18n) Fill() {
	for locale, mapValues := range schema.structure {
		for key, value := range mapValues {
			app.I18n.SaveTranslation(&i18n.Translation{Key: key, Locale: locale, Value: value})
		}
	}
}

func mapLanguage() map[string]string {
	return map[string]string{
		"en-US": "en-US",
		"ru-RU": "ru-RU",
	}
}

// temp functions determines type locale
func parserLocale(header []string) string {
	splitSign := strings.Split(header[0], ";")
	for _, multiLocal := range splitSign {
		local := strings.Split(multiLocal, ",")
		for _, possible := range local {
			if _, isExist := mapLanguage()[possible]; isExist {
				return mapLanguage()[possible]
			}
		}
	}
	return ""
}
