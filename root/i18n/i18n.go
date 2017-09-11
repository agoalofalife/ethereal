package i18n

type FillLocale interface {
	EstablishConnection(configuration interface{}) FillLocale
	Add(StorageLocale)
	Merge(map[string]map[string]string, StorageLocale) FillLocale
}

type Locale interface {
	Get(key string) interface{}
}

type StorageLocale struct {
	Structure map[string]map[string]string
}
