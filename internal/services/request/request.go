package request

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var ErrInvalidStatusCode = errors.New("invalid_status_code")

type IRequestService interface {
	GetDataFromUrls(ctx context.Context, urls []*url.URL) (map[string]string, error)
}

type Service struct {
	maxRequestsCount int
	requestTimeout   time.Duration
}

func NewRequestService(maxRequestsCount int, requestTimeout time.Duration) *Service {
	return &Service{
		maxRequestsCount: maxRequestsCount,
		requestTimeout:   requestTimeout,
	}
}

func (s *Service) DoRequest(ctx context.Context, target *url.URL, dataMap *ResponseData, endCh chan struct{}) error {
	defer func() {
		<-endCh
	}()

	client := http.Client{
		Timeout: s.requestTimeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	/*
		В тексте задания указано "если в процессе обработки хотя бы одного из url произошла ошибка",
		но непонятно что считать ошибкой, поэтому ориентируюсь на статус код и если вернулось не 200, то считаю ошибкой.
	*/
	if resp.StatusCode != http.StatusOK {
		return ErrInvalidStatusCode
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	dataMap.SetValue(target.String(), string(body))

	return nil
}

func (s *Service) GetDataFromUrls(requestCtx context.Context, urls []*url.URL) (map[string]string, error) {
	var groupErr error

	dataMap := NewResponseData(len(urls))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/*
		Не понял можно ли по условиям задания использовать golang.org/x/sync/semaphore, поэтому сделаю "руками".
		Так же круто было бы заюзать golang.org/x/sync/errgroup, но я не уверен, что разрешено.
		Я считаю эти пакеты стандартными, но "из коробки" их нету.
	*/
	semCh := make(chan struct{}, s.maxRequestsCount)
	defer close(semCh)

	wg := sync.WaitGroup{}

	for _, v := range urls {
		select {
		case semCh <- struct{}{}:
			wg.Add(1)

			go func(target *url.URL) {
				defer wg.Done()

				if err := s.DoRequest(ctx, target, dataMap, semCh); err != nil {
					log.Println(err.Error())

					// Если 1 реквест зафэйлился: то останавливаем вообще все
					cancel()

					//
					if groupErr == nil && !errors.Is(err, context.Canceled) {
						groupErr = err
					}
				}
			}(v)
		case <-requestCtx.Done():
			// Если реквест отменен, то отменяем все запросы и больше не посылаем.
			cancel()
			break
		case <-ctx.Done():
			// Если какой-то из-запросов зафэйлился.
			break
		}
	}

	wg.Wait()

	return dataMap.GetData(), groupErr
}
