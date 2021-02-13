package request

import "sync"

// На небольшом кол-ве данных map с mutex быстрее, чем sync.Map.
type ResponseData struct {
	mu   sync.Mutex
	data map[string]string
}

func NewResponseData(size int) *ResponseData {
	return &ResponseData{
		data: make(map[string]string, size),
	}
}

// По условиям задания не ясно что делать, если 1 и тот же url передан дважды, поэтому перезаписываю без проверки.
func (d *ResponseData) SetValue(key, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.data[key] = value
}

func (d *ResponseData) GetData() map[string]string {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.data
}
