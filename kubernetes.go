package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func GetNodes() (*NodeList, error) {
	var nodeList NodeList

	req, err := http.NewRequest("GET", NodesEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+AuthToken)
	req.Header.Set("Accept", "application/json, */*")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d, response: %s", resp.StatusCode, body)
	}

	if err := json.Unmarshal(body, &nodeList); err != nil {
		return nil, fmt.Errorf("Failed to parse JSON: %v, response: %s", err, body)
	}

	return &nodeList, nil

}

func GetRunningPods() (*PodList, error) {
	var podList PodList

	v := url.Values{}
	v.Set("fieldSelector", "status.phase=Running")

	newPodsEndpoint := fmt.Sprintf("%s?%s", PodsEndpoint, v.Encode())

	req, err := http.NewRequest("GET", newPodsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+AuthToken)
	req.Header.Set("Accept", "application/json, */*")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d, response: %s", resp.StatusCode, body)
	}

	if err := json.Unmarshal(body, &podList); err != nil {
		return nil, fmt.Errorf("Failed to parse JSON: %v, response: %s", err, body)
	}

	return &podList, nil

}

func GetUnschedulePods() ([]*Pod, error) {
	var podList PodList
	unscheduledPods := make([]*Pod, 0)

	v := url.Values{}
	v.Set("fieldSelector", "spec.nodeName=")

	newPodsEndpoint := fmt.Sprintf("%s?%s", PodsEndpoint, v.Encode())

	req, err := http.NewRequest("GET", newPodsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+AuthToken)
	req.Header.Set("Accept", "application/json, */*")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d, response: %s", resp.StatusCode, body)
	}

	err = json.Unmarshal(body, &podList)
	if err != nil {
		return unscheduledPods, fmt.Errorf("Failed to parse JSON: %v, response: %s", err, body)
	}

	for _, pod := range podList.Items {
		if pod.Spec.SchedulerName == SchedulerName {
			unscheduledPods = append(unscheduledPods, &pod)
		}
	}

	return unscheduledPods, nil

}

func Fit(pod *Pod) ([]Node, error) {
	nodeList, err := GetNodes()
	if err != nil {
		return nil, err
	}

	podList, err := GetRunningPods()
	if err != nil {
		return nil, err
	}

	resourceUsage := make(map[string]*ResourceUsage)
	for _, node := range nodeList.Items {

		resourceUsage[node.Metadata.Name] = &ResourceUsage{}

	}

	for _, p := range podList.Items {
		for _, c := range p.Spec.Containers {
			if strings.HasSuffix(c.Resources.Requests["cpu"], "m") {
				milliCores := strings.TrimSuffix(c.Resources.Requests["cpu"], "m")
				cores, err := strconv.Atoi(milliCores)
				if err != nil {
					return nil, err
				}
				ru := resourceUsage[p.Spec.NodeName]
				ru.CPU += cores
			}
		}
	}

	var nodes []Node
	var spaceRequaierd int

	for _, c := range pod.Spec.Containers {
		if strings.HasSuffix(c.Resources.Requests["cpu"], "m") {
			milliCores := strings.TrimSuffix(c.Resources.Requests["cpu"], "m")
			cores, err := strconv.Atoi(milliCores)
			if err != nil {
				return nil, err
			}
			spaceRequaierd += cores
		}
	}

	for _, node := range nodeList.Items {
		cpu := node.Status.Allocatable["cpu"]
		cpuFloat, err := strconv.ParseFloat(cpu, 32)
		if err != nil {
			return nil, err
		}

		freeSpace := (int(cpuFloat*1000) - resourceUsage[node.Metadata.Name].CPU)
		if freeSpace > spaceRequaierd {

			nodes = append(nodes, node)
		}
	}

	return nodes, nil

}

func Bind(pod *Pod, node Node) error {
	binding := Binding{
		ApiVersion: "v1",
		Kind:       "Binding",
		Metadata:   Metadata{Name: pod.Metadata.Name},
		Target: Target{
			ApiVersion: "v1",
			Kind:       "Node",
			Name:       node.Metadata.Name,
		},
	}

	b := make([]byte, 0)
	body := bytes.NewBuffer(b)
	err := json.NewEncoder(body).Encode(binding)
	if err != nil {
		return err
	}

	request := &http.Request{
		Body:          ioutil.NopCloser(body),
		ContentLength: int64(body.Len()),
		Header:        make(http.Header),
		Method:        http.MethodPost,
		URL: &url.URL{
			Host:   ApiHost,
			Path:   fmt.Sprintf(BindingsEndpoint, pod.Metadata.Name),
			Scheme: "https",
		},
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+AuthToken)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 201 {
		return errors.New("Binding: Unexpected HTTP status code" + response.Status)
	} else {
		return errors.New("Pod is successfully binded")
	}

}
