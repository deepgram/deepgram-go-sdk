package deepgram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Project struct {
	ProjectId string `json:"project_id"`
	Name      string `json:"name,omitempty"`
	Company   string `json:"company,omitempty"`
}

type ProjectResponse struct {
	Projects []Project `json:"projects"`
}

type ProjectUpdateOptions struct {
	Name    string `json:"name,omitempty"`
	Company string `json:"company,omitempty"`
}

func (dg *Client) ListProjects() (ProjectResponse, error) {
	client := new(http.Client)
	// path := fmt.Sprintf("%s", dg.Path)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: dg.Path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result ProjectResponse
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error getting projects: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) GetProject(projectId string) (Project, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s", dg.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result Project
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error getting project %s: %s\n", projectId, jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) UpdateProject(projectId string, options ProjectUpdateOptions) (Message, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s", dg.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}
	jsonStr, err := json.Marshal(options)
	if err != nil {
		log.Fatal(err)
		return Message{}, err
	}
	req, err := http.NewRequest("PATCH", u.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result Message
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error updating project %s: %s\n", projectId, jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) DeleteProject(projectId string) (Message, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s", dg.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result Message
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error deleting project %s: %s\n", projectId, jsonErr.Error())
		return result, jsonErr
	} else {
		return Message{
			Message: "Project Successfully Deleted",
		}, nil
	}
}
