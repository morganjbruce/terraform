package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
)

type Project struct {
	Name        string `json:"name,omitempty"`
	Key         string `json:"key,omitempty"`
	Description string `json:"description,omitempty"`
	IsPrivate   bool   `json:"is_private,omitempty"`
}

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Update: resourceProjectUpdate,
		Read:   resourceProjectRead,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"is_private": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func newProjectFromResource(d *schema.ResourceData) *Project {
	repo := &Project{
		Name:        d.Get("name").(string),
		Key:         d.Get("key").(string),
		IsPrivate:   d.Get("is_private").(bool),
		Description: d.Get("description").(string),
	}

	return repo
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	project := newProjectFromResource(d)

	bytedata, err := json.Marshal(project)

	if err != nil {
		return err
	}

	_, err = client.Post(fmt.Sprintf("2.0/teams/%s/projects/", d.Get("owner").(string)), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	d.SetId(string(fmt.Sprintf("%s/%s", d.Get("owner").(string), d.Get("key").(string))))

	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	req, _ := client.Get(fmt.Sprintf("2.0/teams/%s/projects/%s",
		d.Get("owner").(string),
		d.Get("key").(string),
	))

	if req.StatusCode == 200 {
		var project Project

		body, readerr := ioutil.ReadAll(req.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &project)
		if decodeerr != nil {
			return decodeerr
		}

		d.Set("name", project.Name)
		d.Set("is_private", project.IsPrivate)
		d.Set("description", project.Description)
		d.Set("key", project.Key)
	}

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	project := newProjectFromResource(d)

	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(project)

	_, err := client.Put(fmt.Sprintf("2.0/teams/%s/projects/%s",
		d.Get("owner").(string),
		d.Get("key").(string),
	), jsonpayload)

	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	_, err := client.Delete(fmt.Sprintf("2.0/teams/%s/projects/%s",
		d.Get("owner").(string),
		d.Get("key").(string),
	))

	return err
}
