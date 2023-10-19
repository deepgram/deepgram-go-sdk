package manage

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

func (dg *ManageClient) ListProjects() (ProjectResponse, error) {
	client := new(http.Client)
	// path := fmt.Sprintf("%s", dg.Client.Path)
	u := url.URL{Scheme: "https", Host: dg.Client.Host, Path: dg.Client.Path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Panic(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Client.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"User-Agent":    []string{dgAgent},
	}

	var result ProjectResponse
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Panic(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error getting projects: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *ManageClient) GetProject(projectId string) (Project, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s", dg.Client.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Client.Host, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Panic(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Client.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"User-Agent":    []string{dgAgent},
	}

	var result Project
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Panic(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error getting project %s: %s\n", projectId, jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *ManageClient) UpdateProject(projectId string, options ProjectUpdateOptions) (Message, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s", dg.Client.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Client.Host, Path: path}
	jsonStr, err := json.Marshal(options)
	if err != nil {
		log.Panic(err)
		return Message{}, err
	}
	req, err := http.NewRequest("PATCH", u.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		//Handle Error
		log.Panic(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Client.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"User-Agent":    []string{dgAgent},
	}

	var result Message
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Panic(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error updating project %s: %s\n", projectId, jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *ManageClient) DeleteProject(projectId string) (Message, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s", dg.Client.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Client.Host, Path: path}

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Panic(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Client.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"User-Agent":    []string{dgAgent},
	}

	var result Message
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Panic(string(b))
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
