package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
)

var (
	ErrLimitReached = errors.New("api limit reached")
)

func (c *Caller) GetPlanContents(ctx context.Context, projectPlanIDs map[int][]int, outDir *string) (map[int]map[int][]int, error) {
	const endpoint = "projects/%d/plans/%d/plan_contents.json"
	const name = "plan_contents"
	var multierr error
	planContentIDs := make(map[int]map[int][]int, 0)

	var wg sync.WaitGroup
	var mu sync.Mutex
	for projectID, planIDs := range projectPlanIDs {
		for _, planID := range planIDs {
			wg.Add(1)
			endpoint := fmt.Sprintf(endpoint, projectID, planID)
			go func(endpoint string, projectID, planID int) {
				defer wg.Done()

				out, ids, err := c.getAllPages(ctx, endpoint, name)
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					multierr = multierror.Append(multierr, err)
				}

				if _, ok := planContentIDs[projectID]; !ok {
					planContentIDs[projectID] = make(map[int][]int)
				}
				planContentIDs[projectID][planID] = append(planContentIDs[projectID][planID], ids...)

				if outDir == nil {
					return
				}

				for page, data := range out {
					path := filepath.Join(*outDir, fmt.Sprintf("%s-%d-%d-%d", name, projectID, planID, page))
					err := os.WriteFile(path, data, 0644)
					if err != nil {
						multierr = multierror.Append(multierr, err)
						return
					}
				}
			}(endpoint, projectID, planID)
		}

	}
	wg.Wait()

	return planContentIDs, multierr
}

func (c *Caller) GetPlans(ctx context.Context, projectIDs []int, outDir *string) (map[int][]int, error) {
	const endpoint = "projects/%d/plans.json"
	const name = "plans"
	var multierr error
	planIDs := make(map[int][]int)

	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, id := range projectIDs {
		wg.Add(1)
		endpoint := fmt.Sprintf(endpoint, id)
		go func(endpoint string, id int) {
			defer wg.Done()

			out, ids, err := c.getAllPages(ctx, endpoint, name)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				multierr = multierror.Append(multierr, err)
			}
			planIDs[id] = append(planIDs[id], ids...)

			if outDir == nil {
				return
			}

			for page, data := range out {
				path := filepath.Join(*outDir, fmt.Sprintf("%s-%d-%d", name, id, page))
				err := os.WriteFile(path, data, 0644)
				if err != nil {
					multierr = multierror.Append(multierr, err)
					return
				}
			}
		}(endpoint, id)
	}
	wg.Wait()

	return planIDs, multierr
}

func (c *Caller) GetTestCases(ctx context.Context, projectSuiteSectionIDs map[int]map[int][]int, outDir *string) error {
	const endpoint = "projects/%d/suites/%d/suite_sections/%d/steps.json"
	const name = "steps"
	var multierr error

	var wg sync.WaitGroup
	var mu sync.Mutex
	const maxGoroutines = 5
	ch := make(chan struct{}, maxGoroutines)
	for projectID, suiteSectionIDs := range projectSuiteSectionIDs {
		for suiteID, sectionIDs := range suiteSectionIDs {
			for _, sectionID := range sectionIDs {

				wg.Add(1)
				endpoint := fmt.Sprintf(endpoint, projectID, suiteID, sectionID)
				ch <- struct{}{}
				go func(endpoint string, projectID, suiteID, sectionID int) {
					defer wg.Done()
					defer func() { <-ch }()

					out, _, err := c.getAllPages(ctx, endpoint, name)
					mu.Lock()
					defer mu.Unlock()
					if err != nil {
						multierr = multierror.Append(multierr, err)
					}

					if outDir == nil {
						return
					}

					for page, data := range out {
						path := filepath.Join(*outDir, fmt.Sprintf("%s-%d-%d-%d-%d", name, projectID, suiteID, sectionID, page))
						err := os.WriteFile(path, data, 0644)
						if err != nil {
							multierr = multierror.Append(multierr, err)
							return
						}
					}
				}(endpoint, projectID, suiteID, sectionID)
			}
		}

	}
	wg.Wait()

	return multierr
}
func (c *Caller) GetRequirements(ctx context.Context, projectDocumentIDs map[int][]int, outDir *string) (map[int]map[int][]int, error) {
	const endpoint = "projects/%d/requirement_documents/%d/requirements.json"
	const name = "requirements"
	var multierr error
	requirementIDs := make(map[int]map[int][]int, 0)

	var wg sync.WaitGroup
	var mu sync.Mutex
	for projectID, docIDs := range projectDocumentIDs {
		for _, docID := range docIDs {
			wg.Add(1)
			endpoint := fmt.Sprintf(endpoint, projectID, docID)
			go func(endpoint string, projectID, docID int) {
				defer wg.Done()

				out, ids, err := c.getAllPages(ctx, endpoint, name)
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					multierr = multierror.Append(multierr, err)
				}

				if _, ok := requirementIDs[projectID]; !ok {
					requirementIDs[projectID] = make(map[int][]int)
				}
				requirementIDs[projectID][docID] = append(requirementIDs[projectID][docID], ids...)

				if outDir == nil {
					return
				}
				for page, data := range out {
					path := filepath.Join(*outDir, fmt.Sprintf("%s-%d-%d-%d", name, projectID, docID, page))
					err := os.WriteFile(path, data, 0644)
					if err != nil {
						multierr = multierror.Append(multierr, err)
						return
					}
				}
			}(endpoint, projectID, docID)
		}

	}
	wg.Wait()

	return requirementIDs, multierr
}
func (c *Caller) GetSuiteSections(ctx context.Context, projectSuiteIDs map[int][]int, outDir *string) (map[int]map[int][]int, error) {
	const endpoint = "projects/%d/suites/%d/suite_sections.json"
	const name = "suite_sections"
	var multierr error
	suiteSectionIDs := make(map[int]map[int][]int)

	var wg sync.WaitGroup
	var mu sync.Mutex
	for projectID, suiteIDs := range projectSuiteIDs {
		for _, suiteID := range suiteIDs {
			wg.Add(1)
			endpoint := fmt.Sprintf(endpoint, projectID, suiteID)
			go func(endpoint string, projectID, suiteID int) {
				defer wg.Done()

				out, ids, err := c.getAllPages(ctx, endpoint, name)
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					multierr = multierror.Append(multierr, err)
				}
				if _, ok := suiteSectionIDs[projectID]; !ok {
					suiteSectionIDs[projectID] = make(map[int][]int)
				}
				suiteSectionIDs[projectID][suiteID] = append(suiteSectionIDs[projectID][suiteID], ids...)

				if outDir == nil {
					return
				}
				for page, data := range out {
					path := filepath.Join(*outDir, fmt.Sprintf("%s-%d-%d-%d", name, projectID, suiteID, page))
					err := os.WriteFile(path, data, 0644)
					if err != nil {
						multierr = multierror.Append(multierr, err)
						return
					}
				}
			}(endpoint, projectID, suiteID)
		}

	}
	wg.Wait()

	return suiteSectionIDs, multierr
}

func (c *Caller) GetRequirementDocuments(ctx context.Context, projectIDs []int, outDir *string) (map[int][]int, error) {
	const endpoint = "projects/%d/requirement_documents.json"
	const name = "requirement_documents"
	var multierr error
	rdIDs := make(map[int][]int)

	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, id := range projectIDs {
		wg.Add(1)
		endpoint := fmt.Sprintf(endpoint, id)
		go func(endpoint string, id int) {
			defer wg.Done()

			out, ids, err := c.getAllPages(ctx, endpoint, name)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				multierr = multierror.Append(multierr, err)
			}
			rdIDs[id] = append(rdIDs[id], ids...)
			if outDir == nil {
				return
			}
			for page, data := range out {
				path := filepath.Join(*outDir, fmt.Sprintf("%s-%d-%d", name, id, page))
				err := os.WriteFile(path, data, 0644)
				if err != nil {
					multierr = multierror.Append(multierr, err)
					return
				}
			}
		}(endpoint, id)
	}
	wg.Wait()

	return rdIDs, multierr
}

func (c *Caller) GetSuites(ctx context.Context, projectIDs []int, outDir *string) (map[int][]int, error) {
	const endpoint = "projects/%d/suites.json"
	const name = "suites"
	var multierr error

	suiteIDs := make(map[int][]int, 0)

	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, id := range projectIDs {
		wg.Add(1)
		endpoint := fmt.Sprintf(endpoint, id)
		go func(endpoint string, id int) {
			defer wg.Done()

			out, ids, err := c.getAllPages(ctx, endpoint, name)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				multierr = multierror.Append(multierr, err)
			}
			suiteIDs[id] = append(suiteIDs[id], ids...)

			if outDir == nil {
				return
			}
			for page, data := range out {
				path := filepath.Join(*outDir, fmt.Sprintf("%s-%d-%d", name, id, page))
				err := os.WriteFile(path, data, 0644)
				if err != nil {
					multierr = multierror.Append(multierr, err)
					return
				}
			}
		}(endpoint, id)
	}
	wg.Wait()

	return suiteIDs, multierr
}

func (c *Caller) GetProjects(ctx context.Context, outDir *string) ([]int, error) {
	const endpoint = "projects.json"
	const name = "projects"
	var res map[int][]byte
	var ids []int
	var err error

	res, ids, err = c.getAllPages(ctx, endpoint, name)

	if outDir == nil {
		return ids, err
	}
	for page, data := range res {
		path := filepath.Join(*outDir, fmt.Sprintf("%s-%d", name, page))
		err := os.WriteFile(path, data, 0644)
		if err != nil {
			return nil, err
		}
	}

	return ids, err
}

func (c *Caller) getAllPages(ctx context.Context, path, name string) (map[int][]byte, []int, error) {
	fullPath, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, nil, err
	}

	result, err := c.genericCallWithBackoff(ctx, fullPath, name, 1)
	if err != nil {
		return nil, nil, err
	}

	out := make(map[int][]byte)
	var a map[string]interface{}
	err = json.Unmarshal(result, &a)
	if err != nil {
		return nil, nil, err
	}
	pages := a["pagination"].(map[string]interface{})["total_pages"].(float64)
	data := a[name].([]interface{})
	out[1], err = json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil, nil, err
	}

	var ids []int
	for _, entry := range data {
		id := entry.(map[string]interface{})["id"].(float64)
		ids = append(ids, int(id))
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	const maxGoroutines = 10
	ch := make(chan struct{}, maxGoroutines)
	for i := 2; i <= int(pages); i++ {
		wg.Add(1)
		ch <- struct{}{}
		go func(page int) {
			defer wg.Done()
			defer func() { <-ch }()

			result, err := c.genericCallWithBackoff(ctx, fullPath, name, page)
			if err != nil {
				log.Error().Err(err).Send()
				return
			}
			mu.Lock()
			defer mu.Unlock()
			var a map[string]interface{}
			err = json.Unmarshal(result, &a)
			if err != nil {
				log.Error().Err(err).Send()
				return
			}
			data := a[name].([]interface{})
			out[page], err = json.MarshalIndent(data, "", "\t")
			if err != nil {
				log.Error().Err(err).Send()
			}
			for _, entry := range data {
				id := entry.(map[string]interface{})["id"].(float64)
				ids = append(ids, int(id))
			}
		}(i)
	}
	wg.Wait()

	return out, ids, nil
}

func (c *Caller) genericCallWithBackoff(ctx context.Context, fullPath, name string, page int) (res []byte, err error) {
	err = backoff.Retry(func() error {
		res, err = c.genericCall(ctx, fullPath, name, page)
		if errors.Is(err, ErrLimitReached) {
			log.Info().Msg("api limit reached sleeping for 2 minutes")
			return err
		} else if err != nil {
			log.Error().Err(err).Send()
			return backoff.Permanent(err)
		}
		return nil
	}, backoff.WithContext(backoff.NewConstantBackOff(2*time.Minute), ctx))
	return
}

func (c *Caller) genericCall(ctx context.Context, fullPath, name string, page int) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullPath, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.email, c.key)
	q := req.URL.Query()
	q.Set("page", fmt.Sprint(page))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		return nil, ErrLimitReached
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status code %d and was expecting 200", resp.StatusCode)
	}
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}
