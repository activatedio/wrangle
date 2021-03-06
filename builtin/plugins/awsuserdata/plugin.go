package awsuserdata

import (
	"sync"

	"text/template"

	"os"

	user2 "os/user"

	"fmt"
	"io/ioutil"

	"strings"

	"errors"

	"github.com/activatedio/wrangle/plugin"
)

type AwsUserDataPluginConfig struct {
	Nameservers []string
}

type AwsUserDataPlugin struct {
	ConfigLock sync.Mutex
	Config     *AwsUserDataPluginConfig
}

func (self *AwsUserDataPlugin) GetConfig() interface{} {

	self.ConfigLock.Lock()
	defer self.ConfigLock.Unlock()

	if self.Config == nil {
		self.Config = &AwsUserDataPluginConfig{}
	}

	return self.Config
}

// TODO - Constant for now, will change in the future
const USER_DATA_TEMPLATE = `#!/bin/bash

set -e

{{if .Config.Nameservers }}
{{range .Config.Nameservers}}
echo "nameserver {{.}}" >> /etc/resolvconf/resolv.conf.d/head
{{end}}
resolvconf -u
{{end}}

apt-get update
apt-get install -y python

username={{ .Username }}
userid={{ .UID }}
ssh_public_key="{{ .SshPublicKey }}"

adduser $username --uid $userid --shell /bin/bash --disabled-password
usermod -a -G sudo $username
ssh_dir=/home/$username/.ssh
ssh_authorized_keys=$ssh_dir/authorized_keys2
mkdir $ssh_dir
chmod 700 $ssh_dir
echo $ssh_public_key > $ssh_authorized_keys
chmod 600 $ssh_authorized_keys

chown -R $username:$username $ssh_dir

sudoers_file=/etc/sudoers.d/01-sudo-nopasswd
echo '%sudo ALL=(ALL) NOPASSWD:ALL' > $sudoers_file
chmod 400 $sudoers_file

`

type userDataData struct {
	Username     string
	SshPublicKey string
	Config       *AwsUserDataPluginConfig
	UID          string
}

func (self *AwsUserDataPlugin) Filter(c plugin.Context) error {

	f, err := os.Create(".user-data.sh")
	defer f.Close()

	if err != nil {
		return err
	}

	t, err := template.New("user-data").Parse(USER_DATA_TEMPLATE)

	if err != nil {
		return err
	}

	user, err := user2.Current()

	if err != nil {
		return err
	}

	k, err := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_rsa.pub", user.HomeDir))

	if err != nil {
		return err
	}

	// TODO This is a tempoaray hack until we have a better config option
	uid := os.Getenv("AWS_USER_DATA_UID")

	if uid == "" {
		return errors.New("Please specify userid via AWS_USER_DATA_UID environment variable.")
	}

	d := &userDataData{
		Username:     user.Username,
		UID:          uid,
		SshPublicKey: strings.TrimSuffix(string(k), "\n"),
		Config:       self.Config,
	}

	err = t.Execute(f, d)

	if err != nil {
		return err
	}

	return c.Next()
}
